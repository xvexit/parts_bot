ALTER TABLE users
    ADD COLUMN IF NOT EXISTS email VARCHAR(255),
    ADD COLUMN IF NOT EXISTS password_hash TEXT NOT NULL;

ALTER TABLE users
    DROP COLUMN IF EXISTS telegram_id;

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_unique
    ON users(email)
    WHERE email IS NOT NULL;