package persistence

import (
	"context"
	"database/sql"
	"time"

	"github.com/LiliBeta/auth-service/internal/domain"

	"github.com/google/uuid"
)

type MySQLUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *MySQLUserRepository {
	return &MySQLUserRepository{db: db}
}

func (r *MySQLUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, email, password, name, created_at, updated_at, verified 
	          FROM users WHERE email = ?`

	row := r.db.QueryRowContext(ctx, query, email)

	var user domain.User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Verified,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *MySQLUserRepository) Save(ctx context.Context, user *domain.User) error {
	if user.ID == "" {
		user.ID = uuid.New().String()
		user.CreatedAt = time.Now()
		user.Verified = false

		query := `INSERT INTO users (id, email, password, name, created_at, verified) 
		          VALUES (?, ?, ?, ?, ?, ?)`
		_, err := r.db.ExecContext(ctx, query,
			user.ID, user.Email, user.Password, user.Name, user.CreatedAt, user.Verified)
		return err
	}

	user.UpdatedAt = time.Now()
	query := `UPDATE users SET email=?, password=?, name=?, updated_at=?, verified=? 
	          WHERE id=?`
	_, err := r.db.ExecContext(ctx, query,
		user.Email, user.Password, user.Name, user.UpdatedAt, user.Verified, user.ID)
	return err
}

// Implementar otros métodos del repositorio...
