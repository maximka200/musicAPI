-- +goose Up
-- +goose StatementBegin
CREATE TABLE songs (
    id SERIAL PRIMARY KEY,
    group_name VARCHAR(255) NOT NULL UNIQUE,
    song_name VARCHAR(255) NOT NULL,
    release_date DATE,
    text TEXT[],
    link VARCHAR(2048)
    );

CREATE INDEX idx_songs_group_date ON songs (group_name, release_date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_songs_group_date;
DROP TABLE songs;
-- +goose StatementEnd
