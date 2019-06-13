package store

import (
	"github.com/jinzhu/gorm"
	"github.com/param108/grpc-chat-server/errors"
	"github.com/param108/grpc-chat-server/models"
	"time"
)

func (store *ChatStore) CreateChatGroup(chatName string) (*models.ChatGroup, errors.GrpcChatError) {
	chatGroup := models.ChatGroup{ChatName: chatName}
	err := store.db.Create(&chatGroup).Error
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	return &chatGroup, nil
}

func (store *ChatStore) FindChatGroup(chatID string) (*models.ChatGroup, errors.GrpcChatError) {
	chatGroup := models.ChatGroup{ID: chatID}
	err := store.db.Find(&chatGroup).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewNotFoundError(err)
		}
		return nil, errors.NewInternalError(err)
	}
	return &chatGroup, nil
}

// Assumes both userID and chatID are valid
func (store *ChatStore) AddUserToChat(userID, chatID string) errors.GrpcChatError {
	userGroup := models.UserGroup{UserID: userID, ChatGroupID: chatID}
	err := store.db.Create(&userGroup).Error
	if err != nil {
		findErr := store.db.Find(&userGroup, "user_id = ? and chat_group_id = ?", userID, chatID).Error
		if findErr != nil {
			return errors.NewInternalError(err)
		}
		return errors.NewAlreadyExistsError(err)
	}
	return nil
}

func (store *ChatStore) ListUsersOfChat(chatID string) ([]string, errors.GrpcChatError) {
	ret := []string{}
	userGroups := []models.UserGroup{}
	err := store.db.Table("user_groups").Select([]string{"user_id"}).Where("chat_group_id = ?", chatID).Find(&userGroups).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return ret, errors.NewNotFoundError(err)
		}
		return ret, errors.NewInternalError(err)
	}

	for i := 0; i < len(userGroups); i++ {
		ret = append(ret, userGroups[i].UserID)
	}
	return ret, nil
}

func (store *ChatStore) AddMessageToChat(userID, chatID, clientTime, message string) errors.GrpcChatError {
	ct, err := time.Parse(time.RFC1123Z, clientTime)
	if err != nil {
		return errors.NewInvalidInputError(err)
	}
	chatMessage := models.ChatMessage{ClientTime: ct, Message: message, UserID: userID, ChatGroupID: chatID,
		ServerTime: time.Now().UTC()}
	err = store.db.Create(&chatMessage).Error
	if err != nil {
		return errors.NewInternalError(err)
	}
	return nil
}

func (store *ChatStore) GetMessagesAfter(chatID, clientTime string) ([]models.ChatMessage, errors.GrpcChatError) {
	ret := []models.ChatMessage{}
	ct, err := time.Parse(time.RFC1123Z, clientTime)
	if err != nil {
		return ret, errors.NewInvalidInputError(err)
	}

	if err := store.db.Where("client_time > ?", ct).Find(&ret).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return ret, errors.NewNotFoundError(err)
		}
		return ret, errors.NewInternalError(err)
	}
	return ret, nil
}
