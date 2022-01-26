package storagefactory

import (
	"fmt"

	"github.com/alexMolokov/hw-go-otus/hw12_13_14_15_calendar/internal/config"
	"github.com/alexMolokov/hw-go-otus/hw12_13_14_15_calendar/internal/storage"
	ms "github.com/alexMolokov/hw-go-otus/hw12_13_14_15_calendar/internal/storage/memory"
	sqls "github.com/alexMolokov/hw-go-otus/hw12_13_14_15_calendar/internal/storage/sql"
)

func NewStorage(config *config.Config) (storage.CalendarStorage, error) {
	switch config.DriverDB {
	case "postgres":
		cdb := config.DB
		s := sqls.New(cdb.MaxConnectionPool)
		dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=%s",
			cdb.User, cdb.Password, cdb.Name, cdb.Host, cdb.Port, cdb.SslMode)
		err := s.Connect(config.DriverDB, dsn)
		return s, err
	default:
		return ms.New(), nil
	}
}
