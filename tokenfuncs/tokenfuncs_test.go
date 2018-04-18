package tokenfuncs

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	tokenstore "github.com/sas-fe/grpc-token-middleware/pb"
	tsmock "github.com/sas-fe/grpc-token-middleware/pb/mock_tokenstore"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

func TestAsyncIncrementUsage(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mts := tsmock.NewMockTokenStoreClient(mockCtrl)
	mtsd := tsmock.NewMockTSDaemonClient(mockCtrl)
	// tf := NewTokenFuncs(mts, mtsd, "test", WithAsync())
	tf := NewTokenFuncs(mts, mtsd, "test")

	nTests := 100
	usage := 0
	mts.EXPECT().IncUsage(
		gomock.Any(),
		gomock.Any(),
	).Do(func(_ interface{}, _ interface{}) {
		usage++
	}).Return(&tokenstore.RpcStatus{Success: true}, nil).Times(nTests)

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

func TestTokenFuncs_GetToken(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mts := tsmock.NewMockTokenStoreClient(mockCtrl)
	mtsd := tsmock.NewMockTSDaemonClient(mockCtrl)
	tf := NewTokenFuncs(mts, mtsd, "test")

	token := "key1"
	noAuth := metadata.Pairs()
	withAuth := map[string][]string{
		"authorization": []string{token},
	}

	type args struct {
		md metadata.MD
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "NoAuth",
			args:    args{md: noAuth},
			want:    "",
			wantErr: true,
		},
		{
			name:    "WithAuth",
			args:    args{md: withAuth},
			want:    "test:key1",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tf.GetToken(tt.args.md)
			if (err != nil) != tt.wantErr {
				t.Errorf("TokenFuncs.GetToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TokenFuncs.GetToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenFuncs_CheckValidity(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mts := tsmock.NewMockTokenStoreClient(mockCtrl)
	mtsd := tsmock.NewMockTSDaemonClient(mockCtrl)
	tf := NewTokenFuncs(mts, mtsd, "test")

	mtsd.EXPECT().CheckValidity(
		gomock.Any(),
		&tokenstore.Token{Id: "test:valid"},
	).Return(&tokenstore.Validity{Valid: true}, nil)

	mtsd.EXPECT().CheckValidity(
		gomock.Any(),
		&tokenstore.Token{Id: "test:invalid"},
	).Return(&tokenstore.Validity{Valid: false}, nil)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "NoAuth",
			args:    args{ctx: context.Background()},
			want:    false,
			wantErr: true,
		},
		{
			name: "Valid",
			args: args{ctx: metadata.NewIncomingContext(context.Background(), map[string][]string{
				"authorization": []string{"valid"},
			})},
			want:    true,
			wantErr: false,
		},
		{
			name: "Invalid",
			args: args{ctx: metadata.NewIncomingContext(context.Background(), map[string][]string{
				"authorization": []string{"invalid"},
			})},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tf.CheckValidity(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("TokenFuncs.CheckValidity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TokenFuncs.CheckValidity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenFuncs_IncrementUsage(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mts := tsmock.NewMockTokenStoreClient(mockCtrl)
	mtsd := tsmock.NewMockTSDaemonClient(mockCtrl)
	tf := NewTokenFuncs(mts, mtsd, "test")

	usage := 0
	mts.EXPECT().IncUsage(
		gomock.Any(),
		gomock.Any(),
	).Do(func(_ interface{}, _ interface{}) {
		usage++
	}).Return(&tokenstore.RpcStatus{Success: true}, nil)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "NoAuth",
			args:    args{ctx: context.Background()},
			want:    false,
			wantErr: true,
		},
		{
			name: "Valid",
			args: args{ctx: metadata.NewIncomingContext(context.Background(), map[string][]string{
				"authorization": []string{"valid"},
			})},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tf.IncrementUsage(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("TokenFuncs.IncrementUsage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TokenFuncs.IncrementUsage() = %v, want %v", got, tt.want)
			}
		})
	}

	if usage != 1 {
		t.Errorf("TokenFuncs.IncrementUsage() incremented %v times, want %v", usage, 1)
	}
}
