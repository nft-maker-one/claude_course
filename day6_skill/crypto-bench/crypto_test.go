package cryptobench

import (
	"bytes"
	"fmt"
	"testing"
)

// ─── shared fixture ──────────────────────────────────────────────────────────

var (
	suite    *CryptoSuite
	messages = map[string][]byte{
		"small":  []byte("hello world"),
		"medium": bytes.Repeat([]byte("A"), 1024),
		"large":  bytes.Repeat([]byte("B"), 32*1024),
	}
)

func TestMain(m *testing.M) {
	var err error
	suite, err = NewCryptoSuite()
	if err != nil {
		panic(fmt.Sprintf("NewCryptoSuite: %v", err))
	}
	m.Run()
}

// ─── AES correctness ─────────────────────────────────────────────────────────

func TestAES128RoundTrip(t *testing.T) {
	for name, msg := range messages {
		t.Run(name, func(t *testing.T) {
			nonce, ct, err := suite.AES128Encrypt(msg)
			if err != nil {
				t.Fatalf("encrypt: %v", err)
			}
			pt, err := suite.AES128Decrypt(nonce, ct)
			if err != nil {
				t.Fatalf("decrypt: %v", err)
			}
			if !bytes.Equal(pt, msg) {
				t.Error("plaintext mismatch after AES-128 round-trip")
			}
		})
	}
}

func TestAES192RoundTrip(t *testing.T) {
	for name, msg := range messages {
		t.Run(name, func(t *testing.T) {
			nonce, ct, err := suite.AES192Encrypt(msg)
			if err != nil {
				t.Fatalf("encrypt: %v", err)
			}
			pt, err := suite.AES192Decrypt(nonce, ct)
			if err != nil {
				t.Fatalf("decrypt: %v", err)
			}
			if !bytes.Equal(pt, msg) {
				t.Error("plaintext mismatch after AES-192 round-trip")
			}
		})
	}
}

func TestAES256RoundTrip(t *testing.T) {
	for name, msg := range messages {
		t.Run(name, func(t *testing.T) {
			nonce, ct, err := suite.AES256Encrypt(msg)
			if err != nil {
				t.Fatalf("encrypt: %v", err)
			}
			pt, err := suite.AES256Decrypt(nonce, ct)
			if err != nil {
				t.Fatalf("decrypt: %v", err)
			}
			if !bytes.Equal(pt, msg) {
				t.Error("plaintext mismatch after AES-256 round-trip")
			}
		})
	}
}

func TestAESGCMTamperDetection(t *testing.T) {
	nonce, ct, err := suite.AES256Encrypt([]byte("secret"))
	if err != nil {
		t.Fatal(err)
	}
	ct[0] ^= 0xFF // flip bits
	_, err = suite.AES256Decrypt(nonce, ct)
	if err == nil {
		t.Error("expected authentication failure on tampered ciphertext")
	}
}

// ─── ECDSA correctness ───────────────────────────────────────────────────────

func testECDSARoundTrip(
	t *testing.T,
	sign func([]byte) (*ECDSASignature, error),
	verify func([]byte, *ECDSASignature) bool,
	curve string,
) {
	t.Helper()
	msg := []byte("test message for " + curve)
	sig, err := sign(msg)
	if err != nil {
		t.Fatalf("%s sign: %v", curve, err)
	}
	if !verify(msg, sig) {
		t.Errorf("%s: signature verification failed", curve)
	}
}

func TestECDSAP256(t *testing.T) {
	testECDSARoundTrip(t, suite.ECDSASignP256, suite.ECDSAVerifyP256, "P-256")
}

func TestECDSAP384(t *testing.T) {
	testECDSARoundTrip(t, suite.ECDSASignP384, suite.ECDSAVerifyP384, "P-384")
}

func TestECDSAP521(t *testing.T) {
	testECDSARoundTrip(t, suite.ECDSASignP521, suite.ECDSAVerifyP521, "P-521")
}

func TestECDSATamperedMessage(t *testing.T) {
	msg := []byte("authentic message")
	sig, err := suite.ECDSASignP256(msg)
	if err != nil {
		t.Fatal(err)
	}
	if suite.ECDSAVerifyP256([]byte("tampered message"), sig) {
		t.Error("P-256: tampered message must not verify")
	}
}

func TestECDSACrossKeyRejection(t *testing.T) {
	msg := []byte("cross-key test")
	sig, err := suite.ECDSASignP256(msg)
	if err != nil {
		t.Fatal(err)
	}
	// P-384 key must reject a P-256 signature
	if suite.ECDSAVerifyP384(msg, sig) {
		t.Error("P-384 must reject a signature produced by P-256 key")
	}
}

// ─── RSA correctness ─────────────────────────────────────────────────────────

