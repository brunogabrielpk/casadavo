package service

import (
	"errors"

	"casadavo/internal/model"
	"casadavo/internal/repository"
)

type ReservationService struct {
	reservations *repository.ReservationRepo
	slots        *repository.TimeSlotRepo
	avail        *repository.AvailabilityRepo
	layout       *repository.LayoutRepo
}

func NewReservationService(
	reservations *repository.ReservationRepo,
	slots *repository.TimeSlotRepo,
	avail *repository.AvailabilityRepo,
	layout *repository.LayoutRepo,
) *ReservationService {
	return &ReservationService{reservations, slots, avail, layout}
}

func (s *ReservationService) Create(res *model.Reservation) error {
	slot, err := s.slots.FindByID(res.SlotID)
	if err != nil || slot == nil {
		return errors.New("time slot not found")
	}

	av, err := s.avail.FindByID(slot.AvailabilityID)
	if err != nil || av == nil || !av.IsOpen {
		return errors.New("date not available for reservations")
	}

	excluded, err := s.layout.IsExcluded(res.TableID, res.Date)
	if err != nil {
		return err
	}
	if excluded {
		return errors.New("table not available on this date")
	}

	conflict, err := s.reservations.ExistsConflict(res.TableID, res.SlotID, res.Date, 0)
	if err != nil {
		return err
	}
	if conflict {
		return errors.New("table already reserved for this slot")
	}

	if av.AutoConfirm {
		res.Status = "confirmed"
	} else {
		res.Status = "pending"
	}

	return s.reservations.Create(res)
}

func (s *ReservationService) Update(res *model.Reservation, requesterID int64, isManager bool) error {
	existing, err := s.reservations.FindByID(res.ID)
	if err != nil || existing == nil {
		return errors.New("reservation not found")
	}
	if !isManager && existing.UserID != requesterID {
		return errors.New("forbidden")
	}

	conflict, err := s.reservations.ExistsConflict(res.TableID, res.SlotID, res.Date, res.ID)
	if err != nil {
		return err
	}
	if conflict {
		return errors.New("table already reserved for this slot")
	}

	return s.reservations.Update(res)
}

func (s *ReservationService) UpdateStatus(id int64, status string) error {
	if status != "confirmed" && status != "refused" {
		return errors.New("invalid status")
	}
	return s.reservations.UpdateStatus(id, status)
}

func (s *ReservationService) Delete(id, requesterID int64, isManager bool) error {
	existing, err := s.reservations.FindByID(id)
	if err != nil || existing == nil {
		return errors.New("reservation not found")
	}
	if !isManager && existing.UserID != requesterID {
		return errors.New("forbidden")
	}
	return s.reservations.Delete(id)
}
