package chatserver

import (
	"context"
	"fmt"
	"github.com/param108/grpc-chat-server/chat"
	"github.com/param108/grpc-chat-server/store"
	"google.golang.org/grpc"
	"net"
	"time"
)

type ChatServerImpl struct {
	Port int
	Msgs []string
	DB   store.Store
}

func NewChatServer(port int) (*ChatServerImpl, error) {
	db, err := store.NewChatStore()
	if err != nil {
		return nil, err
	}
	return &ChatServerImpl{Port: port, DB: db}, nil

}

func (server *ChatServerImpl) WriteMessage(ctx context.Context, msg *chat.Message) (*chat.Response, error) {
	response := &chat.Response{}
	server.Msgs = append(server.Msgs, msg.Data)
	return response, nil
}

func (server *ChatServerImpl) ReadMessages(t *chat.TimeDesc, chServer chat.Chat_ReadMessagesServer) error {
	fmt.Println(len(server.Msgs))
	for i := 0; i < len(server.Msgs); i++ {
		chServer.Send(&chat.Message{Data: server.Msgs[i]})
		time.Sleep(5 * time.Second)
	}
	return nil
}

func (server *ChatServerImpl) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", server.Port))
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return err
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	chat.RegisterChatServer(grpcServer, server)
	grpcServer.Serve(lis)
	return nil
}

func (server *ChatServerImpl) Login(ctx context.Context, loginRequest *chat.LoginRequest) (*chat.LoginResponse, error) {
	_, err := server.DB.FindUser(loginRequest.Username)
	if err != nil {

	}

	return nil, nil
}

func (server *ChatServerImpl) CreateChat(context.Context, *chat.CreateChatRequest) (*chat.CreateChatResponse, error) {
	return nil, nil
}

func (server *ChatServerImpl) ListChats(*chat.ListChatRequest, chat.Chat_ListChatsServer) error {
	return nil
}

func (server *ChatServerImpl) JoinChat(context.Context, *chat.JoinChatRequest) (*chat.JoinChatResponse, error) {
	return nil, nil
}

func (server *ChatServerImpl) SetAvailability(context.Context, *chat.SetAvailableRequest) (*chat.SetAvailableResponse, error) {
	return nil, nil
}

func (server *ChatServerImpl) ListAvailability(*chat.ListAvailableRequest, chat.Chat_ListAvailabilityServer) error {
	return nil
}
