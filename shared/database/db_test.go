package database

import (
	"testing"
)

func TestGetAccountByUsername(t *testing.T) {
	// Create a client under test
	c := New()
	// Query for a known account
	acc, err := c.GetAccountByUsername("briggySmalls90@gmail.com")
	// Ensure no err
	if err != nil {
		t.Error(err)
	}
	// Assert account fields
	if acc.AccountId != "aac45d02-c97d-442c-8431-336d578fdcf7" {
		t.Errorf("Unexpected accountId: %s", acc.AccountId)
	}
}