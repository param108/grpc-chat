ALTER TABLE user_groups ADD CONSTRAINT user_group_uniq UNIQUE(user_id, chat_group_id);