type rsaVariant struct {
	name    string
	encrypt func([]byte) ([]byte, error)
	decrypt func([]byte) ([]byte, error)
	sign    func([]byte) ([]byte, error)
	verify  func([]byte, []byte) error
	maxMsg  int // OAEP limit = (keyBytes - 2*hashBytes - 2)
}

func rsaVariants() []rsaVariant {
	return []rsaVariant{
		{
			"RSA-1024",
			suite.RSA1024Encrypt, suite.RSA1024Decrypt,
			suite.RSA1024Sign, suite.RSA1024Verify,
			62, // 128 - 2*32 - 2
		},
		{
			"RSA-2048",
			suite.RSA2048Encrypt, suite.RSA2048Decrypt,
			suite.RSA2048Sign, suite.RSA2048Verify,
			190, // 256 - 2*32 - 2
		},
		{
			"RSA-4096",
			suite.RSA4096Encrypt, suite.RSA4096Decrypt,
			suite.RSA4096Sign, suite.RSA4096Verify,
			446, // 512 - 2*32 - 2
		},
	}
}

func TestRSAEncryptDecrypt(t *testing.T) {
	for _, v := range rsaVariants() {
		v := v
		t.Run(v.name, func(t *testing.T) {
			msg := bytes.Repeat([]byte("X"), v.maxMsg)
			ct, err := v.encrypt(msg)
			if err != nil {
				t.Fatalf("encrypt: %v", err)
			}
			pt, err := v.decrypt(ct)
			if err != nil {
				t.Fatalf("decrypt: %v", err)
			}
			if !bytes.Equal(pt, msg) {
				t.Error("plaintext mismatch after RSA round-trip")
			}
		})
	}
}

func TestRSASignVerify(t *testing.T) {
	msg := []byte("rsa sign/verify test message")
	for _, v := range rsaVariants() {
		v := v
		t.Run(v.name, func(t *testing.T) {
			sig, err := v.sign(msg)
			if err != nil {
				t.Fatalf("sign: %v", err)
			}
			if err := v.verify(msg, sig); err != nil {
				t.Errorf("verify: %v", err)
			}
		})
	}
}

func TestRSATamperedSignature(t *testing.T) {
	msg := []byte("original message")
	for _, v := range rsaVariants() {
		v := v
		t.Run(v.name, func(t *testing.T) {
			sig, err := v.sign(msg)
			if err != nil {
				t.Fatal(err)
			}
			sig[0] ^= 0xFF
			if v.verify(msg, sig) == nil {
				t.Errorf("%s: tampered signature must not verify", v.name)
			}
		})
	}
}

// ─── AES benchmarks ──────────────────────────────────────────────────────────

var msgSizes = []struct {
	label string
	size  int
}{
	{"1KB", 1024},
	{"16KB", 16 * 1024},
	{"256KB", 256 * 1024},
	{"1MB", 1024 * 1024},
}

