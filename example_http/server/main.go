package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/golang/glog"
	"github.com/sas-fe/grpc-token-middleware/tokenfuncs"
	"google.golang.org/grpc"
)

func do(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func main() {
	flag.Parse()

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

	http.Handle("/", tf.Middleware(http.HandlerFunc(do)))

	glog.V(0).Info("Serving on localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
