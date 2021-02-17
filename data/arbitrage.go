package data

/*
	Developed by Güray Gurkan & Kaan Taha Köken
	Contact us via https://github.com/gurkanguray/ & https://github.com/kaankoken/
*/

type ArbitrageDataClass arbitrage

type arbitrage struct {
	Bid       float64 `json:"BID,omitempty"`
	BidVolume float64 `json:"BIDVOLUME,omitempty"`
	Ask       float64 `json:"ASK,omitempty"`
	AskVolume float64 `json:"ASKVOLUME,omitempty"`
}