func benchAESEncrypt(b *testing.B, encrypt func([]byte) ([]byte, []byte, error), size int) {
	b.Helper()
	msg := bytes.Repeat([]byte("A"), size)
	b.SetBytes(int64(size))
	b.ResetTimer()
	for b.Loop() {
		_, _, err := encrypt(msg)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchAESDecrypt(b *testing.B, encrypt func([]byte) ([]byte, []byte, error), decrypt func([]byte, []byte) ([]byte, error), size int) {
	b.Helper()
	msg := bytes.Repeat([]byte("A"), size)
	nonce, ct, err := encrypt(msg)
	if err != nil {
		b.Fatal(err)
	}
	b.SetBytes(int64(size))
	b.ResetTimer()
	for b.Loop() {
		_, err := decrypt(nonce, ct)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkAES128Encrypt(b *testing.B) {
	for _, s := range msgSizes {
		b.Run(s.label, func(b *testing.B) { benchAESEncrypt(b, suite.AES128Encrypt, s.size) })
	}
}

func BenchmarkAES128Decrypt(b *testing.B) {
	for _, s := range msgSizes {
		b.Run(s.label, func(b *testing.B) { benchAESDecrypt(b, suite.AES128Encrypt, suite.AES128Decrypt, s.size) })
	}
}

func BenchmarkAES192Encrypt(b *testing.B) {
	for _, s := range msgSizes {
		b.Run(s.label, func(b *testing.B) { benchAESEncrypt(b, suite.AES192Encrypt, s.size) })
	}
}

func BenchmarkAES192Decrypt(b *testing.B) {
	for _, s := range msgSizes {
		b.Run(s.label, func(b *testing.B) { benchAESDecrypt(b, suite.AES192Encrypt, suite.AES192Decrypt, s.size) })
	}
}

func BenchmarkAES256Encrypt(b *testing.B) {
	for _, s := range msgSizes {
		b.Run(s.label, func(b *testing.B) { benchAESEncrypt(b, suite.AES256Encrypt, s.size) })
	}
}

func BenchmarkAES256Decrypt(b *testing.B) {
	for _, s := range msgSizes {
		b.Run(s.label, func(b *testing.B) { benchAESDecrypt(b, suite.AES256Encrypt, suite.AES256Decrypt, s.size) })
	}
}

// ─── ECDSA benchmarks ────────────────────────────────────────────────────────

func benchECDSASign(b *testing.B, sign func([]byte) (*ECDSASignature, error)) {
	b.Helper()
	msg := []byte("benchmark message")
	b.ResetTimer()
	for b.Loop() {
		_, err := sign(msg)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchECDSAVerify(b *testing.B, sign func([]byte) (*ECDSASignature, error), verify func([]byte, *ECDSASignature) bool) {
	b.Helper()
	msg := []byte("benchmark message")
	sig, err := sign(msg)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for b.Loop() {
		if !verify(msg, sig) {
			b.Fatal("verify returned false")
		}
	}
}

func BenchmarkECDSAP256Sign(b *testing.B)   { benchECDSASign(b, suite.ECDSASignP256) }
func BenchmarkECDSAP256Verify(b *testing.B) { benchECDSAVerify(b, suite.ECDSASignP256, suite.ECDSAVerifyP256) }
func BenchmarkECDSAP384Sign(b *testing.B)   { benchECDSASign(b, suite.ECDSASignP384) }
func BenchmarkECDSAP384Verify(b *testing.B) { benchECDSAVerify(b, suite.ECDSASignP384, suite.ECDSAVerifyP384) }
func BenchmarkECDSAP521Sign(b *testing.B)   { benchECDSASign(b, suite.ECDSASignP521) }
func BenchmarkECDSAP521Verify(b *testing.B) { benchECDSAVerify(b, suite.ECDSASignP521, suite.ECDSAVerifyP521) }

// ─── RSA benchmarks ──────────────────────────────────────────────────────────

func benchRSAEncrypt(b *testing.B, encrypt func([]byte) ([]byte, error), msgSize int) {
	b.Helper()
	msg := bytes.Repeat([]byte("X"), msgSize)
	b.ResetTimer()
	for b.Loop() {
		_, err := encrypt(msg)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchRSADecrypt(b *testing.B, encrypt func([]byte) ([]byte, error), decrypt func([]byte) ([]byte, error), msgSize int) {
	b.Helper()
	msg := bytes.Repeat([]byte("X"), msgSize)
	ct, err := encrypt(msg)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for b.Loop() {
		_, err := decrypt(ct)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchRSASign(b *testing.B, sign func([]byte) ([]byte, error)) {
	b.Helper()
	msg := []byte("rsa benchmark sign message")
	b.ResetTimer()
	for b.Loop() {
		_, err := sign(msg)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func benchRSAVerify(b *testing.B, sign func([]byte) ([]byte, error), verify func([]byte, []byte) error) {
	b.Helper()
	msg := []byte("rsa benchmark sign message")
	sig, err := sign(msg)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for b.Loop() {
		if err := verify(msg, sig); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRSA1024Encrypt(b *testing.B) { benchRSAEncrypt(b, suite.RSA1024Encrypt, 32) }
func BenchmarkRSA1024Decrypt(b *testing.B) { benchRSADecrypt(b, suite.RSA1024Encrypt, suite.RSA1024Decrypt, 32) }
func BenchmarkRSA1024Sign(b *testing.B)    { benchRSASign(b, suite.RSA1024Sign) }
func BenchmarkRSA1024Verify(b *testing.B)  { benchRSAVerify(b, suite.RSA1024Sign, suite.RSA1024Verify) }

func BenchmarkRSA2048Encrypt(b *testing.B) { benchRSAEncrypt(b, suite.RSA2048Encrypt, 32) }
func BenchmarkRSA2048Decrypt(b *testing.B) { benchRSADecrypt(b, suite.RSA2048Encrypt, suite.RSA2048Decrypt, 32) }
func BenchmarkRSA2048Sign(b *testing.B)    { benchRSASign(b, suite.RSA2048Sign) }
func BenchmarkRSA2048Verify(b *testing.B)  { benchRSAVerify(b, suite.RSA2048Sign, suite.RSA2048Verify) }

func BenchmarkRSA4096Encrypt(b *testing.B) { benchRSAEncrypt(b, suite.RSA4096Encrypt, 32) }
func BenchmarkRSA4096Decrypt(b *testing.B) { benchRSADecrypt(b, suite.RSA4096Encrypt, suite.RSA4096Decrypt, 32) }
func BenchmarkRSA4096Sign(b *testing.B)    { benchRSASign(b, suite.RSA4096Sign) }
func BenchmarkRSA4096Verify(b *testing.B)  { benchRSAVerify(b, suite.RSA4096Sign, suite.RSA4096Verify) }
