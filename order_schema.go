package tdam

type Order struct {
	Session                  string               `json:"session"`   // 'NORMAL' or 'AM' or 'PM' or 'SEAMLESS'
	Duration                 string               `json:"duration"`  //: "'DAY' or 'GOOD_TILL_CANCEL' or 'FILL_OR_KILL'",
	OrderType                string               `json:"orderType"` //: "'MARKET' or 'LIMIT' or 'STOP' or 'STOP_LIMIT' or 'TRAILING_STOP' or 'MARKET_ON_CLOSE' or 'EXERCISE' or 'TRAILING_STOP_LIMIT' or 'NET_DEBIT' or 'NET_CREDIT' or 'NET_ZERO'",
	CancelTime               interface{}          `json:"cancelTime"`
	ComplexOrderStrategyType string               `json:"complexOrderStrategyType"` // 'NONE' or 'COVERED' or 'VERTICAL' or 'BACK_RATIO' or 'CALENDAR' or 'DIAGONAL' or 'STRADDLE' or 'STRANGLE' or 'COLLAR_SYNTHETIC' or 'BUTTERFLY' or 'CONDOR' or 'IRON_CONDOR' or 'VERTICAL_ROLL' or 'COLLAR_WITH_STOCK' or 'DOUBLE_DIAGONAL' or 'UNBALANCED_BUTTERFLY' or 'UNBALANCED_CONDOR' or 'UNBALANCED_IRON_CONDOR' or 'UNBALANCED_VERTICAL_ROLL' or 'CUSTOM'
	Quantity                 int                  `json:"quantity"`
	FilledQuantity           int                  `json:"filledQuantity"`
	RemainingQuantity        int                  `json:"remainingQuantity"`
	RequestedDestination     string               `json:"requestedDestination"` // 'INET' or 'ECN_ARCA' or 'CBOE' or 'AMEX' or 'PHLX' or 'ISE' or 'BOX' or 'NYSE' or 'NASDAQ' or 'BATS' or 'C2' or 'AUTO'
	DestinationLinkName      string               `json:"destinationLinkName"`  // string
	ReleaseTime              string               `json:"releaseTime"`          // string
	StopPrice                int                  `json:"stopPrice"`
	StopPriceLinkBasis       string               `json:"stopPriceLinkBasis"` // 'MANUAL' or 'BASE' or 'TRIGGER' or 'LAST' or 'BID' or 'ASK' or 'ASK_BID' or 'MARK' or 'AVERAGE'
	StopPriceLinkType        string               `json:"stopPriceLinkType"`  // 'VALUE' or 'PERCENT' or 'TICK'
	StopPriceOffset          int                  `json:"stopPriceOffset"`
	StopType                 string               `json:"stopType"`       // 'STANDARD' or 'BID' or 'ASK' or 'LAST' or 'MARK'
	PriceLinkBasis           string               `json:"priceLinkBasis"` // 'MANUAL' or 'BASE' or 'TRIGGER' or 'LAST' or 'BID' or 'ASK' or 'ASK_BID' or 'MARK' or 'AVERAGE'
	PriceLinkType            string               `json:"priceLinkType"`  // 'VALUE' or 'PERCENT' or 'TICK'
	Price                    int                  `json:"price"`
	TaxLotMethod             string               `json:"taxLotMethod"` // 'FIFO' or 'LIFO' or 'HIGH_COST' or 'LOW_COST' or 'AVERAGE_COST' or 'SPECIFIC_LOT'
	OrderLegCollection       []OrderLegCollection `json:"orderLegCollection"`
	ActivationPrice          int                  `json:"activationPrice"`
	SpecialInstruction       string               `json:"specialInstruction"` // 'ALL_OR_NONE' or 'DO_NOT_REDUCE' or 'ALL_OR_NONE_DO_NOT_REDUCE'
	OrderStrategyType        string               `json:"orderStrategyType"`  // 'SINGLE' or 'OCO' or 'TRIGGER'
	OrderId                  int                  `json:"orderId"`
	Cancelable               bool                 `json:"cancelable"`
	Editable                 bool                 `json:"editable"`
	Status                   string               `json:"status"`      // 'AWAITING_PARENT_ORDER' or 'AWAITING_CONDITION' or 'AWAITING_MANUAL_REVIEW' or 'ACCEPTED' or 'AWAITING_UR_OUT' or 'PENDING_ACTIVATION' or 'QUEUED' or 'WORKING' or 'REJECTED' or 'PENDING_CANCEL' or 'CANCELED' or 'PENDING_REPLACE' or 'REPLACED' or 'FILLED' or 'EXPIRED'
	EnteredTime              string               `json:"enteredTime"` // string
	CloseTime                string               `json:"closeTime"`   // string
	Tag                      string               `json:"tag"`         // string
	AccountId                int                  `json:"accountId"`
	/*
	  "orderActivityCollection": [
	    "The type <OrderActivity> has the following subclasses [Execution] descriptions are listed below"
	  ],
	  "replacingOrderCollection": [
	    {}
	  ],
	  "childOrderStrategies": [
	    {}
	  ],
	*/
	StatusDescription string `json:"statusDescription"` // string
}

