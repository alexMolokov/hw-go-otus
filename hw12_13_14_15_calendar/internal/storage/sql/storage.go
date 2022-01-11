package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	st "github.com/alexMolokov/hw-go-otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/jmoiron/sqlx"

	// init postgres.
	_ "github.com/lib/pq"
)

var ErrConnectDB = errors.New("can't connect to database")

type Storage struct {
	db      *sqlx.DB
	maxConn int
	timeout time.Duration
}

func New(maxConn int) *Storage {
	return &Storage{
		maxConn: maxConn,
		timeout: 3 * time.Second,
	}
}

func (s *Storage) Connect(driverName string, dataSourceName string) error {
	db, err := sqlx.Connect(driverName, dataSourceName)
	if err != nil {
		return fmt.Errorf("%s %w", ErrConnectDB, err)
	}

	s.db = db
	if s.maxConn > 0 {
		s.db.SetMaxOpenConns(s.maxConn)
	}

	return nil
}

func (s *Storage) Close() error {
	if s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Storage) CreateEvent(e st.Event) (st.Event, error) {
	query := `INSERT INTO calendar.event 
		(title, description, owner_id, start_date, end_date, notify_for) 
		VALUES (:title, :description, :owner_id, :start_date, :end_date, :notify_for) 
		RETURNING event_id`
	rows, err := s.db.NamedQueryContext(s.getContext(), query, e)
	defer func() {
		if rows == nil {
			return
		}
		rows.Close()
	}()

	if err != nil {
		return e, fmt.Errorf("can't create event: %w", err)
	}

	var id int64
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return e, fmt.Errorf("can't get ID event: %w", err)
		}
	}

	return s.GetEvent(id)
}

func (s *Storage) UpdateEvent(e st.Event) (st.Event, error) {
	query := `UPDATE calendar.event SET 
		title = :title, 
        description = :description, 
        owner_id = :owner_id, 
        start_date = :start_date, end_date = :start_date, 
        notify_for = :notify_for,
		updated_at = CURRENT_TIMESTAMP
		WHERE event_id = :event_id`

	result, err := s.db.NamedExecContext(s.getContext(), query, e)
	if err != nil {
		return e, fmt.Errorf("can't update event id = %d : %w", e.ID, err)
	}

	count, _ := result.RowsAffected()
	if count == 0 {
		return e, fmt.Errorf("can't update event id = %d : %w", e.ID, err)
	}

	return e, nil
}

func (s *Storage) DeleteEvent(id int64) error {
	result, err := s.db.ExecContext(s.getContext(), "DELETE FROM calendar.event WHERE event_id = $1", id)
	if err != nil {
		return fmt.Errorf("can't delete event id = %d %w", id, err)
	}
	count, _ := result.RowsAffected()

	if count == 0 {
		return st.ErrEventNotFound
	}
	return nil
}

func (s *Storage) GetEvent(id int64) (st.Event, error) {
	e := st.Event{}
	row := s.db.QueryRowxContext(s.getContext(), "SELECT * FROM calendar.Event WHERE event_id = $1", id)
	if err := row.Err(); err != nil {
		return e, fmt.Errorf("can't get event id %d: %w", id, err)
	}
	if err := row.StructScan(&e); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e, st.ErrEventNotFound
		}
		return e, fmt.Errorf("can't get event id %d row scan : %w", id, err)
	}
	return e, nil
}

func (s *Storage) GetDayEvents(date time.Time) ([]st.Event, error) {
	return s.getEventsInRange(date, date)
}

func (s *Storage) GetWeekEvents(date time.Time) ([]st.Event, error) {
	return s.getEventsInRange(date, date.AddDate(0, 0, 6))
}

func (s *Storage) GetMonthEvents(date time.Time) ([]st.Event, error) {
	return s.getEventsInRange(date, date.AddDate(0, 1, 0))
}

func (s *Storage) CreateUser(u st.User) (st.User, error) {
	query := `INSERT INTO calendar.user 
		(email, first_name, last_name) 
		VALUES (:email, :first_name, :last_name) 
		RETURNING user_id`

	rows, err := s.db.NamedQueryContext(s.getContext(), query, u)
	defer func() {
		if rows == nil {
			return
		}
		rows.Close()
	}()

	if err != nil {
		return u, fmt.Errorf("can't create user: %w", err)
	}

	var id int64
	for rows.Next() {
		err = rows.Scan(&id)
	}

	if err != nil {
		return u, fmt.Errorf("can't get  user get last insert id: %w", err)
	}

	u.ID = id
	return u, nil
}

func (s *Storage) GetUser(id int64) (st.User, error) {
	u := st.User{}
	row := s.db.QueryRowxContext(s.getContext(), "SELECT * FROM calendar.user WHERE user_id = $1", id)
	if err := row.Err(); err != nil {
		return u, fmt.Errorf("can't get user id %d: %w", id, err)
	}

	if err := row.StructScan(&u); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return u, st.ErrUserNotFound
		}

		return u, fmt.Errorf("can't get event id %d row scan : %w", id, err)
	}
	return u, nil
}

func (s *Storage) UpdateUser(u st.User) (st.User, error) {
	query := `UPDATE calendar.user SET 
		email = :email, first_name = :first_name, last_name = :last_name
		WHERE user_id = :user_id`

	result, err := s.db.NamedExecContext(s.getContext(), query, u)
	if err != nil {
		return u, fmt.Errorf("can't update user id = %d : %w", u.ID, err)
	}

	count, cancel := result.RowsAffected()
	_ = cancel
	if count == 0 {
		return u, fmt.Errorf("can't update user id = %d : %w", u.ID, err)
	}

	return u, nil
}

func (s *Storage) getContext() context.Context {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	_ = cancel
	return ctx
}

func (s *Storage) getEventsInRange(start time.Time, end time.Time) ([]st.Event, error) {
	result := make([]st.Event, 0)
	query := `SELECT * FROM calendar.event
		WHERE (start_date <= $1 and end_date >= $1) OR 
		      (start_date <= $2 and end_date >= $2) OR 
		      (start_date >= $1 and end_date <= $2) 
		ORDER BY start_date ASC`
	err := s.db.Select(&result, query, start, end)
	return result, err
}
