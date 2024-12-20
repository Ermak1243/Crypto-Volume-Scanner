package tests

import (
	"cvs/internal/mocks"
	"cvs/internal/service/exchange"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewBybit tests the NewBybit function
func TestNewBybit(t *testing.T) {
	t.Parallel() // Allows this test to run in parallel with other tests

	// Create mocks for services
	mockUserService := mocks.NewUserService(t)
	mockUserPairsService := mocks.NewUserPairsService(t)
	mockHttpRequestService := mocks.NewHttpRequest(t)
	mockFoundVolumeService := mocks.NewFoundVolumesService(t)
	mockLogger := mocks.NewLogger(t)

	// Call NewBybit with mocked services
	bybits := exchange.NewBybit(
		mockUserService,
		mockUserPairsService,
		mockHttpRequestService,
		mockFoundVolumeService,
		mockLogger,
	)

	// Assert that the returned slice of exchanges is not nil and has expected length
	assert.NotNil(t, bybits)
	assert.Equal(t, 2, len(bybits)) // Assuming there are three initialization functions
}
