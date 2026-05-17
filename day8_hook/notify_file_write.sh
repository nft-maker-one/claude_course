#!/usr/bin/env bash
# Reads the PostToolUse JSON payload from stdin and POSTs file info to the local server.
set -euo pipefail

payload=$(cat)

file_path=$(echo "$payload" | jq -r '.tool_input.file_path // empty')
content=$(echo "$payload" | jq -r '.tool_input.content // empty')
tool_name=$(echo "$payload" | jq -r '.tool_name // "Write"')
timestamp=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build the JSON body
body=$(jq -n \
  --arg tool "$tool_name" \
  --arg path "$file_path" \
  --arg ts "$timestamp" \
  --argjson size "${#content}" \
  '{tool: $tool, file_path: $path, timestamp: $ts, content_length: $size}')

curl -s -o /dev/null \
  -X POST http://localhost:8080/claude/create \
  -H "Content-Type: application/json" \
  -d "$body" || true
