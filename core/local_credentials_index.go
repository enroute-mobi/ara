package core

import (
	"strings"
)

type LocalCredentialsIndex struct {
	byCredentials map[string]PartnerId
	byIdentifier  map[PartnerId][]string
}

func NewLocalCredentialsIndex() *LocalCredentialsIndex {
	return &LocalCredentialsIndex{
		byCredentials: make(map[string]PartnerId),
		byIdentifier:  make(map[PartnerId][]string),
	}
}

func (index *LocalCredentialsIndex) Index(modelId PartnerId, localCredentials string) {
	splitCredentials := splitCredentials(localCredentials)

	// Delete from the index all elements of modelId not in the new localCredentials
	currentCredentials, _ := index.byIdentifier[modelId]
	obsoleteCredentials := difference(currentCredentials, splitCredentials)
	for i := range obsoleteCredentials {
		delete(index.byCredentials, obsoleteCredentials[i])
	}

	for i := range splitCredentials {
		index.byCredentials[splitCredentials[i]] = modelId
	}
	index.byIdentifier[modelId] = splitCredentials
}

func (index *LocalCredentialsIndex) Find(c string) (modelId PartnerId, ok bool) {
	modelId, ok = index.byCredentials[c]
	return
}

func (index *LocalCredentialsIndex) Delete(modelId PartnerId) {
	currentCredentials, ok := index.byIdentifier[modelId]
	if !ok {
		return
	}

	for i := range currentCredentials {
		delete(index.byCredentials, currentCredentials[i])
	}
	delete(index.byIdentifier, modelId)
}

func (index *LocalCredentialsIndex) UniqCredentials(modelId PartnerId, localCredentials string) bool {
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
