package chatmanager

import (
	"github.com/param108/grpc-chat/pubsub"
	"sync"
)

/*SimpleChatManager A Simple Chat Manager
 *Implements ChatManager interface
 */
type SimpleChatManager struct {
	d map[string]*pubsub.PubSub
	m sync.Mutex
}

func NewSimpleChatManager() *SimpleChatManager {
	return &SimpleChatManager{d: make(map[string]*pubsub.PubSub)}
}

func (cm *SimpleChatManager) createPubSub(chatID string) *pubsub.PubSub {
	cm.m.Lock()
	defer cm.m.Unlock()
	ps := pubsub.NewPubSub()
	cm.d[chatID] = ps
	go ps.PipeRoutine()
	return ps
}

//Write Assumes ChatID is valid
func (cm *SimpleChatManager) Write(ChatID string, m pubsub.Message) error {
	if ps, ok := cm.d[ChatID]; ok {
		ps.Send(m)
	} else {
		ps = cm.createPubSub(ChatID)
		ps.Send(m)
	}
	return nil
}

func (cm *SimpleChatManager) Read(ChatID string) (pubsub.Sub, error) {
	if ps, ok := cm.d[ChatID]; ok {
		ret := ps.Subscribe()
		return ret, nil
	}

	ps := cm.createPubSub(ChatID)
	ret := ps.Subscribe()
	return ret, nil
}
