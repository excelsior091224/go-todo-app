
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS tasks(
	id INTEGER PRIMARY KEY,
	title TEXT,
	text TEXT,
    status INTEGER DEFAULT 0,
    created TIMESTAMP DEFAULT (datetime(CURRENT_TIMESTAMP, 'localtime')),
    updated TIMESTAMP DEFAULT (datetime(CURRENT_TIMESTAMP, 'localtime'))
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE tasks;
