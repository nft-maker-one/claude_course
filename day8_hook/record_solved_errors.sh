#!/usr/bin/env bash
# record_solved_errors.sh
# Fires on Stop. Reads pending tool failures logged by log_error.sh and appends
# a formatted reason+solution entry to .claude/hooks/dev_errors.md.
#
# Strategy:
#   1. If a transcript is available, ask Claude to infer how each error was resolved.
#   2. If no transcript (or Claude writes nothing), fall back to a structured entry
#      from the raw error data alone — never leave errors as "unprocessed".

set -euo pipefail

INPUT=$(cat)

SESSION_ID=$(echo "$INPUT" | jq -r '.session_id // "unknown"')
TRANSCRIPT_PATH=$(echo "$INPUT" | jq -r '.transcript_path // ""')
STOP_HOOK_ACTIVE=$(echo "$INPUT" | jq -r '.stop_hook_active // false')

if [ "$STOP_HOOK_ACTIVE" = "true" ]; then
  exit 0
fi

PENDING_DIR="${TMPDIR:-/tmp}/claude-pending-errors"
PENDING_FILE="$PENDING_DIR/${SESSION_ID}.jsonl"

if [ ! -s "$PENDING_FILE" ]; then
  exit 0
fi

PROJECT_DIR="${CLAUDE_PROJECT_DIR:-$(pwd)}"
OUT_DIR="$PROJECT_DIR/.claude/hooks"
OUT_FILE="$OUT_DIR/dev_errors.md"
mkdir -p "$OUT_DIR"

PENDING_JSON=$(jq -s '.' "$PENDING_FILE")
BEFORE_SIZE=$(wc -c < "$OUT_FILE" 2>/dev/null || echo 0)

# --- Attempt 1: transcript-aware summarisation ---
if [ -n "$TRANSCRIPT_PATH" ] && [ -f "$TRANSCRIPT_PATH" ]; then
  PROMPT=$(cat <<EOF
You are summarizing development errors that occurred during a Claude Code session
and how they were resolved, for a running log at $OUT_FILE.

Below is a JSON array of tool failures that occurred in this session:

$PENDING_JSON

The full conversation transcript is at: $TRANSCRIPT_PATH
Read it to determine how each error was diagnosed and resolved.

For EVERY error in the array — whether resolved or not — append a Markdown section
to $OUT_FILE using EXACTLY this template (use the Edit tool to append;
do NOT overwrite existing content):

## <ISO-8601 timestamp> — <one-line error summary>

- **Tool:** <tool name>
- **Context:** <what Claude was trying to do>
- **Error:** <verbatim or trimmed error message>
- **Reason:** <root cause, 1–3 sentences>
- **Solution:** <concrete fix, or "Unresolved — session ended without a fix" if not fixed>

---

Rules:
- Write one entry per distinct error (combine duplicates with the same root cause).
- If the file does not yet exist, start it with "# Dev Errors Log" then a blank line.
- Be concise. No preamble, no closing remarks — just the entries.
EOF
)

  printf '%s' "$PROMPT" | claude -p \
    --permission-mode acceptEdits \
    --allowedTools "Read,Write,Edit,Bash" \
    >/dev/null 2>&1 || true
fi

AFTER_SIZE=$(wc -c < "$OUT_FILE" 2>/dev/null || echo 0)

# --- Attempt 2: fallback — write structured entries from raw error data alone ---
if [ "$AFTER_SIZE" -le "$BEFORE_SIZE" ]; then
  {
    if [ ! -s "$OUT_FILE" ] 2>/dev/null; then
      echo "# Dev Errors Log"
      echo ""
    fi

    echo "$PENDING_JSON" | jq -r '.[] |
      "## \(.timestamp) — \(.tool) error in \(.cwd | split("/") | last)\n" +
      "\n" +
      "- **Tool:** \(.tool)\n" +
      "- **Context:** \(.tool_input | if type == "object" then (.command // .description // tostring) else tostring end | .[0:200])\n" +
      "- **Error:** `\(.error | gsub("\n"; " ") | .[0:300])`\n" +
      "- **Reason:** Recorded automatically — transcript unavailable or session ended before resolution was confirmed.\n" +
      "- **Solution:** Review the error message above and rerun with the correct invocation.\n" +
      "\n---\n"
    '
  } >> "$OUT_FILE"
fi

rm -f "$PENDING_FILE"
exit 0
