/*
	Developed by Güray Gurkan & Kaan Taha Köken
	Contact us via https://github.com/gurkanguray/ & https://github.com/kaankoken/
*/

package main

import (
	"Arbitrage-Trading-Bot/data"
	"Arbitrage-Trading-Bot/helper"
	"Arbitrage-Trading-Bot/network"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

const (
	SELL = "sell"
	BUY  = "buy"
)

func CompareExchanges(whiteBitData data.ArbitrageDataClass, bitFinexData data.ArbitrageDataClass, minDiffRatio float64, market string) {
	whitebit_selected_pair, whitebit_usdt := network.GetWhitebitWallet(market)
	bitfinex_selected_pair, bitfinex_usdt := network.GetBitfinexWallet(market)
	bitfinex_usdt_float, _ := strconv.ParseFloat(bitfinex_usdt, 64)
	whitebit_usdt_float, _ := strconv.ParseFloat(whitebit_usdt, 64)
	whitebit_selected_pair_float, _ := strconv.ParseFloat(whitebit_selected_pair, 64)
	bitfinex_selected_pair_float, _ := strconv.ParseFloat(bitfinex_selected_pair, 64)

	//TODO: SAYILARI MANUEL DEĞİŞTİREBİLİRSİNİZ
	// İŞLEM AÇMAK İÇİN MİNİMUM HESAP BAKİYE LİMİTİ
	if whitebit_usdt_float < 5.0 || bitfinex_usdt_float < 5.0 {
		fmt.Println("Bakiye yetersiz!")
		fmt.Println("Whitebit " + whitebit_selected_pair + " " + whitebit_usdt)
		fmt.Println("Bitfinex " + bitfinex_selected_pair + " " + bitfinex_usdt)
	} else if whitebit_selected_pair_float == 0.0 || bitfinex_selected_pair_float == 0.0 { 
		fmt.Println("Bakiye yetersiz!")
	} else {

		if (whiteBitData.Bid > bitFinexData.Ask) && (((whiteBitData.Bid-bitFinexData.Ask)*float64(100))/whiteBitData.Bid >= minDiffRatio) {

			amount := fmt.Sprintf("%v", whiteBitData.BidVolume)

			if whiteBitData.BidVolume > whitebit_selected_pair_float {
				amount = whitebit_selected_pair
			}

			whitebit, _ := json.Marshal(map[string]string{
				"market":  market + "_USDT",
				"side":    SELL,
				"amount":  amount,
				"request": "/api/v4/main-account/balance",
				"nonce":   fmt.Sprint(helper.CreateTimestamp()),
			})

			amount_float, _ := strconv.ParseFloat(amount, 64)

			if bitFinexData.Ask*amount_float > bitfinex_usdt_float {
				amount = fmt.Sprintf("%v", (bitfinex_usdt_float / bitFinexData.Ask))
			}

			bitfinex, _ := json.Marshal(map[string]string{
				"type":   "MARKET",
				"symbol": "t" + market + "USD",
				"amount": "+" + amount,
				"price":  "0.0",
				"cid":    fmt.Sprint(helper.CreateTimestamp()),
			})

			network.PostWhitebit(whitebit)
			network.PostBitfinex(bitfinex)
			fmt.Println("Condition Satisfied - BitFinex Buy - Whitebit Sell")
		}

		if (bitFinexData.Bid > whiteBitData.Ask) && (((bitFinexData.Bid-whiteBitData.Ask)*float64(100))/bitFinexData.Bid >= minDiffRatio) {
			amount := fmt.Sprintf("%v", bitFinexData.BidVolume)

			if bitFinexData.BidVolume > bitfinex_selected_pair_float {
				amount = bitfinex_selected_pair
			}

			bitfinex, _ := json.Marshal(map[string]string{
				"type":   "MARKET",
				"symbol": "t" + market + "USD",
				"amount": "-" + fmt.Sprintf("%v", math.Abs(bitFinexData.BidVolume)),
				"price":  "0.0",
				"cid":    fmt.Sprint(helper.CreateTimestamp()),
			})

			amount_float, _ := strconv.ParseFloat(amount, 64)

			if whiteBitData.Ask*amount_float > whitebit_usdt_float {
				amount = fmt.Sprintf("%v", (whitebit_usdt_float / whiteBitData.Ask))
			}

			whitebit, _ := json.Marshal(map[string]string{
				"market":  market + "_USDT",
				"side":    BUY,
				"amount":  fmt.Sprintf("%v", math.Abs(bitFinexData.BidVolume)),
				"request": "/api/v4/main-account/balance",
				"nonce":   fmt.Sprint(helper.CreateTimestamp()),
			})

			network.PostWhitebit(whitebit)
			network.PostBitfinex(bitfinex)
			fmt.Println("Condition Satisfied - BitFinex Sell - Whitebit Buy")
		}
		fmt.Println("Condition NOT Satisfied")
	}
	
}
