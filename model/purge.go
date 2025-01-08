package model

import (
	"fmt"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/logger"
)

type Purifier struct {
	date string
}

func NewPurifier(days int) *Purifier {
	return &Purifier{
		date: clock.DefaultClock().Now().AddDate(0, 0, -days).Format("2006-01-02"),
	}
}

func (p *Purifier) Purge() error {
	if p.date == "" {
		return fmt.Errorf("purifier date is empty")
	}

	tx, err := Database.Begin()
	if err != nil {
		return fmt.Errorf("database error: %v", err)
	}

	table_names := []string{
		"stop_areas",
		"stop_area_groups",
		"lines",
		"line_groups",
		"vehicle_journeys",
		"stop_visits",
		"operators",
	}
	for i := range table_names {
		r, err := tx.Exec(p.query(table_names[i]))
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("database error: %v", err)
		}
		ra, err := r.RowsAffected()
		if err != nil {
			logger.Log.Debugf("Error with Result.RowsAffected: %v", err)
		}
		logger.Log.Printf("Purged %v %v from database", ra, table_names[i])
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("database error: %v", err)
	}

	return nil
}

func (p Purifier) query(table_name string) string {
	return fmt.Sprintf("delete from %v where model_name < '%v';", table_name, p.date)
}
