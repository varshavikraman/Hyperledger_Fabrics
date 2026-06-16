package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	gwInstance *client.Gateway
	grpcConn   *grpc.ClientConn
	gwOnce     sync.Once
	gwErr      error
)

// initGateway initializes the Fabric Gateway connection
func initGateway(profile *Profile) (*client.Gateway, *grpc.ClientConn, error) {
	if profile.UseMock {
		return nil, nil, nil
	}

	var creds credentials.TransportCredentials
	if profile.TLSCertPath != "" {
		certPool, err := loadCertificatePool(profile)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to load TLS certificate pool: %w", err)
		}
		creds = credentials.NewClientTLSFromCert(certPool, profile.PeerHostOverride)
	} else {
		creds = insecure.NewCredentials()
	}

	// Dial the peer gRPC endpoint
	conn, err := grpc.Dial(profile.PeerEndpoint, grpc.WithTransportCredentials(creds), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to dial peer endpoint %s: %w", profile.PeerEndpoint, err)
	}

	id, err := newIdentity(profile)
	if err != nil {
		conn.Close()
		return nil, nil, fmt.Errorf("failed to create client identity: %w", err)
	}

	signer, err := newSign(profile)
	if err != nil {
		conn.Close()
		return nil, nil, fmt.Errorf("failed to create client signer: %w", err)
	}

	// Connect to gateway using identity, signer, and gRPC connection
	gw, err := client.Connect(
		id,
		client.WithSign(signer),
		client.WithClientConnection(conn),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(15*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		conn.Close()
		return nil, nil, fmt.Errorf("failed to connect to gateway: %w", err)
	}

	return gw, conn, nil
}

// getGateway returns the global Gateway instance or initializes it
func getGateway() (*client.Gateway, error) {
	profile := loadProfile()
	if profile.UseMock {
		return nil, nil
	}

	gwOnce.Do(func() {
		gwInstance, grpcConn, gwErr = initGateway(profile)
	})

	return gwInstance, gwErr
}

// closeGateway closes the gateway connection
func closeGateway() {
	if grpcConn != nil {
		grpcConn.Close()
	}
	if gwInstance != nil {
		gwInstance.Close()
	}
}
