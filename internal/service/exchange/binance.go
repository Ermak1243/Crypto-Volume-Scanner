package exchange

import (
	"strings"
	"time"

	"cvs/internal/models"
	"cvs/internal/service"
	"cvs/internal/service/logger"
	"cvs/internal/service/orderbook"

	"github.com/goccy/go-json"
	cmap "github.com/orcaman/concurrent-map/v2"
)

// Overall data for all sections of the Binance exchange
var (
	binanceTimeBetweenRequests = 3 * time.Second                       // Time interval between requests to the Binance API
	binancePairsJsonModel      = models.BinancePairsJSONResponse{}     // Model for Binance pairs JSON response
	binanceOrderbookJsonModel  = models.BinanceOrderbookJSONResponse{} // Model for Binance order book JSON response
	binanceOrderbookService    = orderbook.NewOrderbook()              // Instance of the order book service for managing order data

	// Function to parse order book JSON response from Binance
	binanceOrderbookJsonParse = func(bodyBytes []byte) ([][]interface{}, [][]interface{}, error) {
		var model models.BinanceOrderbookJSONResponse

		// Unmarshal the response body into jsonData to inspect the response
		err := json.Unmarshal(bodyBytes, &model)

		return model.Asks, model.Bids, err
	}

	// Function to format Binance API URLs with the trading pair
	binanceUrlFormatter = func(url, pair string) string {
		pairFormatted := strings.Replace(pair, "/", "", -1)                 // Remove slashes from the pair string
		replacer := strings.NewReplacer("symbol=", "symbol="+pairFormatted) // Replace "symbol=" in the URL with the formatted pair

		return replacer.Replace(url) // Return the formatted URL
	}

	// Function to parse exchange pairs from Binance API response
	binanceExchangePairsJsonParse = func(exchangeName string, bodyBytes []byte) ([]models.ExchangePairs, error) {
		var model models.BinancePairsJSONResponse
		// Unmarshal the response body into jsonData to inspect the response
		err := json.Unmarshal(bodyBytes, &model)
		if err != nil {
			return []models.ExchangePairs{}, errUnmarshal("exchange pairs", exchangeName) // Return error if unmarshalling fails
		}

		var exchangePairsSlice []models.ExchangePairs // Slice to hold parsed exchange pairs

		for i := 0; i < len(model.Symbols); i++ { // Iterate over all symbols in pairs data
			if model.Symbols[i].QuoteAsset == "BUSD" { // Skip pairs with BUSD as quote asset
				continue
			}

			exchangePairsSlice = append(exchangePairsSlice, models.ExchangePairs{
				Pair:     model.Symbols[i].BaseAsset + "/" + model.Symbols[i].QuoteAsset, // Construct pair string
				Exchange: exchangeName,                                                   // Set exchange name
			})
		}

		return exchangePairsSlice, nil // Return the slice of exchange pairs
	}
)

// NewBinance initializes instances of different Binance exchanges.
//
// This function creates and returns a slice of Exchange instances for various Binance exchanges,
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
//   - []Exchange: A slice containing instances of different Binance exchanges.
func NewBinance(
	userService service.UserService,
	userPairsService service.UserPairsService,
	httpRequestService service.HttpRequest,
	foundVolumeService service.FoundVolumesService,
	logger logger.Logger,
) []Exchange {
	var binances []Exchange // Slice to hold instances of different Binance exchanges
	initFunctions := []func(exchangesData *ExchangeData) *ExchangeData{
		setBinanceSpotData,
		setBinanceFuturesData,
		setBinanceUsData,
	}

	for _, function := range initFunctions {
		exchangeData := setBinanceOverallData(
			userService,
			userPairsService,
			httpRequestService,
			foundVolumeService,
			logger,
		)

		binances = append(binances, function(exchangeData))
	}

	return binances // Return the slice of Binance exchanges
}

