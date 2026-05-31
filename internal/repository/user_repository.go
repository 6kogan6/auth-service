package repository

import (
	"auth-service/internal/model"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

func (userRepo *UserRepository) CreateUser(newUser model.User) (model.User, error) {
	query := `
		INSERT INTO users (name, email, password_hash, role, is_active)
		VALUES($1, $2, $3, $4, $5)
		RETURNING id, name, email, password_hash, role, is_active, created_at, updated_at;
	`
	var user model.User

	err := userRepo.pool.QueryRow(
		context.Background(), query,
		newUser.Name,
		newUser.Email,
		newUser.PasswordHash,
		newUser.Role,
		newUser.IsActive,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (userRepo *UserRepository) ExistsByEmail(email string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM users WHERE email = $1
		)
	`

	var exists bool

	err := userRepo.pool.QueryRow(context.Background(), query, email).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
