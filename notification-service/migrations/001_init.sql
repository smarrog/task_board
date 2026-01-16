-- +goose Up

CREATE TABLE IF NOT EXISTS notifications (
    outbox_id TEXT PRIMARY KEY,
    event_type TEXT NOT NULL,
    aggregate_type TEXT NOT NULL,
    aggregate_id TEXT NOT NULL,
    event_created_at TIMESTAMPTZ NOT NULL,
    version INT NOT NULL,
    payload JSONB NOT NULL,
    text TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_notifications_event_type ON notifications(event_type);
CREATE INDEX IF NOT EXISTS idx_notifications_aggregate ON notifications(aggregate_type, aggregate_id);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at);

-- +goose Down
DROP TABLE IF EXISTS notifications;
