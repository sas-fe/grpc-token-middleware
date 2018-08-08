package tokenfuncs

import (
	"errors"

	"github.com/golang/glog"
	"github.com/sas-fe/grpc-token-middleware"
	"github.com/sas-fe/grpc-token-middleware/pb"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

const serviceSepString = ":"

// TokenFuncs implements tokenapi.TokenAPI
type TokenFuncs struct {
	TSClient       tokenstore.TokenStoreClient
	TSDaemonClient tokenstore.TSDaemonClient
	ServiceName    string
	AsyncIncChan   chan string
	Async          bool
}

// TFOptions configures how TokenFuncs are set up.
type TFOptions func(*TokenFuncs)

// CheckValidity returns true if token is valid, false otherwise.
func (t *TokenFuncs) CheckValidity(ctx context.Context) (bool, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return false, errors.New("Request metadata doesn't exist")
	}
	tokenID, err := t.GetToken(md)
	if err != nil {
		return false, err
	}

	res, err := t.TSDaemonClient.CheckValidity(ctx, &tokenstore.Token{Id: tokenID})
	if err != nil {
		glog.Errorf("TSDaemonClient.CheckValidity() Error: %v", err)
		return false, err
	}
	if !res.Valid {
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
		glog.Errorf("TSClient.IncUsage() Error: %v", err)
		return false, err
	}
	return true, nil
}

// AsyncIncrementUsage increments token usage via a channel.
func (t *TokenFuncs) AsyncIncrementUsage(ctx context.Context) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		glog.Errorf("AsyncIncrementUsage() Error: No incoming metadata")
		return
	}
	tokenID, err := t.GetToken(md)
	if err != nil {
		glog.Errorf("AsyncIncrementUsage() Error: %v", err)
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
			glog.V(2).Infof("Received: %v", tokenID)
			ctx := context.Background()
			_, err := t.TSClient.IncUsage(ctx, &tokenstore.Token{Id: tokenID})
			if err != nil {
				glog.Errorf("TSClient.IncUsage() Error: %v", err)
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

	return t.ServiceName + serviceSepString + apiToken[0], nil
}

// NewTokenStoreClient namespaces the function from tokenstore.
func NewTokenStoreClient(tsConn *grpc.ClientConn) tokenstore.TokenStoreClient {
	return tokenstore.NewTokenStoreClient(tsConn)
}

// NewTSDaemonClient namespaces the function from tokenstore.
func NewTSDaemonClient(tdConn *grpc.ClientConn) tokenstore.TSDaemonClient {
	return tokenstore.NewTSDaemonClient(tdConn)
}

// NewTokenFuncs creates a new TokenFuncs using a connection to the tokenstore and service name.
func NewTokenFuncs(tsClient tokenstore.TokenStoreClient, tdClient tokenstore.TSDaemonClient, serviceName string, options ...TFOptions) *TokenFuncs {
	tokenFuncs := &TokenFuncs{
		TSClient:       tsClient,
		TSDaemonClient: tdClient,
		ServiceName:    serviceName,
	}

	// Force async if no other options
	if len(options) == 0 {
		options = append(options, WithAsync())
	}

	for _, option := range options {
		option(tokenFuncs)
	}

	return tokenFuncs
}

// WithAsync sets marks async incrementation.
func WithAsync() TFOptions {
	return func(tf *TokenFuncs) {
		asyncChan := make(chan string)
		tf.AsyncIncChan = asyncChan
		tf.Async = true
		go tf.ListenAndInc()
	}
}

var _ tokenapi.TokenAPI = (*TokenFuncs)(nil)
