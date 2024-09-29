package exchange

import (
	"context"
	"cvs/internal/models"  // Importing models for domain-specific data structures
	"cvs/internal/service" // Importing service layer for user and order book services
	"cvs/internal/service/orderbook"
	"fmt"
	"io"
	"log"

	"strconv"
	"sync"
	"time"

	// Importing JSON library for marshaling and unmarshaling
	cmap "github.com/orcaman/concurrent-map/v2" // Importing concurrent map for thread-safe storage
)

const directoryPath = "internal.service.exchange." // Path for logging operations

var (
	errUnmarshal = func(dataType, exchange string) error {
		return fmt.Errorf("response unmarshal error: %s %s", exchange, dataType) // Error for unmarshalling failures
	}
)

// Exchange defines the interface for managing exchange operations.
// It includes methods for retrieving pairs, getting order books, and finding volumes.
type Exchange interface {
	StartWork()                                                         // Method to start the exchange's work
	GetAllPairsOfExchange()                                             // Method to retrieve all pairs available on the exchange
	GetOrderbookPeriodically()                                          // Method to fetch order book data periodically
	FindVolumeInOrderbookPeriodically()                                 // Method to find volume in the order book periodically
	FillPairsSubscribedStorage()                                        // Method to fill exchange pairs subscribed to pairs subscribed storage
	ExchangeName() string                                               // Method to get the name of the exchange
	AddPairToSubscribedPairs(pair string)                               // Method to add a pair to the list of subscribed pairs
	ClearSubscribedPairsStorage()                                       // Method to clear the list of subscribed pairs
	DeletePairFromSubscribedPairs(pair string)                          // Method to delete a pair from the list of subscribed pairs
	SetExchangeIntoAllExchangesStorage(exchange Exchange)               // Method to set the exchange into the AllExchanges storage
	SetEchangePairsToStorage(exchangePairsSlice []models.ExchangePairs) // Method to set the exchange pairs into the allPairsOfExchange storage
	GetOrderbookDataFromExchange(pair string)                           // Method to get the order book data from the exchange
}

// exchange is a concrete implementation of the Exchange interface.
// It holds various services and data related to an exchange.
type ExchangeData struct {
	userService         service.UserService         // User service for managing user data
	userPairsService    service.UserPairsService    // User pairs service for managing user pairs data
	foundVolumesService service.FoundVolumesService // Service for managing found volumes
	httpRequestService  service.HttpRequest         // HTTP request service for making API calls
	allExchangesStorage AllExchanges                // All exchanges storage service

	orderbookService orderbook.Orderbook // Order book service for managing order data

	allPairsOfExchange cmap.ConcurrentMap[string, models.ExchangePairs] // Concurrent map storing all pairs available on this exchange

	pairsSubscribed cmap.ConcurrentMap[string, bool] // List of pairs that are subscribed to updates

	timeBetweenRequests time.Duration // Duration between requests to the exchange API

	pairsUrlForGetRequest     string                                                                      // URL for getting pairs information from the exchange
	orderbookUrlForGetRequest string                                                                      // URL for getting order book data from the exchange
	exchangeName              string                                                                      // Name of the exchange
	pairsJsonModel            interface{}                                                                 // Model for pairs JSON response
	orderbookJsonModel        interface{}                                                                 // Model for order book JSON response
	urlFormatter              func(url, pair string) string                                               // Function to format URLs with trading pairs
	orderbookJsonParse        func(bodyBytes []byte) ([][]interface{}, [][]interface{}, error)            // Function to parse order book JSON response
	exchangePairsJsonParse    func(exchangeName string, bodyBytes []byte) ([]models.ExchangePairs, error) // Function to parse exchange pairs from JSON response
}

