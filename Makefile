GO_FILES=$(shell find . -name '*.go')
TEST_DB_DRIVER="postgres"
TEST_DB_USERNAME="test_grpc_chat"
TEST_DB_PASSWORD="test_grpc_chat"
TEST_DB_HOST="localhost"
TEST_DB_PORT="5432"
TEST_DB_NAME="test_grpc_chat"
REDIS_HOST="localhost"
REDIS_PORT="6379"
MIGRATION_PATH=file://db/migrations

LOCAL_DB_DRIVER="postgres"
LOCAL_DB_USERNAME="grpc_chat"
LOCAL_DB_PASSWORD="grpc_chat"
LOCAL_DB_HOST="localhost"
LOCAL_DB_PORT="5432"
LOCAL_DB_NAME="grpc_chat"
LOCAL_MIGRATION_PATH=file://db/migrations

chat/chat.pb.go:	chat/chat.proto
	protoc -I chat/ chat/chat.proto --go_out=plugins=grpc:chat

build:	chat/chat.pb.go ${GO_FILES}
	mkdir -p out
	go build  -o out/grpc_chat

bundle: build
	cp ${PRODUCTION_CONFIG} out/.env
	cp -R db out/
	tar -zcvf bundle.tgz out

build-java: chat/chat.proto
	rm -rf chat/java
	mkdir -p chat/java
	protoc -I chat/ chat/chat.proto --java_out=chat/java/

test: build
	REDIS_HOST=${REDIS_HOST} REDIS_PORT=${REDIS_PORT} \
	DB_DRIVER=${TEST_DB_DRIVER} DB_USERNAME=${TEST_DB_USERNAME} \
	DB_PASSWORD=${TEST_DB_PASSWORD} DB_HOST=${TEST_DB_HOST} \
	DB_PORT=${TEST_DB_PORT} DB_NAME=${TEST_DB_NAME} MIGRATION_PATH=${MIGRATION_PATH} out/grpc_chat migrate
	REDIS_HOST=${REDIS_HOST} REDIS_PORT=${REDIS_PORT} \
	TEST_DB_DRIVER=${TEST_DB_DRIVER} TEST_DB_USERNAME=${TEST_DB_USERNAME} \
	TEST_DB_PASSWORD=${TEST_DB_PASSWORD} TEST_DB_HOST=${TEST_DB_HOST} \
	TEST_DB_PORT=${TEST_DB_PORT} TEST_DB_NAME=${TEST_DB_NAME} go test -v ./... --cover

test.rollback: build
	REDIS_HOST=${REDIS_HOST} REDIS_PORT=${REDIS_PORT} \
	DB_DRIVER=${TEST_DB_DRIVER} DB_USERNAME=${TEST_DB_USERNAME} \
	DB_PASSWORD=${TEST_DB_PASSWORD} DB_HOST=${TEST_DB_HOST} \
	DB_PORT=${TEST_DB_PORT} DB_NAME=${TEST_DB_NAME} MIGRATION_PATH=${MIGRATION_PATH} out/grpc_chat rollback

migrate: build
	REDIS_HOST=${REDIS_HOST} REDIS_PORT=${REDIS_PORT} \
	DB_DRIVER=${LOCAL_DB_DRIVER} DB_USERNAME=${LOCAL_DB_USERNAME} \
	DB_PASSWORD=${LOCAL_DB_PASSWORD} DB_HOST=${LOCAL_DB_HOST} \
	DB_PORT=${LOCAL_DB_PORT} DB_NAME=${LOCAL_DB_NAME} MIGRATION_PATH=${LOCAL_MIGRATION_PATH} out/grpc_chat migrate

rollback: build
	REDIS_HOST=${REDIS_HOST} REDIS_PORT=${REDIS_PORT} \
	DB_DRIVER=${LOCAL_DB_DRIVER} DB_USERNAME=${LOCAL_DB_USERNAME} \
	DB_PASSWORD=${LOCAL_DB_PASSWORD} DB_HOST=${LOCAL_DB_HOST} \
	DB_PORT=${LOCAL_DB_PORT} DB_NAME=${LOCAL_DB_NAME} MIGRATION_PATH=${LOCAL_MIGRATION_PATH} out/grpc_chat rollback
