#!/usr/bin/env bash
# log_error.sh
# Captures developer errors during a Claude Code session into a per-session
# pending-errors file. The Stop hook later turns these into reason+solution
# entries in dev_errors.md.
#
# Fires on two distinct events:
#   1. PostToolUseFailure — the tool itself failed (timeout, missing binary,
#      permission denied, MCP error, etc).
#   2. PostToolUse — the tool succeeded as a tool, but what it *ran* errored.
#      Mainly: Bash commands that exited non-zero (failed tests, Python
#      tracebacks, compile errors, etc). Claude Code reports these as
#      "tool success" because Bash ran fine — we inspect the output
#      to spot the real failure.

set -euo pipefail

INPUT=$(cat)

EVENT=$(echo "$INPUT" | jq -r '.hook_event_name // "unknown"')
SESSION_ID=$(echo "$INPUT" | jq -r '.session_id // "unknown"')
CWD=$(echo "$INPUT" | jq -r '.cwd // empty')
TOOL_NAME=$(echo "$INPUT" | jq -r '.tool_name // "unknown"')
TOOL_INPUT=$(echo "$INPUT" | jq -c '.tool_input // {}')
TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Decide whether this event represents an error worth logging.
ERROR_MSG=""

if [ "$EVENT" = "PostToolUseFailure" ]; then
  # Tool-level failure — always log.
  ERROR_MSG=$(echo "$INPUT" | jq -r '
    .tool_response.error
    // .error
    // (.tool_response | tostring)
    // "tool failure (no message captured)"
  ')

elif [ "$EVENT" = "PostToolUse" ]; then
  # Tool succeeded as a tool — but did what it ran actually work?
  # Only inspect Bash; Edit/Write etc. surface failures via PostToolUseFailure.
  if [ "$TOOL_NAME" != "Bash" ]; then
    exit 0
  fi

  EXIT_CODE=$(echo "$INPUT" | jq -r '
    .tool_response.exitCode
    // .tool_response.exit_code
    // .tool_response.returncode
    // 0
  ')
  STDERR=$(echo "$INPUT" | jq -r '.tool_response.stderr // ""')
  STDOUT=$(echo "$INPUT" | jq -r '.tool_response.stdout // ""')

  if [ "$EXIT_CODE" != "0" ] && [ "$EXIT_CODE" != "null" ]; then
    ERROR_MSG=$(printf 'exit code: %s\nstderr:\n%s\nstdout (tail):\n%s' \
      "$EXIT_CODE" \
      "$(printf '%s' "$STDERR" | tail -c 4000)" \
      "$(printf '%s' "$STDOUT" | tail -c 2000)")
  else
    # Some scripts print tracebacks but still exit 0. Heuristic: stderr
    # mentions a traceback-like pattern.
    if printf '%s' "$STDERR" | grep -qE 'Traceback \(most recent call last\)|Error:|Exception:|FATAL|PANIC'; then
      ERROR_MSG=$(printf 'exit code: 0 but stderr contains error markers\nstderr:\n%s' \
        "$(printf '%s' "$STDERR" | tail -c 4000)")
    fi
  fi
fi

if [ -z "$ERROR_MSG" ]; then
  exit 0
fi

PENDING_DIR="${TMPDIR:-/tmp}/claude-pending-errors"
mkdir -p "$PENDING_DIR"
PENDING_FILE="$PENDING_DIR/${SESSION_ID}.jsonl"

jq -nc \
  --arg ts "$TIMESTAMP" \
  --arg event "$EVENT" \
  --arg cwd "$CWD" \
  --arg tool "$TOOL_NAME" \
  --argjson input "$TOOL_INPUT" \
  --arg error "$ERROR_MSG" \
  '{timestamp: $ts, event: $event, cwd: $cwd, tool: $tool, tool_input: $input, error: $error}' \
  >> "$PENDING_FILE"

exit 0