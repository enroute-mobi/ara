package model

type StopAreaLineIds []LineId

func (ids *StopAreaLineIds) Add(id LineId) {
	if !ids.Contains(id) {
		*ids = append(*ids, id)
	}
}

func (ids StopAreaLineIds) Contains(id LineId) bool {
	for _, lineId := range ids {
		if lineId == id {
			return true
		}
	}
	return false
}

func (ids StopAreaLineIds) Copy() (t StopAreaLineIds) {
	t = append(t, ids...)
	return t
}
