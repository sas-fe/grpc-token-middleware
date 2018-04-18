package main

import (
	"log"
	"net"

	"golang.org/x/net/context"

	"github.com/sas-fe/grpc-token-middleware"
	pb "github.com/sas-fe/grpc-token-middleware/example/pb"
	"github.com/sas-fe/grpc-token-middleware/tokenfuncs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func (s *server) Do(ctx context.Context, in *pb.DoRequest) (*pb.DoResponse, error) {
	return &pb.DoResponse{Success: true}, nil
}

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	tdConn, err := grpc.Dial("localhost:50050", opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer tdConn.Close()

	tsConn, err := grpc.Dial("localhost:50051", opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer tsConn.Close()

	tsClient := tokenfuncs.NewTokenStoreClient(tsConn)
	tdClient := tokenfuncs.NewTSDaemonClient(tdConn)
	tf := tokenfuncs.NewTokenFuncs(tsClient, tdClient, "test", tokenfuncs.WithAsync())

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(tokenapi.UnaryServerInterceptor(tf)),
		grpc.StreamInterceptor(tokenapi.StreamServerInterceptor(tf)),
	)
	pb.RegisterTokenAPIExampleServer(s, &server{})

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v\n", err)
	}
}
