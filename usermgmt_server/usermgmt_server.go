package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jackc/pgx/v4"
	pb "github.com/storyofhis/basic-grpc/usermgmt"
	"google.golang.org/grpc"
)

const (
	port = ":8080"
)

type UserManagementServer struct {
	conn *pgx.Conn
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
	createSql := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name TEXT,
			age INT
		);
	`
	_, err := server.conn.Exec(context.Background(), createSql)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Table creation failed: %v\n", err)
		os.Exit(1)
	}

	// readBytes, err := ioutil.ReadFile("users.json")
	// var user_list *pb.UserList = &pb.UserList{}
	// var user_id int32 = int32(rand.Intn(1000))
	// created_user := &pb.User{Name: in.GetName(), Age: in.GetAge(), Id: user_id}
	created_user := &pb.User{Name: in.GetName(), Age: in.GetAge()}
	// server.user_list.Users = append(server.user_list.Users, created_user)
	// if err != nil {
	// 	if os.IsNotExist(err) {
	// 		log.Print("File not found. Creating a new file")
	// 		user_list.Users = append(user_list.Users, created_user)
	// 		jsonBytes, err := protojson.Marshal(
	// 			user_list,
	// 		)
	// 		if err != nil {
	// 			log.Fatalf("JSON Marshaling failed: %v", err)
	// 		}
	// 		if err := ioutil.WriteFile("users.json", jsonBytes, 0664); err != nil {
	// 			log.Fatalf("Failed write to file: %v", err)
	// 		}
	// 		return created_user, nil
	// 	} else {
	// 		log.Fatalln("Error reading File : ", err)
	// 	}
	// }

	// if err := protojson.Unmarshal(readBytes, user_list); err != nil {
	// 	log.Fatalf("Failed to parse user list: %v", err)
	// }

	// user_list.Users = append(user_list.Users, created_user)
	// jsonBytes, err := protojson.Marshal(
	// 	user_list,
	// )
	// if err != nil {
	// 	log.Fatalf("JSON Marshaling failed: %v", err)
	// }
	// if err := ioutil.WriteFile("users.json", jsonBytes, 0664); err != nil {
	// 	log.Fatalf("Failed write to file: %v", err)
	// }

	tx, err := server.conn.Begin(context.Background())
	if err != nil {
		log.Fatalf("conn.Begin failed: %v", err)
	}
	_, err = tx.Exec(
		context.Background(),
		`INSERT INTO users (name, age) VALUES ($1, $2)`,
		created_user.Name, created_user.Age,
	)
	if err != nil {
		log.Fatalf("tx.exec failed: %v", err)
	}
	tx.Commit(context.Background())
	return created_user, nil
}

func (server *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UserList, error) {
	var user_list *pb.UserList = &pb.UserList{}
	// jsonBytes, err := ioutil.ReadFile("users.json")
	// if err != nil {
	// 	log.Fatalf("Failed to read from file: %v", err)
	// }
	rows, err := server.conn.Query(
		context.Background(),
		`SELECT * FROM users`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		user := pb.User{}
		err = rows.Scan(&user.Id, &user.Name, &user.Age)
		if err != nil {
			return nil, err
		}
		user_list.Users = append(user_list.Users, &user)
	}

	// if err := protojson.Unmarshal(jsonBytes, user_list); err != nil {
	// 	log.Fatalf("Unmarshaling failed: %v", err)
	// }
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
	database_url := "postgres://postgres:0000@localhost:5432/postgres"
	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil {
		log.Fatalf("Unable to establish connection: %v", err)
	}
	defer conn.Close(context.Background())
	var user_mgmt_server *UserManagementServer = NewUserManagementServer()
	user_mgmt_server.conn = conn
	if err := user_mgmt_server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
