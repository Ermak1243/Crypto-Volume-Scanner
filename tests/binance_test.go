package tests

import (
	"main/internal/mocks"
	"main/internal/service/exchange"

	"testing" // Importing the testing package for writing tests

	"github.com/stretchr/testify/assert" // Importing assert from testify for easier assertions
)

// TestNewBinance tests the NewBinance function
func TestNewBinance(t *testing.T) {
	t.Parallel() // Allows this test to run in parallel with other tests

	// Create mocks for services
	mockUserService := mocks.NewUserService(t)
	mockUserPairsService := mocks.NewUserPairsService(t)
	mockHttpRequestService := mocks.NewHttpRequest(t)
	mockFoundVolumeService := mocks.NewFoundVolumesService(t)
	mockAllExchangesStorage := mocks.NewAllExchanges(t)

	// Call NewBinance with mocked services
	binances := exchange.NewBinance(
		mockUserService,
		mockUserPairsService,
		mockHttpRequestService,
		mockFoundVolumeService,
		mockAllExchangesStorage,
	)

	// Assert that the returned slice of exchanges is not nil and has expected length
	assert.NotNil(t, binances)
	assert.Equal(t, 3, len(binances)) // Assuming there are three initialization functions
}
