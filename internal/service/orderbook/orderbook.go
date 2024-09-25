package orderbook

import (
	"fmt"
	"main/internal/domain/models"
	"sort"
	"sync"
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/spf13/cast"
)

// Orderbook defines the interface for managing an order book.
// It includes methods for retrieving asks and bids, upserting data, and searching for volumes.
type Orderbook interface {
	Asks(pair string) map[string]interface{}                                 // Method to retrieve all ask orders for a given pair
	Bids(pair string) map[string]interface{}                                 // Method to retrieve all bid orders for a given pair
	Upsert(pair string, asks, bids [][]interface{})                          // Method to update or insert ask and bid orders
	SearchVolume(pair, exchange string, search float64) []models.FoundVolume // Method to search for volumes based on a specified value
}

// orderbook is a concrete implementation of the Orderbook interface.
// It holds a concurrent map to store order book data by pairs.
type orderbook struct {
	cmap.ConcurrentMap[string, orderbookData] // Concurrent map storing order book data by pair
}

// orderbookData holds the details of an order book entry.
// It includes the pair, asks, bids, and sorted lists of found volumes.
type orderbookData struct {
	Pair               string                                  // The trading pair (e.g., "BTC/USD")
	asks               cmap.ConcurrentMap[string, interface{}] // Concurrent map for ask orders
	bids               cmap.ConcurrentMap[string, interface{}] // Concurrent map for bid orders
	asksSortedByVolume []models.FoundVolume                    // Sorted list of asks by volume
	bidsSortedByVolume []models.FoundVolume                    // Sorted list of bids by volume
	asksSortedByPrice  []models.FoundVolume                    // Sorted list of asks by price
	bidsSortedByPrice  []models.FoundVolume                    // Sorted list of bids by price
}

// sortedSlice holds two slices of FoundVolume sorted by volume and price.
type sortedSlice struct {
	ByVolume []models.FoundVolume // Slice of volumes sorted by volume
	ByPrice  []models.FoundVolume // Slice of volumes sorted by price
}

// NewOrderbook creates a new instance of orderbook.
// It initializes the concurrent map for storing order book data.
func NewOrderbook() Orderbook {
	level2Data := &orderbook{
		cmap.New[orderbookData](), // Initialize the concurrent map for order book data
	}

	return level2Data // Return the new orderbook instance
}

// Asks retrieves all ask orders for a given trading pair.
func (o *orderbook) Asks(pair string) map[string]interface{} {
	orderbook, _ := o.Get(pair) // Get the order book data for the specified pair

	return orderbook.asks.Items() // Return all ask orders as a map
}

// Bids retrieves all bid orders for a given trading pair.
func (o *orderbook) Bids(pair string) map[string]interface{} {
	orderbook, _ := o.Get(pair) // Get the order book data for the specified pair

	return orderbook.bids.Items() // Return all bid orders as a map
}

// Upsert updates or inserts ask and bid orders into the order book.
// It organizes the data in a nested concurrent map structure based on user ID, pair, exchange, and side.
func (o *orderbook) Upsert(pair string, asks, bids [][]interface{}) {
	var wg sync.WaitGroup

	o.Remove(pair) // Remove any existing data for the specified pair

	level2Data := orderbookData{
		Pair: pair,
		asks: cmap.New[interface{}](), // Initialize concurrent map for asks
		bids: cmap.New[interface{}](), // Initialize concurrent map for bids
	}

	wg.Add(2) // Prepare to wait for two goroutines

	go func() {
		defer wg.Done() // Decrement WaitGroup counter when done

		buffAsks := cmap.New[interface{}]() // Temporary concurrent map for incoming asks

		for _, val := range asks { // Iterate over incoming asks
			buffAsks.Set(fmt.Sprintf("%v", val[0]), val[1]) // Store each ask in the temporary map
		}

		level2Data.asks = buffAsks                                             // Assign temporary asks to level2Data
		level2Data.asksSortedByPrice = sortHashMap(buffAsks.Items()).ByPrice   // Sort asks by price
		level2Data.asksSortedByVolume = sortHashMap(buffAsks.Items()).ByVolume // Sort asks by volume
	}()
	go func() {
		defer wg.Done() // Decrement WaitGroup counter when done

		buffBids := cmap.New[interface{}]() // Temporary concurrent map for incoming bids

		for _, val := range bids { // Iterate over incoming bids
			buffBids.Set(fmt.Sprintf("%v", val[0]), val[1]) // Store each bid in the temporary map
		}

		level2Data.bids = buffBids                                             // Assign temporary bids to level2Data
		level2Data.bidsSortedByPrice = sortHashMap(buffBids.Items()).ByPrice   // Sort bids by price
		level2Data.bidsSortedByVolume = sortHashMap(buffBids.Items()).ByVolume // Sort bids by volume
	}()

	wg.Wait() // Wait for both goroutines to finish

	o.Set(pair, level2Data) // Store the updated level2Data in the main order book structure
}

