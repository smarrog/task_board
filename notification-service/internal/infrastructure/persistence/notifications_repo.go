package persistence

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/smarrog/task-board/notification-service/internal/domain/notification"
)

type NotificationsRepo struct {
	pg  *pgxpool.Pool
	log *zerolog.Logger
}

func NewNotificationsRepo(pg *pgxpool.Pool, log *zerolog.Logger) *NotificationsRepo {
	return &NotificationsRepo{pg: pg, log: log}
}

func (r *NotificationsRepo) Save(ctx context.Context, n notification.HistoryRecord) error {
	_, err := r.pg.Exec(ctx, `
		INSERT INTO notifications (
			outbox_id,
			event_type,
			aggregate_type,
			aggregate_id,
			event_created_at,
			version,
			payload,
			text
		) VALUES ($1,$2,$3,$4,$5,$6,$7::jsonb,$8)
		ON CONFLICT (outbox_id) DO NOTHING
	`, n.OutboxId, n.EventType, n.AggregateType, n.AggregateId, n.EventCreatedAt, n.Version, n.Payload, n.Text)
	return err
}
