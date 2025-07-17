package repository

import (
	"context"
	"database/sql"

	"github.com/Yagshymyradov/subscriptions-service/internal/models"
)

const createQuery = `
INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
VALUES ($1,$2,$3,$4,$5) RETURNING id;
`
const getByIDQuery = `
SELECT id, service_name, price, user_id, start_date, end_date
FROM subscriptions
WHERE id = $1;
`
const listByUserQuery = `
SELECT id, service_name, price, user_id, start_date, end_date
FROM subscriptions
WHERE user_id = $1
ORDER BY id;
`

const updateQuery = `
UPDATE subscriptions
SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5
WHERE id = $6;
`

const deleteQuery = `
DELETE FROM subscriptions
WHERE id = $1;
`

const totalCostquery = `
SELECT COALESCE(SUM(price), 0) AS total
FROM subscriptions
WHERE user_id = $1
	AND ($2 = '' OR service_name ILIKE '%' || $2 || '%')
	AND start_date <= (make_date($4, $3, 1) + INTERVAL '1 month' - INTERVAL '1 day')
	AND (end_date IS NULL OR end_date >= make_date($4, $3, 1));
`

type PostgresSubscriptionRepository struct {
	db *sql.DB
}

func NewPostgresSubscriptionRepository(db *sql.DB) *PostgresSubscriptionRepository {
	return &PostgresSubscriptionRepository{db: db}
}

func (r *PostgresSubscriptionRepository) Create(ctx context.Context, s *models.Subscription) error {
	return r.db.QueryRowContext(ctx, createQuery,
		s.ServiceName, s.Price, s.UserID, s.StartDate, s.EndDate,
	).Scan(&s.ID)
}

func (r *PostgresSubscriptionRepository) GetByID(ctx context.Context, id int) (*models.Subscription, error) {
	var s models.Subscription
	err := r.db.QueryRowContext(ctx, getByIDQuery, id).
		Scan(&s.ID, &s.ServiceName, &s.Price, &s.UserID, &s.StartDate, &s.EndDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

func (r *PostgresSubscriptionRepository) List(ctx context.Context, userID string) ([]models.Subscription, error) {
	rows, err := r.db.QueryContext(ctx, listByUserQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []models.Subscription
	for rows.Next() {
		var s models.Subscription
		if err := rows.Scan(&s.ID, &s.ServiceName, &s.Price, &s.UserID, &s.StartDate, &s.EndDate); err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return subs, nil
}

func (r *PostgresSubscriptionRepository) Update(ctx context.Context, s *models.Subscription) error {
	res, err := r.db.ExecContext(ctx, updateQuery,
		s.ServiceName, s.Price, s.UserID, s.StartDate, s.EndDate, s.ID,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *PostgresSubscriptionRepository) Delete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *PostgresSubscriptionRepository) TotalCost(ctx context.Context, userID string, month, year int, serviceFilter string) (int, error) {
	var total int
	err := r.db.QueryRowContext(ctx, totalCostquery, userID, serviceFilter, month, year).Scan(&total)
	return total, err
}
