package src

import (
	"fmt"
	"net/http"
)

type errorsNamespace struct {}
var Err = errorsNamespace {}

func (errorsNamespace) UnknownPlatform(platformName string) error {
	return ExpectedError {
		Status: http.StatusBadRequest,
		Message: fmt.Sprintf("Unknown platform: %q", platformName),
	}
}
func (errorsNamespace) OpenDB(cause error) error {
	return fmt.Errorf("Failed to open Database: %w", cause)
}
func (errorsNamespace) PingDB(cause error) error {
    return fmt.Errorf("Failed to ping db: %w", cause)
}
func (errorsNamespace) NewUser(cause error) error {
    return fmt.Errorf("Failed to create user: %w", cause)
}
func (errorsNamespace) CheckUserExistence(cause error) error {
    return fmt.Errorf("Failed to check user existence: %w", cause)
}
func (errorsNamespace) NewUserByPlatform(cause error) error {
    return fmt.Errorf("Failed to create new user by platform: %w", cause)
}
func (errorsNamespace) BeginTx(cause error) error {
    return fmt.Errorf("Failed to begin transaction: %w", cause)
}
func (errorsNamespace) CommitTx(cause error) error {
    return fmt.Errorf("Failed to commit transaction: %w", cause)
}
