ALTER TABLE users 
    DROP COLUMN IF EXISTS email,
    DROP COLUMN IF EXISTS password_hash;

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS telegram_id BIGINT;

ALTER TABLE users
    ALTER COLUMN telegram_id SET NOT NULL;


DROP INDEX IF EXISTS idx_users_email_unique;

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_telegram_unique
    ON users(telegram_id);
