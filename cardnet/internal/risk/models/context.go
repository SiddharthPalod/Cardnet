package models

type RiskContext struct {
	TxCount1Min int
	TxCount5Min int

	LastCountry string
	LastTxTime  int64

	UserAvgAmount int64
	UserP95Amount int64

	MerchantRiskScore int
	BINRiskScore      int
}
