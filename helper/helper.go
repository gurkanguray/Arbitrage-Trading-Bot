/*
	Developed by Güray Gurkan & Kaan Taha Köken
	Contact us via https://github.com/gurkanguray/ & https://github.com/kaankoken/
*/

package helper

import (
	"crypto/hmac"
	"crypto/sha512"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func ConvertToBase64(data []byte) string {
	uEnc := b64.StdEncoding.EncodeToString(data)
	return uEnc
}

func CreateHMAC(data []byte, key string, exchange string) string {
	var hmacKey hash.Hash
	switch exchange {
	case "whitebit":
		hmacKey = hmac.New(sha512.New, []byte(key))
	case "bitfinex":
		hmacKey = hmac.New(sha512.New384, []byte(key))
	}
	hmacKey.Write(data)

	sha := fmt.Sprintf("%x", hmacKey.Sum(nil))

	return sha
}

func CreateTimestamp() int {
	return int(time.Now().Unix())
}

func ReadJSONFile(key string) (string, string) {
	jsonFile, err := os.Open("src/keys.json")
	if err != nil {
		fmt.Println(err)
	}

	var embeddedObject map[string]interface{}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	if err := json.Unmarshal(byteValue, &embeddedObject); err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()
	result := embeddedObject[key].(map[string]interface{})

	return result["API_KEY"].(string), result["API_SECRET"].(string)
}

func Filter(ss []string, test func(string) bool) (ret []string) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}
