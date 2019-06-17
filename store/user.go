package store

import (
	"github.com/jinzhu/gorm"
	"github.com/param108/grpc-chat/errors"
	"github.com/param108/grpc-chat/models"
	"github.com/satori/go.uuid"
)

func (store *ChatStore) FindUser(username string) (*models.User, errors.GrpcChatError) {
	user := models.User{}
	err := store.db.Where("name = ?", username).Find(&user).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewNotFoundError(err)
		}
		return nil, errors.NewInternalError(err)
	}
	return &user, nil
}

func (store *ChatStore) CreateUser(username string, firebaseToken string, userRole string) (*models.User, errors.GrpcChatError) {
	createdUser := &models.User{Name: username, FirebaseToken: firebaseToken, Role: userRole}
	createErr := store.db.Create(createdUser).Error
	if createErr != nil {
		// check if user already exists
		_, err := store.FindUser(username)
		if err != nil {
			return nil, errors.NewInternalError(createErr)
		}
		return nil, errors.NewAlreadyExistsError(createErr)
	}
	return createdUser, nil
}

func (store *ChatStore) UpdateFirebaseToken(userID string, firebaseToken string) errors.GrpcChatError {
	user := models.User{ID: userID}
	if err := store.db.Find(&user).Error; err != nil {
		return errors.NewNotFoundError(err)
	}

	user.FirebaseToken = firebaseToken
	if err := store.db.Save(&user).Error; err != nil {
		return errors.NewInternalError(err)
	}
	return nil
}

func (store *ChatStore) CreateUserToken(userID string) (*models.UserToken, errors.GrpcChatError) {
	user := models.User{ID: userID}
	if err := store.db.Find(&user).Error; err != nil {
		return nil, errors.NewNotFoundError(err)
	}
	userToken := models.UserToken{}
	store.db.Delete(&userToken, "user_id = ?", user.ID)
	userToken = models.UserToken{UserID: user.ID, UserToken: uuid.NewV4().String(),
		FirebaseToken: user.FirebaseToken, UserRole: user.Role}
	err := store.db.Create(&userToken).Error
	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	return &userToken, nil
}

func (store *ChatStore) FindUserToken(token string) (*models.UserToken, errors.GrpcChatError) {
	userToken := models.UserToken{}
	err := store.db.Find(&userToken, "user_token = ?", token).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, errors.NewNotFoundError(err)
		}
		return nil, errors.NewInternalError(err)
	}

	return &userToken, nil
}

func (store *ChatStore) FindUserFromToken(token string) (*models.User, errors.GrpcChatError) {
	userToken, err := store.FindUserToken(token)
	if err != nil {
		return nil, err
	}

	user := models.User{ID: userToken.UserID}
	findErr := store.db.Find(&user).Error
	if findErr != nil {
		if gorm.IsRecordNotFoundError(findErr) {
			return nil, errors.NewNotFoundError(findErr)
		}
		return nil, errors.NewInternalError(findErr)
	}

	return &user, nil
}
