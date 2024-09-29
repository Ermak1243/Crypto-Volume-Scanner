package tests

import (
	"cvs/internal/service/orderbook"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert" // Import assert package for better assertions
)

// TestOrderbook_Upsert tests the Upsert function of the Orderbook.
func TestOrderbook_Upsert(t *testing.T) {
	t.Parallel() // Run tests in parallel for efficiency

	tests := []struct {
		name      string          // Name of the test case
		pair      string          // Trading pair to be tested
		asks      [][]interface{} // Asks data to be inserted
		bids      [][]interface{} // Bids data to be inserted
		expectErr bool            // Expected outcome: true if an error is expected
	}{
		{
			name: "Valid upsert", // Test case for valid upsert operation
			pair: "BTC/USD",
			asks: [][]interface{}{
				{"50000", "1"}, // Ask price and volume
				{"51000", "2"},
			},
			bids: [][]interface{}{
				{"49000", "1"}, // Bid price and volume
				{"48000", "2"},
			},
			expectErr: false, // No error expected for valid input
		},
		{
			name: "Empty asks", // Test case for empty asks
			pair: "BTC/USD",
			asks: [][]interface{}{}, // No asks provided
			bids: [][]interface{}{
				{"49000", "1"},
			},
			expectErr: false, // No error expected; empty asks are valid
		},
		{
			name: "Empty bids", // Test case for empty bids
			pair: "BTC/USD",
			asks: [][]interface{}{
				{"50000", "1"},
			},
			bids:      [][]interface{}{}, // No bids provided
			expectErr: false,             // No error expected; empty bids are valid
		},
	}

	for _, tc := range tests {
		tc := tc // Capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel() // Run this test case in parallel

			ob := orderbook.NewOrderbook()       // Create a new orderbook instance
			ob.Upsert(tc.pair, tc.asks, tc.bids) // Perform the upsert operation

			// Check if asks and bids are set correctly
			asks := ob.Asks(tc.pair) // Retrieve asks for the trading pair
			bids := ob.Bids(tc.pair) // Retrieve bids for the trading pair

			assert.Equal(t, len(tc.asks), len(asks), "Expected %d asks, got %d", len(tc.asks), len(asks)) // Validate number of asks
			assert.Equal(t, len(tc.bids), len(bids), "Expected %d bids, got %d", len(tc.bids), len(bids)) // Validate number of bids
		})
	}
}

// TestOrderbook_Asks tests the Asks function of the Orderbook.
func TestOrderbook_Asks(t *testing.T) {
	t.Parallel() // Run tests in parallel for efficiency

	ob := orderbook.NewOrderbook()                                                         // Create a new orderbook instance
	ob.Upsert("BTC/USD", [][]interface{}{{"50000", "1"}}, [][]interface{}{{"49000", "1"}}) // Insert test data

	asks := ob.Asks("BTC/USD") // Retrieve asks for the trading pair

	assert.Equal(t, 1, len(asks), "Expected 1 ask, got %d", len(asks))                                      // Validate number of asks retrieved
	assert.Equal(t, "1", asks["50000"], "Expected ask price 50000 to have volume 1, got %s", asks["50000"]) // Validate ask volume
}

// TestOrderbook_Bids tests the Bids function of the Orderbook.
func TestOrderbook_Bids(t *testing.T) {
	t.Parallel() // Run tests in parallel for efficiency

	ob := orderbook.NewOrderbook()                                                         // Create a new orderbook instance
	ob.Upsert("BTC/USD", [][]interface{}{{"50000", "1"}}, [][]interface{}{{"49000", "1"}}) // Insert test data

	bids := ob.Bids("BTC/USD") // Retrieve bids for the trading pair

	assert.Equal(t, 1, len(bids), "Expected 1 bid, got %d", len(bids))                                      // Validate number of bids retrieved
	assert.Equal(t, "1", bids["49000"], "Expected bid price 49000 to have volume 1, got %s", bids["49000"]) // Validate bid volume
}

// TestOrderbook_SearchVolume tests the SearchVolume function of the Orderbook.
func TestOrderbook_SearchVolume(t *testing.T) {
	t.Parallel() // Run tests in parallel for efficiency

	ob := orderbook.NewOrderbook()                                                         // Create a new orderbook instance
	ob.Upsert("BTC/USD", [][]interface{}{{"50000", "1"}}, [][]interface{}{{"49000", "1"}}) // Insert test data

	volumes := ob.SearchVolume("BTC/USD", "binance", 1) // Search volumes based on criteria

	assert.Equal(t, 2, len(volumes), "Expected 2 volumes, got %d", len(volumes)) // Validate total volumes retrieved
}

// TestOrderbook_ConcurrentAccess tests concurrent access to the Orderbook.
func TestOrderbook_ConcurrentAccess(t *testing.T) {
	t.Parallel() // Run tests in parallel for efficiency

	ob := orderbook.NewOrderbook() // Create a new orderbook instance

	var wg sync.WaitGroup

	// Simulate concurrent upsert operations
	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			ob.Upsert("BTC/USD", [][]interface{}{{"50000", "1"}}, [][]interface{}{{"49000", "1"}})
		}()
	}
	wg.Wait()

	// Check if asks and bids are set correctly after concurrent operations
	asks := ob.Asks("BTC/USD")
	bids := ob.Bids("BTC/USD")

	assert.Greater(t, len(asks), 0, "Expected at least 1 ask, got %d", len(asks))
	assert.Greater(t, len(bids), 0, "Expected at least 1 bid, got %d", len(bids))
}
