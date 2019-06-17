# GRPC-Chat Server and Client

## Setup

1. cp env.txt .env
2. update the db details in the .env file
- if you want tests to work setup the TEST... variables as well
3. 

`make build`

4. Executable is out/grpc-chat
5. Migrate the db

`./out/grpc-chat migrate`

6. Run the server

./out/grpc-chat server

7. Run the client in console mode

./out/grpc-chat client console

8. You should see

`-> `

prompt


## Prompt commands for client

1. `login <username> <some key>`
- this logs you in. The server doesnt do authentication rt now, so the key is useless

2. `create_chat <chat name>`
- creates a new chat

3. `list_chats`
- lists all the chats available on the server

```
-> list_chats
ChatName	ChatID	UserIDs
blublu	5fa89b3d-e154-44c4-ad41-785d8a330535	a39b2fdf-e5d0-41a1-8893-c4cfb7ae6bb0,fee1adb3-d969-4258-8907-224fa2bc284b
bumbum	8cd61367-8556-4bfb-a3aa-940b7f209700	a39b2fdf-e5d0-41a1-8893-c4cfb7ae6bb0
baby	8d8c3079-c679-4f29-81a7-114eb4e8fe44	a39b2fdf-e5d0-41a1-8893-c4cfb7ae6bb0
```

4. `start_chat <chat id>`
- use the second column from `list_chats` output.

5. now just type and you should see the chat on all clients who have connected.
- other users will have their username prefixed

```
-> start_chat 5fa89b3d-e154-44c4-ad41-785d8a330535
Hi I am param
anil:hello
anil:I am anil
```

6. `quit` to exit a chat
