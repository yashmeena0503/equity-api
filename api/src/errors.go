package src

import (
	"fmt"
	"net/http"
)

func ErrUnknownPlatform(platformName string) error {
	return ExpectedError {
		Status: http.StatusBadRequest,
		Message: fmt.Sprintf("Unknown platform: %q", platformName),
	}
}
func ErrOpenDB(cause error) error {
	return fmt.Errorf("Failed to open Database: %w", cause)
}
func ErrPingDB(cause error) error {
    return fmt.Errorf("Failed to ping db: %w", cause)
}
func ErrNewUser(cause error) error {
    return fmt.Errorf("Failed to create user: %w", cause)
}
func ErrCheckUserExistence(cause error) error {
    return fmt.Errorf("Failed to check user existence: %w", cause)
}
func ErrNewUserByPlatform(cause error) error {
    return fmt.Errorf("Failed to create new user by platform: %w", cause)
}
func ErrBeginTx(cause error) error {
    return fmt.Errorf("Failed to begin transaction: %w", cause)
}
func ErrCommitTx(cause error) error {
    return fmt.Errorf("Failed to commit transaction: %w", cause)
}
func ErrNewUnid(cause error) error {
    return fmt.Errorf("Failed to generate new unid: %w", cause)
}
func ErrLastInsertId(cause error) error {
	return fmt.Errorf("Failed to fetch last insert id: %w", cause)
}
func ErrNewTransaction(cause error) error {
	return fmt.Errorf("Failed to create new transaction row: %w", cause)
}
func ErrNewCoin(cause error) error {
	return fmt.Errorf("Failed to create new coin row: %w", cause)
}
func ErrNewTransactionLink(cause error) error {
	return fmt.Errorf("Failed to create new transaction_links row: %w", cause)
}
func ErrNewPayout(cause error) error {
	return fmt.Errorf("Failed to create new payouts row: %w", cause)
}
func ErrGetUnidByPlatform(cause error) error {
	return fmt.Errorf("Failed to get unid by platform: %w", cause)
}
func ErrEmptyPayout() error {
	return fmt.Errorf("Empty Payout Buffer")
}