// SearchVolume retrieves found volumes based on a specified search value.
// It searches both asks and bids concurrently.
func (o *orderbook) SearchVolume(pair, exchange string, search float64) []models.FoundVolume {
	var volumes []models.FoundVolume // Slice to hold found volumes results
	level2Data, exist := o.Get(pair) // Get the order book data for the specified pair
	if !exist {                      // Check if data exists for the pair
		return volumes // Return empty slice if not found
	}

	asksSlice := level2Data.asksSortedByVolume // Get sorted asks by volume from level2Data
	bidsSlice := level2Data.bidsSortedByVolume // Get sorted bids by volume from level2Data

	var wg sync.WaitGroup // WaitGroup to synchronize goroutines

	wg.Add(2) // Prepare to wait for two goroutines

	go func() {
		defer wg.Done() // Decrement WaitGroup counter when done

		foundVolumeData := binarySearch(pair, asksSlice, search) // Perform binary search on asks slice
		if foundVolumeData.Price != 0 {                          // Check if found volume has a valid price
			percentDistance := (foundVolumeData.Price - level2Data.asksSortedByPrice[0].Price) / foundVolumeData.Price * 100 // Calculate percentage distance from first ask price

			foundVolumeData.Difference = percentDistance // Store calculated difference in found volume data
			foundVolumeData.VolumeTimeFound = time.Now()
		}

		foundVolumeData.Side = "asks" // Set found volume side to "ask"
		foundVolumeData.Pair = pair
		foundVolumeData.Exchange = exchange

		volumes = append(volumes, foundVolumeData) // Append found volume data to results
	}()
	go func() {
		defer wg.Done() // Decrement WaitGroup counter when done

		foundVolumeData := binarySearch(pair, bidsSlice, search) // Perform binary search on bids slice
		if foundVolumeData.Price != 0 {                          // Check if found volume has a valid price
			percentDistance := (level2Data.bidsSortedByPrice[len(level2Data.bidsSortedByPrice)-1].Price - foundVolumeData.Price) / level2Data.bidsSortedByPrice[len(level2Data.bidsSortedByPrice)-1].Price * 100 // Calculate percentage distance from last bid price

			foundVolumeData.Difference = percentDistance // Store calculated difference in found volume data
			foundVolumeData.VolumeTimeFound = time.Now()
		}

		foundVolumeData.Side = "bids" // Set found volume side to "bid"
		foundVolumeData.Pair = pair
		foundVolumeData.Exchange = exchange

		volumes = append(volumes, foundVolumeData) // Append found volume data to results
	}()

	wg.Wait() // Wait for both goroutines to finish

	return volumes // Return all found volumes retrieved
}