// InitAllExchanges initializes instances of all exchanges and starts their operations.
//
// This function creates and initializes instances of various exchanges (Binance and Bybit) by
// utilizing the provided services. It sets up goroutines to manage the retrieval of trading pairs,
// order book data, and volume finding processes for each exchange concurrently.
//
// Parameters:
//   - userService: The service for managing user data.
//   - userPairsService: The service for managing user pairs data.
//   - httpRequestService: The service for making HTTP requests.
//   - foundVolumesStorage: The service for managing found volumes data.
//   - allExchangesStorage: The storage that holds all exchanges, allowing access to exchange-related operations.
//
// This function does not return any values. It manages concurrency using goroutines and waits for
// all initialization tasks to complete before returning.
func InitAllExchanges(
	userService service.UserService,
	userPairsService service.UserPairsService,
	httpRequestService service.HttpRequest,
	foundVolumesStorage service.FoundVolumesService,
	allExchangesStorage AllExchanges,
) {
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()

		// Create instances of Binance exchanges
		binances := NewBinance(
			userService,
			userPairsService,
			httpRequestService,
			foundVolumesStorage,
			allExchangesStorage,
		)

		var binanceWg sync.WaitGroup

		for _, binance := range binances {
			binanceWg.Add(1)

			go func(binance Exchange) {
				defer binanceWg.Done()

				binance.StartWork()
			}(binance)
		}

		binanceWg.Wait() // Wait for all Binance goroutines to finish
	}()

	go func() {
		defer wg.Done()

		// Create instances of Bybit exchanges
		bybits := NewBybit(
			userService,
			userPairsService,
			httpRequestService,
			foundVolumesStorage,
			allExchangesStorage,
		)

		var bybitWg sync.WaitGroup

		for _, bybit := range bybits {
			bybitWg.Add(1)

			go func(bybit Exchange) {
				defer bybitWg.Done()

				bybit.StartWork()
			}(bybit)
		}

		bybitWg.Wait() // Wait for all Binance goroutines to finish
	}()

	wg.Wait() // Wait for the initial goroutine to finish
}

// StartWork starts the exchange's work by filling the pairs subscribed storage, retrieving all
// pairs available on the exchange, and starting the periodic fetching of order book data and
// finding volume in the order book. This method calls the following methods in order: FillPairsSubscribedStorage,
// GetAllPairsOfExchange, FindVolumeInOrderbookPeriodically, and GetOrderbookPeriodically.
func (e *ExchangeData) StartWork() {
	e.SetExchangeIntoAllExchangesStorage(e) // Add each exchange instance to the AllExchanges storage

	e.FillPairsSubscribedStorage()        // Fill pairs subscribed storage
	e.GetAllPairsOfExchange()             // Retrieve all pairs available on exchange instance
	e.FindVolumeInOrderbookPeriodically() // Start finding volume in the order book periodically
	e.GetOrderbookPeriodically()          // Start fetching order book data periodically
}

// GetAllPairsOfExchange retrieves all trading pairs available on the exchange.
//
// This method makes a GET request to the exchange's API to fetch the trading pairs
// information. It reads the response body, parses the JSON data into a slice of
// ExchangePairs, and stores this data in the exchange's storage.
//
// The method performs the following steps:
// 1. Sends an HTTP GET request to the URL specified by pairsUrlForGetRequest.
// 2. Reads the response body into bytes.
// 3. Parses the JSON response into a slice of ExchangePairs.
// 4. Logs any errors encountered during parsing.
// 5. Calls SetEchangePairsToStorage to store the retrieved pairs in storage.
//
// This method does not return any values and does not produce errors directly.
// However, it logs any errors encountered during the HTTP request or JSON parsing.
//
// Example usage:
//
//	e.GetAllPairsOfExchange()
func (e *ExchangeData) GetAllPairsOfExchange() {
	const op = directoryPath + "GetAllPairsOfExchange"

	resp, _ := e.httpRequestService.Get(e.pairsUrlForGetRequest) // Make a GET request to retrieve pairs information
	defer resp.Body.Close()                                      // Ensure response body is closed after reading

	bodyBytes, _ := io.ReadAll(resp.Body) // Read response body into bytes

	exchangePairsSlice, err := e.exchangePairsJsonParse(e.exchangeName, bodyBytes) // Parse JSON response into exchange pairs slice
	if err != nil {
		log.Println(e.exchangeName, op, ": ", err) // Log any errors encountered during parsing
	}

	e.SetEchangePairsToStorage(exchangePairsSlice) // Store the retrieved pairs in storage
}

