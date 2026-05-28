import { createClient } from "https://esm.sh/@supabase/supabase-js@2";

const STABLECOINS = new Set([
  "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v", // USDC
  "Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB", // USDT
]);

const MONITORED = new Set([
  "3BEaKXn847KD7VGkK2sfcrdTSD51DWXWJW6tqubHjU4g",
  "FdHEV7sRsHJch6LmnEBtzLtU28zK3Dogb9tbhXLNyyGT",
  "FdRZFLAAP1kgGjTVeKgdD8pKwqT8freLPNo5Ur8u14hJ",
  "BtHnhjh8Ta2vGZkrL76AmXMur1ADo3tpRo2kQMe2BCSU",
  "DNio5ketGczp2dsB3jJCm3XNsRETaAM87q3YVMPA5F9z",
  "9dW7y6yWSHuk2HCS2keR14h2NyhKPKnz5CTMC9SMdX45",
  "A2MwjTFz4jzT1mY4xrqkwm1vAbZDrqnA6QJoyTAU8Djw",
  "5CVjp216wUnqUuzPkWHWRiA8ovbfiE3MbnxADuehUotK",
  "DK6sTE5yz4xFT4eravZX16HKLvkKnfynYirZAtXRA9jQ",
  "HXbwX1oGic2A5skRxVrsrQDrbAiq21u4JEzvKaFTdJBm",
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
    if (tx.type !== "SWAP") continue;
    const swap = tx.events?.swap ?? {};
    const valueUsd = computeValueUsd(swap, solPrice);

    for (const transfer of tx.tokenTransfers ?? []) {
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
