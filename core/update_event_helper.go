package core

func GetModelReferenceSlice(refs map[string]struct{}) []string {
	refSlice := make([]string, len(refs))
	i := 0
	for ref := range refs {
		refSlice[i] = ref
		i++
	}
	return refSlice
}
