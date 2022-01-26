package memorystorage

import (
	"sync"
	"sync/atomic"
	"time"

	st "github.com/alexMolokov/hw-go-otus/hw12_13_14_15_calendar/internal/storage"
)

var (
	eventID = generateID()
	userID  = generateID()
)

func generateID() func() int64 {
	var id int64
	return func() int64 {
		return atomic.AddInt64(&id, 1)
	}
}

func isDateBefore(t time.Time, u time.Time) bool {
	return t.Equal(u) || t.Before(u)
}

func isDateAfter(t time.Time, u time.Time) bool {
	return t.Equal(u) || t.After(u)
}

type Storage struct {
	events map[int64]st.Event
	users  map[int64]st.User
	sync.RWMutex
}

func New() *Storage {
	return &Storage{
		events: make(map[int64]st.Event),
		users:  make(map[int64]st.User),
	}
}

func (s *Storage) CreateEvent(e st.Event) (st.Event, error) {
	s.Lock()
	defer s.Unlock()

	e.ID = eventID()
	s.events[e.ID] = e
	return e, nil
}

func (s *Storage) GetEvent(id int64) (st.Event, error) {
	s.RLock()
	defer s.RUnlock()

	e, ok := s.events[id]

	if !ok {
		return e, st.ErrEventNotFound
	}

	return e, nil
}

func (s *Storage) UpdateEvent(e st.Event) (st.Event, error) {
	s.RLock()
	defer s.RUnlock()

	if _, ok := s.events[e.ID]; !ok {
		return e, st.ErrEventNotFound
	}

	s.events[e.ID] = e
	return e, nil
}

func (s *Storage) DeleteEvent(id int64) error {
	s.RLock()
	if _, ok := s.events[id]; !ok {
		return st.ErrEventNotFound
	}
	s.RUnlock()

	s.Lock()
	delete(s.events, id)
	s.Unlock()

	return nil
}

func (s *Storage) GetDayEvents(date time.Time) ([]st.Event, error) {
	return s.getEventsInRange(date, date), nil
}

func (s *Storage) GetWeekEvents(date time.Time) ([]st.Event, error) {
	return s.getEventsInRange(date, date.AddDate(0, 0, 6)), nil
}

func (s *Storage) GetMonthEvents(date time.Time) ([]st.Event, error) {
	return s.getEventsInRange(date, date.AddDate(0, 1, 0)), nil
}

func (s *Storage) getEventsInRange(start time.Time, end time.Time) []st.Event {
	var result []st.Event

	s.RLock()
	defer s.RUnlock()

	for _, v := range s.events {
		isStartInRange := isDateAfter(start, v.StartDate) && isDateBefore(start, v.EndDate)
		isEndInRange := isDateAfter(end, v.StartDate) && isDateBefore(end, v.EndDate)
		isCoverRange := isDateAfter(end, v.EndDate) && isDateBefore(start, v.StartDate)
		if isStartInRange || isEndInRange || isCoverRange {
			result = append(result, v)
		}
	}
	return result
}

func (s *Storage) CreateUser(u st.User) (st.User, error) {
	s.Lock()
	defer s.Unlock()

	u.ID = userID()
	s.users[u.ID] = u
	return u, nil
}

func (s *Storage) UpdateUser(u st.User) (st.User, error) {
	s.Lock()
	defer s.Unlock()

	_, ok := s.users[u.ID]
	if !ok {
		return u, st.ErrUserNotFound
	}
	s.users[u.ID] = u

	return u, nil
}

func (s *Storage) GetUser(id int64) (st.User, error) {
	s.RLock()
	defer s.RUnlock()

	u, ok := s.users[id]
	if !ok {
		return u, st.ErrUserNotFound
	}

	return u, nil
}

func (s *Storage) Close() error {
	return nil
}
