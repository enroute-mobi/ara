package model

import (
	"fmt"

	"github.com/af83/edwig/logger"
)

type Purifier struct {
	date string
}

func NewPurifier(days int) *Purifier {
	return &Purifier{
		date: DefaultClock().Now().AddDate(0, 0, -days).Format("2006-01-02"),
	}
}

func (p *Purifier) Purge() error {
	if p.date == "" {
		return fmt.Errorf("purifier date is empty")
	}

	_, err := Database.Exec("BEGIN;")
	if err != nil {
		return fmt.Errorf("database error: %v", err)
	}

	table_names := []string{"stop_areas", "lines", "vehicle_journeys", "stop_visits", "operators"}
	for i := range table_names {
		r, err := Database.Exec(p.query(table_names[i]))
		if err != nil {
			Database.Exec("ROLLBACK;")
			return fmt.Errorf("database error: %v", err)
		}
		ra, err := r.RowsAffected()
		if err != nil {
			logger.Log.Debugf("Error with Result.RowsAffected: %v", err)
		}
		logger.Log.Debugf("Purged %v %v from database", ra, table_names[i])
	}

	// Commit transaction
	_, err = Database.Exec("COMMIT;")
	if err != nil {
		return fmt.Errorf("database error: %v", err)
	}

	return nil
}

func (p Purifier) query(table_name string) string {
	return fmt.Sprintf("delete from %v where model_name < '%v';", table_name, p.date)
}
