package repository

import (
	"context"
	"log"

	"zarinpal-platform/core/trace"
	"zarinpal-platform/services/auth/data/model"
	"zarinpal-platform/services/auth/errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OtpRepository interface {
	GenerateOTP(ctx context.Context, mobile string) error
	VerifyOTP(ctx context.Context, mobile string, otp string) error
	NewUser(ctx context.Context, name string) (*model.User, error)
	GetUser(ctx context.Context, id int32) (*model.User, error)
}

type otpRepository struct {
	db *pgxpool.Pool
}

func (r *otpRepository) NewUser(ctx context.Context, name string) (*model.User, error) {
	_, span := trace.GetTracer().Start(ctx, "OtpRepository.NewUser")
	defer span.End()

	var user model.User
	row := r.db.QueryRow(ctx, "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id, name, email, created_at, updated_at", name, name+"@example.com")
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *otpRepository) GetUser(ctx context.Context, id int32) (*model.User, error) {
	_, span := trace.GetTracer().Start(ctx, "OtpRepository.GetUser")
	defer span.End()

	var user model.User
	row := r.db.QueryRow(ctx, "SELECT id, name, email, created_at, updated_at FROM users WHERE id = $1", id)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func NewOtpRepository(dbConnection *pgxpool.Pool) OtpRepository {

	repo := &otpRepository{
		db: dbConnection,
	}

	repo.initSchema()
	return repo
}

func (r *otpRepository) GenerateOTP(ctx context.Context, mobile string) error {
	_, span := trace.GetTracer().Start(ctx, "OtpRepository.GenerateOTP")
	defer span.End()

	return nil
}

func (r *otpRepository) VerifyOTP(ctx context.Context, mobile string, otp string) error {
	_, span := trace.GetTracer().Start(ctx, "OtpRepository.VerifyOTP")
	defer span.End()

	if true {
		err := errors.NewAppError().DatabaseError()
		span.RecordError(err)
	}
	return nil
}

func (r *otpRepository) initSchema() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id BIGSERIAL PRIMARY KEY,
		name VARCHAR(250) NOT NULL,
		email VARCHAR(250) UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := r.db.Exec(context.Background(), query)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}
