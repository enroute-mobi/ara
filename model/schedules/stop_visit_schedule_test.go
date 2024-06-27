package schedules

import (
	"testing"
	"time"
)

func Test_StopVisitsSchedule_Eq(t *testing.T) {
	svs1 := NewStopVisitSchedules()
	svs2 := NewStopVisitSchedules()

	if !svs1.Eq(svs2) {
		t.Error("The two schedules should be equal")
	}

	t1 := time.Now()
	svs1.SetArrivalTime(Aimed, t1)
	if svs1.Eq(svs2) {
		t.Error("The two schedules should not be equal")
	}

	t2 := time.Now().Add(10 * time.Second)
	svs2.SetArrivalTime(Expected, t2)
	if svs1.Eq(svs2) {
		t.Error("The two schedules should not be equal")
	}

	svs2.SetArrivalTime(Aimed, t1)
	svs1.SetArrivalTime(Expected, t2)

	if !svs1.Eq(svs2) {
		t.Error("The two schedules should be equal")
	}
}

func Test_StopVisitsSchedule_Include(t *testing.T) {
	svs1 := NewStopVisitSchedules()
	svs2 := NewStopVisitSchedules()

	if !svs1.Include(svs2) {
		t.Error("The first schedule should include the second")
	}

	t1 := time.Now()
	svs1.SetArrivalTime(Aimed, t1)
	if !svs1.Include(svs2) {
		t.Error("The first schedule should include the second")
	}

	t2 := time.Now().Add(10 * time.Second)
	svs2.SetArrivalTime(Expected, t2)
	if svs1.Include(svs2) {
		t.Error("The first schedule should not include the second")
	}

	svs1.SetArrivalTime(Expected, t2)
	if !svs1.Include(svs2) {
		t.Error("The first schedule should include the second")
	}

	svs2.SetArrivalTime(Aimed, t1)
	if !svs1.Include(svs2) {
		t.Error("The first schedule should include the second")
	}
}
