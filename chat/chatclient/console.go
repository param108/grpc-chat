package chatclient

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var (
	UserToken string
	Username  string
	ChatID    string
)

func (m *MemoryChatClient) console() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("-> ")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)

		words := strings.Fields(text)

		m.execute(words)
	}

}

func (m *MemoryChatClient) execute(words []string) {
	switch words[0] {
	case "login":
		if len(words) != 3 {
			fmt.Println("Usage: login <username> <firebasekey>")
			return
		}
		r, err := m.Login(words[1], words[2])
		if err != nil {
			fmt.Println("Failed to login:" + err.Error())
			return
		}

		Username = words[1]
		UserToken = r
		fmt.Println(UserToken)
	case "create_chat":
		if len(words) != 2 {
			fmt.Println("Usage: create_chat <chat_name>")
			return
		}
		r, err := m.CreateChat(UserToken, words[1])
		if err != nil {
			fmt.Println("Failed to login:" + err.Error())
			return
		}
		fmt.Println(r)
	case "list_chats":
		if len(words) != 1 {
			fmt.Println("Usage: list_chats")
			return
		}
		details, err := m.ListChats(UserToken)
		if err != nil {
			fmt.Println("Failed to list chats:" + err.Error())
			return
		}
		fmt.Println("ChatName\tChatID\tUserIDs")
		for _, detail := range details {
			fmt.Printf("%s\t%s\t%s\n", detail.ChatName, detail.ChatID,
				strings.Join(detail.UserList, ","))
		}
	case "join_chat":
		if len(words) != 2 {
			fmt.Println("Usage: join_chat <chat_id>")
			return
		}

		err := m.JoinChat(UserToken, words[1])

		if err != nil {
			fmt.Println("Failed to join chat:" + err.Error())
			return
		}
	case "start_chat":
		if len(words) != 2 {
			fmt.Println("Usage: start_chat <chat_id>")
		}
		m.StartChat(Username, UserToken, words[1])
	}
}

func (m *MemoryChatClient) StartChat(username string, userToken string, chatID string) {
	reader := bufio.NewReader(os.Stdin)
	quit := make(chan int)
	done := make(chan int)
	go m.read(userToken, chatID, quit, done)

	for {
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)

		m.write(username, userToken, chatID, text)
		if text == "quit" {
			quit <- 1
			<-done
			return
		}
	}
}
