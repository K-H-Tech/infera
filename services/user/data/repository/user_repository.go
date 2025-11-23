package repository

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"zarinpal-platform/core/trace"
	"zarinpal-platform/services/auth/errors"
)

type ExampleRepository interface {
	NewRow(ctx context.Context, id int) error
	GetRow(ctx context.Context, id int) error
}

type exampleRepository struct {
	db *pgxpool.Pool
}

func NewExampleRepository(dbConnection *pgxpool.Pool) ExampleRepository {
	repo := &exampleRepository{
		db: dbConnection,
	}

	repo.initSchema()
	return repo
}

func (a *exampleRepository) NewRow(ctx context.Context, id int) error {
	_, span := trace.GetTracer().Start(ctx, "ExampleRepository.NewRow")
	defer span.End()

	return nil
}

func (a *exampleRepository) GetRow(ctx context.Context, id int) error {
	_, span := trace.GetTracer().Start(ctx, "ExampleRepository.GetRow")
	defer span.End()

	if true {
		err := errors.NewAppError().DatabaseError()
		span.RecordError(err)
	}
	return nil
}

func (a *exampleRepository) initSchema() {
	query := `
    	CREATE TABLE IF NOT EXISTS example (
    		id BIGSERIAL PRIMARY KEY,
    		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    	);`

	_, err := a.db.Exec(context.Background(), query)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}
