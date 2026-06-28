package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/HackMateGolang/user-service/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) (string, error) {
	_, err := r.db.Exec(ctx, "INSERT INTO users (login, username) VALUES ($1, $2)", user.Login, user.Username)
	if err != nil {
		return "", fmt.Errorf("REPO: inserting user error: %w", err)
	}

	return user.Login, nil
}

func (r *UserRepository) ReadUser(ctx context.Context, req *models.ReadUserRequest) (*models.User, error) {
	if req.Login == "" {
		return nil, fmt.Errorf("REPO: Login is empty")
	}

	var user models.User

	query := `SELECT 
		login, username, first_name, second_name, patronymic,
		stack, description, contacts, short_desc, avatar 
	FROM users
	WHERE login=$1`

	err := r.db.QueryRow(ctx, query, req.Login).Scan(
		&user.Login,
		&user.Username,
		&user.FirstName,
		&user.SecondName,
		&user.Patronymic,
		&user.Stack,
		&user.Description,
		&user.Contacts,
		&user.ShortDesc,
		&user.Avatar,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("REPO: User not found: %w", err)
		}
		return nil, fmt.Errorf("REPO: read user db error: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) ReplaceUser(ctx context.Context, req *models.UpdateUserRequest) (bool, error) {
	query := `UPDATE users
		SET
		username = $2, 
		first_name = $3, 
		second_name = $4, 
		patronymic = $5,
		stack = $6, 
		description = $7, 
		contacts = $8, 
		short_desc = $9, 
		avatar = $10
		WHERE login=$1`

	res, err := r.db.Exec(ctx, query,
		req.Login, req.Username, req.FirstName, req.SecondName, req.Patronymic,
		req.Stack, req.Description, req.Contacts, req.ShortDesc, req.Avatar,
	)

	if err != nil {
		return false, fmt.Errorf("REPO: Replacing user error: %w", err)
	}

	if res.RowsAffected() == 0 {
		return false, fmt.Errorf("REPO: User not found")
	}

	return true, nil
}

func (r *UserRepository) PatchUser(ctx context.Context, req *models.PatchUserRequest) (bool, error) {
	query := `UPDATE users
		SET
		username = COALESCE($2, username), 
		first_name = COALESCE($3, first_name), 
		second_name = COALESCE($4, second_name), 
		patronymic = COALESCE($5, patronymic),
		stack = COALESCE($6, stack), 
		description = COALESCE($7, description), 
		contacts = COALESCE($8, contacts),
		short_desc = COALESCE($9, short_desc), 
		avatar = COALESCE($10, avatar)
		WHERE login=$1`

	res, err := r.db.Exec(ctx, query,
		req.Login, req.Username, req.FirstName, req.SecondName, req.Patronymic,
		req.Stack, req.Description, req.Contacts, req.ShortDesc, req.Avatar,
	)

	if err != nil {
		return false, fmt.Errorf("REPO: User patching error: %w", err)
	}

	if res.RowsAffected() == 0 {
		return false, fmt.Errorf("REPO: User not found")
	}

	return true, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, req *models.DeleteUserRequest) (bool, error) {
	query := `DELETE FROM users WHERE login=$1`
	res, err := r.db.Exec(ctx, query, req.Login)

	if err != nil {
		return false, fmt.Errorf("REPO: User deleting failed: %w", err)
	}

	if res.RowsAffected() == 0 {
		return false, fmt.Errorf("REPO: User not found")
	}

	return true, nil
}
