package chatclient

import (
	"context"
	"errors"
	"fmt"
	"github.com/param108/grpc-chat-server/chat"
	"google.golang.org/grpc"
)

type MemoryChatClient struct {
	grpcClient chat.ChatClient
	Port       int
}

func NewMemoryChatClient(port int) *MemoryChatClient {
	return &MemoryChatClient{
		Port: port,
	}
}

func (m *MemoryChatClient) WriteMessage(ctx context.Context, in *chat.Message, opts ...grpc.CallOption) (*chat.Response, error) {

	return nil, nil
}

func (m *MemoryChatClient) ReadMessages(ctx context.Context, in *chat.TimeDesc, opts ...grpc.CallOption) (chat.Chat_ReadMessagesClient, error) {
	return nil, nil
}

func (m *MemoryChatClient) read() {
	ctx := context.TODO()
	recv, err := m.grpcClient.ReadMessages(ctx, &chat.TimeDesc{})
	if err != nil {
		fmt.Printf("Failed to read %v", err)
		return
	}
	msg, err := recv.Recv()
	for err == nil {
		fmt.Println(msg.Data)
		msg, err = recv.Recv()
	}

	if err != nil {
		fmt.Printf("Failed to read %v", err)
		return
	}

}

func (m *MemoryChatClient) write(msg string) {
	ctx := context.TODO()
	wrappedMsg := &chat.Message{Data: msg}
	_, err := m.grpcClient.WriteMessage(ctx, wrappedMsg)
	if err != nil {
		fmt.Printf("Failed to write: %v", err)
		return
	}
}

func (m *MemoryChatClient) connect(conn *grpc.ClientConn) {
	m.grpcClient = chat.NewChatClient(conn)
}

func (m *MemoryChatClient) Start(cmd []string) error {
	serverAddr := fmt.Sprintf("localhost:%d", m.Port)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("fail to dial: %v", err)
	}
	defer conn.Close()

	m.connect(conn)

	switch cmd[0] {
	case "read":
		m.read()
	case "write":
		m.write(cmd[1])
	default:
		fmt.Printf("Invalid command %s", cmd[0])
	}
	return nil

}

// Returns the UserToken for the user.
func (m *MemoryChatClient) Login(username string, firebaseKey string) (string, error) {
	ctx := context.TODO()
	response, err := m.grpcClient.Login(ctx, &chat.LoginRequest{Username: username, FirebaseKey: firebaseKey})
	if err != nil {
		return "", err
	}

	if !response.Success {
		return "", errors.New(response.Error)
	}

	return response.UserToken, nil
}

// returns the chat ID
func (m *MemoryChatClient) CreateChat(userToken string, chatName string) (string, error) {
	ctx := context.TODO()
	response, err := m.grpcClient.CreateChat(ctx,
		&chat.CreateChatRequest{UserToken: userToken, ChatName: chatName})
	if err != nil {
		return "", err
	}

	if !response.Success {
		return "", errors.New(response.Error)
	}

	return response.ChatID, nil
}

func (m *MemoryChatClient) ListChats(userToken string) ([]*chat.ChatDetail, error) {
	ctx := context.TODO()

	ret := []*chat.ChatDetail{}
	chatListSocket, err := m.grpcClient.ListChats(ctx, &chat.ListChatRequest{UserToken: userToken})
	if err != nil {
		return ret, err
	}

	chatDetail, err := chatListSocket.Recv()
	for err != nil {
		ret = append(ret, chatDetail)
		chatDetail, err = chatListSocket.Recv()
	}

	return ret, nil
}

func (m *MemoryChatClient) JoinChat(userToken string, chatID string) error {
	ctx := context.TODO()
	req := &chat.JoinChatRequest{UserToken: userToken, ChatID: chatID}
	response, err := m.grpcClient.JoinChat(ctx, req)
	if err != nil {
		return err
	}

	if !response.Success {
		return errors.New(response.Error)
	}

	return nil
}

func (m *MemoryChatClient) SetAvailability(context.Context, *chat.SetAvailableRequest) (*chat.SetAvailableResponse, error) {
	return nil, nil
}

func (m *MemoryChatClient) ListAvailability(*chat.ListAvailableRequest, chat.Chat_ListAvailabilityServer) error {
	return nil
}
