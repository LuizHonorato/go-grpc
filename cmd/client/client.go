package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/codeedu/fc2-grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect to gRPC Server: %v", err)
	}

	defer connection.Close()

	client := pb.NewUserServiceClient(connection)
	// AddUser(client)
	// AddUserVerbose(client)
	// AddUsers(client)
	AddUserStreamBoth(client)
}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Luizão",
		Email: "l@l.com",
	}

	res, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make to gRPC request: %v", err)
	}

	log.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Luizão",
		Email: "l@l.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not make to gRPC request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not receive the msg: %v", err)
		}

		fmt.Println("Status: ", stream.Status, " - ", stream.GetUser())
	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		&pb.User{
			Id:    "l1",
			Name:  "Luiz",
			Email: "luiz@luiz.com",
		},
		&pb.User{
			Id:    "l2",
			Name:  "Luiz2",
			Email: "luiz2@luiz.com",
		},
		&pb.User{
			Id:    "l3",
			Name:  "Luiz3",
			Email: "luiz3@luiz.com",
		},
		&pb.User{
			Id:    "l4",
			Name:  "Luiz4",
			Email: "luiz4@luiz.com",
		},
		&pb.User{
			Id:    "l5",
			Name:  "Luiz5",
			Email: "luiz5@luiz.com",
		},
		&pb.User{
			Id:    "l6",
			Name:  "Luiz6",
			Email: "luiz6@luiz.com",
		},
	}

	stream, err := client.AddUsers(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient) {
	stream, err := client.AddUserStreamBoth(context.Background())
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	reqs := []*pb.User{
		&pb.User{
			Id:    "l1",
			Name:  "Luiz",
			Email: "luiz@luiz.com",
		},
		&pb.User{
			Id:    "l2",
			Name:  "Luiz2",
			Email: "luiz2@luiz.com",
		},
		&pb.User{
			Id:    "l3",
			Name:  "Luiz3",
			Email: "luiz3@luiz.com",
		},
		&pb.User{
			Id:    "l4",
			Name:  "Luiz4",
			Email: "luiz4@luiz.com",
		},
		&pb.User{
			Id:    "l5",
			Name:  "Luiz5",
			Email: "luiz5@luiz.com",
		},
		&pb.User{
			Id:    "l6",
			Name:  "Luiz6",
			Email: "luiz6@luiz.com",
		},
	}

	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user: ", req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}

		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving data: %v", err)
				break
			}
			fmt.Printf("Receiving user %v com status: %v\n", res.GetUser().GetName(), res.GetStatus())
		}

		close(wait)
	}()

	<-wait
}
