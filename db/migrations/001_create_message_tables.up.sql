CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE chat_messages (
       id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
       server_time TIMESTAMPTZ DEFAULT NOW(),
       client_time TIMESTAMPTZ NOT NULL,
       message text NOT NULL,
       user_id uuid NOT NULL,
       chat_group_id uuid NOT NULL
);

CREATE TABLE chat_groups (
       id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
       chat_name varchar(50) NOT NULL
);

CREATE TABLE users (
       id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
       name varchar(50) NOT NULL
);

CREATE TABLE user_groups (
       id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
       chat_group_id uuid NOT NULL,
       user_id uuid NOT NULL
);

CREATE TABLE chat_message_archives (
       id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
       server_time TIMESTAMPTZ DEFAULT NOW(),
       client_time TIMESTAMPTZ NOT NULL,
       message text NOT NULL,
       user_id uuid NOT NULL,
       chat_group_id uuid NOT NULL
)
