package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	ErrorHandler(err)
	public := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, public
}

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func NewWallet() *Wallet {
	private, public := NewKeyPair()
	wallet := Wallet{private, public}
	return &wallet
}

func (w *Wallet) EncodeWalletKeys() {
	x509Encoded, _ := x509.MarshalECPrivateKey(&w.PrivateKey)
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(&w.PrivateKey.PublicKey)
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

	f, err := os.OpenFile("private.pem", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	ErrorHandler(err)
	f.WriteString(string(pemEncoded))

	f, err = os.OpenFile("public.pem", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	ErrorHandler(err)
	f.WriteString(string(pemEncodedPub))
}

func DecodeWalletKeys(pemEncoded string, pemEncodedPub string) Wallet {
	block, _ := pem.Decode([]byte(pemEncoded))
	x509Encoded := block.Bytes
	privateKey, _ := x509.ParseECPrivateKey(x509Encoded)

	blockPub, _ := pem.Decode([]byte(pemEncodedPub))
	x509EncodedPub := blockPub.Bytes
	// genericPublicKey, _ := x509.ParsePKIXPublicKey(x509EncodedPub)
	// publicKey := genericPublicKey.(*ecdsa.PublicKey)

	wallet := Wallet{PrivateKey: *privateKey, PublicKey: x509EncodedPub}

	return wallet
}
