package grpcclient

import (
    "context"
    "fmt"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
)

type BackendClient struct {
    conn   *grpc.ClientConn
    target string
}

func NewBackendClient(target, certPath string) (*BackendClient, error) {
    creds, err := credentials.NewClientTLSFromFile(certPath, "")
    if err != nil {
        return nil, err
    }
    conn, err := grpc.Dial(target, grpc.WithTransportCredentials(creds))
    if err != nil {
        return nil, err
    }
    return &BackendClient{conn: conn, target: target}, nil
}

func (c *BackendClient) Close() error {
    return c.conn.Close()
}

func (c *BackendClient) Ping(ctx context.Context) error {
    // Placeholder for actual generated gRPC client call.
    fmt.Println("connecting to backend", c.target)
    return nil
}
