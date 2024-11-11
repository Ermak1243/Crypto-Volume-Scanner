package tests

import (
	"bytes"
	"cvs/internal/mocks"
	"cvs/internal/service/exchange"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInitAllExchanges(t *testing.T) {
	t.Parallel() // Allows this test to run in parallel with other tests

	// Create mocks for services
	mockUserService := mocks.NewUserService(t)
	mockUserPairsService := mocks.NewUserPairsService(t)
	mockHttpRequestService := mocks.NewHttpRequest(t)
	mockFoundVolumeService := mocks.NewFoundVolumesService(t)
	mockLogger := mocks.NewLogger(t)
	allExchangesStorage := exchange.NewAllExchangesService(mockLogger)

	mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockHttpRequestService.On("Get", mock.Anything).Return(http.Response{Body: io.NopCloser(bytes.NewReader([]byte("test")))}, nil)
	mockUserPairsService.On("GetPairsByExchange", mock.Anything, mock.Anything).Return(nil, nil)

	allExchanges := exchange.InitAllExchanges(
		mockUserService,
		mockUserPairsService,
		mockHttpRequestService,
		mockFoundVolumeService,
		allExchangesStorage,
		mockLogger,
	)

	assert.EqualValues(t, 5, len(allExchanges.All()))
}
