import { createClient } from "https://esm.sh/@supabase/supabase-js@2";

const STABLECOINS = new Set([
  "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v", // USDC
  "Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB", // USDT
]);

const MONITORED = new Set([
  "CyaE1VxvBrahnPWkqm5VsdCvyS2QmNht2UFrKJHga54o",
  "6S8GezkxYUfZy9JPtYnanbcZTMB87Wjt1qx3c6ELajKC",
  "4BdKaxN8G6ka4GYtQQWk4G4dZRUTX2vQH9GcXdBREFUk",
  "Bi4rd5FH5bYEN8scZ7wevxNZyNmKHdaBcvewdPFxYdLt",
  "DxM1hfY8FQ8dNGrucuJzhJcF8KRbjk8WBwrgKvQ9spPv",
]);

const SOL_MINT = "So11111111111111111111111111111111111111112";

async function getSolPrice(): Promise<number> {
  const apiKey = Deno.env.get("JUPITER_API_KEY");
  if (!apiKey) return 0;
  try {
    const resp = await fetch(
      `https://api.jup.ag/price/v3?ids=${SOL_MINT}`,
      {
        headers: { "x-api-key": apiKey },
        signal: AbortSignal.timeout(5000),
      },
    );
    if (!resp.ok) return 0;
    const json = await resp.json();
    const entry = json?.data?.[SOL_MINT];
    if (!entry || entry.confidenceLevel === "low") return 0;
    return parseFloat(entry.price) || 0;
  } catch {
    return 0;
  }
}

function computeValueUsd(swap: Record<string, any>, solPrice: number): number | null {
  const lamports = swap?.nativeInput?.amount;
  if (lamports && solPrice) {
    return Math.round((lamports / 1e9) * solPrice * 10000) / 10000;
  }
  for (const input of swap?.tokenInputs ?? []) {
    if (STABLECOINS.has(input.mint)) return input.tokenAmount ?? null;
  }
  return null;
}

Deno.serve(async (req) => {
  if (req.method !== "POST") {
    return new Response("Method not allowed", { status: 405 });
  }

  const supabase = createClient(
    Deno.env.get("SUPABASE_URL")!,
    Deno.env.get("SUPABASE_SERVICE_ROLE_KEY")!,
  );

  let body: unknown;
  try {
    body = await req.json();
  } catch {
    return new Response("Invalid JSON", { status: 400 });
  }

  const txList: Record<string, any>[] = Array.isArray(body) ? body : [body];
  const solPrice = await getSolPrice();
  const rows: Record<string, unknown>[] = [];

  for (const tx of txList) {
    console.log("[debug] tx.type:", tx.type, "sig:", tx.signature?.slice(0, 16));

    if (tx.type !== "SWAP") {
      console.log("[debug] skipped: type is not SWAP");
      continue;
    }
    const swap = tx.events?.swap ?? {};
    const valueUsd = computeValueUsd(swap, solPrice);

    const transfers = tx.tokenTransfers ?? [];
    console.log("[debug] tokenTransfers count:", transfers.length);
    transfers.forEach((t: Record<string, any>, i: number) => {
      console.log(`[debug]   transfer[${i}] toUserAccount=${t.toUserAccount} mint=${t.mint} amount=${t.tokenAmount}`);
    });

    for (const transfer of transfers) {
      const buyer: string = transfer.toUserAccount ?? "";
      if (!MONITORED.has(buyer)) continue;
      const tokenAmt = parseFloat(transfer.tokenAmount);
      if (!tokenAmt) continue;
      rows.push({
        address: buyer,
        token: transfer.mint ?? null,
        token_amt: tokenAmt,
        value_usd: valueUsd,
      });
    }
  }

  if (rows.length > 0) {
    const { error } = await supabase.from("smart_money").insert(rows);
    if (error) {
      console.error("Supabase insert failed:", error.message);
      return new Response(JSON.stringify({ error: error.message }), {
        status: 500,
        headers: { "Content-Type": "application/json" },
      });
    }
  }

  console.log(`Processed ${txList.length} txs, inserted ${rows.length} buys`);
  return new Response(JSON.stringify({ inserted: rows.length }), {
    headers: { "Content-Type": "application/json" },
  });
});
