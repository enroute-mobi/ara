package partners

import (
	"strings"
)

type LocalCredentialsIndex struct {
	byCredentials map[string]Id
	byIdentifier  map[Id][]string
}

func NewLocalCredentialsIndex() *LocalCredentialsIndex {
	return &LocalCredentialsIndex{
		byCredentials: make(map[string]Id),
		byIdentifier:  make(map[Id][]string),
	}
}

func (index *LocalCredentialsIndex) Index(modelId Id, localCredentials string) {
	splitCredentials := splitCredentials(localCredentials)

	// Delete from the index all elements of modelId not in the new localCredentials
	currentCredentials := index.byIdentifier[modelId]
	obsoleteCredentials := difference(currentCredentials, splitCredentials)
	for i := range obsoleteCredentials {
		delete(index.byCredentials, obsoleteCredentials[i])
	}

	for i := range splitCredentials {
		index.byCredentials[splitCredentials[i]] = modelId
	}
	index.byIdentifier[modelId] = splitCredentials
}

func (index *LocalCredentialsIndex) Find(c string) (modelId Id, ok bool) {
	modelId, ok = index.byCredentials[c]
	return
}

func (index *LocalCredentialsIndex) Delete(modelId Id) {
	currentCredentials, ok := index.byIdentifier[modelId]
	if !ok {
		return
	}

	for i := range currentCredentials {
		delete(index.byCredentials, currentCredentials[i])
	}
	delete(index.byIdentifier, modelId)
}

func (index *LocalCredentialsIndex) UniqCredentials(modelId Id, localCredentials string) bool {
	splitCredentials := splitCredentials(localCredentials)

	for i := range splitCredentials {
		x, ok := index.byCredentials[splitCredentials[i]]
		if ok && x != modelId {
			return false
		}
	}
	return true
}

// Split and trim spaces from the local_credentials string
// We expect something like local_credential + "," + local_credentials
func splitCredentials(c string) (r []string) {
	if c == "" {
		return []string{""}
	}
	sc := strings.Split(c, ",")
	for i := range sc {
		sc[i] = strings.TrimSpace(sc[i])
		if sc[i] != "" {
			r = append(r, sc[i])
		}
	}
	return
}

// Return all the element of slice a not in slice b
func difference(a, b []string) (diff []string) {
	mb := make(map[string]struct{}, len(b))
	for i := range b {
		mb[b[i]] = struct{}{}
	}

	for i := range a {
		if _, found := mb[a[i]]; !found {
			diff = append(diff, a[i])
		}
	}

	return diff
}
