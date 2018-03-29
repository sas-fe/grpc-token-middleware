package tokenfuncs

import (
	"errors"
	"log"

	"github.com/sas-fe/grpc-token-middleware"
	"github.com/sas-fe/grpc-token-middleware/pb"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

// TokenFuncs implements tokenapi.TokenAPI
type TokenFuncs struct {
	TSClient     tokenstore.TokenStoreClient
	ServiceName  string
	AsyncIncChan chan string
	Async        bool
}

// CheckValidity returns true if token is valid, false otherwise.
func (t *TokenFuncs) CheckValidity(ctx context.Context) (bool, error) {
	md, _ := metadata.FromIncomingContext(ctx)
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
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return false, errors.New("Request metadata doesn't exist")
	}
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

// AsyncIncrementUsage increments token usage via a channel.
func (t *TokenFuncs) AsyncIncrementUsage(ctx context.Context) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return
	}
	tokenID, err := t.GetToken(md)
	if err != nil {
		return
	}

	t.AsyncIncChan <- tokenID
}

// IsAsync checks if we are using async incrementation.
func (t *TokenFuncs) IsAsync() bool {
	return t.Async
}

// ListenAndInc listens for tokenID on the AsyncIncChan channel and performs the incrementation.
func (t *TokenFuncs) ListenAndInc() {
	for {
		select {
		case tokenID := <-t.AsyncIncChan:
			log.Printf("Received %v\n", tokenID)
			ctx := context.Background()
			_, err := t.TSClient.IncUsage(ctx, &tokenstore.Token{Id: tokenID})
			if err != nil {
				log.Println(err)
			}
		}
	}
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
func NewTokenFuncs(tsConn *grpc.ClientConn, serviceName string, asyncInc bool) *TokenFuncs {
	tsClient := tokenstore.NewTokenStoreClient(tsConn)

	asyncChan := make(chan string)

	tokenFuncs := &TokenFuncs{
		TSClient:     tsClient,
		ServiceName:  serviceName,
		AsyncIncChan: asyncChan,
		Async:        asyncInc,
	}

	if asyncInc {
		go tokenFuncs.ListenAndInc()
	}

	return tokenFuncs
}

var _ tokenapi.TokenAPI = (*TokenFuncs)(nil)
