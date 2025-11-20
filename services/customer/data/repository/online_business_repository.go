package repository

import (
	"context"
	"log"

	"zarinpal-platform/core/trace"
	"zarinpal-platform/services/customer/data/model"
	"zarinpal-platform/services/customer/errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OnlineBusinessRepository interface {
	Create(ctx context.Context, business *model.OnlineBusiness) (*model.OnlineBusiness, error)
	GetByID(ctx context.Context, id string) (*model.OnlineBusiness, error)
	GetByUserID(ctx context.Context, userID string) ([]*model.OnlineBusiness, error)
}

type onlineBusinessRepository struct {
	db *pgxpool.Pool
}

func NewOnlineBusinessRepository(dbConnection *pgxpool.Pool) OnlineBusinessRepository {
	repo := &onlineBusinessRepository{
		db: dbConnection,
	}

	repo.initSchema()
	return repo
}

func (r *onlineBusinessRepository) Create(ctx context.Context, business *model.OnlineBusiness) (*model.OnlineBusiness, error) {
	_, span := trace.GetTracer().Start(ctx, "OnlineBusinessRepository.Create")
	defer span.End()

	// Generate UUIDs
	business.ID = uuid.New().String()
	business.CustomerID = uuid.New().String()

	query := `
		INSERT INTO online_businesses (
			id, customer_id, website_name, url, enamad_id, user_id, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, NOW(), NOW()
		) RETURNING created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		business.ID,
		business.CustomerID,
		business.WebsiteName,
		business.URL,
		business.EnamadID,
		business.UserID,
	).Scan(&business.CreatedAt, &business.UpdatedAt)

	if err != nil {
		span.RecordError(err)
		return nil, errors.NewAppError(ctx).DatabaseError()
	}

	return business, nil
}

func (r *onlineBusinessRepository) GetByID(ctx context.Context, id string) (*model.OnlineBusiness, error) {
	_, span := trace.GetTracer().Start(ctx, "OnlineBusinessRepository.GetByID")
	defer span.End()

	business := &model.OnlineBusiness{}
	query := `
		SELECT id, customer_id, website_name, url, enamad_id, user_id, created_at, updated_at
		FROM online_businesses
		WHERE id = $1
	`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&business.ID,
		&business.CustomerID,
		&business.WebsiteName,
		&business.URL,
		&business.EnamadID,
		&business.UserID,
		&business.CreatedAt,
		&business.UpdatedAt,
	)

	if err != nil {
		span.RecordError(err)
		return nil, errors.NewAppError(ctx).DatabaseError()
	}

	return business, nil
}

func (r *onlineBusinessRepository) GetByUserID(ctx context.Context, userID string) ([]*model.OnlineBusiness, error) {
	_, span := trace.GetTracer().Start(ctx, "OnlineBusinessRepository.GetByUserID")
	defer span.End()

	query := `
		SELECT id, customer_id, website_name, url, enamad_id, user_id, created_at, updated_at
		FROM online_businesses
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		span.RecordError(err)
		return nil, errors.NewAppError(ctx).DatabaseError()
	}
	defer rows.Close()

	businesses := make([]*model.OnlineBusiness, 0)
	for rows.Next() {
		business := &model.OnlineBusiness{}
		err := rows.Scan(
			&business.ID,
			&business.CustomerID,
			&business.WebsiteName,
			&business.URL,
			&business.EnamadID,
			&business.UserID,
			&business.CreatedAt,
			&business.UpdatedAt,
		)
		if err != nil {
			span.RecordError(err)
			return nil, errors.NewAppError(ctx).DatabaseError()
		}
		businesses = append(businesses, business)
	}

	return businesses, nil
}

func (r *onlineBusinessRepository) initSchema() {
	query := `
		CREATE TABLE IF NOT EXISTS online_businesses (
			id UUID PRIMARY KEY,
			customer_id UUID NOT NULL,
			website_name VARCHAR(255) NOT NULL,
			url VARCHAR(500) NOT NULL,
			enamad_id VARCHAR(100),
			user_id VARCHAR(100) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			CONSTRAINT unique_website_url UNIQUE (url)
		);

		CREATE INDEX IF NOT EXISTS idx_online_businesses_user_id ON online_businesses(user_id);
		CREATE INDEX IF NOT EXISTS idx_online_businesses_customer_id ON online_businesses(customer_id);
	`

	_, err := r.db.Exec(context.Background(), query)
	if err != nil {
		log.Fatalf("Failed to create online_businesses table: %v", err)
	}
}
