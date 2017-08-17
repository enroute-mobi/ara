package model

type StopAreaLineIds []LineId

func (ids *StopAreaLineIds) Add(id LineId) {
	if ids.Contains(id) {
		return
	}
	*ids = append(*ids, id)
}

func (ids StopAreaLineIds) Contains(id LineId) bool {
	for _, lineId := range ids {
		if lineId == id {
			return true
		}
	}
	return false
}
