package storage

import (
	"errors"
	"time"

	"github.com/emvi/null"
)

var (
	ErrEventNotFound = errors.New("event not found")
	ErrUserNotFound  = errors.New("user not found")
)

type Event struct {
	ID          int64       `db:"event_id"`
	Title       string      `db:"title"`
	Description null.String `db:"description"`
	StartDate   time.Time   `db:"start_date"`
	EndDate     time.Time   `db:"end_date"`
	OwnerID     int64       `db:"owner_id"`
	NotifyFor   int16       `db:"notify_for"`
	CreatedAt   time.Time   `db:"created_at"`
	UpdatedAt   null.Time   `db:"updated_at"`
}

func (e *Event) SetDescription(description string) *Event {
	e.Description.String = description
	e.Description.Valid = true
	return e
}

type User struct {
	ID        int64       `db:"user_id"`
	FirstName null.String `db:"first_name"`
	LastName  null.String `db:"last_name"`
	Email     string      `db:"email"`
}

func (u *User) SetFirstName(firstName string) *User {
	u.FirstName.String = firstName
	u.FirstName.Valid = true
	return u
}

func (u *User) SetLastName(lastName string) *User {
	u.LastName.String = lastName
	u.LastName.Valid = true
	return u
}

type CalendarEvent interface {
	CreateEvent(e Event) (Event, error)
	UpdateEvent(e Event) (Event, error)
	DeleteEvent(id int64) error
	GetEvent(id int64) (Event, error)

	GetDayEvents(date time.Time) ([]Event, error)
	GetWeekEvents(time.Time) ([]Event, error)
	GetMonthEvents(time.Time) ([]Event, error)
}

type CalendarUser interface {
	CreateUser(u User) (User, error)
	GetUser(id int64) (User, error)
	UpdateUser(u User) (User, error)
}

type Closer interface {
	Close() error
}

type CalendarStorage interface {
	CalendarEvent
	CalendarUser
	Closer
}
