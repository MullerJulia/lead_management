package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateUUID(t *testing.T) {
	tests := []struct {
		testID int
	}{
		{testID: 1},
		{testID: 2},
	}

	uuids := make(map[string]bool) // to store and check for uniqueness

	for _, tc := range tests {
		t.Run(fmt.Sprintf("TestUUIDGeneration%d", tc.testID), func(t *testing.T) {
			uuid := GenerateUUID()
			assert.NotEmpty(t, uuid, "UUID should not be empty")

			// Check for uniqueness among all generated UUIDs in this test suite
			if _, exists := uuids[uuid]; exists {
				t.Errorf("UUID is not unique: %s", uuid)
			}
			uuids[uuid] = true // Mark this UUID as seen
		})
	}
}
