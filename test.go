package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type whitebitStatus struct {
	Success  bool `json:"bid,omitempty"`
	Message string `json:"message,omitempty"`
	Result whitebitResult `json:"result,omitempty"`
}

type whitebitResult struct {
	Bid float64 `json:"bid,omitempty"`
	Ask float64 `json:"ask,omitempty"`
	Open float64 `json:"open,omitempty"`
	High float64 `json:"high,omitempty"`
	Low float64 `json:"low,omitempty"`
	Last float64 `json:"last,omitempty"`
	Volume float64 `json:"volume,omitempty"`
	Change float64 `json:"change,omitempty"`
}

func main() {
	a, b := whitebit("BTC")

	fmt.Println(a, b)

  
}

func whitebit(pair string) (float64, float64)  {
	url := "https://whitebit.com/api/v1/public/ticker?market="+"pair"+"_USDT"
	method := "GET"

	client := &http.Client {
  	}
  	req, err := http.NewRequest(method, url, nil)

	if err != nil {
    	fmt.Println(err)
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
  	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
	
	var embeddedObject EmbeddedObject

	json.Unmarshal(body, &embeddedObject)

	result := embeddedObject["result"].(EmbeddedObject)

	return result.GetAsk(), result.GetBid()

	

	// fmt.Printf("success: %v\n\n", v.Result)

	fmt.Println("ASK:",result["ask"])

	ask := result["ask"].(string)
	a, err := strconv.ParseFloat(ask, 64)
	if err != nil {
    	fmt.Println(a)
	}

	bid := result["bid"].(string)
	b, err := strconv.ParseFloat(bid, 64)
	if err != nil {
    	fmt.Println(b)
	}

	return result.GetAsk(), result.GetBid()

}

// func find(obj interface{}, key string) (interface{}, bool) {
//     //if the argument is not a map, ignore it
//     mobj, ok := obj.(map[string]interface{})
//     if !ok {
//         return nil, false
//     }

//     for k, v := range mobj {
//         //key match, return value
//         if k == key {
//             return v, true
//         }

//         //if the value is a map, search recursively
//         if m, ok := v.(map[string]interface{}); ok {
//             if res, ok := find(m, key); ok {
//                 return res, true
//             }
//         }
//         //if the value is an array, search recursively 
//         //from each element
//         if va, ok := v.([]interface{}); ok {
//             for _, a := range va {
//                 if res, ok := find(a, key); ok {
//                     return res,true
//                 }
//             }
//         }
//     }

//     //element not found
//     return nil,false
// }