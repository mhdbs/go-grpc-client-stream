package main

import (
	"fmt"
	"io"
	"log"
	greet "mhdbs/go-grpc-client-stream/pb"
	"net"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) GreetManyTimes(req *greet.GreetManyTimesRequest, stream greet.GreetService_GreetManyTimesServer) error {
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello" + firstName + "Number" + strconv.Itoa(i)
		res := &greet.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)

	}
	return nil
}

func (*server) LongGreet(stream greet.GreetService_LongGreetServer) error {
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&greet.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("error %v", err)
		}
		firstname := req.GetGreeting().GetFirstName()
		result += "Helo " + firstname + "! "
		fmt.Println("GOt the data", result)
	}
}

func main() {
	fmt.Println("Server program")
	lis, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		fmt.Println("Failed to listen ", err)
	}
	s := grpc.NewServer()
	greet.RegisterGreetServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
