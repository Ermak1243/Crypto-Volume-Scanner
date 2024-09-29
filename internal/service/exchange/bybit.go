package exchange

import (
	"strings"
	"time"

	"cvs/internal/models"
	"cvs/internal/service"
	"cvs/internal/service/orderbook"

	"github.com/goccy/go-json"
	cmap "github.com/orcaman/concurrent-map/v2"
)

// Overall data for all sections of the Bybit exchange
var (
	bybitTimeBetweenRequests = 3 * time.Second                     // Time interval between requests to the Bybit API
	bybitPairsJsonModel      = models.BybitPairsJSONResponse{}     // Model for Bybit pairs JSON response
	bybitOrderbookJsonModel  = models.BybitOrderbookJSONResponse{} // Model for Bybit order book JSON response
	bybitOrderbookService    = orderbook.NewOrderbook()            // Instance of the order book service for managing order data

	// Function to parse order book JSON response from Bybit
	bybitOrderbookJsonParse = func(bodyBytes []byte) ([][]interface{}, [][]interface{}, error) {
		var model models.BybitOrderbookJSONResponse

		// Unmarshal the response body into jsonData to inspect the response
		err := json.Unmarshal(bodyBytes, &model)

		return model.Result.Asks, model.Result.Bids, err
	}

	// Function to format Bybit API URLs with the trading pair
	bybitUrlFormatter = func(url, pair string) string {
		pairFormatted := strings.Replace(pair, "/", "", -1)                 // Remove slashes from the pair string
		replacer := strings.NewReplacer("symbol=", "symbol="+pairFormatted) // Replace "symbol=" in the URL with the formatted pair

		return replacer.Replace(url) // Return the formatted URL
	}

	// Function to parse exchange pairs from Bybit API response
	bybitExchangePairsJsonParse = func(exchangeName string, bodyBytes []byte) ([]models.ExchangePairs, error) {
		var model models.BybitPairsJSONResponse

		// Unmarshal the response body into jsonData to inspect the response
		err := json.Unmarshal(bodyBytes, &model)
		if err != nil {
			return []models.ExchangePairs{}, errUnmarshal("exchange pairs", exchangeName) // Return error if unmarshalling fails
		}

		var exchangePairsSlice []models.ExchangePairs // Slice to hold parsed exchange pairs

		for i := 0; i < len(model.Result.List); i++ { // Iterate over all symbols in pairs data
			exchangePairsSlice = append(exchangePairsSlice, models.ExchangePairs{
				Pair:     model.Result.List[i].BaseCoin + "/" + model.Result.List[i].BaseCoin, // Construct pair string
				Exchange: exchangeName,                                                        // Set exchange name
			})
		}

		return exchangePairsSlice, nil // Return the slice of exchange pairs
	}
)

// NewBybit initializes instances of different Bybit exchanges.
//
// This function creates and returns a slice of Exchange instances for various Bybit exchanges,
// including Spot and Futures exchanges. It uses the provided user service, user pairs service,
// HTTP request service, and found volume service to set up each exchange's data.
//
// Parameters:
//   - userService: The service for managing user data.
//   - userPairsService: The service for managing user pairs data.
//   - httpRequestService: The service for making HTTP requests.
//   - foundVolumeService: The service for managing found volumes.
//
// Returns:
//   - []Exchange: A slice containing instances of different Bybit exchanges.
func NewBybit(
	userService service.UserService,
	userPairsService service.UserPairsService,
	httpRequestService service.HttpRequest,
	foundVolumeService service.FoundVolumesService,
	allExchangesStorage AllExchanges,
) []Exchange {
	var bybits []Exchange // Slice to hold instances of different Bybit exchanges
	initFunctions := []func(exchangesData *ExchangeData) *ExchangeData{
		setBybitSpotData,
		setBybitFuturesData,
	}

	for _, function := range initFunctions {
		exchangeData := setBybitOverallData(
			userService,
			userPairsService,
			httpRequestService,
			foundVolumeService,
			allExchangesStorage,
		)

		bybits = append(bybits, function(exchangeData))
	}

	return bybits // Return the slice of Bybit exchanges
}

