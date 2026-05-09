package repository

import (
	"database/sql"
	"fmt"

	"casadavo/internal/model"
)

type ReservationRepo struct{ db *sql.DB }

func NewReservationRepo(db *sql.DB) *ReservationRepo { return &ReservationRepo{db} }

const detailQuery = `
SELECT r.id, r.user_id, r.table_id, r.slot_id, r.date, r.party_size, r.status,
       COALESCE(r.notes,''), r.created_at, r.updated_at,
       u.name, u.email, u.phone,
       t.number, t.location,
       ts.slot_time
FROM reservations r
JOIN users     u  ON u.id  = r.user_id
JOIN tables    t  ON t.id  = r.table_id
JOIN time_slots ts ON ts.id = r.slot_id`

func scanDetail(rows *sql.Rows) (model.ReservationDetail, error) {
	var d model.ReservationDetail
	err := rows.Scan(
		&d.ID, &d.UserID, &d.TableID, &d.SlotID, &d.Date, &d.PartySize, &d.Status,
		&d.Notes, &d.CreatedAt, &d.UpdatedAt,
		&d.UserName, &d.UserEmail, &d.UserPhone,
		&d.TableNum, &d.Location,
		&d.SlotTime,
	)
	return d, err
}

func (r *ReservationRepo) ListAll() ([]model.ReservationDetail, error) {
	rows, err := r.db.Query(detailQuery + ` ORDER BY r.date, ts.slot_time`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanDetails(rows)
}

func (r *ReservationRepo) ListByUser(userID int64) ([]model.ReservationDetail, error) {
	rows, err := r.db.Query(detailQuery+` WHERE r.user_id=? ORDER BY r.date, ts.slot_time`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanDetails(rows)
}

func scanDetails(rows *sql.Rows) ([]model.ReservationDetail, error) {
	var list []model.ReservationDetail
	for rows.Next() {
		d, err := scanDetail(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, d)
	}
	return list, rows.Err()
}

func (r *ReservationRepo) Create(res *model.Reservation) error {
	sqlRes, err := r.db.Exec(
		`INSERT INTO reservations (user_id, table_id, slot_id, date, party_size, status, notes)
		 VALUES (?,?,?,?,?,?,?)`,
		res.UserID, res.TableID, res.SlotID, res.Date, res.PartySize, res.Status, res.Notes,
	)
	if err != nil {
		return fmt.Errorf("create reservation: %w", err)
	}
	res.ID, _ = sqlRes.LastInsertId()
	return nil
}

func (r *ReservationRepo) Update(res *model.Reservation) error {
	_, err := r.db.Exec(
		`UPDATE reservations SET table_id=?, slot_id=?, date=?, party_size=?, notes=?, updated_at=CURRENT_TIMESTAMP WHERE id=?`,
		res.TableID, res.SlotID, res.Date, res.PartySize, res.Notes, res.ID,
	)
	return err
}

func (r *ReservationRepo) UpdateStatus(id int64, status string) error {
	_, err := r.db.Exec(
		`UPDATE reservations SET status=?, updated_at=CURRENT_TIMESTAMP WHERE id=?`,
		status, id,
	)
	return err
}

func (r *ReservationRepo) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM reservations WHERE id=?`, id)
	return err
}

func (r *ReservationRepo) FindByID(id int64) (*model.Reservation, error) {
	res := &model.Reservation{}
	err := r.db.QueryRow(
		`SELECT id, user_id, table_id, slot_id, date, party_size, status, COALESCE(notes,''), created_at, updated_at FROM reservations WHERE id=?`, id,
	).Scan(&res.ID, &res.UserID, &res.TableID, &res.SlotID, &res.Date, &res.PartySize, &res.Status, &res.Notes, &res.CreatedAt, &res.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return res, err
}

func (r *ReservationRepo) ExistsConflict(tableID, slotID int64, date string, excludeID int64) (bool, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM reservations WHERE table_id=? AND slot_id=? AND date=? AND id!=? AND status!='refused'`,
		tableID, slotID, date, excludeID,
	).Scan(&count)
	return count > 0, err
}
