package main

import (
	"context"
	"fmt"
	"log"
	greet "mhdbs/go-grpc-client-stream/pb"
	"time"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("client program")
	cc, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
		defer cc.Close()
	}
	c := greet.NewGreetServiceClient(cc)
	// doServerStreaming(c)
	doClientStreaming(c)

}

func doClientStreaming(c greet.GreetServiceClient) {
	requests := []*greet.LongGreetRequest{
		&greet.LongGreetRequest{
			Greeting: &greet.Greeting{
				FirstName: "John",
			},
		},
		&greet.LongGreetRequest{
			Greeting: &greet.Greeting{
				FirstName: "Doe",
			},
		},
		&greet.LongGreetRequest{
			Greeting: &greet.Greeting{
				FirstName: "Done",
			},
		},
	}
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	for _, req := range requests {
		fmt.Println("Sending data", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Final response", res)
}

// func doServerStreaming(c greet.GreetServiceClient) {
// 	fmt.Println("server streaming roc")
// 	req := &greet.GreetManyTimesRequest{
// 		Greeting: &greet.Greeting{
// 			FirstName: "bilal",
// 			LastName:  "mohammed",
// 		},
// 	}
// 	resStream, err := c.GreetManyTimes(context.Background(), req)
// 	if err != nil {
// 		log.Fatalf("Error while calling greetmany times RPC: %v", err)
// 	}
// 	for {
// 		msg, err := resStream.Recv()
// 		if err == io.EOF {
// 			fmt.Println("End of the stream ")
// 			break
// 		}
// 		if err != nil {
// 			log.Fatalf("Error while  getting the stream %v", err)
// 		}
// 		fmt.Println("Response from greet many times %v", msg.GetResult())
// 	}

// }
