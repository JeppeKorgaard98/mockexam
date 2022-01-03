package main

import (
	"MockExam/Increment/protobuf"
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

type server struct {
	protobuf.UnimplementedIncrementServer
}

const amountOfServers int = 3

var FrontEndConns [amountOfServers]protobuf.IncrementClient
var Conns [amountOfServers]*grpc.ClientConn

func main() {
	log.Print("Welcome Frontend.")

	go FrontendServerStart()
	for i := 0; i < amountOfServers; i++ {
		portToDial := 8080 + i + 1
		FrontEndConns[i], Conns[i] = Dial(portToDial)
		defer Conns[i].Close()
	}

	time.Sleep(1000 * time.Second)
}

func FrontendServerStart() {
	lis, err := net.Listen("tcp", ":8085")

	if err != nil { //error before listening
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer() //we create a new server
	protobuf.RegisterIncrementServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func Dial(port int) (protobuf.IncrementClient, *grpc.ClientConn) {
	conn, err := grpc.Dial(":"+strconv.Itoa(port), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil { //error can not establish connection
		log.Fatalf("did not connect: %v", err)
	}

	frontend := protobuf.NewIncrementClient(conn)
	message, userAlreadyExistsError := frontend.NewNode(context.Background(), &protobuf.NewNodeRequest{})
	if userAlreadyExistsError != nil {
		if message == nil {
			fmt.Println("Username is already in use")
		}
	} else {
		fmt.Println("Dial to " + strconv.Itoa(port) + " was succesful")
		return frontend, conn
	}
	return nil, nil
}

func (s *server) NewIncrement(ctx context.Context, in *protobuf.NewIncrementRequest) (*protobuf.NewIncrementReply, error) {
	responsesFromServers := make([]int64, amountOfServers)
	for i := 0; i < amountOfServers; i++ {
		var response, err = FrontEndConns[i].NewIncrement(context.Background(), &protobuf.NewIncrementRequest{})
		if err == nil {
			responsesFromServers[i] = response.Answer
		}
	}
	var valueToReturn int64 = validatedResponse(responsesFromServers)
	fmt.Println("received request, highest number from validated from server is " + strconv.Itoa(int(valueToReturn)))
	return &protobuf.NewIncrementReply{Answer: valueToReturn}, nil
}

func validatedResponse(list []int64) int64 {
	var highestNumber int64
	for i := 0; i < len(list); i++ {
		if list[i] > highestNumber {
			highestNumber = list[i]
		}
	}
	updateAllServers(highestNumber)
	return highestNumber
}

func updateAllServers(number int64) {
	for i := 0; i < amountOfServers; i++ {
		var _, _ = FrontEndConns[i].NewUpdateNumbers(context.Background(), &protobuf.NewUpdateNumbersRequest{UpdatedNumber: number})
	}
}
