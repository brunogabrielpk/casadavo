package repository

import (
	"database/sql"

	"casadavo/internal/model"
)

type LayoutRepo struct{ db *sql.DB }

func NewLayoutRepo(db *sql.DB) *LayoutRepo { return &LayoutRepo{db} }

func (r *LayoutRepo) ListByDate(date string) ([]model.TableExclusion, error) {
	rows, err := r.db.Query(
		`SELECT id, table_id, date FROM table_exclusions WHERE date=? ORDER BY table_id`, date,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.TableExclusion
	for rows.Next() {
		var e model.TableExclusion
		if err := rows.Scan(&e.ID, &e.TableID, &e.Date); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func (r *LayoutRepo) Add(e *model.TableExclusion) error {
	res, err := r.db.Exec(
		`INSERT OR IGNORE INTO table_exclusions (table_id, date) VALUES (?,?)`,
		e.TableID, e.Date,
	)
	if err != nil {
		return err
	}
	e.ID, _ = res.LastInsertId()
	return nil
}

func (r *LayoutRepo) Remove(id int64) error {
	_, err := r.db.Exec(`DELETE FROM table_exclusions WHERE id=?`, id)
	return err
}

func (r *LayoutRepo) IsExcluded(tableID int64, date string) (bool, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM table_exclusions WHERE table_id=? AND date=?`, tableID, date,
	).Scan(&count)
	return count > 0, err
}
