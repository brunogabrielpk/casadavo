package repository

import (
	"database/sql"
	"fmt"

	"casadavo/internal/model"
)

type UserRepo struct{ db *sql.DB }

func NewUserRepo(db *sql.DB) *UserRepo { return &UserRepo{db} }

func (r *UserRepo) Create(u *model.User) error {
	res, err := r.db.Exec(
		`INSERT INTO users (name, email, phone, password, role) VALUES (?,?,?,?,?)`,
		u.Name, u.Email, u.Phone, u.Password, u.Role,
	)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	u.ID, _ = res.LastInsertId()
	return nil
}

func (r *UserRepo) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	err := r.db.QueryRow(
		`SELECT id, name, email, phone, password, role, created_at FROM users WHERE email=?`, email,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.Password, &u.Role, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return u, err
}

func (r *UserRepo) EnsureAdminRole(email string) error {
	_, err := r.db.Exec(`UPDATE users SET role='gerente' WHERE email=?`, email)
	return err
}

func (r *UserRepo) FindByID(id int64) (*model.User, error) {
	u := &model.User{}
	err := r.db.QueryRow(
		`SELECT id, name, email, phone, password, role, created_at FROM users WHERE id=?`, id,
	).Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.Password, &u.Role, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return u, err
}
