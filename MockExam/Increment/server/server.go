package main

import (
	"MockExam/Increment/protobuf"
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"
)

var number int64 = 0

type server struct {
	protobuf.UnimplementedIncrementServer
}

func main() {
	log.Print("Welcome Server. You need to provide a name:")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	name := strings.Replace(text, "\n", "", 1)
	port := strings.Replace(name, "S", "", 1)

	lis, err := net.Listen("tcp", ":808"+port)

	if err != nil { //error before listening
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer() //we create a new server
	protobuf.RegisterIncrementServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil { //error while listening
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) NewNode(ctx context.Context, in *protobuf.NewNodeRequest) (*protobuf.NewNodeReply, error) {
	return &protobuf.NewNodeReply{}, nil
}

func (s *server) NewIncrement(ctx context.Context, in *protobuf.NewIncrementRequest) (*protobuf.NewIncrementReply, error) {
	number += 1
	fmt.Println("my current value is: " + strconv.Itoa(int(number)))
	return &protobuf.NewIncrementReply{Answer: number}, nil
}

func (s *server) NewUpdateNumbers(ctx context.Context, in *protobuf.NewUpdateNumbersRequest) (*protobuf.NewUpdateNumbersReply, error) {
	number = in.UpdatedNumber
	fmt.Println("updated this to: " + strconv.Itoa(int(number)))
	return &protobuf.NewUpdateNumbersReply{}, nil
}
