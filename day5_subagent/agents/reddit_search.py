#!/usr/bin/env python3
"""
Standalone Reddit search — no API key required.
Uses Reddit public JSON endpoints directly (free, no auth).

Usage:
    python3 reddit_search.py "PENGU solana" --days 14 --depth deep
    python3 reddit_search.py "dogwifhat" --subreddit solana --days 7
    python3 reddit_search.py "bitcoin" --json
"""

import argparse
import json
import sys
import time
from datetime import datetime, timezone, timedelta
from urllib.parse import quote_plus

import requests

HEADERS = {
    "User-Agent": "Mozilla/5.0 (compatible; token-analyze-agent/1.0; +research)"
}

DEPTH_MAP = {"quick": 10, "default": 25, "deep": 50}


def search_reddit(
    query: str,
    days: int = 30,
    depth: str = "default",
    subreddit: str | None = None,
) -> list[dict]:
    """
    Search Reddit posts via the public JSON API.
    Returns a list of post dicts sorted by score descending.
    """
    limit = DEPTH_MAP.get(depth, 25)
    cutoff = datetime.now(tz=timezone.utc) - timedelta(days=days)
    cutoff_ts = cutoff.timestamp()

    results: list[dict] = []

    if subreddit:
        urls = [
            f"https://www.reddit.com/r/{subreddit}/search.json"
            f"?q={quote_plus(query)}&restrict_sr=1&sort=relevance&limit={limit}&t=month"
        ]
    else:
        urls = [
            f"https://www.reddit.com/search.json"
            f"?q={quote_plus(query)}&sort=relevance&limit={limit}&t=month",
            f"https://www.reddit.com/search.json"
            f"?q={quote_plus(query)}&sort=hot&limit={limit}&t=month",
        ]

    seen: set[str] = set()

    for url in urls:
        try:
            resp = requests.get(url, headers=HEADERS, timeout=12)
            if resp.status_code == 429:
                print("Rate limited by Reddit — waiting 5s", file=sys.stderr)
                time.sleep(5)
                resp = requests.get(url, headers=HEADERS, timeout=12)
            resp.raise_for_status()
            posts = resp.json()["data"]["children"]
        except Exception as e:
            print(f"Reddit request failed ({url}): {e}", file=sys.stderr)
            continue

        for post in posts:
            d = post["data"]
            post_id = d.get("id", "")
            if post_id in seen:
                continue
            created = d.get("created_utc", 0)
            if created < cutoff_ts:
                continue
            if d.get("stickied"):
                continue

            link = d.get("url", "")
            if link.startswith("/r/"):
                link = "https://www.reddit.com" + link

            seen.add(post_id)
            results.append(
                {
                    "id": post_id,
                    "title": d.get("title", ""),
                    "subreddit": d.get("subreddit", ""),
                    "score": d.get("score", 0),
                    "num_comments": d.get("num_comments", 0),
                    "date": datetime.fromtimestamp(created, tz=timezone.utc).strftime("%Y-%m-%d"),
                    "url": link,
                    "selftext": d.get("selftext", "")[:300].strip(),
                }
            )

        time.sleep(0.5)  # be polite between requests

    results.sort(key=lambda x: x["score"], reverse=True)
    return results[:limit]


def main() -> None:
    p = argparse.ArgumentParser(description="Reddit search (no API key)")
    p.add_argument("topic", nargs="+", help="Search query")
    p.add_argument("--days", type=int, default=30, help="Days to look back (default: 30)")
    p.add_argument(
        "--depth",
        choices=["quick", "default", "deep"],
        default="default",
        help="Result depth: quick=10, default=25, deep=50",
    )
    p.add_argument("--subreddit", default=None, help="Restrict to a subreddit")
    p.add_argument("--json", action="store_true", help="Output raw JSON")
    args = p.parse_args()

    query = " ".join(args.topic)
    now = datetime.now(tz=timezone.utc)
    from_date = (now - timedelta(days=args.days)).strftime("%Y-%m-%d")
    to_date = now.strftime("%Y-%m-%d")

    print(f"Searching Reddit: '{query}' ({from_date} → {to_date})", file=sys.stderr)

    results = search_reddit(query, days=args.days, depth=args.depth, subreddit=args.subreddit)

    if args.json:
        print(json.dumps(results, indent=2, ensure_ascii=False))
        return

    if not results:
        print("No results found.")
        return

    print(f"\n{'='*60}")
    print(f"Reddit results for: {query}")
    print(f"Period: {from_date} → {to_date}  |  Found: {len(results)}")
    print(f"{'='*60}\n")

    for item in results:
        print(f"[{item['id']}] {item['title']}")
        print(
            f"  r/{item['subreddit']} | ↑{item['score']} | "
            f"{item['num_comments']} comments | {item['date']}"
        )
        print(f"  {item['url']}")
        if item["selftext"]:
            print(f"  Preview: {item['selftext'][:200]}...")
        print()


if __name__ == "__main__":
    main()
