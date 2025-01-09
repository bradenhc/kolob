// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"time"
)

// GenerateSelfSignedCert generates a self-signed ECDSA certificate and private key.
// It returns the certificate and key in PEM format, or an error if the generation fails.
func GenerateSelfSignedCert() ([]byte, []byte, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour)

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, nil, err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Kolob"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		return nil, nil, err
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyPEM, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return nil, nil, err
	}

	keyPEMBlock := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyPEM})

	return certPEM, keyPEMBlock, nil
}

// SaveCertAndKey saves the provided certificate and key to the specified file paths.
// The certificate is saved with read permissions for everyone, while the key is saved with
// read/write permissions for the owner only. Returns an error if saving the files fails.
func SaveCertAndKey(certPEM, keyPEM []byte, certPath, keyPath string) error {
	err := os.WriteFile(certPath, certPEM, 0644)
	if err != nil {
		return err
	}

	err = os.WriteFile(keyPath, keyPEM, 0600)
	if err != nil {
		return err
	}

	return nil
}
