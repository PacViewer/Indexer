package client

import (
	"context"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Pactus struct {
	Blockchain  pactus.BlockchainClient
	Transaction pactus.TransactionClient
	Network     pactus.NetworkClient
}

func NewPactus(ctx context.Context, rpc string) (*Pactus, error) {
	dialOpts := make([]grpc.DialOption, 0)
	dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	dialOpts = append(dialOpts, grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor()))

	conn, err := grpc.DialContext(ctx, rpc, dialOpts...)
	if err != nil {
		return nil, err
	}

	return &Pactus{
		pactus.NewBlockchainClient(conn),
		pactus.NewTransactionClient(conn),
		pactus.NewNetworkClient(conn),
	}, nil
}
