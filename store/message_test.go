package store

import (
	"github.com/param108/grpc-chat-server/errors"
	"github.com/param108/grpc-chat-server/models"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateChatGroup(t *testing.T) {
	defer teardown()
	chatGroup, err := testStore.CreateChatGroup("newChat")
	assert.Nil(t, err, "Failed to create chat group")
	assert.Equal(t, "newChat", chatGroup.ChatName, "Invalid name seen")
}

func TestCreateDuplicateChatGroup(t *testing.T) {
	defer teardown()
	testStore.CreateChatGroup("newChat")
	_, err := testStore.CreateChatGroup("newChat")
	assert.Nil(t, err, "Failed to throw error for create chat group")

	chatGroups := []models.ChatGroup{}
	testStore.db.Find(&chatGroups)
	assert.Equal(t, 2, len(chatGroups), "Invalid name seen")
}

func TestFindChatGroup(t *testing.T) {
	defer teardown()
	chatGroup, _ := testStore.CreateChatGroup("newChat")
	foundChatGroup, err := testStore.FindChatGroup(chatGroup.ID)
	assert.Nil(t, err, "Failed to find ChatGroup")
	assert.Equal(t, "newChat", foundChatGroup.ChatName, "Invalid chat name found")
}

func TestFindChatGroupNotFound(t *testing.T) {
	defer teardown()
	_, err := testStore.FindChatGroup(uuid.NewV4().String())
	assert.NotNil(t, err, "Failed to throw error for notfind ChatGroup")
	assert.Equal(t, errors.NotFoundError, err.Code(), "Invalid error seen")
}

func TestAddUserToChat(t *testing.T) {
	defer teardown()
	user, _ := testStore.CreateUser("user_1", "firebase_token_1", "user")
	chat, _ := testStore.CreateChatGroup("chat_1")
	err := testStore.AddUserToChat(user.ID, chat.ID)
	assert.Nil(t, err, "Failed to add user to chat")
}

func TestAddUserToChatDuplicate(t *testing.T) {
	defer teardown()
	user, _ := testStore.CreateUser("user_1", "firebase_token_1", "user")
	chat, _ := testStore.CreateChatGroup("chat_1")
	testStore.AddUserToChat(user.ID, chat.ID)
	err := testStore.AddUserToChat(user.ID, chat.ID)
	assert.NotNil(t, err, "Should throw error if associate twice")
	assert.Equal(t, errors.AlreadyExistsError, err.Code(), "Invalid error seen")
}

func TestListUsersOfChat(t *testing.T) {
	defer teardown()
	user1, _ := testStore.CreateUser("user_1", "firebase_token_1", "user")
	user2, _ := testStore.CreateUser("user_2", "firebase_token_2", "user")
	chat, _ := testStore.CreateChatGroup("the_chat_group")
	testStore.AddUserToChat(user1.ID, chat.ID)
	testStore.AddUserToChat(user2.ID, chat.ID)
	userList, err := testStore.ListUsersOfChat(chat.ID)
	assert.Nil(t, err, "Failed to find userList")
	assert.Equal(t, 2, len(userList), "Invalid number of users found")
	assert.True(t, (userList[0] == user1.ID && userList[1] == user2.ID) ||
		(userList[0] == user2.ID && userList[1] == user1.ID),
		"Invalid userList found")
}

func TestAddMessageToChat(t *testing.T) {
	defer teardown()
	user1, _ := testStore.CreateUser("user_1", "firebase_token_1", "user")
	chat, _ := testStore.CreateChatGroup("the_chat_group")
	testStore.AddUserToChat(user1.ID, chat.ID)

	err := testStore.AddMessageToChat(user1.ID, chat.ID, time.Now().Format(time.RFC1123Z),
		"Hello my name is bogie")
	messages := []models.ChatMessage{}
	testStore.db.Find(&messages)
	assert.Equal(t, 1, len(messages), "Invalid number of messages seen")
	assert.Nil(t, err, "Failed to add message")

	etime := messages[0].ServerTime.Sub(messages[0].ClientTime)
	if messages[0].ClientTime.After(messages[0].ServerTime) {
		etime = messages[0].ClientTime.Sub(messages[0].ServerTime)
	}
	assert.Condition(t, func() bool { return etime < (5 * time.Minute) }, "Times are off")
}

func TestGetMessagesAfter(t *testing.T) {
	defer teardown()
	user1, _ := testStore.CreateUser("user_1", "firebase_token_1", "user")
	chat, _ := testStore.CreateChatGroup("the_chat_group")
	testStore.AddUserToChat(user1.ID, chat.ID)

	firstTimeStamp := time.Now().Format(time.RFC1123Z)
	time.Sleep(time.Second)
	testStore.AddMessageToChat(user1.ID, chat.ID, time.Now().Format(time.RFC1123Z),
		"line 1")
	time.Sleep(10 * time.Second)

	secondTimeStamp := time.Now().Format(time.RFC1123Z)

	time.Sleep(time.Second)

	testStore.AddMessageToChat(user1.ID, chat.ID, time.Now().Format(time.RFC1123Z),
		"line 2")

	messages, err := testStore.GetMessagesAfter(chat.ID, firstTimeStamp)
	assert.Nil(t, err, "Failed to get messages")

	assert.Equal(t, 2, len(messages), "Incorrect number of messages read")

	assert.Equal(t, "line 1", messages[0].Message, "Invalid first message")
	assert.Equal(t, "line 2", messages[1].Message, "Invalid first message")

	messages, err = testStore.GetMessagesAfter(chat.ID, secondTimeStamp)
	assert.Nil(t, err, "Failed to get messages")

	assert.Equal(t, 1, len(messages), "Incorrect number of messages read")
	assert.Equal(t, "line 2", messages[0].Message, "Invalid first message")

}