type OrderLegCollection struct {
	OrderLegType   string `json:"orderLegType"` // 'EQUITY' or 'OPTION' or 'INDEX' or 'MUTUAL_FUND' or 'CASH_EQUIVALENT' or 'FIXED_INCOME' or 'CURRENCY'
	LegId          int    `json:"legId"`
	Instrument     string `json:"instrument"`     // The type <Instrument> has the following subclasses [Option, MutualFund, CashEquivalent, Equity, FixedIncome] descriptions are listed below\"
	Instruction    string `json:"instruction"`    // 'BUY' or 'SELL' or 'BUY_TO_COVER' or 'SELL_SHORT' or 'BUY_TO_OPEN' or 'BUY_TO_CLOSE' or 'SELL_TO_OPEN' or 'SELL_TO_CLOSE' or 'EXCHANGE'
	PositionEffect string `json:"positionEffect"` // 'OPENING' or 'CLOSING' or 'AUTOMATIC'
	Quantity       int    `json:"quantity"`
	QuantityType   string `json:"quantityType"` // 'ALL_SHARES' or 'DOLLARS' or 'SHARES'
}

/*
//The class <Instrument> has the
//following subclasses:
//-Option
//-MutualFund
//-CashEquivalent
//-Equity
//-FixedIncome
//JSON for each are listed below:

//Option:
{
  "assetType": "'EQUITY' or 'OPTION' or 'INDEX' or 'MUTUAL_FUND' or 'CASH_EQUIVALENT' or 'FIXED_INCOME' or 'CURRENCY'",
  "cusip": "string",
  "symbol": "string",
  "description": "string",
  "type": "'VANILLA' or 'BINARY' or 'BARRIER'",
  "putCall": "'PUT' or 'CALL'",
  "underlyingSymbol": "string",
  "optionMultiplier": 0,
  "optionDeliverables": [
    {
      "symbol": "string",
      "deliverableUnits": 0,
      "currencyType": "'USD' or 'CAD' or 'EUR' or 'JPY'",
      "assetType": "'EQUITY' or 'OPTION' or 'INDEX' or 'MUTUAL_FUND' or 'CASH_EQUIVALENT' or 'FIXED_INCOME' or 'CURRENCY'"
    }
  ]
}

//OR

//MutualFund:
{
  "assetType": "'EQUITY' or 'OPTION' or 'INDEX' or 'MUTUAL_FUND' or 'CASH_EQUIVALENT' or 'FIXED_INCOME' or 'CURRENCY'",
  "cusip": "string",
  "symbol": "string",
  "description": "string",
  "type": "'NOT_APPLICABLE' or 'OPEN_END_NON_TAXABLE' or 'OPEN_END_TAXABLE' or 'NO_LOAD_NON_TAXABLE' or 'NO_LOAD_TAXABLE'"
}

//OR

//CashEquivalent:
{
  "assetType": "'EQUITY' or 'OPTION' or 'INDEX' or 'MUTUAL_FUND' or 'CASH_EQUIVALENT' or 'FIXED_INCOME' or 'CURRENCY'",
  "cusip": "string",
  "symbol": "string",
  "description": "string",
  "type": "'SAVINGS' or 'MONEY_MARKET_FUND'"
}

//OR

//Equity:
{
  "assetType": "'EQUITY' or 'OPTION' or 'INDEX' or 'MUTUAL_FUND' or 'CASH_EQUIVALENT' or 'FIXED_INCOME' or 'CURRENCY'",
  "cusip": "string",
  "symbol": "string",
  "description": "string"
}

//OR

//FixedIncome:
{
  "assetType": "'EQUITY' or 'OPTION' or 'INDEX' or 'MUTUAL_FUND' or 'CASH_EQUIVALENT' or 'FIXED_INCOME' or 'CURRENCY'",
  "cusip": "string",
  "symbol": "string",
  "description": "string",
  "maturityDate": "string",
  "variableRate": 0,
  "factor": 0
}

//The class <OrderActivity> has the
//following subclasses:
//-Execution
//JSON for each are listed below:

//Execution:
{
  "activityType": "'EXECUTION' or 'ORDER_ACTION'",
  "executionType": "'FILL'",
  "quantity": 0,
  "orderRemainingQuantity": 0,
  "executionLegs": [
    {
      "legId": 0,
      "quantity": 0,
      "mismarkedQuantity": 0,
      "price": 0,
      "time": "string"
    }
  ]
}
*/
