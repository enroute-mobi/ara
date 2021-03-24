package cache

import "fmt"

type ErrNotFound string

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("Cannot find key %s in cache table", string(e))
}
