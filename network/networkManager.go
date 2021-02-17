/*
	Developed by Güray Gurkan & Kaan Taha Köken
	Contact us via https://github.com/gurkanguray/ & https://github.com/kaankoken/
*/

package network

import (
	"Arbitrage-Trading-Bot/data"
	"Arbitrage-Trading-Bot/helper"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

const (
	GET               = "GET"
	POST              = "POST"
	WHITEBIT          = "whitebit"
	BITFENIX          = "bitfinex"
	WHITEBIT_LINK     = "https://whitebit.com/api/v4/public/orderbook/"
	BITFENIX_OLD_LINK = "https://api-pub.bitfinex.com/v2/ticker/t"
	BITFENIX_LINK     = "https://api-pub.bitfinex.com/v2/book/t"
	WHITEBIT_CURRENCY = "https://whitebit.com/api/v1/public/symbols"
	WHITEBIT_ORDER    = "https://whitebit.com/api/v4/order/market"
	BITFENIX_ORDER    = "https://api.bitfinex.com/v2/auth/w/order/submit"
	BITFENIX_WALLET   = "https://api.bitfinex.com/v2/auth/r/wallets"
	WHITEBIT_WALLET   = "https://whitebit.com/api/v1/account/balances"
)

var client = &http.Client{}
var arbitrageWhitebit = data.ArbitrageDataClass{}
var arbitrageBitfinex = data.ArbitrageDataClass{}

func SendRequestExchange(request_type string, provider string, link string) data.ArbitrageDataClass {
	switch request_type {
	case GET:
		if provider == BITFENIX {
			return getBitfenix(link)
		} else {
			return getWhitebit(link)
		}
	case POST:
	}
	return data.ArbitrageDataClass{}
}

func SendRequest(request_type string, link string) []string {
	switch request_type {
	case GET:
		return get(link)
	case POST:
		//post()
	}
	return nil
}

func getWhitebit(url string) data.ArbitrageDataClass {

	req, err := http.NewRequest(GET, url, nil)

	if err != nil {
		fmt.Println(err)
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	var embeddedObject map[string]interface{}

	if err := json.Unmarshal(body, &embeddedObject); err != nil {
		log.Fatal(err)
	}
	//fmt.Println(embeddedObject)
	//bid := (embeddedObject["bid"]).(string)
	//bid := fmt.Sprintf("%f", EmbeddedObject.GetAsk(embeddedObject).(string))
	//arbitrageWhitebit.Bid = bid

	askRaw := embeddedObject["asks"].([]interface{})
	bidRaw := embeddedObject["bids"].([]interface{})

	ask := askRaw[0].([]interface{})
	bid := bidRaw[0].([]interface{})

	i, err := strconv.ParseFloat(ask[0].(string), 64)
	arbitrageWhitebit.Ask = i
	j, _ := strconv.ParseFloat(ask[1].(string), 64)
	arbitrageWhitebit.AskVolume = j
	k, _ := strconv.ParseFloat(bid[0].(string), 64)
	arbitrageWhitebit.Bid = k
	l, _ := strconv.ParseFloat(bid[1].(string), 64)
	arbitrageWhitebit.BidVolume = l

	return arbitrageWhitebit
}

func getBitfenix(url string) data.ArbitrageDataClass {
	req, err := http.NewRequest(GET, url, nil)

	if err != nil {
		fmt.Println(err)
	}
	query := req.URL.Query()
	query.Add("len", "1")
	req.URL.RawQuery = query.Encode()

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()

	var decoded []interface{}
	err = json.NewDecoder(res.Body).Decode(&decoded)
	if err != nil {
		log.Fatal(err)
	}

	if len(decoded) > 0 {
		if reflect.TypeOf(decoded[0]).Kind() == reflect.String {
			return data.ArbitrageDataClass{}
		}
		bidRaw := decoded[0].([]interface{})
		askRaw := decoded[1].([]interface{})

		bid := bidRaw[0].(float64)
		bidVolume := bidRaw[2].(float64)

		ask := askRaw[0].(float64)
		askVolume := askRaw[2].(float64)

		arbitrageBitfinex.Bid = bid
		arbitrageBitfinex.BidVolume = bidVolume
		arbitrageBitfinex.Ask = ask
		arbitrageBitfinex.AskVolume = math.Abs(askVolume)

		return arbitrageBitfinex
	}
	return data.ArbitrageDataClass{}
	//bitfinexStruct, _ := json.Marshal(tempStruct)
}

func get(url string) []string {
	req, err := http.NewRequest(GET, url, nil)
	if err != nil {
		fmt.Println(err)
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()
	var decoded map[string]interface{}

	err = json.NewDecoder(res.Body).Decode(&decoded)
	if err != nil {
		log.Fatal(err)
	}

	mappedData := decoded["result"].([]interface{})
	s := make([]string, len(mappedData))
	for i, v := range mappedData {
		s[i] = fmt.Sprint(v)
	}

	filterList := func(s string) bool { return strings.HasSuffix(s, "_USDT") }
	filteredList := helper.Filter(s, filterList)

	return filteredList
}

func PostWhitebit(responseBody []byte) {
	payload := helper.ConvertToBase64(responseBody)
	key, secret := helper.ReadJSONFile(WHITEBIT)

	req, err := http.NewRequest(POST, WHITEBIT_ORDER, bytes.NewBuffer(responseBody))
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Set("Content-type", "application/json")
	req.Header.Set("X-TXC-APIKEY", key)
	req.Header.Set("X-TXC-PAYLOAD", payload)
	req.Header.Set("X-TXC-SIGNATURE", helper.CreateHMAC([]byte(payload), secret, WHITEBIT))

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
}

func PostBitfinex(responseBody []byte) {
	nonce := strconv.Itoa(helper.CreateTimestamp())

	req, err := http.NewRequest(POST, BITFENIX_ORDER, bytes.NewBuffer(responseBody))
	if err != nil {
		fmt.Println(err)
	}

	signature := "/api/v2/auth/r/wallets" + nonce + string(responseBody)
	key, secret := helper.ReadJSONFile(BITFENIX)
	sig := helper.CreateHMAC([]byte(signature), secret, BITFENIX)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("bfx-nonce", nonce)
	req.Header.Add("bfx-apikey", key)
	req.Header.Add("bfx-signature", sig)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
}

func GetBitfinexWallet(pair string) (string, string) {
	balance, _ := json.Marshal(map[string]string{})

	nonce := strconv.Itoa(helper.CreateTimestamp())
	req, err := http.NewRequest(POST, BITFENIX_WALLET, bytes.NewBuffer(balance))
	if err != nil {
		fmt.Println(err)
	}

	signature := "/api/v2/auth/r/wallets" + nonce + string(balance)
	key, secret := helper.ReadJSONFile(BITFENIX)
	sig := helper.CreateHMAC([]byte(signature), secret, BITFENIX)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("bfx-nonce", nonce)
	req.Header.Add("bfx-apikey", key)
	req.Header.Add("bfx-signature", sig)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	var usdt string
	var selectedPair string
	if len(body) > 0 {
		var embeddedObject []interface{}
		if err := json.Unmarshal(body, &embeddedObject); err != nil {
			log.Fatal(err)
		}

		for i, v := range embeddedObject {
			value := embeddedObject[i].([]interface{})
			if value[1].(string) == "USD" {
				usdt = value[4].(string)
				fmt.Println(v)
			}

			if value[1].(string) == pair {
				selectedPair = value[4].(string)
			}
		}
	}
	return selectedPair, usdt
}

func GetWhitebitWallet(pair string) (string, string) {
	responseBody, _ := json.Marshal(map[string]string{
		"request": "/api/v1/account/balances",
		"nonce":   strconv.Itoa(helper.CreateTimestamp()),
	})

	payload := helper.ConvertToBase64(responseBody)
	key, secret := helper.ReadJSONFile(WHITEBIT)

	req, err := http.NewRequest(POST, WHITEBIT_WALLET, bytes.NewBuffer(responseBody))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-type", "application/json")
	req.Header.Set("X-TXC-APIKEY", key)
	req.Header.Set("X-TXC-PAYLOAD", payload)
	req.Header.Set("X-TXC-SIGNATURE", helper.CreateHMAC([]byte(payload), secret, WHITEBIT))

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	var embeddedObject map[string]interface{}

	if err := json.Unmarshal(body, &embeddedObject); err != nil {
		log.Fatal(err)
	}

	result := embeddedObject["result"].(map[string]interface{})
	usdt := result["USDT"].(map[string]interface{})
	selected_pair := result[pair].(map[string]interface{})

	return selected_pair["available"].(string), usdt["available"].(string)
}
