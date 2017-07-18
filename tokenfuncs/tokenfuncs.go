package tokenfuncs

import (
	"github.com/sas-fe/grpc-token-middleware/pb"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

// TokenFuncs implements tokenapi.TokenAPI
type TokenFuncs struct {
	TSClient    tokenstore.TokenStoreClient
	ServiceName string
}

// CheckValidity returns true if token is valid, false otherwise.
func (t *TokenFuncs) CheckValidity(ctx context.Context) (bool, error) {
	md, _ := metadata.FromContext(ctx)
	tokenID, err := t.GetToken(md)
	if err != nil {
		return false, err
	}

	res, err := t.TSClient.TokenStatus(ctx, &tokenstore.Token{Id: tokenID})
	if err != nil {
		return false, err
	}
	if !res.Allowed {
		return false, nil
	}
	return true, nil
}

// IncrementUsage performs the incrementing of token usage.
func (t *TokenFuncs) IncrementUsage(ctx context.Context) (bool, error) {
	md, _ := metadata.FromContext(ctx)
	tokenID, err := t.GetToken(md)
	if err != nil {
		return false, err
	}

	_, err = t.TSClient.IncUsage(ctx, &tokenstore.Token{Id: tokenID})
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetToken parses the token from metadata and creates the tokenID.
func (t *TokenFuncs) GetToken(md metadata.MD) (string, error) {
	apiToken, ok := md["authorization"]
	if !ok {
		return "", grpc.Errorf(codes.InvalidArgument, "Missing 'authorization' from request metadata")
	}

	return t.ServiceName + ":" + apiToken[0], nil
}

// NewTokenFuncs creates a new TokenFuncs using a connection to the tokenstore and service name.
func NewTokenFuncs(tsConn *grpc.ClientConn, serviceName string) *TokenFuncs {
	tsClient := tokenstore.NewTokenStoreClient(tsConn)

	return &TokenFuncs{
		TSClient:    tsClient,
		ServiceName: serviceName,
	}
}
