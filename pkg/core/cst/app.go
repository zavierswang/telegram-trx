package cst

const (
	AppName           = "telegram-trx"
	BaseName          = "telegram"
	DateTimeFormatter = "2006-01-02 15:04:05"
	TimeFormatter     = "15:04:05"
	PublicKey         = "./public.key"
	PrivateKey        = "./private.pem"
	UserAgent         = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"
)

const (
	OrderStatus = iota
	OrderStatusSuccess
	OrderStatusRunning
	OrderStatusReceived
	OrderStatusApiSuccess
	OrderStatusApiFailure
	OrderStatusFailure
	OrderStatusNotSufficientFunds
	OrderStatusCancel
)

const (
	OkxMarketTradesApi  = "https://www.okx.com/priapi/v5/market/trades"
	OkxTradingOrdersApi = "https://www.okx.com/v3/c2c/tradingOrders/books"
	LicenseApi          = "http://127.0.0.1/api/license"
)
