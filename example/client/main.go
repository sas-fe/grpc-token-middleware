package main

import (
	"log"

	"golang.org/x/net/context"

	pb "github.com/sas-fe/grpc-token-middleware/example/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial("localhost:50052", opts...)
	if err != nil {
		log.Fatal("Error connecting to server: ", err)
	}
	defer conn.Close()
	client := pb.NewTokenAPIExampleClient(conn)

	md := map[string][]string{
		"authorization": []string{"key1"},
	}
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	for i := 0; i < 100; i++ {
		_, err = client.Do(ctx, &pb.DoRequest{Do: true})
		if err != nil {
			log.Fatal("Error executing: ", err)
		}
	}
}
