package tokenfuncs

import (
	"testing"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"

	"github.com/golang/mock/gomock"
	pb "github.com/sas-fe/grpc-token-middleware/pb"
	tsmock "github.com/sas-fe/grpc-token-middleware/pb/mock_tokenstore"
)

func TestAsyncIncrementUsage(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mts := tsmock.NewMockTokenStoreClient(mockCtrl)
	tf := NewTokenFuncs(mts, "test", WithAsync())

	nTests := 100
	usage := 0
	mts.EXPECT().IncUsage(
		gomock.Any(),
		gomock.Any(),
	).Do(func(_ interface{}, _ interface{}) {
		usage++
	}).Return(&pb.RpcStatus{Success: true}, nil).Times(nTests)

	md := map[string][]string{
		"authorization": []string{"key1"},
	}
	ctx := metadata.NewIncomingContext(context.Background(), md)

	for i := 0; i < nTests; i++ {
		go tf.AsyncIncrementUsage(ctx)
	}

	time.Sleep(100 * time.Millisecond)

	currUsage := usage
	if currUsage != nTests {
		t.Errorf("Expected AsyncIncrementUsage to increment %v times, but got %v times", nTests, currUsage)
	}
}
