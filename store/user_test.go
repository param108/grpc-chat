package store

import (
	"github.com/param108/grpc-chat/errors"
	"github.com/param108/grpc-chat/models"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserNotFound(t *testing.T) {
	defer teardown()
	_, err := testStore.FindUser("BlaBla")
	assert.NotNil(t, err, "Failed to raise error when not found")
	assert.Equal(t, errors.NotFoundError, err.Code(), "Wrong error type found:"+err.Code())
}

func TestUserFound(t *testing.T) {
	defer teardown()
	createdUser := &models.User{Name: "user_1", FirebaseToken: "firebase_token_1", Role: "user"}
	err := testStore.db.Create(createdUser).Error
	user, err := testStore.FindUser("user_1")
	assert.Nil(t, err, "Failed to find user_1")
	assert.Equal(t, createdUser.Name, user.Name, "Wrong name found:"+user.Name)
	assert.Equal(t, createdUser.FirebaseToken, user.FirebaseToken, "Wrong firebase token found:"+user.FirebaseToken)
	assert.Equal(t, createdUser.Role, user.Role, "Wrong role found:"+user.Role)
}

func TestCreateUser(t *testing.T) {
	defer teardown()
	_, err := testStore.CreateUser("user_1", "firebase_token_1", "user")
	if err != nil {
		assert.Nil(t, err, "Could not create user:"+err.Error())
	}
	user, err := testStore.FindUser("user_1")
	if err != nil {
		assert.Nil(t, err, "Could not find user:"+err.Error())
	}
	assert.Equal(t, "user_1", user.Name, "Invalid name seen"+user.Name)
	assert.Equal(t, "firebase_token_1", user.FirebaseToken, "Invalid firebase token seen"+user.FirebaseToken)
	assert.Equal(t, "user", user.Role, "Invalid role seen"+user.Role)

}

func TestCreateDuplicateUser(t *testing.T) {
	defer teardown()
	testStore.CreateUser("user_1", "firebase_token_1", "user")
	_, err := testStore.CreateUser("user_1", "firebase_token_2", "user")
	assert.NotNil(t, err, "Should not be able to create user")
	if err != nil {
		assert.Equal(t, errors.AlreadyExistsError, err.Code(), "Incorrect error code")
	}
}

func TestCreateDuplicateFirebaseToken(t *testing.T) {
	defer teardown()
	testStore.CreateUser("user_1", "firebase_token_1", "user")
	_, err := testStore.CreateUser("user_2", "firebase_token_1", "user")
	assert.NotNil(t, err, "Should not be able to create user")
	if err != nil {
		assert.Equal(t, errors.InternalError, err.Code(), "Incorrect error code")
	}
}

func TestUpdateFirebaseToken(t *testing.T) {
	defer teardown()
	user, _ := testStore.CreateUser("user_1", "firebase_token_1", "user")
	err := testStore.UpdateFirebaseToken(user.ID, "new_firebase_token")
	if err != nil {
		assert.Nil(t, err, "Failed to update token:"+err.Error())
	}

	foundUser, err := testStore.FindUser("user_1")
	if err != nil {
		assert.Nil(t, err, "Failed to find updated user:"+err.Error())
	}

	assert.Equal(t, "new_firebase_token", foundUser.FirebaseToken,
		"Invalid FirebaseToken:"+foundUser.FirebaseToken)
}

func TestUpdateFirebaseTokenNotFound(t *testing.T) {
	defer teardown()
	err := testStore.UpdateFirebaseToken(uuid.NewV4().String(), "new_firebase_token")
	assert.NotNil(t, err, "Should return error if user doesnt exist")
	assert.Equal(t, errors.NotFoundError, err.Code(), "Incorred code seen:"+err.Code())
}

func TestCreateUserToken(t *testing.T) {
	defer teardown()
	user, _ := testStore.CreateUser("user_1", "firebase_token_1", "user")
	userToken, err := testStore.CreateUserToken(user.ID)
	if err != nil {
		assert.Nil(t, err, "Failed to create user token:"+err.Error())
	}
	assert.Equal(t, user.ID, userToken.UserID, "Invalid user id seen in token:"+userToken.UserID)
	assert.Equal(t, user.Role, userToken.UserRole, "Invalid user id seen in token:"+userToken.UserRole)
	assert.Equal(t, user.FirebaseToken, userToken.FirebaseToken, "Invalid user id seen in token:"+userToken.FirebaseToken)
}

func TestCreateSecondUserToken(t *testing.T) {
	defer teardown()
	user, _ := testStore.CreateUser("user_1", "firebase_token_1", "user")
	oldUserToken, _ := testStore.CreateUserToken(user.ID)
	_, err := testStore.CreateUserToken(user.ID)
	if err != nil {
		assert.Nil(t, err, "Failed to create User Token")
	}
	foundUserTokens := []models.UserToken{}
	testStore.db.Where("user_id = ?", user.ID).Find(&foundUserTokens)
	assert.Equal(t, 1, len(foundUserTokens), "Incorrect number of tokens seen")
	assert.NotEqual(t, oldUserToken, foundUserTokens[0].UserToken, "Token not updated")
}

func TestCreateSecondUserTokenWithOtherUser(t *testing.T) {
	defer teardown()
	user, _ := testStore.CreateUser("user_1", "firebase_token_1", "user")
	user2, _ := testStore.CreateUser("user_2", "firebase_token_2", "user")
	testStore.CreateUserToken(user2.ID)
	oldUserToken, _ := testStore.CreateUserToken(user.ID)
	_, err := testStore.CreateUserToken(user.ID)
	if err != nil {
		assert.Nil(t, err, "Failed to create User Token")
	}
	foundUserTokens := []models.UserToken{}
	testStore.db.Where("user_id = ?", user.ID).Find(&foundUserTokens)
	assert.Equal(t, 1, len(foundUserTokens), "Incorrect number of tokens seen")
	assert.NotEqual(t, oldUserToken, foundUserTokens[0].UserToken, "Token not updated")

	foundUserTokens = []models.UserToken{}
	testStore.db.Where("user_id = ?", user2.ID).Find(&foundUserTokens)
	assert.Equal(t, 1, len(foundUserTokens), "Incorrect number of tokens seen")
}

func TestFindUserToken(t *testing.T) {
	defer teardown()
	user, _ := testStore.CreateUser("user_1", "firebase_token_1", "user")
	userToken, _ := testStore.CreateUserToken(user.ID)
	foundUserToken, err := testStore.FindUserToken(userToken.UserToken)
	if err != nil {
		assert.Nil(t, err, "Failed to find UserToken")
	}

	assert.Equal(t, *userToken, *foundUserToken, "Found wrong userToken")
}

func TestFindUserTokenNotFound(t *testing.T) {
	defer teardown()
	_, err := testStore.FindUserToken("junk")
	assert.NotNil(t, err, "Failed to find UserToken")
	assert.Equal(t, errors.NotFoundError, err.Code(), "Invalid code seen")
}

func TestFindUserFromToken(t *testing.T) {
	defer teardown()
	user, _ := testStore.CreateUser("user_1", "firebase_token_1", "user")
	userToken, _ := testStore.CreateUserToken(user.ID)
	foundUser, err := testStore.FindUserFromToken(userToken.UserToken)
	if err != nil {
		assert.Nil(t, err, "Failed to find User")
	}

	assert.Equal(t, *user, *foundUser, "Found wrong user")
}

func TestFindUserFromTokenNotFound(t *testing.T) {
	defer teardown()
	_, err := testStore.FindUserFromToken("junk")
	assert.NotNil(t, err, "Failed to find User")
	assert.Equal(t, errors.NotFoundError, err.Code(), "Invalid code seen")
}
