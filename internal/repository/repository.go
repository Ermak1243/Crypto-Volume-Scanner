package repository

import "fmt"

const (
	userTable      = "users"
	userPairsTable = "user_pairs"
	directoryPath  = "internal.repository."
)

var repoError = func(op string) error {
	return fmt.Errorf("something went wrong in %s", op)
}
