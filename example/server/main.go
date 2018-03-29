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
	tsConn, err := grpc.Dial("localhost:50051", opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer tsConn.Close()
	tf := tokenfuncs.NewTokenFuncs(tsConn, "test", tokenfuncs.WithAsync())
	// tf := tokenfuncs.NewTokenFuncs(tsConn, "test")

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
