package methods

import (
	"context"

	"github.com/rs/zerolog/log"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

// SayHello implements helloworld.GreeterServer.SayHello
func SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}
