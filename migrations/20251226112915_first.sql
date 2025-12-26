-- +goose Up
CREATE TABLE IF NOT EXISTS profiles (
                                        user_id INTEGER PRIMARY KEY, -- Беремо цей ID з Auth-сервісу
                                        username TEXT NOT NULL UNIQUE,
                                        first_name TEXT,
                                        last_name TEXT,
                                        birth_day TIMESTAMP,
                                        phone_number TEXT
);
-- Індекс для швидкого пошуку по нікнейму
CREATE INDEX idx_profiles_username ON profiles(username);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
