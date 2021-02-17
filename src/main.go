/*
	Developed by Güray Gurkan & Kaan Taha Köken
	Contact us via https://github.com/gurkanguray/ & https://github.com/kaankoken/
*/

package main

import (
	"Arbitrage-Trading-Bot/data"
	network "Arbitrage-Trading-Bot/network"
	"fmt"

	"github.com/chebyrash/promise"
	"github.com/jasonlvhit/gocron"
)

func main() {

	arbj := menu()
	pair := "NEO"

	//"https://api-pub.bitfinex.com/v2/book/tBTCUSD/P0"
	//"https://api-pub.bitfinex.com/v2/book/t"
	url_bitfinex := network.BITFENIX_LINK + pair + "USD/" + "P0"
	url_whitebit := network.WHITEBIT_LINK + pair + "_USDT?depth=5&level=1"

	cronJob := gocron.NewScheduler()
	cronJob.Every(2).Second().Do(promiseWork, url_bitfinex, url_whitebit, arbj, pair)
	<-cronJob.Start()
}

func whiteBitGoRoutine(link string) data.ArbitrageDataClass {
	return network.SendRequestExchange(network.GET, network.WHITEBIT, link)
}

func bitFinexGoRoutine(link string) data.ArbitrageDataClass {
	return network.SendRequestExchange(network.GET, network.BITFENIX, link)
}

func promiseWork(link1 string, link2 string, arbitrage float64, market string) {
	var bitFinexPromise = promise.Resolve(bitFinexGoRoutine(link1))
	var whiteBitPromise = promise.Resolve(whiteBitGoRoutine(link2))

	results, _ := promise.All(bitFinexPromise, whiteBitPromise).Await()
	values := results.([]promise.Any)

	fmt.Println()
	fmt.Printf("    BITFINEX \t| VOLUME \t| WHITEBIT \t| VOLUME\n")
	fmt.Printf("Bid %f   | %f   | %f   | %f\n", values[0].(data.ArbitrageDataClass).Bid, values[0].(data.ArbitrageDataClass).BidVolume,
		values[1].(data.ArbitrageDataClass).Bid, values[1].(data.ArbitrageDataClass).BidVolume)
	fmt.Printf("Ask %f   | %f   | %f   | %f\n", values[0].(data.ArbitrageDataClass).Ask, values[0].(data.ArbitrageDataClass).AskVolume,
		values[1].(data.ArbitrageDataClass).Ask, values[1].(data.ArbitrageDataClass).AskVolume)
	fmt.Println()
	CompareExchanges(values[1].(data.ArbitrageDataClass), values[0].(data.ArbitrageDataClass), arbitrage, market)
	fmt.Println()
}

func menu() float64 {
	//reader := bufio.NewReader(os.Stdin)
	var arbj float64

	fmt.Print("Arbitrajı Girin (Örnek: 10.7): ")
	//temp, _ := reader.ReadString('\n')
	fmt.Scanf("%f ", &arbj)

	//fmt.Print("Minimum Alış/Satış Miktar ($): ")
	//temp, _ = reader.ReadString('\n')
	//fmt.Scanf("%f ", &temp)

	//fmt.Println("Coin Listesi: ")
	//for i, s := range currencyList {
	//	fmt.Printf("%s : (%d)\n", s, i)
	//}
	//temp, _ = reader.ReadString('\n')
	//index, _ := strconv.ParseInt(temp, 10, 16)
	//fmt.Print("Tercih: ")
	//fmt.Scanf("%f ", &temp)
	return arbj
}

/*
func createCurrencyTable() (ret []string) {
	data := network.SendRequest(network.GET, network.WHITEBIT_CURRENCY)
	for _, s := range data {
		splitted := strings.Split(s, "_")[0]
		link := network.BITFENIX_LINK + splitted + "USD/" + "P0"
		result := bitFinexGoRoutine(link)

		if result.Bid > 0.0 {
			ret = append(ret, splitted)
		}
	}
	return ret
}
*/
