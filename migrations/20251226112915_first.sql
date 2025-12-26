-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS profiles (
                                        user_id INTEGER PRIMARY KEY,
                                        username TEXT NOT NULL UNIQUE,
                                        first_name TEXT,
                                        last_name TEXT,
                                        birth_day TIMESTAMP,
                                        phone_number TEXT
);

CREATE INDEX idx_profiles_username ON profiles(username);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_profiles_username;
DROP TABLE IF EXISTS profiles;
-- +goose StatementEnd