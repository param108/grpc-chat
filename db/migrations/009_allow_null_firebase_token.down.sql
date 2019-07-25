ALTER TABLE users DROP COLUMN firebase_token;
ALTER TABLE users ADD COLUMN firebase_token text NOT NULL;
