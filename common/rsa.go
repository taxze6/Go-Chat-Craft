package common

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"os"
)

func RsaDecoder(data string) (string, error) {
	privateKeyPEM, err := os.ReadFile("common/private.pem")
	if err != nil {
		return "", errors.New("unable to read private key file")
	}

	// Parse PEM data
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		return "", errors.New("invalid PEM data")
	}

	var key interface{}
	if block.Type == "PRIVATE KEY" {
		// Parse PKCS8 private key
		key, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return "", errors.New("unable to parse private key")
		}
	} else {
		return "", errors.New("unsupported private key format")
	}

	// Convert the encrypted string to a byte array
	encryptedData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", errors.New("unable to decode data")
	}
	// Decrypt the data using the private key
	privateKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return "", errors.New("invalid private key type")
	}
	decryptedData, err := rsa.DecryptPKCS1v15(nil, privateKey, encryptedData)
	if err != nil {
		return "", errors.New("decryption failed")
	}

	return string(decryptedData), nil
}
