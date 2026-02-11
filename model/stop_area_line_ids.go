package model

import "slices"

type StopAreaLineIds []LineId

func (ids *StopAreaLineIds) Add(id LineId) {
	if !ids.Contains(id) {
		*ids = append(*ids, id)
	}
}

func (ids StopAreaLineIds) Contains(id LineId) bool {
	return slices.Contains(ids, id)
}

func (ids StopAreaLineIds) Copy() (t StopAreaLineIds) {
	t = append(t, ids...)
	return t
}
