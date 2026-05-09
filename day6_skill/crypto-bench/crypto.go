package cryptobench

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
	"math/big"
)

// CryptoSuite holds key material for AES, ECDSA, and RSA operations.
type CryptoSuite struct {
	// AES
	aesKey128 []byte
	aesKey192 []byte
	aesKey256 []byte

	// ECDSA
	ecP256Key *ecdsa.PrivateKey
	ecP384Key *ecdsa.PrivateKey
	ecP521Key *ecdsa.PrivateKey

	// RSA
	rsa1024Key *rsa.PrivateKey
	rsa2048Key *rsa.PrivateKey
	rsa4096Key *rsa.PrivateKey
}

// NewCryptoSuite generates all key material. Expensive — call once.
func NewCryptoSuite() (*CryptoSuite, error) {
	s := &CryptoSuite{}

	// AES keys
	s.aesKey128 = make([]byte, 16)
	s.aesKey192 = make([]byte, 24)
	s.aesKey256 = make([]byte, 32)
	for _, k := range [][]byte{s.aesKey128, s.aesKey192, s.aesKey256} {
		if _, err := rand.Read(k); err != nil {
			return nil, fmt.Errorf("aes key gen: %w", err)
		}
	}

	// ECDSA keys
	curves := []struct {
		curve elliptic.Curve
		dst   **ecdsa.PrivateKey
	}{
		{elliptic.P256(), &s.ecP256Key},
		{elliptic.P384(), &s.ecP384Key},
		{elliptic.P521(), &s.ecP521Key},
	}
	for _, c := range curves {
		k, err := ecdsa.GenerateKey(c.curve, rand.Reader)
		if err != nil {
			return nil, fmt.Errorf("ecdsa key gen %s: %w", c.curve.Params().Name, err)
		}
		*c.dst = k
	}

	// RSA keys
	rsaBits := []struct {
		bits int
		dst  **rsa.PrivateKey
	}{
		{1024, &s.rsa1024Key},
		{2048, &s.rsa2048Key},
		{4096, &s.rsa4096Key},
	}
	for _, r := range rsaBits {
		k, err := rsa.GenerateKey(rand.Reader, r.bits)
		if err != nil {
			return nil, fmt.Errorf("rsa key gen %d: %w", r.bits, err)
		}
		*r.dst = k
	}

	return s, nil
}

// ─── AES-GCM ────────────────────────────────────────────────────────────────

func aesGCMEncrypt(key, plaintext []byte) (nonce, ciphertext []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}
	nonce = make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, nil, err
	}
	ct := gcm.Seal(nil, nonce, plaintext, nil)
	return nonce, ct, nil
}

func aesGCMDecrypt(key, nonce, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// AES128Encrypt encrypts plaintext with AES-128-GCM.
func (s *CryptoSuite) AES128Encrypt(plaintext []byte) (nonce, ct []byte, err error) {
	return aesGCMEncrypt(s.aesKey128, plaintext)
}

// AES128Decrypt decrypts AES-128-GCM ciphertext.
func (s *CryptoSuite) AES128Decrypt(nonce, ct []byte) ([]byte, error) {
	return aesGCMDecrypt(s.aesKey128, nonce, ct)
}

// AES192Encrypt encrypts plaintext with AES-192-GCM.
func (s *CryptoSuite) AES192Encrypt(plaintext []byte) (nonce, ct []byte, err error) {
	return aesGCMEncrypt(s.aesKey192, plaintext)
}

// AES192Decrypt decrypts AES-192-GCM ciphertext.
func (s *CryptoSuite) AES192Decrypt(nonce, ct []byte) ([]byte, error) {
	return aesGCMDecrypt(s.aesKey192, nonce, ct)
}

// AES256Encrypt encrypts plaintext with AES-256-GCM.
func (s *CryptoSuite) AES256Encrypt(plaintext []byte) (nonce, ct []byte, err error) {
	return aesGCMEncrypt(s.aesKey256, plaintext)
}

// AES256Decrypt decrypts AES-256-GCM ciphertext.
func (s *CryptoSuite) AES256Decrypt(nonce, ct []byte) ([]byte, error) {
	return aesGCMDecrypt(s.aesKey256, nonce, ct)
}

// ─── ECDSA ──────────────────────────────────────────────────────────────────

type ECDSASignature struct {
	R, S *big.Int
}

func ecdsaSign(key *ecdsa.PrivateKey, msg []byte) (*ECDSASignature, error) {
	h := sha256.Sum256(msg)
	r, s, err := ecdsa.Sign(rand.Reader, key, h[:])
	if err != nil {
		return nil, err
	}
	return &ECDSASignature{R: r, S: s}, nil
}

func ecdsaVerify(pub *ecdsa.PublicKey, msg []byte, sig *ECDSASignature) bool {
	h := sha256.Sum256(msg)
	return ecdsa.Verify(pub, h[:], sig.R, sig.S)
}

// ECDSASignP256 signs msg with the P-256 key.
func (s *CryptoSuite) ECDSASignP256(msg []byte) (*ECDSASignature, error) {
	return ecdsaSign(s.ecP256Key, msg)
}

// ECDSAVerifyP256 verifies a P-256 signature.
func (s *CryptoSuite) ECDSAVerifyP256(msg []byte, sig *ECDSASignature) bool {
	return ecdsaVerify(&s.ecP256Key.PublicKey, msg, sig)
}

// ECDSASignP384 signs msg with the P-384 key.
func (s *CryptoSuite) ECDSASignP384(msg []byte) (*ECDSASignature, error) {
	return ecdsaSign(s.ecP384Key, msg)
}

// ECDSAVerifyP384 verifies a P-384 signature.
func (s *CryptoSuite) ECDSAVerifyP384(msg []byte, sig *ECDSASignature) bool {
	return ecdsaVerify(&s.ecP384Key.PublicKey, msg, sig)
}

// ECDSASignP521 signs msg with the P-521 key.
func (s *CryptoSuite) ECDSASignP521(msg []byte) (*ECDSASignature, error) {
	return ecdsaSign(s.ecP521Key, msg)
}

// ECDSAVerifyP521 verifies a P-521 signature.
func (s *CryptoSuite) ECDSAVerifyP521(msg []byte, sig *ECDSASignature) bool {
	return ecdsaVerify(&s.ecP521Key.PublicKey, msg, sig)
}

// ─── RSA ────────────────────────────────────────────────────────────────────

func rsaEncrypt(pub *rsa.PublicKey, plaintext []byte) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, plaintext, nil)
}

