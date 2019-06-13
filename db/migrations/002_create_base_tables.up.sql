CREATE TABLE user_tokens (
       id SERIAL NOT NULL PRIMARY KEY,
       user_token TEXT NOT NULL,
       user_id uuid NOT NULL,
       user_role TEXT NOT NULL,
       firebase_token TEXT NOT NULL
);

ALTER TABLE users ADD COLUMN role TEXT NOT NULL;

ALTER TABLE users ADD COLUMN firebase_token TEXT NOT NULL;

CREATE INDEX user_id_index on user_tokens (user_token);

CREATE TABLE availables (
       id SERIAL NOT NULL PRIMARY KEY,
       user_id uuid NOT NULL,
       start_time TIMESTAMPTZ NOT NULL,
       end_time TIMESTAMPTZ NOT NULL
);

ALTER TABLE user_tokens ADD CONSTRAINT user_tokens_user_id_fk FOREIGN KEY (user_id) REFERENCES users (id);

ALTER TABLE availables ADD CONSTRAINT availables_user_id_fk FOREIGN KEY (user_id) REFERENCES users (id);

ALTER TABLE user_groups ADD CONSTRAINT user_groups_user_id_fk FOREIGN KEY (user_id) REFERENCES users (id);

ALTER TABLE chat_messages ADD CONSTRAINT chat_messages_user_id_fk FOREIGN KEY (user_id) REFERENCES users (id);

ALTER TABLE chat_message_archives ADD CONSTRAINT chat_message_archives_user_id_fk FOREIGN KEY (user_id) REFERENCES users (id);

ALTER TABLE chat_messages ADD CONSTRAINT chat_messages_chat_group_id_fk FOREIGN KEY (chat_group_id) REFERENCES chat_groups (id);

ALTER TABLE chat_message_archives ADD CONSTRAINT chat_message_archives_chat_group_id_fk FOREIGN KEY (chat_group_id) REFERENCES chat_groups (id);