// FillPairsSubscribedStorage retrieves and stores the subscribed trading pairs for the exchange.
//
// This method fetches the pairs associated with the current exchange from the userPairsService.
// It uses the exchange's name to get the relevant pairs and stores them in the
// pairsSubscribed field of the exchange struct.
//
// If an error occurs while retrieving the pairs, it logs the error with context about the operation.
//
// Example usage:
//
//	e.FillPairsSubscribedStorage()
func (e *ExchangeData) FillPairsSubscribedStorage() {
	const op = directoryPath + "FillPairsSubscribedStorage"

	pairs, err := e.userPairsService.GetPairsByExchange(context.Background(), e.exchangeName)
	if err != nil {
		log.Println(e.exchangeName, op, ": ", err) // Log any errors encountered during retrieval
	}

	for _, pair := range pairs {
		e.pairsSubscribed.Set(pair, true) // Store each pair in the exchange's pairsSubscribed field
	}
}

// GetOrderbookDataFromExchange retrieves order book data for a specific trading pair from the exchange.
//
// This method makes a GET request to the exchange's API to fetch the order book data
// for the specified trading pair. It reads the response body, parses the JSON data
// into asks and bids, and updates the order book service with this data.
//
// Parameters:
//   - pair: A string representing the trading pair for which to retrieve order book data.
//
// This method does not return any values and does not produce errors directly.
// However, it logs any errors encountered during the HTTP request or JSON parsing.
// If an error occurs during parsing, it will be logged with the exchange name and operation context.
//
// Example usage:
//
//	e.GetOrderbookDataFromExchange("BTC/USD")
func (e *ExchangeData) GetOrderbookDataFromExchange(pair string) {
	const op = directoryPath + "GetOrderbookDataFromExchange"

	// Make a GET request to retrieve order book data using formatted URL
	resp, err := e.httpRequestService.Get(e.urlFormatter(e.orderbookUrlForGetRequest, pair))
	if err != nil {
		// Log any errors encountered during the GET request
		log.Println(e.exchangeName, op, ": ", err)
	}

	defer resp.Body.Close() // Ensure response body is closed after reading

	// Read the response body into bytes
	bodyBytes, _ := io.ReadAll(resp.Body)

	// Parse JSON response into asks and bids slices
	asks, bids, err := e.orderbookJsonParse(bodyBytes)
	if len(asks) == 0 || len(bids) == 0 || err != nil {
		// Log any errors encountered during JSON parsing
		log.Println(e.exchangeName, op, ": ", "empty asks or bids or error parsing JSON. Error: ", err)
	}

	// Update or insert order book data into the order book service
	e.orderbookService.Upsert(pair, asks, bids) // Update or insert order book data into the order book service
}

// GetOrderbookPeriodically fetches order book data from the exchange for subscribed pairs at regular intervals.
//
// This method runs as a goroutine and continuously checks for subscribed pairs.
// If there are subscribed pairs, it iterates over each pair and retrieves the order book data
// from the exchange using the GetOrderbookDataFromExchange method.
// It sleeps for timeBetweenRequests variable  value milliseconds between requests to avoid hitting rate limits imposed by the exchange API.
// If there are no subscribed pairs, it waits for 1 second before checking again.
//
// This method will run indefinitely until the application is terminated or the goroutine is stopped.
//
// Possible Errors:
//   - Errors may occur during the retrieval of order book data, but these errors are logged
//     and do not interrupt the execution of this method.
func (e *ExchangeData) GetOrderbookPeriodically() {
	go func() {
		for {
			pairsSubscribed := e.pairsSubscribed.Keys() // Get all subscribed pairs keys

			if len(pairsSubscribed) != 0 { // Check if there are any subscribed pairs
				for _, pair := range pairsSubscribed { // Iterate over each subscribed pair
					e.GetOrderbookDataFromExchange(pair) // Fetch order book data from the exchange

					time.Sleep(e.timeBetweenRequests) // Sleep briefly between requests to avoid rate limiting
				}
			}

			time.Sleep(time.Second) // Sleep before checking again
		}
	}()
}

