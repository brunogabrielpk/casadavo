package repository

import (
	"database/sql"
	"fmt"

	"casadavo/internal/model"
)

type TableRepo struct{ db *sql.DB }

func NewTableRepo(db *sql.DB) *TableRepo { return &TableRepo{db} }

func (r *TableRepo) List() ([]model.Table, error) {
	rows, err := r.db.Query(`SELECT id, number, location, capacity, is_active FROM tables ORDER BY number`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tables []model.Table
	for rows.Next() {
		var t model.Table
		var active int
		if err := rows.Scan(&t.ID, &t.Number, &t.Location, &t.Capacity, &active); err != nil {
			return nil, err
		}
		t.IsActive = active == 1
		tables = append(tables, t)
	}
	return tables, rows.Err()
}

func (r *TableRepo) Create(t *model.Table) error {
	res, err := r.db.Exec(
		`INSERT INTO tables (number, location, capacity, is_active) VALUES (?,?,?,?)`,
		t.Number, t.Location, t.Capacity, boolToInt(t.IsActive),
	)
	if err != nil {
		return fmt.Errorf("create table: %w", err)
	}
	t.ID, _ = res.LastInsertId()
	return nil
}

func (r *TableRepo) Update(t *model.Table) error {
	_, err := r.db.Exec(
		`UPDATE tables SET number=?, location=?, capacity=?, is_active=? WHERE id=?`,
		t.Number, t.Location, t.Capacity, boolToInt(t.IsActive), t.ID,
	)
	return err
}

func (r *TableRepo) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM tables WHERE id=?`, id)
	return err
}

func (r *TableRepo) FindByID(id int64) (*model.Table, error) {
	t := &model.Table{}
	var active int
	err := r.db.QueryRow(
		`SELECT id, number, location, capacity, is_active FROM tables WHERE id=?`, id,
	).Scan(&t.ID, &t.Number, &t.Location, &t.Capacity, &active)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	t.IsActive = active == 1
	return t, err
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
