package auth

import "fmt"

type ErrPermission struct {
	Reason string
}

func (e ErrPermission) Error() string {
	return fmt.Sprintf("permission denied: %s", e.Reason)
}
