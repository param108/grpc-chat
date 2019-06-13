GO_FILES=$(shell find . -name '*.go')
TEST_DB_DRIVER="postgres"
TEST_DB_USERNAME="test_grpc_chat"
TEST_DB_PASSWORD="test_grpc_chat"
TEST_DB_HOST="localhost"
TEST_DB_PORT="5432"
TEST_DB_NAME="test_grpc_chat"
MIGRATION_PATH=file://db/migrations

chat/chat.pb.go:	chat/chat.proto
	protoc -I chat/ chat/chat.proto --go_out=plugins=grpc:chat

build:	chat/chat.pb.go ${GO_FILES}
	mkdir -p out
	go build  -o out/grpc_chat


test: build
	DB_DRIVER=${TEST_DB_DRIVER} DB_USERNAME=${TEST_DB_USERNAME} \
	DB_PASSWORD=${TEST_DB_PASSWORD} DB_HOST=${TEST_DB_HOST} \
	DB_PORT=${TEST_DB_PORT} DB_NAME=${TEST_DB_NAME} MIGRATION_PATH=${MIGRATION_PATH} out/grpc_chat migrate
	TEST_DB_DRIVER=${TEST_DB_DRIVER} TEST_DB_USERNAME=${TEST_DB_USERNAME} \
	TEST_DB_PASSWORD=${TEST_DB_PASSWORD} TEST_DB_HOST=${TEST_DB_HOST} \
	TEST_DB_PORT=${TEST_DB_PORT} TEST_DB_NAME=${TEST_DB_NAME} go test -v ./... --cover

test.rollback: build
	DB_DRIVER=${TEST_DB_DRIVER} DB_USERNAME=${TEST_DB_USERNAME} \
	DB_PASSWORD=${TEST_DB_PASSWORD} DB_HOST=${TEST_DB_HOST} \
	DB_PORT=${TEST_DB_PORT} DB_NAME=${TEST_DB_NAME} MIGRATION_PATH=${MIGRATION_PATH} out/grpc_chat rollback
