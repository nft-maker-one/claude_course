"""
Run once to register a Helius webhook for all addresses in smart_money.txt.
Usage:
    WEBHOOK_URL=https://xxxx.ngrok.io python register_webhook.py
"""

import os
import sys
import httpx
from dotenv import load_dotenv

load_dotenv()

HELIUS_KEY = os.environ["HELIUS_KEY"]
SUPABASE_PROJECT_URL = os.environ.get("SUPABASE_PROJECT_URL", "").rstrip("/")

# Default to the Supabase Edge Function URL; override with WEBHOOK_URL env var if needed.
_default_url = f"{SUPABASE_PROJECT_URL}/functions/v1/smart-money-hook" if SUPABASE_PROJECT_URL else ""
WEBHOOK_URL = os.environ.get("WEBHOOK_URL") or _default_url

if not WEBHOOK_URL:
    print("ERROR: set SUPABASE_PROJECT_URL in .env or WEBHOOK_URL env var explicitly.")
    sys.exit(1)

addresses = [
    line.strip()
    for line in open("smart_money.txt")
    if line.strip()
]

payload = {
    "webhookURL": WEBHOOK_URL.rstrip("/"),
    "accountAddresses": addresses,
    "transactionTypes": ["SWAP"],
    "webhookType": "enhanced",
}

resp = httpx.post(
    f"https://api.helius.xyz/v0/webhooks?api-key={HELIUS_KEY}",
    json=payload,
    timeout=15,
)
resp.raise_for_status()
data = resp.json()

print("Webhook registered successfully.")
print(f"  webhookID : {data.get('webhookID')}")
print(f"  webhook   : {data.get('webhookURL')}")
print(f"  addresses : {len(addresses)}")
print()
print("Save the webhookID — you need it to update or delete this webhook.")
