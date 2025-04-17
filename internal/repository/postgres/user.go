package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliskhannn/pvz-service/internal/auth"
	"github.com/aliskhannn/pvz-service/internal/domain"
	"github.com/aliskhannn/pvz-service/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	if user.Email == "" || user.Password == "" || user.Role == "" {
		return fmt.Errorf("email, password and role are required")
	}

	if user.Role != "employee" && user.Role != "moderator" {
		return fmt.Errorf("role must be 'employee' or 'moderator'")
	}

	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}
	user.Password = hashedPassword

	query := `INSERT INTO users (email, password, role) VALUES ($1, $2, $3)`
	_, err = r.db.Exec(ctx, query, user.Email, user.Password, user.Role)
	if err != nil {
		return fmt.Errorf("user could not be created: %w", err)
	}

	return nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}

	var user domain.User

	query := `SELECT id, email, password, role FROM users WHERE email = $1`
	err := r.db.QueryRow(ctx, query, email).Scan(&user.Id, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf(`user not found`)
		}
		return nil, fmt.Errorf("error getting user by email: %w", err)
	}

	return &user, nil
}