// sortHashMap sorts a hashmap of interface values into slices sorted by volume and price.
// It returns a sortedSlice containing both sorted slices.
//
// Parameters:
//   - hashmap: A map where the key is a string (representing price) and the value is an interface{} (representing volume).
//
// Returns:
//   - A sortedSlice containing two slices: one sorted by volume and another sorted by price.
func sortHashMap(hashmap map[string]interface{}) sortedSlice {
	sortedByVolume := make([]models.FoundVolume, 0, len(hashmap)) // Slice to hold volumes sorted by volume
	sortedByPrice := make([]models.FoundVolume, 0, len(hashmap))  // Slice to hold volumes sorted by price

	index := 0 // Index for tracking the position in the slices

	// Iterate over each key in the hashmap to populate the sorted slices
	for k := range hashmap {
		// Append a new FoundVolume to the sortedByVolume slice
		sortedByVolume = append(sortedByVolume, models.FoundVolume{
			Index:  index,
			Price:  cast.ToFloat64(k),          // Convert key (price) from string to float64
			Volume: cast.ToFloat64(hashmap[k]), // Convert value (volume) from interface{} to float64
		})

		// Append a new FoundVolume to the sortedByPrice slice
		sortedByPrice = append(sortedByPrice, models.FoundVolume{
			Index:  index,
			Price:  cast.ToFloat64(k),          // Convert key (price) from string to float64
			Volume: cast.ToFloat64(hashmap[k]), // Convert value (volume) from interface{} to float64
		})

		index++ // Increment index for the next entry
	}

	var wg sync.WaitGroup // WaitGroup to synchronize goroutines

	wg.Add(2) // Add two goroutines to the WaitGroup

	// Goroutine for sorting by volume
	go func() {
		defer wg.Done() // Decrement WaitGroup counter when done

		sort.SliceStable(sortedByVolume, func(i, j int) bool { // Sort the slice by volume using a stable sort
			return sortedByVolume[i].Volume < sortedByVolume[j].Volume // Compare volumes for sorting order
		})
	}()

	// Goroutine for sorting by price
	go func() {
		defer wg.Done() // Decrement WaitGroup counter when done

		sort.SliceStable(sortedByPrice, func(i, j int) bool { // Sort the slice by price using a stable sort
			return sortedByPrice[i].Price < sortedByPrice[j].Price // Compare prices for sorting order
		})
	}()

	wg.Wait() // Wait for both sorting goroutines to finish

	return sortedSlice{ // Return a struct containing both sorted slices
		ByVolume: sortedByVolume,
		ByPrice:  sortedByPrice,
	}
}

// binarySearch performs a binary search on a slice of FoundVolumes to find a volume matching the search criteria.
// It returns the FoundVolume that matches or is closest to the specified search value.
//
// Parameters:
//   - pair: The trading pair being searched (not used in this implementation but could be relevant for logging or context).
//   - slice: A slice of FoundVolume objects sorted by volume.
//   - search: The volume value to search for in the slice.
//
// Returns:
//   - A FoundVolume object that matches the search criteria.
func binarySearch(pair string, slice []models.FoundVolume, search float64) models.FoundVolume {
	mid := len(slice) / 2                  // Calculate the midpoint index of the slice.
	var foundVolumeData models.FoundVolume // Variable to hold the found volume data.

	switch { // Determine which case to execute based on the length of the slice and the value at the midpoint.
	case len(slice) == 0: // Base case: If the slice is empty,
		foundVolumeData = models.FoundVolume{} // Return an empty FoundVolume.
	case slice[mid].Volume >= search: // If the volume at the midpoint is greater than or equal to the search value,
		foundVolumeData = slice[mid] // Set foundVolumeData to the midpoint volume (potential match).
	case slice[mid].Volume < search: // If the volume at the midpoint is less than the search value,
		// Recursively search in the right half of the slice (elements after mid).
		foundVolumeData = binarySearch(pair, slice[mid+1:], search)
	default: // This case handles any unexpected scenarios (though it should not be reached).
		foundVolumeData = slice[mid] // Fallback to returning the midpoint volume.
	}

	return foundVolumeData // Return the found volume data (or an empty FoundVolume if not found).
}
