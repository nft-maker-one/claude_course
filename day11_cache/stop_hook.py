import os, httpx
from dotenv import load_dotenv
load_dotenv()
wid = ""
r = httpx.delete(f'https://api.helius.xyz/v0/webhooks/{wid}?api-key={os.environ["HELIUS_KEY"]}')
print('deleted' if r.is_success else r.text)