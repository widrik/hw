-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events (
  id uuid primary key,
  title text,
  description text,
  start_at text,
  finished_at text,
  created_at text,
  updated_at text,
  user_id VARCHAR(100) NOT NULL,
  notify_at int DEFAULT 900
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd
