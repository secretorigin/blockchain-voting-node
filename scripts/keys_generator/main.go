package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
)

func main() {
	// Generate ECDSA private key
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	// Print the private key in hex format
	privateKeyHex := hex.EncodeToString(privateKey.D.Bytes())
	fmt.Println("Private Key (hex):", privateKeyHex)

	// Get the public key corresponding to the private key
	publicKey := privateKey.PublicKey

	// Serialize the public key to uncompressed format
	publicKeyBytes := elliptic.Marshal(publicKey.Curve, publicKey.X, publicKey.Y)

	// Print the public key in hex format
	publicKeyHex := hex.EncodeToString(publicKeyBytes)
	fmt.Println("Public Key (hex):", publicKeyHex)
}
