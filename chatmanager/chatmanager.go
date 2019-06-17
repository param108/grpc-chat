package chatmanager

import (
	"github.com/param108/grpc-chat-server/pubsub"
)

//ChatManager This component creates the distributor for chat messages to recipients
// Tasks are
// 1. TODO On connect read from db and send to the recipient
// 2. TODO Once all old messages are sent, then subscribe the recipient to all new messages
type ChatManager interface {
	Write(ChatID string, m pubsub.Message) error
	Read(ChatID string) (pubsub.Sub, error)
}
