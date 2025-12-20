-- +goose Up
CREATE TABLE IF NOT EXISTS outbox_events (
   id UUID PRIMARY KEY,
   event_type TEXT NOT NULL,
   aggregate_type TEXT NOT NULL,   -- например: "task"
   aggregate_id UUID NOT NULL,     -- id задачи
   payload JSONB NOT NULL,
   created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
   published_at TIMESTAMPTZ
);

CREATE INDEX idx_outbox_unpublished ON outbox_events(published_at) WHERE published_at IS NULL;
CREATE INDEX idx_outbox_created ON outbox_events(created_at);

-- +goose Down
DROP TABLE IF EXISTS outbox_events;
