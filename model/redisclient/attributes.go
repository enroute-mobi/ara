package redisclient

// We need to define all common strings as consts to optimize memory allocation.
// Here we should define all attributes we need to find by, and all model names.
//
// This maybe should be shared in the whole application by creating an ara package or something
const (
	ModelID      = "Id"
	ReferentID   = "ReferentId"
	ByReferentID = ReferentID

	Line = "line"
)
