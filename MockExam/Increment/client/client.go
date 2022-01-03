package main

import (
	"MockExam/Increment/protobuf"
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
)

func main() {
	log.Print("Welcome Client")

	conn, err := grpc.Dial(":8085", grpc.WithInsecure(), grpc.WithBlock())

	if err != nil { //error can not establish connection
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := protobuf.NewIncrementClient(conn)

	go TakeInput(client)
	time.Sleep(1000 * time.Second)
}

func TakeInput(client protobuf.IncrementClient) {
	for {
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		inputParsed := strings.Replace(input, "\n", "", 1)
		if inputParsed == "inc" {
			var result, _ = client.NewIncrement(context.Background(), &protobuf.NewIncrementRequest{})
			fmt.Println(result.Answer)
		}
	}
}
