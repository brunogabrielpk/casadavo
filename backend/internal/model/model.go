package model

import "time"

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Role      string    `json:"role"` // "cliente" | "gerente"
	CreatedAt time.Time `json:"created_at"`
}

type Table struct {
	ID       int64  `json:"id"`
	Number   int    `json:"number"`
	Location string `json:"location"` // "frente" | "fundos"
	Capacity int    `json:"capacity"`
	IsActive bool   `json:"is_active"`
}

type Availability struct {
	ID          int64  `json:"id"`
	Date        string `json:"date"` // YYYY-MM-DD
	IsOpen      bool   `json:"is_open"`
	AutoConfirm bool   `json:"auto_confirm"`
}

type TimeSlot struct {
	ID             int64  `json:"id"`
	AvailabilityID int64  `json:"availability_id"`
	SlotTime       string `json:"slot_time"` // HH:MM
}

type Reservation struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	TableID   int64     `json:"table_id"`
	SlotID    int64     `json:"slot_id"`
	Date      string    `json:"date"`
	PartySize int       `json:"party_size"`
	Status    string    `json:"status"` // "pending" | "confirmed" | "refused"
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TableExclusion struct {
	ID      int64  `json:"id"`
	TableID int64  `json:"table_id"`
	Date    string `json:"date"`
}

type ReservationDetail struct {
	Reservation
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
	TableNum  int    `json:"table_number"`
	Location  string `json:"location"`
	SlotTime  string `json:"slot_time"`
}
