package repository

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/khaldeezal/subscriptions-service/internal/domain"
	"github.com/khaldeezal/subscriptions-service/internal/utils"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewSubscriptionRepository(ctx context.Context, dns string) (*Repository, error) {
	pool, err := newPool(ctx, dns)
	if err != nil {
		return nil, err
	}

	return &Repository{pool: pool}, nil
}

func newPool(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	return pgxpool.NewWithConfig(ctx, cfg)
}

func (r *Repository) ClosePool() {
	r.pool.Close()
}

func (r *Repository) Create(ctx context.Context, in *domain.CreateInput) (string, error) {
	id := uuid.New().String()
	start, err := utils.ParseYearMonth(in.StartDate)
	if err != nil {
		return "", err
	}

	var end *time.Time
	if in.EndDate != nil && *in.EndDate != "" {
		et, err := utils.ParseYearMonth(*in.EndDate)
		if err != nil {
			return "", err
		}
		end = &et
	}
	if in.Price < 0 {
		return "", errors.New("price must be >= 0")
	}

	q := `INSERT INTO subscriptions(id,user_id,service_name,price,start_date,end_date)
		VALUES ($1,$2,$3,$4,$5,$6)`

	_, err = r.pool.Exec(ctx, q, id, in.UserID, in.ServiceName, in.Price, start, end)
	return id, err
}

func (r *Repository) Get(ctx context.Context, id string) (*domain.Subscription, error) {
	q := `SELECT id,user_id,service_name,price,start_date,end_date,created_at,updated_at
		FROM subscriptions WHERE id=$1`
	row := r.pool.QueryRow(ctx, q, id)

	var s domain.Subscription
	var start, end *time.Time

	err := row.Scan(&s.ID, &s.UserID, &s.ServiceName, &s.Price, &start, &end, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}

	s.StartDate = utils.YmString(*start)
	if end != nil {
		e := utils.YmString(*end)
		s.EndDate = &e
	}
	return &s, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	cmd, err := r.pool.Exec(ctx, `DELETE FROM subscriptions WHERE id=$1`, id)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (r *Repository) List(ctx context.Context, f *domain.ListFilter) ([]domain.Subscription, error) {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.PageSize <= 0 || f.PageSize > 100 {
		f.PageSize = 20
	}
	offset := (f.Page - 1) * f.PageSize

	where := "WHERE 1=1"
	args := []any{}
	idx := 1

	if f.UserID != nil && *f.UserID != "" {
		where += " AND user_id=$" + itoa(idx)
		args = append(args, *f.UserID)
		idx++
	}

	if f.ServiceName != nil && *f.ServiceName != "" {
		where += " AND service_name=$" + itoa(idx)
		args = append(args, *f.ServiceName)
		idx++
	}

	q := `SELECT id,user_id,service_name,price,start_date,end_date,created_at,updated_at
		FROM subscriptions ` + where + ` ORDER BY created_at DESC LIMIT ` + itoa(f.PageSize) + ` OFFSET ` + itoa(offset)

	rows, err := r.pool.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Subscription

	for rows.Next() {
		var s domain.Subscription
		var start, end *time.Time

		if err := rows.Scan(&s.ID, &s.UserID, &s.ServiceName, &s.Price, &start, &end, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}

		s.StartDate = utils.YmString(*start)
		if end != nil {
			e := utils.YmString(*end)
			s.EndDate = &e
		}
		out = append(out, s)
	}

	return out, rows.Err()
}

// Update updates all fields.
func (r *Repository) Update(ctx context.Context, id string, in *domain.CreateInput) error {
	start, err := utils.ParseYearMonth(in.StartDate)
	if err != nil {
		return err
	}

	var end *time.Time
	if in.EndDate != nil && *in.EndDate != "" {
		et, err := utils.ParseYearMonth(*in.EndDate)
		if err != nil {
			return err
		}
		end = &et
	}

	cmd, err := r.pool.Exec(ctx, `UPDATE subscriptions
		SET user_id=$1, service_name=$2, price=$3, start_date=$4, end_date=$5, updated_at=now()
		WHERE id=$6`, in.UserID, in.ServiceName, in.Price, start, end, id)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func itoa(n int) string { return fmtInt(n) }

func fmtInt(n int) string { return strconv.Itoa(n) }