// setBinanceOverallData initializes and sets up overall data for all Binance exchanges.
//
// This function creates an instance of the exchange struct and populates it with the necessary services,
// models, and configurations required for interacting with Binance exchanges. It prepares the exchange
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
func setBinanceOverallData(
	userService service.UserService,
	userPairsService service.UserPairsService,
	httpRequestService service.HttpRequest,
	foundVolumeService service.FoundVolumesService,
	logger logger.Logger,
) *ExchangeData {
	binanceExchangesData := ExchangeData{
		userService:            userService,
		userPairsService:       userPairsService,
		httpRequestService:     httpRequestService,
		foundVolumesService:    foundVolumeService,
		logger:                 logger,
		pairsJsonModel:         binancePairsJsonModel,            // Set pairs JSON model for exchanges
		orderbookJsonModel:     binanceOrderbookJsonModel,        // Set orderbook JSON model for exchanges
		urlFormatter:           binanceUrlFormatter,              // Set URL formatter function for exchanges
		timeBetweenRequests:    binanceTimeBetweenRequests,       // Set time between requests for exchanges
		orderbookService:       binanceOrderbookService,          // Assign order book service instance to exchanges data
		pairsSubscribed:        cmap.New[bool](),                 // Initialize subscribed pairs list as empty
		allPairsOfExchange:     cmap.New[models.ExchangePairs](), // Initialize concurrent map for all pairs of the exchange
		orderbookJsonParse:     binanceOrderbookJsonParse,        // Set order book JSON parsing function for exchanges
		exchangePairsJsonParse: binanceExchangePairsJsonParse,    // Set exchange pairs JSON parsing function for exchanges
	}

	return &binanceExchangesData
}

// setBinanceSpotData sets up data specific to the Binance Spot exchange.
//
// This function configures the exchange struct with settings specific to the Binance Spot exchange,
// including URLs for API calls and initializing necessary fields.
//
// Parameters:
//   - exchangesData: A pointer to the exchange struct to be configured.
//
// Returns:
//   - *exchange: A pointer to the updated exchange struct.
func setBinanceSpotData(exchangesData *ExchangeData) *ExchangeData {
	exchangesData.exchangeName = "binance_spot"                                                        // Set the name of the exchange to "binanceSpot"
	exchangesData.pairsUrlForGetRequest = "https://api.binance.com/api/v3/exchangeInfo"                // URL for getting pairs information
	exchangesData.orderbookUrlForGetRequest = "https://api.binance.com/api/v1/depth?symbol=&limit=500" // URL for getting order book data

	return exchangesData // Return updated exchanges data
}

// setBinanceUsData sets up data specific to the Binance US exchange.
//
// This function configures the exchange struct with settings specific to the Binance US exchange,
// including URLs for API calls and initializing necessary fields.
//
// Parameters:
//   - exchangesData: A pointer to the exchange struct to be configured.
//
// Returns:
//   - *exchange: A pointer to the updated exchange struct.
func setBinanceUsData(exchangesData *ExchangeData) *ExchangeData {
	exchangesData.exchangeName = "binance_us"                                                         // Set the name of the exchange to "binanceUs"
	exchangesData.pairsUrlForGetRequest = "https://api.binance.us/api/v3/exchangeInfo"                // URL for getting pairs information from Binance US
	exchangesData.orderbookUrlForGetRequest = "https://api.binance.us/api/v3/depth?symbol=&limit=500" // URL for getting order book data from Binance US

	return exchangesData // Return updated exchanges data
}

// setBinanceFuturesData sets up data specific to the Binance Futures exchange.
//
// This function configures the exchange struct with settings specific to the Binance Futures exchange,
// including URLs for API calls and initializing necessary fields.
//
// Parameters:
//   - exchangesData: A pointer to the exchange struct to be configured.
//
// Returns:
//   - *exchange: A pointer to the updated exchange struct.
func setBinanceFuturesData(exchangesData *ExchangeData) *ExchangeData {
	exchangesData.exchangeName = "binance_futures"                                                       // Set the name of the exchange to "binanceFutures"
	exchangesData.pairsUrlForGetRequest = "https://fapi.binance.com/fapi/v1/exchangeInfo"                // URL for getting futures pairs information
	exchangesData.orderbookUrlForGetRequest = "https://fapi.binance.com/fapi/v1/depth?symbol=&limit=500" // URL for getting futures order book data

	return exchangesData // Return updated exchanges data
}
