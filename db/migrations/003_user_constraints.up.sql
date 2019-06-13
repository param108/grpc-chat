ALTER TABLE users ADD CONSTRAINT user_name_uniq UNIQUE (name);
ALTER TABLE users ADD CONSTRAINT user_firebase_token_uniq UNIQUE (firebase_token);
