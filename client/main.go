package main

import (
	"fmt"
	"net"
	"net/http"
	"time"

	pb "github.com/k2wanko/appengine-grpc-example/helloworld"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/socket"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	appengine.Main()
}

func init() {
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	conn, err := grpc.Dial(address,
		grpc.WithDialer(func(addr string, timeout time.Duration) (conn net.Conn, err error) {
			log.Debugf(ctx, "WithDialer: %s %v", addr, timeout)
			return socket.DialTimeout(ctx, "tcp", addr, timeout)
		}),
		grpc.WithInsecure())
	if err != nil {
		log.Criticalf(ctx, "did not connect: %v", err)
		http.Error(w, fmt.Sprintf("did not connect: %v", err), 500)
		return
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	name := appengine.RequestID(ctx)
	res, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Criticalf(ctx, "could not greet: %v", err)
		http.Error(w, fmt.Sprintf("could not greet: %v", err), 500)
		return
	}
	fmt.Fprintf(w, "Greeting: %s", res.Message)
}
