-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users(
			id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			username VARCHAR(255) NOT NULL UNIQUE,
            password_hash VARCHAR NOT NULL,
            user_session JSON,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);
CREATE TABLE IF NOT EXISTS notes(
			id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		);

CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS update_updated_at_notes ON notes;

CREATE TRIGGER update_date_notes
    BEFORE UPDATE ON notes
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
DROP TABLE notes;
-- +goose StatementEnd