func rsaDecrypt(priv *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, ciphertext, nil)
}

func rsaSign(priv *rsa.PrivateKey, msg []byte) ([]byte, error) {
	h := sha512.Sum512(msg)
	return rsa.SignPSS(rand.Reader, priv, crypto.SHA512, h[:], nil)
}

func rsaVerify(pub *rsa.PublicKey, msg, sig []byte) error {
	h := sha512.Sum512(msg)
	return rsa.VerifyPSS(pub, crypto.SHA512, h[:], sig, nil)
}

// RSA1024Encrypt encrypts with RSA-1024 OAEP.
func (s *CryptoSuite) RSA1024Encrypt(plaintext []byte) ([]byte, error) {
	return rsaEncrypt(&s.rsa1024Key.PublicKey, plaintext)
}

// RSA1024Decrypt decrypts RSA-1024 OAEP ciphertext.
func (s *CryptoSuite) RSA1024Decrypt(ct []byte) ([]byte, error) {
	return rsaDecrypt(s.rsa1024Key, ct)
}

// RSA1024Sign signs with RSA-1024 PSS.
func (s *CryptoSuite) RSA1024Sign(msg []byte) ([]byte, error) {
	return rsaSign(s.rsa1024Key, msg)
}

// RSA1024Verify verifies an RSA-1024 PSS signature.
func (s *CryptoSuite) RSA1024Verify(msg, sig []byte) error {
	return rsaVerify(&s.rsa1024Key.PublicKey, msg, sig)
}

// RSA2048Encrypt encrypts with RSA-2048 OAEP.
func (s *CryptoSuite) RSA2048Encrypt(plaintext []byte) ([]byte, error) {
	return rsaEncrypt(&s.rsa2048Key.PublicKey, plaintext)
}

// RSA2048Decrypt decrypts RSA-2048 OAEP ciphertext.
func (s *CryptoSuite) RSA2048Decrypt(ct []byte) ([]byte, error) {
	return rsaDecrypt(s.rsa2048Key, ct)
}

// RSA2048Sign signs with RSA-2048 PSS.
func (s *CryptoSuite) RSA2048Sign(msg []byte) ([]byte, error) {
	return rsaSign(s.rsa2048Key, msg)
}

// RSA2048Verify verifies an RSA-2048 PSS signature.
func (s *CryptoSuite) RSA2048Verify(msg, sig []byte) error {
	return rsaVerify(&s.rsa2048Key.PublicKey, msg, sig)
}

// RSA4096Encrypt encrypts with RSA-4096 OAEP.
func (s *CryptoSuite) RSA4096Encrypt(plaintext []byte) ([]byte, error) {
	return rsaEncrypt(&s.rsa4096Key.PublicKey, plaintext)
}

// RSA4096Decrypt decrypts RSA-4096 OAEP ciphertext.
func (s *CryptoSuite) RSA4096Decrypt(ct []byte) ([]byte, error) {
	return rsaDecrypt(s.rsa4096Key, ct)
}

// RSA4096Sign signs with RSA-4096 PSS.
func (s *CryptoSuite) RSA4096Sign(msg []byte) ([]byte, error) {
	return rsaSign(s.rsa4096Key, msg)
}

// RSA4096Verify verifies an RSA-4096 PSS signature.
func (s *CryptoSuite) RSA4096Verify(msg, sig []byte) error {
	return rsaVerify(&s.rsa4096Key.PublicKey, msg, sig)
}

// ─── helpers ────────────────────────────────────────────────────────────────

var ErrSignatureInvalid = errors.New("signature invalid")
