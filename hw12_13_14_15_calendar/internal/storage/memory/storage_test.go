package memorystorage

import (
	"sort"
	"testing"
	"time"

	st "github.com/alexMolokov/hw-go-otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/emvi/null"
	"github.com/stretchr/testify/require"
)

func getUser() st.User {
	return st.User{
		Email:     "sportmalex@yandex.ru",
		FirstName: null.NewString("Alex", true),
		LastName:  null.NewString("Molokov", true),
	}
}

func getDate(year int, day int, month time.Month) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func TestUser(t *testing.T) {
	t.Run("create user", func(t *testing.T) {
		s := New()
		u := getUser()
		cu, err := s.CreateUser(u)
		require.Nil(t, err, "Ошибка при создании пользователя")
		require.Greater(t, cu.ID, int64(0), "ID пользователя должно быть больше нуля")
		require.Equal(t, cu.LastName.String, u.LastName.String, "Фамилии должны сопадать")
		require.Equal(t, cu.FirstName.String, u.FirstName.String, "Имена должны совпадать")
	})
	t.Run("get user", func(t *testing.T) {
		s := New()
		u := getUser()
		cu, _ := s.CreateUser(u)
		gu, err := s.GetUser(cu.ID)

		require.Nil(t, err, "Пользователь должен существовать")
		require.Equal(t, cu.ID, gu.ID, "ID должны совпадать")

		nextID := userID()
		_, err = s.GetUser(nextID)
		require.ErrorIs(t, err, st.ErrUserNotFound, "Пользователь не должен существовать")
	})
	t.Run("update user", func(t *testing.T) {
		s := New()
		u := getUser()
		cu, _ := s.CreateUser(u)
		cu.FirstName = null.NewString("Alex update", true)
		s.UpdateUser(cu)
		gu, err := s.GetUser(cu.ID)

		require.Nil(t, err, "Пользователь должен существовать")
		require.Equal(t, cu.FirstName.String, gu.FirstName.String, "Имена должны сопадать")
	})
}

func TestEvent(t *testing.T) {
	t.Run("crud event", func(t *testing.T) {
		s := New()

		startDate := getDate(1900, 1, time.January)
		endDate := startDate.AddDate(0, 0, 1)

		event, err := s.CreateEvent(st.Event{OwnerID: 1, StartDate: startDate, EndDate: endDate})
		require.Nil(t, err)
		require.NotNil(t, event)

		gEvent, err := s.GetEvent(event.ID)
		require.Nil(t, err)
		require.Equal(t, event, gEvent)

		gEvent.Description = null.NewString("update event", true)
		_, err = s.UpdateEvent(gEvent)
		require.Nil(t, err)

		uEvent, _ := s.GetEvent(event.ID)
		require.Equal(t, gEvent.Description.String, uEvent.Description.String, "Описания  должны сопадать")

		err = s.DeleteEvent(uEvent.ID)
		require.Nil(t, err)
		_, err = s.GetEvent(uEvent.ID)
		require.ErrorIs(t, err, st.ErrEventNotFound)
	})
	t.Run("get day events", func(t *testing.T) {
		s := New()

		s.events[1] = st.Event{
			ID:        1,
			OwnerID:   1,
			StartDate: getDate(2022, 1, time.January),
			EndDate:   getDate(2022, 2, time.January),
		}
		s.events[2] = st.Event{
			ID:        2,
			OwnerID:   1,
			StartDate: getDate(2022, 1, time.January),
			EndDate:   getDate(2022, 7, time.January),
		}
		s.events[3] = st.Event{
			ID:        3,
			OwnerID:   1,
			StartDate: getDate(2022, 3, time.January),
			EndDate:   getDate(2022, 11, time.January),
		}
		s.events[4] = st.Event{
			ID:        4,
			OwnerID:   1,
			StartDate: getDate(2022, 13, time.January),
			EndDate:   getDate(2022, 16, time.January),
		}
		s.events[5] = st.Event{
			ID:        5,
			OwnerID:   1,
			StartDate: getDate(2020, 13, time.February),
			EndDate:   getDate(2020, 16, time.March),
		}

		result, err := s.GetDayEvents(getDate(2021, 2, time.January))
		require.Nil(t, err)
		require.Equal(t, 0, len(result))

		result, err = s.GetDayEvents(getDate(2022, 2, time.January))
		require.Nil(t, err)

		sort.Slice(result, func(i int, j int) bool {
			return result[i].ID < result[j].ID
		})

		require.Equal(t, []st.Event{s.events[1], s.events[2]}, result)
	})
	t.Run("get week events", func(t *testing.T) {
		s := New()

		s.events[1] = st.Event{
			ID: 1, OwnerID: 1,
			StartDate: getDate(2022, 1, time.January),
			EndDate:   getDate(2022, 2, time.January),
		}
		s.events[2] = st.Event{
			ID:        2,
			OwnerID:   1,
			StartDate: getDate(2022, 1, time.January),
			EndDate:   getDate(2022, 7, time.January),
		}
		s.events[3] = st.Event{
			ID:        3,
			OwnerID:   1,
			StartDate: getDate(2022, 3, time.January),
			EndDate:   getDate(2022, 11, time.January),
		}
		s.events[4] = st.Event{
			ID:        4,
			OwnerID:   1,
			StartDate: getDate(2022, 13, time.January),
			EndDate:   getDate(2022, 16, time.January),
		}

		result, err := s.GetWeekEvents(getDate(2021, 2, time.January))
		require.Nil(t, err)
		require.Equal(t, 0, len(result))

		result, err = s.GetWeekEvents(getDate(2022, 2, time.January))
		require.Nil(t, err)

		sort.Slice(result, func(i int, j int) bool {
			return result[i].ID < result[j].ID
		})

		require.Equal(t, []st.Event{s.events[1], s.events[2], s.events[3]}, result)
	})
	t.Run("get months events", func(t *testing.T) {
		s := New()

		s.events[1] = st.Event{
			ID:        1,
			OwnerID:   1,
			StartDate: getDate(2022, 1, time.January),
			EndDate:   getDate(2022, 2, time.January),
		}
		s.events[2] = st.Event{
			ID: 2, OwnerID: 1,
			StartDate: getDate(2022, 1, time.January),
			EndDate:   getDate(2022, 7, time.January),
		}
		s.events[3] = st.Event{
			ID:        3,
			OwnerID:   1,
			StartDate: getDate(2022, 3, time.January),
			EndDate:   getDate(2022, 11, time.January),
		}
		s.events[4] = st.Event{
			ID:        4,
			OwnerID:   1,
			StartDate: getDate(2022, 13, time.January),
			EndDate:   getDate(2022, 16, time.January),
		}

		result, err := s.GetMonthEvents(getDate(2022, 6, time.January))
		sort.Slice(result, func(i int, j int) bool {
			return result[i].ID < result[j].ID
		})
		require.Nil(t, err)

		require.Equal(t, []st.Event{s.events[2], s.events[3], s.events[4]}, result)
	})
}
