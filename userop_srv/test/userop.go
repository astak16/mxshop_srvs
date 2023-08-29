package main

import (
	"context"
	"mxshow_srvs/userop_srv/proto"

	"google.golang.org/grpc"
)

var userFavClient proto.UserFavClient
var addressClient proto.AddressClient
var messageClient proto.MessageClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	userFavClient = proto.NewUserFavClient(conn)
	addressClient = proto.NewAddressClient(conn)
	messageClient = proto.NewMessageClient(conn)
}

func TestAddressList() {
	_, err := addressClient.GetAddressList(context.Background(), &proto.AddressRequest{UserId: 1})
	if err != nil {
		panic(err)
	}
}

func TestMessageList() {
	_, err := messageClient.MessageList(context.Background(), &proto.MessageRequest{UserId: 1})
	if err != nil {
		panic(err)
	}
}

func TestUserFav() {
	_, err := userFavClient.GetFavList(context.Background(), &proto.UserFavRequest{UserId: 1})
	if err != nil {
		panic(err)
	}
}

func main() {
	Init()
	defer conn.Close()

	TestAddressList()
	// TestMessageList()
	// TestUserFav()
}
