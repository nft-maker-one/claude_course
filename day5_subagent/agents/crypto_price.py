#!/usr/bin/env python3
"""
CoinGecko price fetcher — no API key required.
Returns price (USD) and 24h change for one or more tokens.

Usage:
    python3 crypto_price.py pudgy-penguins
    python3 crypto_price.py bitcoin ethereum solana
    python3 crypto_price.py dogwifhat --json
"""

import argparse
import json
import sys
from urllib.parse import quote

import requests

HEADERS = {
    "User-Agent": "Mozilla/5.0 (compatible; token-analyze-agent/1.0; +research)"
}
BASE = "https://api.coingecko.com/api/v3"


def get_prices(coin_ids: list[str]) -> list[dict]:
    """Fetch current price + 24h change for each CoinGecko coin ID."""
    ids = ",".join(coin_ids)
    url = (
        f"{BASE}/simple/price"
        f"?ids={ids}&vs_currencies=usd&include_24hr_change=true&include_market_cap=true&include_24hr_vol=true"
    )
    try:
        resp = requests.get(url, headers=HEADERS, timeout=12)
        resp.raise_for_status()
        data = resp.json()
    except Exception as e:
        print(f"CoinGecko request failed: {e}", file=sys.stderr)
        return []

    return [
        {
            "coin": coin_id,
            "price_usd": info.get("usd"),
            "change_24h": round(info.get("usd_24h_change", 0) or 0, 2),
            "market_cap_usd": info.get("usd_market_cap"),
            "volume_24h_usd": info.get("usd_24h_vol"),
        }
        for coin_id, info in data.items()
    ]


def resolve_id(ticker_or_name: str) -> str | None:
    """Resolve a ticker/name to a CoinGecko ID. Returns best match or None."""
    url = f"{BASE}/search?query={quote(ticker_or_name)}"
    try:
        resp = requests.get(url, headers=HEADERS, timeout=12)
        resp.raise_for_status()
        coins = resp.json().get("coins", [])
        return coins[0]["id"] if coins else None
    except Exception as e:
        print(f"CoinGecko search failed: {e}", file=sys.stderr)
        return None


def main() -> None:
    p = argparse.ArgumentParser(description="CoinGecko price lookup (no API key)")
    p.add_argument("coins", nargs="+", help="CoinGecko IDs or tickers (e.g. pudgy-penguins, BONK)")
    p.add_argument("--resolve", action="store_true", help="Auto-resolve tickers to CoinGecko IDs")
    p.add_argument("--json", action="store_true", help="Output raw JSON")
    args = p.parse_args()

    coin_ids = args.coins
    if args.resolve:
        resolved = []
        for c in coin_ids:
            cg_id = resolve_id(c)
            if cg_id:
                resolved.append(cg_id)
                print(f"Resolved '{c}' → '{cg_id}'", file=sys.stderr)
            else:
                print(f"Could not resolve '{c}' — skipping", file=sys.stderr)
        coin_ids = resolved

    if not coin_ids:
        print("No valid coin IDs to fetch.", file=sys.stderr)
        sys.exit(1)

    results = get_prices(coin_ids)

    if args.json:
        print(json.dumps(results, indent=2))
        return

    print(f"\n{'='*55}")
    print(f"{'Coin':<25} {'Price (USD)':>12} {'24h Change':>12}")
    print(f"{'='*55}")
    for r in results:
        change = r["change_24h"]
        sign = "+" if change >= 0 else ""
        print(f"{r['coin']:<25} ${r['price_usd']:>11.6f} {sign}{change:>10.2f}%")
    print(f"{'='*55}\n")


if __name__ == "__main__":
    main()
