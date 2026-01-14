package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	do "github.com/smarrog/task-board/auth-service/internal/domain"
)

type UsersRepo struct {
	pg  *pgxpool.Pool
	log *zerolog.Logger
}

func NewUsersRepo(pg *pgxpool.Pool, log *zerolog.Logger) *UsersRepo {
	return &UsersRepo{pg: pg, log: log}
}

func (r *UsersRepo) Create(ctx context.Context, u *do.User) error {
	_, err := r.pg.Exec(ctx, `
		INSERT INTO users (id, email, username, password_hash, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, u.Id(), u.Email(), u.Username(), u.PwdHash(), u.CreatedAt())
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "ux_users_email" { // 23505 - UniqueViolation
			return do.ErrEmailAlreadyExists
		}
		return err
	}
	return nil
}

func (r *UsersRepo) GetByEmail(ctx context.Context, email do.Email) (*do.User, error) {
	var userIdRaw uuid.UUID
	var userNameRaw, pwdHashRaw string
	var createdAtRaw time.Time

	err := r.pg.QueryRow(ctx, `
		SELECT id, username, password_hash, created_at
		FROM users
		WHERE email = $1
	`, email.String()).Scan(&userIdRaw, &userNameRaw, &pwdHashRaw, &createdAtRaw)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, do.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	userId, err := do.UserIdFromUUID(userIdRaw)
	if err != nil {
		return nil, err
	}
	userName, err := do.NewUserName(userNameRaw)
	if err != nil {
		return nil, err
	}
	pwdHash, err := do.NewPwdHash(pwdHashRaw)
	if err != nil {
		return nil, err
	}

	return do.RehydrateUser(userId, email, userName, pwdHash, createdAtRaw), nil
}

func (r *UsersRepo) GetById(ctx context.Context, id do.UserId) (*do.User, error) {
	var emailRaw, userNameRaw, pwdHashRaw string
	var createdAtRaw time.Time

	err := r.pg.QueryRow(ctx, `
		SELECT email, username, password_hash, created_at
		FROM users
		WHERE id = $1
	`, id.UUID()).Scan(&emailRaw, &userNameRaw, &pwdHashRaw, &createdAtRaw)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, do.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	email, err := do.NewEmail(emailRaw)
	if err != nil {
		return nil, err
	}
	userName, err := do.NewUserName(userNameRaw)
	if err != nil {
		return nil, err
	}
	pwdHash, err := do.NewPwdHash(pwdHashRaw)
	if err != nil {
		return nil, err
	}

	return do.RehydrateUser(id, email, userName, pwdHash, createdAtRaw), nil
}
