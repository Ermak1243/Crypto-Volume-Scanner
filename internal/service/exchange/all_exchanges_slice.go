package exchange

import (
	"fmt"

	cmap "github.com/orcaman/concurrent-map/v2"
)

// AllExchanges defines the interface for managing multiple exchange instances.
// It includes methods for adding and retrieving exchanges.
type AllExchanges interface {
	Add(exchange Exchange)            // Method to add a new exchange to the storage
	Get(exchangeName string) Exchange // Method to retrieve an exchange by its name
	All() []Exchange                  // Method to retrieve all exchanges stored in the storage
}

// allExchanges is a concrete implementation of the AllExchanges interface.
// It holds a concurrent map to store exchange instances.
type allExchanges struct {
	exchanges cmap.ConcurrentMap[string, Exchange] // Concurrent map storing exchanges by their names
}

// NewAllExchangesService creates a new instance of allExchanges.
// It initializes the concurrent map for storing exchange instances.
func NewAllExchangesService() AllExchanges {
	return &allExchanges{
		exchanges: cmap.New[Exchange](), // Initialize the concurrent map for exchanges
	}
}

// Add adds a new exchange to the storage.
// It stores the exchange in the concurrent map using its name as the key.
func (ae *allExchanges) Add(exchange Exchange) {
	ae.exchanges.Set(exchange.ExchangeName(), exchange) // Store the exchange in the map
}

// Get retrieves an exchange by its name from the storage.
// If the exchange does not exist, it logs a message and returns a nil value.
func (ae *allExchanges) Get(exchangeName string) Exchange {
	exchange, exists := ae.exchanges.Get(exchangeName) // Attempt to retrieve the exchange from the map
	if !exists {
		fmt.Printf("exchange with name %s does not exist in AllExchanges storage", exchangeName) // Log message if not found
	}

	return exchange // Return the retrieved exchange (or nil if not found)
}

// All retrieves all exchanges stored in the concurrent map.
func (ae *allExchanges) All() []Exchange {
	var exchanges []Exchange // Slice to hold instances of Exchange

	// Iterate over all exchanges in the concurrent map and append them to the list
	for _, exchange := range ae.exchanges.Items() {
		exchanges = append(exchanges, exchange) // Add each exchange to the slice
	}

	return exchanges // Return the list of exchanges
}
