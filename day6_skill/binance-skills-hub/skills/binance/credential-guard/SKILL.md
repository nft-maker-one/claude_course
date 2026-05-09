---
name: credential-guard
description: |
  Security skill for managing Binance API credentials safely.
  Enforces loading all secrets from a local .env file — never from conversation input.
  Masks any credential that appears in tool output before showing it to the user.
  Apply this skill whenever any other Binance skill requires authentication.
  Auto-trigger on: API key setup, profile creation, any authenticated API call.
metadata:
  version: 1.0.0
  author: binance-skills-hub
license: MIT
---

# Credential Guard Skill

This skill governs how **all** API credentials are handled across every Binance skill.
It must be applied **before** any authenticated operation.

---

## Rules (MANDATORY — never override)

### 1. Never accept credentials from conversation

- Do **not** ask the user to paste API keys, secrets, or tokens into the chat.
- Do **not** read credentials from user messages, even if the user provides them.
- If the user pastes a credential in chat, respond:
  > "Please store that in your `.env` file. I'll load it from there — never from the conversation."
- Exception: if `.env` does not exist yet, guide the user to create it (see Setup section below).

### 2. Always load from `.env`

Locate `.env` by searching upward from the current working directory:

```bash
# Load .env into the current shell
set -a
source "$(pwd)/.env"
set +a
```

Required variables:

| Variable | Used by |
|----------|---------|
| `BINANCE_API_KEY` | binance-cli profile, authenticated REST calls |
| `BINANCE_API_SECRET` | binance-cli profile, signed requests |
| `X_SQUARE_OPENAPI_KEY` | square-post skill |

Fail fast if a required variable is missing:

```bash
python3 - <<'EOF'
import os, sys
required = ["BINANCE_API_KEY", "BINANCE_API_SECRET"]
missing = [k for k in required if not os.environ.get(k)]
if missing:
    print(f"ERROR: Missing in .env: {', '.join(missing)}", file=sys.stderr)
    sys.exit(1)
print("Credentials loaded ✅")
EOF
```

### 3. Mask all credentials in output

Before displaying any command output to the user, replace full credential values with masked versions.

**Masking rule**: show first 5 characters + `...` + last 4 characters.

```python
def mask(value: str) -> str:
    if len(value) <= 9:
        return "***"
    return value[:5] + "..." + value[-4:]
```

Examples:

| Raw value | Masked display |
|-----------|---------------|
| `jw3fiP44L6sc3dUIuGmDHkqdX8XNgZt9Epg4rNs8iDb2FGToPt40BMPeLB3ManiW` | `jw3fi...aniW` |
| `b7e1623212c44dc1a5e7aa14fe56f2b5` | `b7e16...2b5` |

Apply masking to:
- Any shell command shown to the user (replace literal key values with masked form)
- Any CLI output that echoes back credentials
- Any log or error messages containing key material

### 4. Never write credentials to files other than `.env`

- Do not write API keys into SKILL.md, CLAUDE.md, config files, or any source-tracked file.
- `.env` must be listed in `.gitignore` — verify this before any git operation.

### 5. Redact before showing commands

When constructing a shell command that embeds a credential, **always** show the user a redacted version while running the real command internally:

**Show user:**
```bash
binance-cli profile create --name default --env prod \
  --api-key "jw3fi...aniW" \
  --api-secret "rZEC2...uxQ" \
  --select --force
```

**Actually execute** (using env var, never literal value in shown text):
```bash
set -a && source .env && set +a
binance-cli profile create --name default --env prod \
  --api-key "$BINANCE_API_KEY" \
  --api-secret "$BINANCE_API_SECRET" \
  --select --force
```

---

## Setup Workflow

Run this when credentials are not yet configured.

### Step 1 — Create `.env` if missing

Tell the user:
> "Please create a `.env` file in `<project-root>/` with the following format, then tell me when it's ready:"

```env
# Binance API
BINANCE_API_KEY=your_api_key_here
BINANCE_API_SECRET=your_api_secret_here

# Binance Square OpenAPI
X_SQUARE_OPENAPI_KEY=your_square_key_here
```

### Step 2 — Verify `.gitignore`

```bash
grep -q "^\.env$" .gitignore || echo ".env" >> .gitignore
echo ".gitignore status: .env is excluded ✅"
```

### Step 3 — Load and validate

```bash
set -a && source .env && set +a
python3 -c "
import os
keys = ['BINANCE_API_KEY', 'BINANCE_API_SECRET']
for k in keys:
    v = os.environ.get(k, '')
    masked = v[:5]+'...'+v[-4:] if len(v) > 9 else '***'
    status = '✅' if v else '❌ MISSING'
    print(f'{k}: {masked} {status}')
"
```

### Step 4 — Configure binance-cli profile

```bash
set -a && source .env && set +a
binance-cli profile create \
  --name default \
  --env prod \
  --api-key "$BINANCE_API_KEY" \
  --api-secret "$BINANCE_API_SECRET" \
  --select --force
```

---

## Integration with Other Skills

Apply this skill **before** any of the following:

| Skill | Credentials required |
|-------|---------------------|
| `binance` | `BINANCE_API_KEY`, `BINANCE_API_SECRET` |
| `p2p` (authenticated) | `BINANCE_API_KEY`, `BINANCE_API_SECRET` |
| `fiat` (SAPI endpoints) | `BINANCE_API_KEY`, `BINANCE_API_SECRET` |
| `square-post` | `X_SQUARE_OPENAPI_KEY` |
| `payment` | Handled by `payment_skill.py` (reads its own config) |
| `onchain-pay` | RSA key path (store path in `.env` as `ONCHAIN_PAY_PRIVATE_KEY_PATH`) |

### Canonical credential-loading snippet (reuse across skills)

```bash
ENV_FILE="$(pwd)/.env"
if [ ! -f "$ENV_FILE" ]; then
  echo "ERROR: .env not found at $ENV_FILE" >&2
  exit 1
fi
set -a && source "$ENV_FILE" && set +a
```

---

## Security Checklist

Before any authenticated operation:

- [ ] `.env` file exists
- [ ] `.env` is in `.gitignore`
- [ ] Credentials loaded via `source .env`, not from conversation
- [ ] Command shown to user uses masked values only
- [ ] No credential appears in any displayed output or log
