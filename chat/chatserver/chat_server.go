package chatserver

import (
	"context"
	//std_errors "errors"
	"fmt"
	"github.com/param108/grpc-chat-server/chat"
	"github.com/param108/grpc-chat-server/chatmanager"
	//"github.com/param108/grpc-chat-server/errors"
	"github.com/param108/grpc-chat-server/store"
	"google.golang.org/grpc"
	"net"
)

type ChatServerImpl struct {
	Port int
	Msgs []string
	DB   store.Store
	CM   chatmanager.ChatManager
}

func NewChatServer(port int) (*ChatServerImpl, error) {
	db, err := store.NewChatStore()
	if err != nil {
		return nil, err
	}
	CM := chatmanager.NewSimpleChatManager()
	return &ChatServerImpl{Port: port, DB: db, CM: CM}, nil
}

func (server *ChatServerImpl) WriteMessage(ctx context.Context, msg *chat.Message) (*chat.Response, error) {
	response := &chat.Response{}
	err := server.CM.Write(msg.ChatID, msg)
	if err != nil {
		response.Success = false
		response.Error = "FailToWrite"
		fmt.Println("Failed to write")
		return response, err
	}
	server.Msgs = append(server.Msgs, msg.Data)
	response.Success = true
	fmt.Println("success write:" + msg.Data)
	return response, nil
}

func (server *ChatServerImpl) ReadMessages(t *chat.TimeDesc, chServer chat.Chat_ReadMessagesServer) error {
	fmt.Println(len(server.Msgs))
	sub, err := server.CM.Read(t.ChatID)
	if err != nil {
		return err
	}

	for {
		select {
		case v := <-sub.ReadChan:
			if v == nil {
				fmt.Println("Failing ReadMessages center closed connection")
				return nil
			}
			fmt.Println(v)
			err := chServer.Send(v.(*chat.Message))
			if err != nil {
				fmt.Println("Failing ReadMessages Cannot Send")
				sub.Quit()
				return nil
			}
		}
	}
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
	user, err := server.DB.FindUser(loginRequest.Username)
	if err != nil {
		userVal, err := server.DB.CreateUser(loginRequest.Username, loginRequest.FirebaseKey,
			"user")
		if err != nil {
			return &chat.LoginResponse{Success: false, Error: err.Code()}, err
		}

		user = userVal
	}

	userToken, err := server.DB.CreateUserToken(user.ID)
	if err != nil {
		return &chat.LoginResponse{Success: false, Error: err.Code()}, err
	}

	return &chat.LoginResponse{UserToken: userToken.UserToken, Success: true, Error: ""}, nil
}

func (server *ChatServerImpl) CreateChat(ctx context.Context, createChatRequest *chat.CreateChatRequest) (*chat.CreateChatResponse, error) {
	user, err := server.DB.FindUserFromToken(createChatRequest.UserToken)
	if err != nil {
		return &chat.CreateChatResponse{Success: false, Error: err.Code()}, err
	}

	chatGroup, err := server.DB.CreateChatGroup(createChatRequest.ChatName)
	if err != nil {
		return &chat.CreateChatResponse{Success: false, Error: err.Code()}, err
	}

	err = server.DB.AddUserToChat(user.ID, chatGroup.ID)
	if err != nil {
		return &chat.CreateChatResponse{Success: false, Error: err.Code()}, err
	}
	return &chat.CreateChatResponse{Success: true, Error: "", ChatID: chatGroup.ID}, nil
}

func (server *ChatServerImpl) ListChats(listChatRequest *chat.ListChatRequest, callbacks chat.Chat_ListChatsServer) error {
	_, err := server.DB.FindUserFromToken(listChatRequest.UserToken)
	if err != nil {
		return err
	}

	// FIXME removing this check for now
	// if user.Role != "admin" {
	// 	err = errors.NewForbiddenError(std_errors.New("Not Admin"))
	// 	return err
	// }

	chatGroups, err := server.DB.ListOpenChatGroups()
	if err != nil {
		return err
	}

	for i := 0; i < len(chatGroups); i++ {
		userList, err := server.DB.ListUsersOfChat(chatGroups[i].ID)
		if err != nil {
			return err
		}
		sendErr := callbacks.Send(&chat.ChatDetail{UserList: userList, ChatName: chatGroups[i].ChatName,
			ChatID: chatGroups[i].ID})
		if sendErr != nil {
			return sendErr
		}
	}
	return nil
}

func (server *ChatServerImpl) JoinChat(ctx context.Context, joinChatRequest *chat.JoinChatRequest) (*chat.JoinChatResponse, error) {
	user, err := server.DB.FindUserFromToken(joinChatRequest.UserToken)
	if err != nil {
		return nil, err
	}

	// FIXME removing this check for now
	// if user.Role != "admin" {
	//	err = errors.NewForbiddenError(std_errors.New("Not Admin"))
	//	return &chat.JoinChatResponse{Success: false, Error: err.Code()}, err
	// }

	err = server.DB.AddUserToChat(user.ID, joinChatRequest.ChatID)
	if err != nil {
		return &chat.JoinChatResponse{Success: false, Error: err.Code()}, err
	}
	return &chat.JoinChatResponse{Success: true, Error: ""}, nil
}

func (server *ChatServerImpl) SetAvailability(context.Context, *chat.SetAvailableRequest) (*chat.SetAvailableResponse, error) {
	return nil, nil
}

func (server *ChatServerImpl) ListAvailability(*chat.ListAvailableRequest, chat.Chat_ListAvailabilityServer) error {
	return nil
}
