package token

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"testing"

	"github.com/braumsmilk/go-token/keys"
	"github.com/stretchr/testify/assert"
)


func getRsaKeyBytes() (priv []byte, pub []byte) {
	k, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}

	// Encode private key to PKCS#1 ASN.1 PEM.
	priv = pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(k),
		},
	)
	pubB, err := x509.MarshalPKIXPublicKey(&k.PublicKey)
	if err != nil {
		panic(err)
	}

	pub = pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubB,
		},
	)

	return
}

func TestNewJwtToken(t *testing.T) {

	assert.Nil(t, keys.Init(getRsaKeyBytes()))
	tkn, err := NewJwtToken("1", "aud", "id", "issuer")
	assert.Nilf(t, err, "should not error when getting new token")

	jwtToken, err := DecryptJwtToken(tkn)
	assert.Nilf(t, err, "should not fail to decrypt token")
	assert.NotNilf(t, jwtToken, "should have gotten a non-nil token")
}

func Benchmark_NewJwtToken(b *testing.B) {
	keys.Init(getRsaKeyBytes())

	b.Run("creating 10000 tokens for 100 different users", func(b *testing.B) {
		users := []string{}
		for i := 0; i < 100; i++ {
			users = append(users, fmt.Sprintf("%d", i))
		}

		for i := 0; i < 1000; i++ {
			for _, u := range users {
				_, err := NewJwtToken(u, "aud", "id", "issuer")
				assert.Nil(b, err)
			}
		}
	})
}