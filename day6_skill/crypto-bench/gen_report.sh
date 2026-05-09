#!/usr/bin/env bash
# Runs Go unit tests + benchmarks and writes a Markdown report.
set -euo pipefail

REPORT="${1:-bench_report.md}"
BENCH_OUT="$(mktemp)"
TEST_OUT="$(mktemp)"
trap 'rm -f "$BENCH_OUT" "$TEST_OUT"' EXIT

TIMESTAMP="$(date '+%Y-%m-%d %H:%M:%S')"
GOARCH="$(go env GOARCH)"
GOOS="$(go env GOOS)"
GOVERSION="$(go version | awk '{print $3}')"
CPU="$(grep -m1 'model name' /proc/cpuinfo 2>/dev/null | cut -d: -f2 | xargs || echo 'unknown')"

# ─── unit tests ──────────────────────────────────────────────────────────────
echo ">>> Running unit tests..."
go test -v -count=1 -bench='^$' ./... 2>&1 | tee "$TEST_OUT" || true
if grep -q '^FAIL' "$TEST_OUT"; then
  TEST_STATUS="FAIL"
else
  TEST_STATUS="PASS"
fi

# ─── benchmarks ──────────────────────────────────────────────────────────────
echo ">>> Running benchmarks (benchtime=3s)..."
go test -bench=. -benchmem -benchtime=3s -count=1 ./... 2>&1 | tee "$BENCH_OUT"

# ─── report ──────────────────────────────────────────────────────────────────
{
cat <<HEADER
# Crypto Benchmark Report

| Field | Value |
|-------|-------|
| Generated | ${TIMESTAMP} |
| Go version | ${GOVERSION} |
| OS / Arch | ${GOOS} / ${GOARCH} |
| CPU | ${CPU} |
| Unit tests | ${TEST_STATUS} |

---

## Unit Test Results

\`\`\`
HEADER

grep -E '^(--- (PASS|FAIL)|FAIL|ok )' "$TEST_OUT" || grep -E '(PASS|FAIL)' "$TEST_OUT" | head -40

cat <<'SECT'
```

---

## Benchmark Results

> Columns: `ns/op` = nanoseconds per operation · `B/op` = bytes allocated · `allocs/op` = heap allocations

### AES-GCM (Encrypt)

```
SECT

grep 'BenchmarkAES.*Encrypt' "$BENCH_OUT" || true

cat <<'SECT'
```

### AES-GCM (Decrypt)

```
SECT

grep 'BenchmarkAES.*Decrypt' "$BENCH_OUT" || true

cat <<'SECT'
```

### ECDSA (Sign)

```
SECT

grep 'BenchmarkECDSA.*Sign[^V]' "$BENCH_OUT" || true

cat <<'SECT'
```

### ECDSA (Verify)

```
SECT

grep 'BenchmarkECDSA.*Verify' "$BENCH_OUT" || true

cat <<'SECT'
```

### RSA (Encrypt / Decrypt)

```
SECT

grep -E 'BenchmarkRSA.*(Encrypt|Decrypt)' "$BENCH_OUT" || true

cat <<'SECT'
```

### RSA (Sign / Verify)

```
SECT

grep -E 'BenchmarkRSA.*(Sign|Verify)' "$BENCH_OUT" || true

cat <<'SECT'
```

---

## Full Raw Output

<details>
<summary>Click to expand</summary>

```
SECT

cat "$BENCH_OUT"

cat <<'FOOTER'
```

</details>
FOOTER

} > "$REPORT"

echo ""
echo ">>> Report written to: $REPORT"
