#!/usr/bin/env python3
"""
Twitter search via TweeterPy — focused data-fetch script for agents.

Usage:
    python3 twitter_search.py --query "$PENGU solana" --filter Top --total 30 \
        --auth-token YOUR_TOKEN

    python3 twitter_search.py --query "dogwifhat" --filter Latest --total 20 \
        --session /path/to/session.json
"""

import sys
import os
import json
import argparse

# Bundled TweeterPy lives at .claude/agents/lib/ — resolved relative to this file
_LIB_DIR = os.path.join(os.path.dirname(os.path.abspath(__file__)), "lib")
sys.path.insert(0, _LIB_DIR)


def _load_client(auth_token: str | None, session_path: str | None, proxy: str | None):
    from tweeterpy import TweeterPy  # type: ignore

    proxies = {"http": proxy, "https": proxy} if proxy else None
    client = TweeterPy(proxies=proxies, log_level="ERROR")

    if auth_token:
        client.generate_session(auth_token=auth_token)
    elif session_path:
        # Accept both a TweeterPy session file and a browser-exported cookie JSON
        if session_path.endswith(".json"):
            with open(session_path) as f:
                cookies = json.load(f)
            # Browser-exported format: list of {name, value, ...}
            if isinstance(cookies, list):
                token = next(
                    (c["value"] for c in cookies if c.get("name") == "auth_token"), None
                )
                if token:
                    client.generate_session(auth_token=token)
                else:
                    raise ValueError("auth_token cookie not found in session file")
            else:
                client.load_session(path=session_path)
        else:
            client.load_session(path=session_path)
    else:
        raise ValueError("Provide --auth-token or --session")

    return client


def search(
    query: str,
    search_filter: str = "Top",
    total: int = 20,
    auth_token: str | None = None,
    session_path: str | None = None,
    proxy: str | None = None,
) -> list[dict]:
    """Return a list of tweet dicts matching the query."""
    from tweeterpy.util import find_nested_key  # type: ignore

    client = _load_client(auth_token, session_path, proxy)
    result = client.search(query, search_filter=search_filter, total=total)
    tweets = find_nested_key(result.get("data", []), "full_text")
    rate = result.get("api_rate_limit", {})

    return {
        "query": query,
        "filter": search_filter,
        "count": len(tweets),
        "tweets": tweets,
        "rate_limit": {
            "remaining": rate.get("remaining_requests_count"),
            "total": rate.get("total_limit"),
            "resets_after": rate.get("resets_after"),
        },
    }


def main() -> None:
    p = argparse.ArgumentParser(description="Twitter search via TweeterPy")
    p.add_argument("--query", required=True, help="Search query string")
    p.add_argument(
        "--filter",
        default="Top",
        choices=["Top", "Latest", "People", "Photos", "Videos"],
        help="Search filter (default: Top)",
    )
    p.add_argument("--total", type=int, default=20, help="Max results (default: 20)")
    p.add_argument("--auth-token", help="Twitter auth_token cookie value")
    p.add_argument(
        "--session",
        help="Path to TweeterPy session file or browser-exported cookie JSON",
    )
    p.add_argument("--proxy", help="Proxy URL, e.g. http://host:port")
    p.add_argument("--json", action="store_true", help="Output raw JSON")
    args = p.parse_args()

    try:
        data = search(
            query=args.query,
            search_filter=args.filter,
            total=args.total,
            auth_token=args.auth_token,
            session_path=args.session,
            proxy=args.proxy,
        )
    except Exception as e:
        print(f"ERROR: {e}", file=sys.stderr)
        sys.exit(1)

    if args.json:
        print(json.dumps(data, indent=2, ensure_ascii=False, default=str))
        return

    print(f"\n=== Twitter Search: '{data['query']}' [{data['filter']}] — {data['count']} results ===\n")
    for i, tweet in enumerate(data["tweets"], 1):
        print(f"  {i:>3}. {tweet[:200]}")
    rl = data["rate_limit"]
    if rl.get("remaining") is not None:
        print(f"\n[Rate Limit] {rl['remaining']}/{rl['total']} remaining | resets: {rl['resets_after']}")


if __name__ == "__main__":
    main()
