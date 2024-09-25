package repository

import "fmt"

const (
	userTable      = "users"
	userPairsTable = "user_pairs"
	directoryPath  = "internal.repository."
)

var repoError = func(method string) error {
	return fmt.Errorf("something went wrong while %s", method)
}
