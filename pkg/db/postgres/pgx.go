package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"zarinpal-platform/core/logger"
)

// Config holds the configuration for PostgreSQL connection
type Config struct {
	Host         string
	Port         int
	User         string
	Password     string
	Database     string
	PoolConfig   PoolConfig
	SSLMode      string
	MaxRetries   int
	RetryTimeout time.Duration
}

// PoolConfig holds the configuration for connection pool
type PoolConfig struct {
	MaxConnections        int32
	MinConnections        int32
	MaxConnLifetime       time.Duration
	MaxConnIdleTime       time.Duration
	HealthCheckPeriod     time.Duration
	MaxConnLifetimeJitter time.Duration
}

// DbConnection represents a PostgreSQL database connection
type DbConnection struct {
	pool   *pgxpool.Pool
	config *Config
}

// NewPgxPoolConnection creates a new PostgreSQL database connection
func NewPgxPoolConnection(ctx context.Context, config *Config) (*pgxpool.Pool, error) {

	// Set default values if not provided
	if config.PoolConfig.MaxConnections == 0 {
		config.PoolConfig.MaxConnections = 10
	}
	if config.PoolConfig.MinConnections == 0 {
		config.PoolConfig.MinConnections = 2
	}
	if config.PoolConfig.MaxConnLifetime == 0 {
		config.PoolConfig.MaxConnLifetime = time.Hour
	}
	if config.PoolConfig.MaxConnIdleTime == 0 {
		config.PoolConfig.MaxConnIdleTime = 30 * time.Minute
	}
	if config.PoolConfig.HealthCheckPeriod == 0 {
		config.PoolConfig.HealthCheckPeriod = time.Minute
	}
	if config.SSLMode == "" {
		config.SSLMode = "disable"
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}
	if config.RetryTimeout == 0 {
		config.RetryTimeout = 5 * time.Second
	}

	db := &DbConnection{
		config: config,
	}

	if err := db.connect(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db.pool, nil
}

// connect establishes the database connection with retry mechanism
func (db *DbConnection) connect(ctx context.Context) error {
	var err error
	for i := 0; i < db.config.MaxRetries; i++ {
		if i > 0 {
			logger.Log.Infof("Retrying database connection (attempt %d/%d)", i+1, db.config.MaxRetries)
			time.Sleep(db.config.RetryTimeout)
		}

		poolConfig, err := db.createPoolConfig()
		if err != nil {
			continue
		}

		db.pool, err = pgxpool.NewWithConfig(ctx, poolConfig)
		if err != nil {
			logger.Log.Errorf("Failed to create connection pool: %v", err)
			continue
		}

		// Test the connection
		if err = db.pool.Ping(ctx); err != nil {
			logger.Log.Errorf("Failed to ping database: %v", err)
			continue
		}

		logger.Log.Info("Successfully connected to PostgreSQL(PGX) database")
		return nil
	}

	return fmt.Errorf("failed to connect to database after %d attempts: %w", db.config.MaxRetries, err)
}

// createPoolConfig creates pgxpool configuration
func (db *DbConnection) createPoolConfig() (*pgxpool.Config, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		db.config.User,
		db.config.Password,
		db.config.Host,
		db.config.Port,
		db.config.Database,
		db.config.SSLMode,
	)

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %w", err)
	}

	// Configure the connection pool
	poolConfig.MaxConns = db.config.PoolConfig.MaxConnections
	poolConfig.MinConns = db.config.PoolConfig.MinConnections
	poolConfig.MaxConnLifetime = db.config.PoolConfig.MaxConnLifetime
	poolConfig.MaxConnIdleTime = db.config.PoolConfig.MaxConnIdleTime
	poolConfig.HealthCheckPeriod = db.config.PoolConfig.HealthCheckPeriod
	poolConfig.MaxConnLifetimeJitter = db.config.PoolConfig.MaxConnLifetimeJitter

	return poolConfig, nil
}