// FindVolumeInOrderbookPeriodically searches for trading volumes in the order book
// for subscribed pairs at regular intervals.
//
// This method runs as a goroutine and continuously checks for subscribed pairs.
// If there are no subscribed pairs, it waits for one second before checking again.
// For each subscribed pair, it retrieves the user IDs from memory and processes
// each user's settings to search for volumes in the order book using the specified
// exact values. The found volumes are then upserted into the found volumes service.
//
// The method utilizes goroutines to handle concurrent processing of user settings
// and volume searches, ensuring that multiple users can be processed simultaneously.
//
// Note: This method will run indefinitely until the application is terminated or
// the goroutine is stopped.
//
// Possible Errors:
//   - Errors may occur during the retrieval of user pairs or while searching for volumes,
//     but these errors are logged and do not interrupt the execution of this method.
func (e *ExchangeData) FindVolumeInOrderbookPeriodically() {
	go func() {
		for {
			pairsSubscribed := e.pairsSubscribed.Keys() // Get all subscribed pairs keys

			if len(pairsSubscribed) != 0 { // Check if there are any subscribed pairs
				for _, pair := range pairsSubscribed { // Iterate over each subscribed pair
					var wg sync.WaitGroup // WaitGroup to manage goroutines

					for _, userID := range e.userService.GetUsersIdFromMemory().Keys() {
						wg.Add(1) // Increment WaitGroup counter

						go func(userID string) { // Start a new goroutine for each user ID
							defer wg.Done() // Decrement counter when done

							userIdInt, _ := strconv.Atoi(userID) // Convert user ID to int

							userSettings, _ := e.userPairsService.GetAllUserPairs(context.Background(), userIdInt)

							for _, pairSettings := range userSettings { // Iterate over each user's pair settings
								foundVolumes := e.orderbookService.SearchVolume(pair, e.exchangeName, pairSettings.ExactValue) // Search for volumes

								for _, volume := range foundVolumes { // Iterate over found volumes
									e.foundVolumesService.UpsertFoundVolume(pairSettings, volume) // Upsert volume into service
								}
							}
						}(userID)

						time.Sleep(100 * time.Millisecond) // Sleep briefly between processing users
					}

					wg.Wait() // Wait for all goroutines to finish before proceeding to the next pair
				}
			}

			time.Sleep(time.Second)
		}
	}()
}

// SetEchangePairsToStorage stores all pairs of an exchange into its storage.
//
// This method takes a slice of ExchangePairs and iterates over each pair.
// For each pair, it adds the pair data to the exchange's storage using a
// concurrent map for thread-safe operations.
//
// Parameters:
//   - exchangePairsSlice: A slice of models.ExchangePairs containing the pairs
//     to be stored in the exchange's storage.
//
// This method does not return any values and does not produce errors.
// It assumes that the provided slice is valid and that the concurrent map
// is initialized properly. If the slice is empty, no operations are performed.
func (e *ExchangeData) SetEchangePairsToStorage(exchangePairsSlice []models.ExchangePairs) {
	for _, pairData := range exchangePairsSlice { // Iterate over each pair data
		e.allPairsOfExchange.Set(pairData.Pair, pairData) // Store each pair in the concurrent map
	}
}

// SetExchangeIntoAllExchangesStorage adds an exchange to the AllExchanges storage.
//
// This method adds the provided exchange to the allExchangesStorage concurrent map.
// It does not return any values and does not produce errors. It assumes that the
// provided exchange is valid and that the concurrent map is initialized properly.
func (e *ExchangeData) SetExchangeIntoAllExchangesStorage(exchange Exchange) {
	e.allExchangesStorage.Add(exchange)
}

// ExchangeName returns the name of the exchange.
func (e *ExchangeData) ExchangeName() string {
	return e.exchangeName //
}

// AddPairToSubscribedPairs adds a trading pair to the set of subscribed pairs for this exchange.
// It takes a string parameter representing the pair to be added and sets the value in the concurrent map to true.
// This method does not return any values and does not produce errors. If the pair is already subscribed, this method has no effect.
func (e *ExchangeData) AddPairToSubscribedPairs(pair string) {
	e.pairsSubscribed.Set(pair, true)
}

func (e *ExchangeData) ClearSubscribedPairsStorage() {
	e.pairsSubscribed.Clear()
}

// DeletePairFromSubscribedPairs deletes a trading pair from the set of subscribed pairs for this exchange.
// It takes a string parameter representing the pair to be deleted and sets the value in the concurrent map to false.
// This method does not return any values and does not produce errors. If the pair is not subscribed, this method has no effect.
func (e *ExchangeData) DeletePairFromSubscribedPairs(pair string) {
	e.pairsSubscribed.Remove(pair)
}
