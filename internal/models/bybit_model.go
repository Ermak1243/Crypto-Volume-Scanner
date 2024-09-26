package models

type BybitPairsJSONResponse struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
	Result  struct {
		Category string `json:"category"`
		List     []struct {
			Symbol        string `json:"symbol"`
			BaseCoin      string `json:"baseCoin"`
			QuoteCoin     string `json:"quoteCoin"`
			Innovation    string `json:"innovation"`
			Status        string `json:"status"`
			LotSizeFilter struct {
				BasePrecision  string `json:"basePrecision"`
				QuotePrecision string `json:"quotePrecision"`
				MinOrderQty    string `json:"minOrderQty"`
				MaxOrderQty    string `json:"maxOrderQty"`
				MinOrderAmt    string `json:"minOrderAmt"`
				MaxOrderAmt    string `json:"maxOrderAmt"`
			} `json:"lotSizeFilter"`
			PriceFilter struct {
				TickSize string `json:"tickSize"`
			} `json:"priceFilter"`
		} `json:"list"`
	} `json:"result"`
	RetExtInfo struct {
	} `json:"retExtInfo"`
	Time int64 `json:"time"`
}

type BybitOrderbookJSONResponse struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
	Result  struct {
		S    string          `json:"s"`
		Asks [][]interface{} `json:"a"`
		Bids [][]interface{} `json:"b"`
		Ts   int64           `json:"ts"`
		U    int             `json:"u"`
	} `json:"result"`
	RetExtInfo struct {
	} `json:"retExtInfo"`
	Time int64 `json:"time"`
}
