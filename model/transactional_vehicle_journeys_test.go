package model

import (
	"testing"
)

func Test_TransactionalVehicleJourneys_Find_NotFound(t *testing.T) {
	model := NewMemoryModel()
	vehicleJourneys := NewTransactionalVehicleJourneys(model)

	_, ok := vehicleJourneys.Find("6ba7b814-9dad-11d1-0-00c04fd430c8")
	if ok {
		t.Errorf("Find should return false when VehicleJourney isn't found")
	}
}

func Test_TransactionalVehicleJourneys_Find_Model(t *testing.T) {
	model := NewMemoryModel()
	vehicleJourneys := NewTransactionalVehicleJourneys(model)

	existingVehicleJourney := model.VehicleJourneys().New()
	model.VehicleJourneys().Save(&existingVehicleJourney)

	vehicleJourneyId := existingVehicleJourney.Id()

	vehicleJourney, ok := vehicleJourneys.Find(vehicleJourneyId)
	if !ok {
		t.Errorf("Find should return true when VehicleJourney is found")
	}
	if vehicleJourney.Id() != vehicleJourneyId {
		t.Errorf("Find should return a VehicleJourney with the given Id")
	}
}

func Test_TransactionalVehicleJourneys_Find_Saved(t *testing.T) {
	model := NewMemoryModel()
	vehicleJourneys := NewTransactionalVehicleJourneys(model)

	existingVehicleJourney := vehicleJourneys.New()
	vehicleJourneys.Save(&existingVehicleJourney)

	vehicleJourneyId := existingVehicleJourney.Id()

	vehicleJourney, ok := vehicleJourneys.Find(vehicleJourneyId)
	if !ok {
		t.Errorf("Find should return true when VehicleJourney is found")
	}
	if vehicleJourney.Id() != vehicleJourneyId {
		t.Errorf("Find should return a VehicleJourney with the given Id")
	}
}

func Test_TransactionalVehicleJourneys_FindAll(t *testing.T) {
	model := NewMemoryModel()
	vehicleJourneys := NewTransactionalVehicleJourneys(model)

	for i := 0; i < 5; i++ {
		existingVehicleJourney := vehicleJourneys.New()
		vehicleJourneys.Save(&existingVehicleJourney)
	}

	foundVehicleJourneys := vehicleJourneys.FindAll()

	if len(foundVehicleJourneys) != 5 {
		t.Errorf("FindAll should return all vehicleJourneys")
	}
}

func Test_TransactionalVehicleJourneys_Save(t *testing.T) {
	model := NewMemoryModel()
	vehicleJourneys := NewTransactionalVehicleJourneys(model)

	vehicleJourney := vehicleJourneys.New()
	objectid := NewObjectID("kind", "value")
	vehicleJourney.SetObjectID(objectid)

	if success := vehicleJourneys.Save(&vehicleJourney); !success {
		t.Errorf("Save should return true")
	}
	if vehicleJourney.Id() == "" {
		t.Errorf("New VehicleJourney identifier shouldn't be an empty string")
	}
	if _, ok := model.VehicleJourneys().Find(vehicleJourney.Id()); ok {
		t.Errorf("VehicleJourney shouldn't be saved before commit")
	}
	if _, ok := model.VehicleJourneys().FindByObjectId(objectid); ok {
		t.Errorf("VehicleJourney shouldn't be found by objectid before commit")
	}
}

func Test_TransactionalVehicleJourneys_Delete(t *testing.T) {
	model := NewMemoryModel()
	vehicleJourneys := NewTransactionalVehicleJourneys(model)

	existingVehicleJourney := model.VehicleJourneys().New()
	objectid := NewObjectID("kind", "value")
	existingVehicleJourney.SetObjectID(objectid)
	model.VehicleJourneys().Save(&existingVehicleJourney)

	vehicleJourneys.Delete(&existingVehicleJourney)

	_, ok := vehicleJourneys.Find(existingVehicleJourney.Id())
	if !ok {
		t.Errorf("VehicleJourney should not be deleted before commit")
	}
	_, ok = vehicleJourneys.FindByObjectId(objectid)
	if !ok {
		t.Errorf("VehicleJourney should be found by objectid before commit")
	}
}

func Test_TransactionalVehicleJourneys_Commit(t *testing.T) {
	model := NewMemoryModel()
	vehicleJourneys := NewTransactionalVehicleJourneys(model)

	// Test Save
	vehicleJourney := vehicleJourneys.New()
	objectid := NewObjectID("kind", "value")
	vehicleJourney.SetObjectID(objectid)
	vehicleJourneys.Save(&vehicleJourney)

	// Test Delete
	existingVehicleJourney := model.VehicleJourneys().New()
	secondObjectid := NewObjectID("kind", "value2")
	existingVehicleJourney.SetObjectID(secondObjectid)
	model.VehicleJourneys().Save(&existingVehicleJourney)
	vehicleJourneys.Delete(&existingVehicleJourney)

	vehicleJourneys.Commit()

	if _, ok := model.VehicleJourneys().Find(vehicleJourney.Id()); !ok {
		t.Errorf("VehicleJourney should be saved after commit")
	}
	if _, ok := model.VehicleJourneys().FindByObjectId(objectid); !ok {
		t.Errorf("VehicleJourney should be found by objectid after commit")
	}

	if _, ok := vehicleJourneys.Find(existingVehicleJourney.Id()); ok {
		t.Errorf("VehicleJourney should be deleted after commit")
	}
	if _, ok := vehicleJourneys.FindByObjectId(secondObjectid); ok {
		t.Errorf("VehicleJourney should not be foundable by objectid after commit")
	}
}

func Test_TransactionalVehicleJourneys_Rollback(t *testing.T) {
	model := NewMemoryModel()
	vehicleJourneys := NewTransactionalVehicleJourneys(model)

	vehicleJourney := vehicleJourneys.New()
	objectid := NewObjectID("kind", "value")
	vehicleJourney.SetObjectID(objectid)
	vehicleJourneys.Save(&vehicleJourney)

	vehicleJourneys.Rollback()
	vehicleJourneys.Commit()

	if _, ok := model.VehicleJourneys().Find(vehicleJourney.Id()); ok {
		t.Errorf("VehicleJourney should not be saved with a rollback")
	}
	if _, ok := model.VehicleJourneys().FindByObjectId(objectid); ok {
		t.Errorf("VehicleJourney should not be foundable by objectid with a rollback")
	}
}
