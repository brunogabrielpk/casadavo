package repository

import (
	"database/sql"
	"fmt"

	"casadavo/internal/model"
)

type AvailabilityRepo struct{ db *sql.DB }

func NewAvailabilityRepo(db *sql.DB) *AvailabilityRepo { return &AvailabilityRepo{db} }

func (r *AvailabilityRepo) List() ([]model.Availability, error) {
	rows, err := r.db.Query(`SELECT id, date, is_open, auto_confirm FROM availability ORDER BY date`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var avs []model.Availability
	for rows.Next() {
		var a model.Availability
		var isOpen, autoConfirm int
		if err := rows.Scan(&a.ID, &a.Date, &isOpen, &autoConfirm); err != nil {
			return nil, err
		}
		a.IsOpen = isOpen == 1
		a.AutoConfirm = autoConfirm == 1
		avs = append(avs, a)
	}
	return avs, rows.Err()
}

func (r *AvailabilityRepo) Create(a *model.Availability) error {
	res, err := r.db.Exec(
		`INSERT INTO availability (date, is_open, auto_confirm) VALUES (?,?,?)`,
		a.Date, boolToInt(a.IsOpen), boolToInt(a.AutoConfirm),
	)
	if err != nil {
		return fmt.Errorf("create availability: %w", err)
	}
	a.ID, _ = res.LastInsertId()
	return nil
}

func (r *AvailabilityRepo) Update(a *model.Availability) error {
	_, err := r.db.Exec(
		`UPDATE availability SET is_open=?, auto_confirm=? WHERE id=?`,
		boolToInt(a.IsOpen), boolToInt(a.AutoConfirm), a.ID,
	)
	return err
}

func (r *AvailabilityRepo) FindByID(id int64) (*model.Availability, error) {
	a := &model.Availability{}
	var isOpen, autoConfirm int
	err := r.db.QueryRow(
		`SELECT id, date, is_open, auto_confirm FROM availability WHERE id=?`, id,
	).Scan(&a.ID, &a.Date, &isOpen, &autoConfirm)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	a.IsOpen = isOpen == 1
	a.AutoConfirm = autoConfirm == 1
	return a, err
}

// --- Time Slots ---

type TimeSlotRepo struct{ db *sql.DB }

func NewTimeSlotRepo(db *sql.DB) *TimeSlotRepo { return &TimeSlotRepo{db} }

func (r *TimeSlotRepo) ListByAvailability(availID int64) ([]model.TimeSlot, error) {
	rows, err := r.db.Query(
		`SELECT id, availability_id, slot_time FROM time_slots WHERE availability_id=? ORDER BY slot_time`,
		availID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var slots []model.TimeSlot
	for rows.Next() {
		var s model.TimeSlot
		if err := rows.Scan(&s.ID, &s.AvailabilityID, &s.SlotTime); err != nil {
			return nil, err
		}
		slots = append(slots, s)
	}
	return slots, rows.Err()
}

func (r *TimeSlotRepo) Create(s *model.TimeSlot) error {
	res, err := r.db.Exec(
		`INSERT INTO time_slots (availability_id, slot_time) VALUES (?,?)`,
		s.AvailabilityID, s.SlotTime,
	)
	if err != nil {
		return fmt.Errorf("create time slot: %w", err)
	}
	s.ID, _ = res.LastInsertId()
	return nil
}

func (r *TimeSlotRepo) Delete(id int64) error {
	_, err := r.db.Exec(`DELETE FROM time_slots WHERE id=?`, id)
	return err
}

func (r *TimeSlotRepo) FindByID(id int64) (*model.TimeSlot, error) {
	s := &model.TimeSlot{}
	err := r.db.QueryRow(
		`SELECT id, availability_id, slot_time FROM time_slots WHERE id=?`, id,
	).Scan(&s.ID, &s.AvailabilityID, &s.SlotTime)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return s, err
}
