package store

import (
	"github.com/param108/grpc-chat/errors"
	"github.com/param108/grpc-chat/models"
)

type Store interface {
	FindUser(username string) (*models.User, errors.GrpcChatError)
	CreateUser(username, firebaseToken, userRole string) (*models.User, errors.GrpcChatError)
	UpdateFirebaseToken(userID string, firebaseToken string) errors.GrpcChatError
	CreateUserToken(userID string) (*models.UserToken, errors.GrpcChatError)
	FindUserToken(userToken string) (*models.UserToken, errors.GrpcChatError)
	FindUserFromToken(userToken string) (*models.User, errors.GrpcChatError)
	CreateChatGroup(chatName string) (*models.ChatGroup, errors.GrpcChatError)
	FindChatGroup(chatID string) (*models.ChatGroup, errors.GrpcChatError)
	AddUserToChat(userID, chatID string) errors.GrpcChatError
	ListUsersOfChat(chatID string) ([]string, errors.GrpcChatError)
	AddMessageToChat(userID, chatID, clientTime, message string) errors.GrpcChatError
	GetMessagesAfter(chatID, clientTime string) ([]models.ChatMessage, errors.GrpcChatError)
	ListOpenChatGroups() ([]models.ChatGroup, errors.GrpcChatError)
}
