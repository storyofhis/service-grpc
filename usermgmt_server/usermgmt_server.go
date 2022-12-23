package main

import (
	"context"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"

	pb "github.com/storyofhis/basic-grpc/usermgmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	port = ":8080"
)

type UserManagementServer struct {
	pb.UnimplementedUserManagementServer
	// user_list *pb.UserList

}

func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{
		// user_list: &pb.UserList{},

	}
}

func (server *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	log.Printf("Received: %v\n", in.GetName())
	readBytes, err := ioutil.ReadFile("users.json")
	var user_list *pb.UserList = &pb.UserList{}
	var user_id int32 = int32(rand.Intn(1000))
	created_user := &pb.User{Name: in.GetName(), Age: in.GetAge(), Id: user_id}
	// server.user_list.Users = append(server.user_list.Users, created_user)
	if err != nil {
		if os.IsNotExist(err) {
			log.Print("File not found. Creating a new file")
			user_list.Users = append(user_list.Users, created_user)
			jsonBytes, err := protojson.Marshal(
				user_list,
			)
			if err != nil {
				log.Fatalf("JSON Marshaling failed: %v", err)
			}
			if err := ioutil.WriteFile("users.json", jsonBytes, 0664); err != nil {
				log.Fatalf("Failed write to file: %v", err)
			}
			return created_user, nil
		} else {
			log.Fatalln("Error reading File : ", err)
		}
	}

	if err := protojson.Unmarshal(readBytes, user_list); err != nil {
		log.Fatalf("Failed to parse user list: %v", err)
	}

	user_list.Users = append(user_list.Users, created_user)
	jsonBytes, err := protojson.Marshal(
		user_list,
	)
	if err != nil {
		log.Fatalf("JSON Marshaling failed: %v", err)
	}
	if err := ioutil.WriteFile("users.json", jsonBytes, 0664); err != nil {
		log.Fatalf("Failed write to file: %v", err)
	}
	return created_user, nil
}

func (server *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UserList, error) {
	jsonBytes, err := ioutil.ReadFile("users.json")
	if err != nil {
		log.Fatalf("Failed to read from file: %v", err)
	}
	var user_list *pb.UserList = &pb.UserList{}
	if err := protojson.Unmarshal(jsonBytes, user_list); err != nil {
		log.Fatalf("Unmarshaling failed: %v", err)
	}
	// return server.user_list, nil
	return user_list, nil
}

func (server *UserManagementServer) Run() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserManagementServer(s, server)
	log.Printf("server listening at %v", lis.Addr())
	return s.Serve(lis)
}

func main() {
	var user_mgmt_server *UserManagementServer = NewUserManagementServer()
	if err := user_mgmt_server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
