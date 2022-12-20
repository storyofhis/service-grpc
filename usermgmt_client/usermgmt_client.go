package main

import (
	"context"
	"log"
	"time"

	pb "github.com/storyofhis/basic-grpc/usermgmt"
	"google.golang.org/grpc"
)

const (
	address = "localhost:8080"
)

func main() {

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect : %v", err)
		return
	}
	defer conn.Close()

	c := pb.NewUserManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var new_users = make(map[string]uint32)
	new_users["Alice "] = 43
	new_users["Bob"] = 30

	for name, age := range new_users {
		r, err := c.CreateNewUser(ctx, &pb.NewUser{Name: name, Age: int32(age)})
		if err != nil {
			log.Fatalf("Could not create user %v", err)
			return
		}
		log.Printf(
			`User Details : 
				NAME : %s
				AGE : %d
				ID : %d
			`,
			r.GetName(), r.GetAge(), r.GetId(),
		)
	}
}