// setBybitOverallData initializes and sets up overall data for all Bybit exchanges.
//
// This function creates an instance of the exchange struct and populates it with the necessary services,
// models, and configurations required for interacting with Bybit exchanges. It prepares the exchange
// with settings for handling trading pairs, order books, and request formatting.
//
// Parameters:
//   - userService: The service for managing user data.
//   - userPairsService: The service for managing user pairs data.
//   - httpRequestService: The service for making HTTP requests.
//   - foundVolumeService: The service for managing found volumes.
//
// Returns:
//   - *exchange: A pointer to the initialized exchange struct, ready for use in API interactions.
func setBybitOverallData(
	userService service.UserService,
	userPairsService service.UserPairsService,
	httpRequestService service.HttpRequest,
	foundVolumeService service.FoundVolumesService,
	allExchangesStorage AllExchanges,
) *ExchangeData {
	bybitExchangesData := ExchangeData{
		userService:            userService,
		userPairsService:       userPairsService,
		httpRequestService:     httpRequestService,
		foundVolumesService:    foundVolumeService,
		allExchangesStorage:    allExchangesStorage,
		pairsJsonModel:         bybitPairsJsonModel,              // Set pairs JSON model for exchanges
		orderbookJsonModel:     bybitOrderbookJsonModel,          // Set orderbook JSON model for exchanges
		urlFormatter:           bybitUrlFormatter,                // Set URL formatter function for exchanges
		timeBetweenRequests:    bybitTimeBetweenRequests,         // Set time between requests for exchanges
		orderbookService:       bybitOrderbookService,            // Assign order book service instance to exchanges data
		pairsSubscribed:        cmap.New[bool](),                 // Initialize subscribed pairs list as empty
		allPairsOfExchange:     cmap.New[models.ExchangePairs](), // Initialize concurrent map for all pairs of the exchange
		orderbookJsonParse:     bybitOrderbookJsonParse,          // Set order book JSON parsing function for exchanges
		exchangePairsJsonParse: bybitExchangePairsJsonParse,      // Set exchange pairs JSON parsing function for exchanges
	}

	return &bybitExchangesData
}

// setBybitSpotData sets up data specific to the Bybit Spot exchange.
//
// This function configures the exchange struct with settings specific to the Bybit Spot exchange,
// including URLs for API calls and initializing necessary fields.
//
// Parameters:
//   - exchangesData: A pointer to the exchange struct to be configured.
//
// Returns:
//   - *exchange: A pointer to the updated exchange struct.
func setBybitSpotData(exchangesData *ExchangeData) *ExchangeData {
	const category = "spot"

	exchangesData.exchangeName = "bybit_spot"                                                                                          // Set the name of the exchange to "bybitSpot"
	exchangesData.pairsUrlForGetRequest = "https://api.bytick.com/v5/market/instruments-info?category=" + category                     // URL for getting pairs information
	exchangesData.orderbookUrlForGetRequest = "https://api.bytick.com/v5/market/orderbook?category=" + category + "&symbol=&limit=200" // URL for getting order book data

	return exchangesData // Return updated exchanges data
}

// setBybitFuturesData sets up data specific to the Bybit Futures exchange.
//
// This function configures the exchange struct with settings specific to the Bybit Futures exchange,
// including URLs for API calls and initializing necessary fields.
//
// Parameters:
//   - exchangesData: A pointer to the exchange struct to be configured.
//
// Returns:
//   - *exchange: A pointer to the updated exchange struct.
func setBybitFuturesData(exchangesData *ExchangeData) *ExchangeData {
	const category = "linear"

	exchangesData.exchangeName = "bybit_futures"                                                                                       // Set the name of the exchange to "bybitFutures"
	exchangesData.pairsUrlForGetRequest = "https://api.bytick.com/v5/market/instruments-info?category=" + category                     // URL for getting futures pairs information
	exchangesData.orderbookUrlForGetRequest = "https://api.bytick.com/v5/market/orderbook?category=" + category + "&symbol=&limit=200" // URL for getting futures order book data

	return exchangesData // Return updated exchanges data
}
