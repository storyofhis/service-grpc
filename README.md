# gRPC with Go (Golang) 🚀
make a simple system for user management service. this project is self study to get knowledge for basic gRPC with GO (golang). 

used tools : 
> make `protoc` as protobuf in this project 

propoc :
```
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative usermgmt/usermgmt.proto
```

* if you want to change protoc file. you just run `script`
> bash paths.sh

- run you server service 
```
go run usermgmt_server/usermgmt_server.go
```
- run your client service
```
go run usermgmt_client/usermgmt_client.go
```
