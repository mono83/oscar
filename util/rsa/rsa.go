package rsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"math/big"
	"time"
)

// Generate generates new RSA pubkey/privkey/rsa
func Generate(len int, ttl time.Duration) (*RSA, error) {
	// generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, len)
	if err != nil {
		return nil, err
	}
	publicKey := &privateKey.PublicKey

	// certificate template
	start := time.Now()
	template := &x509.Certificate{
		IsCA: true,
		BasicConstraintsValid: true,
		SubjectKeyId:          []byte{1, 2, 3},
		SerialNumber:          big.NewInt(1234),
		Subject: pkix.Name{
			Country:      []string{"International"},
			Organization: []string{"Oscar"},
		},
		NotBefore: start.Add(-time.Second),
		NotAfter:  start.Add(ttl),
		// see http://golang.org/pkg/crypto/x509/#KeyUsage
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}

	var parent = template
	cert, err := x509.CreateCertificate(rand.Reader, template, parent, publicKey, privateKey)
	if err != nil {
		return nil, err
	}

	return &RSA{
		key:  privateKey,
		cert: cert,
	}, nil
}

// RSA is complex structure with Lua bindings, used to handle RSA encryption, decryption
// and digital signatures
type RSA struct {
	key  *rsa.PrivateKey
	cert []byte
}

// EncodedCertificate returns PEM-encoded version of RSA certificate
func (r RSA) EncodedCertificate() string {
	encoded := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: r.cert,
	})
	return string(encoded)
}

// SignSHA256 produces digital signature of provided data using RSAwithSHA256 sig algo
func (r RSA) SignSHA256(data []byte) ([]byte, error) {
	hashed := sha256.Sum256(data)
	return rsa.SignPKCS1v15(rand.Reader, r.key, crypto.SHA256, hashed[:])
}

// SignSHA256B64 produces digital signature of provided data using RSAwithSHA256 sig algo
// Result encoded into string using Base64
func (r RSA) SignSHA256B64(data []byte) (string, error) {
	bts, err := r.SignSHA256(data)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(bts), nil
}
