package main

import (
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hyperledger/fabric-gateway/pkg/identity"
)

// Profile configures the connection details for a Fabric network
type Profile struct {
	MspID            string
	CertPath         string
	KeyPath          string
	TLSCertPath      string
	PeerEndpoint     string
	PeerHostOverride string
	UseMock          bool
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func loadProfile() *Profile {
	mockEnv := getEnv("MOCK_FABRIC", "true")
	useMock := mockEnv == "true" || mockEnv == "1"

	return &Profile{
		MspID:            getEnv("FABRIC_MSP_ID", "Org1MSP"),
		CertPath:         getEnv("FABRIC_CERT_PATH", ""),
		KeyPath:          getEnv("FABRIC_KEY_PATH", ""),
		TLSCertPath:      getEnv("FABRIC_TLS_CERT_PATH", ""),
		PeerEndpoint:     getEnv("FABRIC_PEER_ENDPOINT", "localhost:7051"),
		PeerHostOverride: getEnv("FABRIC_PEER_HOST_OVERRIDE", ""),
		UseMock:          useMock,
	}
}

// readFirstFileInDir helper to read the key file which typically has a dynamic name in Fabric's MSP folders
func readFirstFileInDir(dirPath string) ([]byte, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if !file.IsDir() {
			return os.ReadFile(filepath.Join(dirPath, file.Name()))
		}
	}
	return nil, fmt.Errorf("no files found in directory: %s", dirPath)
}

func newIdentity(profile *Profile) (*identity.X509Identity, error) {
	if profile.CertPath == "" {
		return nil, fmt.Errorf("certificate path is not configured")
	}

	var certBytes []byte
	var err error
	info, err := os.Stat(profile.CertPath)
	if err == nil && info.IsDir() {
		certBytes, err = readFirstFileInDir(profile.CertPath)
	} else {
		certBytes, err = os.ReadFile(profile.CertPath)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}

	cert, err := identity.CertificateFromPEM(certBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	id, err := identity.NewX509Identity(profile.MspID, cert)
	if err != nil {
		return nil, fmt.Errorf("failed to create x509 identity: %w", err)
	}

	return id, nil
}

func newSign(profile *Profile) (identity.Sign, error) {
	if profile.KeyPath == "" {
		return nil, fmt.Errorf("private key path is not configured")
	}

	var keyBytes []byte
	var err error
	info, err := os.Stat(profile.KeyPath)
	if err == nil && info.IsDir() {
		keyBytes, err = readFirstFileInDir(profile.KeyPath)
	} else {
		keyBytes, err = os.ReadFile(profile.KeyPath)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	privateKey, err := identity.PrivateKeyFromPEM(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	signer, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create private key signer: %w", err)
	}

	return signer, nil
}

func loadCertificatePool(profile *Profile) (*x509.CertPool, error) {
	if profile.TLSCertPath == "" {
		return nil, fmt.Errorf("TLS certificate path is not configured")
	}

	var pemCert []byte
	var err error
	info, err := os.Stat(profile.TLSCertPath)
	if err == nil && info.IsDir() {
		pemCert, err = readFirstFileInDir(profile.TLSCertPath)
	} else {
		pemCert, err = os.ReadFile(profile.TLSCertPath)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read TLS certificate: %w", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemCert) {
		return nil, fmt.Errorf("failed to append certificate to pool")
	}

	return certPool, nil
}
