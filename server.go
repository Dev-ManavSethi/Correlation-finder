package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"

	"github.com/mongodb/mongo-go-driver/bson"

	"github.com/montanaflynn/stats"
	"golang.org/x/net/websocket"
)

type SocketRequest struct {
	Exchange string   `json:"exchange"`
	Pairs    []string `json : "pairs"`
	Candle   string   `json:"candle"`
	Price    string   `json:"price"`
	Days     int64    `json:"days"`
}

func websocketHandler(ws *websocket.Conn) {

	MongoDBclient := ConnectToMongoDB()

	Request := &SocketRequest{}

	error1 := websocket.JSON.Receive(ws, &Request)
	FatalOnError(error1, "Error recieving request from Websocket")

	FindCorrelation(ws, Request, MongoDBclient)

}

func ConnectToMongoDB() *mongo.Client {
	MongoDBclient, err := mongo.Connect(context.TODO(), "mongodb://34.227.81.97:27017")

	FatalOnError(err, "Error connecting to MongoDB (mongo db go official driver)")
	// Check the connection
	err2 := MongoDBclient.Ping(context.TODO(), nil)
	FatalOnError(err2, "Error pinging mongodb (mongo db go official driver)")

	fmt.Println("Connected to MongoDB!")
	return MongoDBclient

}

func FindCorrelation(ws *websocket.Conn, Request *SocketRequest, MongoDBclient *mongo.Client) {
	Exchange := Request.Exchange

	switch Exchange {
	case "binance":
		FindCorrelationFromBinance(ws, Request, MongoDBclient)
		break

		// case "bitfinex":
		// 	FindCorrelationFromBitfinex(ws, Request, MongoDBclient)
		// 	break

		// case "bitmex":
		// 	FindCorrelationFromBitmex(ws, Request, MongoDBclient)
		// 	break
		// case "poloniex":
		// 	FindCorrelationFromPoloniex(ws, Request, MongoDBclient)
		// 	break
	}
}

//done
func FindCorrelationFromBinance(ws *websocket.Conn, Request *SocketRequest, MongoDBclient *mongo.Client) {

	Pairs := Request.Pairs
	//Days := Request.Days
	Price := Request.Price
	//Candle := Request.Candle

	switch Price {
	case "trade":
		var Correlations []float64

		i := 365

		//	for j := 0; j < 4; j++ {

		timeNow := time.Now().Unix()

		timegap := i * 24 * 60 * 60

		startTime := timeNow - int64(timegap)

		//------------------------------------------------
		var TradePrices1 []float64

		var TradePrices2 []float64

		List1Channel := make(chan bool, 1)
		List2Channel := make(chan bool, 1)

		go func() {

			BinanceTradesCollectionPair1 := MongoDBclient.Database("binance_trades").Collection(Pairs[0])
			log.Println("Connected to db of binanace " + Pairs[0])

			//	Count1, error1 := BinanceTradesCollectionPair1.Find(context.TODO(), bson.M{"time": bson.M{"$gt": startTime}}).Count()
			Cursor1, err1 := BinanceTradesCollectionPair1.Find(context.TODO(), bson.M{"time": bson.M{"$gt": startTime}})
			FatalOnError(err1, "Error getting cursor")

			var Count1 int = 0

			for Cursor1.Next(context.Background()) {
				Count1++
				result := &bson.D{}

				err := Cursor1.Decode(result)
				FatalOnError(err, "Error decoding result")

				for _, E := range *result {
					if E.Key == "price" {

						TradePrice, errr := strconv.ParseFloat(E.Value.(string), 64)
						FatalOnError(errr, "Error parsing float64 from String")
						TradePrices1 = append(TradePrices1, TradePrice)
					}
				}

			}

			log.Printf("Found %d records of %s", Count1, Pairs[0])

			List1Channel <- true
		}()

		go func() {

			BinanceTradesCollectionPair2 := MongoDBclient.Database("binance_trades").Collection(Pairs[1])
			log.Println("Connected to db of binanace " + Pairs[1])

			Cursor2, err1 := BinanceTradesCollectionPair2.Find(context.TODO(), bson.M{"time": bson.M{"$gt": startTime}})
			FatalOnError(err1, "oops")

			var Count2 int = 0

			for Cursor2.Next(context.Background()) {
				result := &bson.D{}

				err := Cursor2.Decode(result)
				FatalOnError(err, "Error decoding result")

				for _, E := range *result {
					if E.Key == "price" {
						TradePrice, errr := strconv.ParseFloat(E.Value.(string), 64)
						FatalOnError(errr, "Error parsing float64 from String")
						TradePrices2 = append(TradePrices2, TradePrice)
					}
				}
				Count2++

			}

			log.Printf("Found %d records of %s", Count2, Pairs[1])

			List2Channel <- true

		}()

		<-List1Channel
		<-List2Channel

		correlation, err := stats.Correlation(TradePrices1, TradePrices2)
		FatalOnError(err, "Error finding correlation")

		Correlations = append(Correlations, correlation)

		// if i == 7 {
		// 	i = 10
		// }
		// if i == 10 {
		// 	i = 15
		// }
		// if i == 15 {
		// 	i = 30
		// }
		// if i == 30 {

		// }

		//}

		//do something with Correlations

		for _, Correlation := range Correlations {
			log.Println(Correlation)
		}

		break

	case "close":
		break
	}

}

//done
// func FindCorrelationFromBitfinex(ws *websocket.Conn, Request *SocketRequest, MongoDBclient *mongo.Client) {

// 	Pairs := Request.Pairs
// 	Days := Request.Days
// 	Candle := Request.Candle

// 	switch Request.Price {

// 	case "trade":

// 		var Correlations []float64

// 		i := 7

// 		for j := 0; j < 4; j++ {

// 			timeNow := time.Now().Unix()

// 			timegap := i * 24 * 60 * 60

// 			startTime := timeNow - int64(timegap)

// 			//------------------------------------------------
// 			var TradePrices1 []float64

// 			var TradePrices2 []float64

// 			List1Channel := make(chan bool, 1)
// 			List2Channel := make(chan bool, 1)

// 			go func() {

// 				BinanceTradesCollectionPair1 := MongoDBclient.Database("bitfinex_trades").Collection(Pairs[0])
// 				log.Println("Connected to db of bitfinex " + Pairs[0])

// 				Cursor1, err1 := BinanceTradesCollectionPair1.Find(context.TODO(), bson.M{"timestamp": bson.M{"$gt": startTime}})
// 				FatalOnError(err1, "Error getting Cursor")

// 				for Cursor1.Next(context.Background()) {
// 					result := &bson.D{}

// 					err := Cursor1.Decode(result)
// 					FatalOnError(err, "Error decoding result")

// 					for _, E := range *result {
// 						if E.Key == "price" {

// 							TradePrice, errr := strconv.ParseFloat(E.Value.(string), 64)
// 							FatalOnError(errr, "Error parsing float64 from String")
// 							TradePrices1 = append(TradePrices1, TradePrice)
// 						}
// 					}

// 				}

// 				List1Channel <- true
// 			}()

// 			go func() {

// 				BinanceTradesCollectionPair2 := MongoDBclient.Database("bitfinex_trades").Collection(Pairs[1])
// 				log.Println("Connected to db of binanace " + Pairs[1])

// 				Cursor2, err1 := BinanceTradesCollectionPair2.Find(context.TODO(), bson.M{"timestamp": bson.M{"$gt": startTime}})
// 				FatalOnError(err1, "oops")

// 				for Cursor2.Next(context.Background()) {
// 					result := &bson.D{}

// 					err := Cursor2.Decode(result)
// 					FatalOnError(err, "Error decoding result")

// 					for _, E := range *result {
// 						if E.Key == "price" {
// 							TradePrice, errr := strconv.ParseFloat(E.Value.(string), 64)
// 							FatalOnError(errr, "Error parsing float64 from String")
// 							TradePrices2 = append(TradePrices2, TradePrice)
// 						}
// 					}

// 				}
// 				List2Channel <- true

// 			}()

// 			<-List1Channel
// 			<-List2Channel

// 			correlation, err := stats.Correlation(TradePrices1, TradePrices2)
// 			FatalOnError(err, "Error finding correlation")

// 			Correlations := append(Correlations, correlation)

// 			if i == 7 {
// 				i = 10
// 			}
// 			if i == 10 {
// 				i = 15
// 			}
// 			if i == 15 {
// 				i = 30
// 			}
// 			if i == 30 {

// 			}

// 		}

// 		//do something with Correlations

// 		break

// 	case "close":
// 		break
// 	}

// }

// //not done
// func FindCorrelationFromBitmex(ws *websocket.Conn, Request *SocketRequest, MongoDBclient *mongo.Client) {

// 	Pairs := Request.Pairs
// 	Days := Request.Days
// 	Price := Request.Price
// 	Candle := Request.Candle

// 	switch Request.Price {
// 	case "trade":
// 		var Correlations []float64

// 		i := 7

// 		for j := 0; j < 4; j++ {

// 			timeNow := time.Now().Unix()

// 			timegap := i * 24 * 60 * 60

// 			startTime := timeNow - int64(timegap)

// 			//------------------------------------------------
// 			var TradePrices1 []float64

// 			var TradePrices2 []float64

// 			List1Channel := make(chan bool, 1)
// 			List2Channel := make(chan bool, 1)

// 			go func() {

// 				BinanceTradesCollectionPair1 := MongoDBclient.Database("bitmex_trades").Collection(Pairs[0])
// 				log.Println("Connected to db of bitfinex " + Pairs[0])

// 				Cursor1, err1 := BinanceTradesCollectionPair1.Find(context.TODO(), bson.M{"time": bson.M{"$gt": startTime}})
// 				FatalOnError(err1, "Error getting Cursor")

// 				for Cursor1.Next(context.Background()) {
// 					result := &bson.D{}

// 					err := Cursor1.Decode(result)
// 					FatalOnError(err, "Error decoding result")

// 					for _, E := range *result {
// 						if E.Key == "price" {

// 							TradePrice, errr := strconv.ParseFloat(E.Value.(string), 64)
// 							FatalOnError(errr, "Error parsing float64 from String")
// 							TradePrices1 = append(TradePrices1, TradePrice)
// 						}
// 					}

// 				}

// 				List1Channel <- true
// 			}()

// 			go func() {

// 				BinanceTradesCollectionPair2 := MongoDBclient.Database("bitmex_trades").Collection(Pairs[1])
// 				log.Println("Connected to db of binanace " + Pairs[1])

// 				Cursor2, err1 := BinanceTradesCollectionPair2.Find(context.TODO(), bson.M{"time": bson.M{"$gt": startTime}})
// 				FatalOnError(err1, "oops")

// 				for Cursor2.Next(context.Background()) {
// 					result := &bson.D{}

// 					err := Cursor2.Decode(result)
// 					FatalOnError(err, "Error decoding result")

// 					for _, E := range *result {
// 						if E.Key == "price" {
// 							TradePrice, errr := strconv.ParseFloat(E.Value.(string), 64)
// 							FatalOnError(errr, "Error parsing float64 from String")
// 							TradePrices2 = append(TradePrices2, TradePrice)
// 						}
// 					}

// 				}
// 				List2Channel <- true

// 			}()

// 			<-List1Channel
// 			<-List2Channel

// 			correlation, err := stats.Correlation(TradePrices1, TradePrices2)
// 			FatalOnError(err, "Error finding correlation")

// 			Correlations := append(Correlations, correlation)

// 			if i == 7 {
// 				i = 10
// 			}
// 			if i == 10 {
// 				i = 15
// 			}
// 			if i == 15 {
// 				i = 30
// 			}
// 			if i == 30 {

// 			}

// 		}

// 		//do something with Correlations
// 		break

// 	case "close":
// 		break
// 	}

// }

// //doubt
// func FindCorrelationFromPoloniex(ws *websocket.Conn, Request *SocketRequest, MongoDBclient *mongo.Client) {

// 	Pairs := Request.Pairs
// 	Days := Request.Days
// 	Price := Request.Price
// 	Candle := Request.Candle
// 	var Correlations []float64

// 	i := 7

// 	for j := 0; j < 4; j++ {

// 		timeNow := time.Now().Unix()

// 		timegap := i * 24 * 60 * 60

// 		startTime := timeNow - int64(timegap)

// 		StartTimeFormatted := time.Unix(0, 0)

// 		FormattedTime := StartTimeFormatted.Format("2006-01-02 15:04:05")

// 		//------------------------------------------------
// 		var TradePrices1 []float64

// 		var TradePrices2 []float64

// 		List1Channel := make(chan bool, 1)
// 		List2Channel := make(chan bool, 1)

// 		go func() {

// 			BinanceTradesCollectionPair1 := MongoDBclient.Database("poloniex_trades").Collection(Pairs[0])
// 			log.Println("Connected to db of bitfinex " + Pairs[0])

// 			Cursor1, err1 := BinanceTradesCollectionPair1.Find(context.TODO(), bson.M{"this.Date > " + FormattedTime})
// 			FatalOnError(err1, "Error getting Cursor")

// 			for Cursor1.Next(context.Background()) {
// 				result := &bson.D{}

// 				err := Cursor1.Decode(result)
// 				FatalOnError(err, "Error decoding result")

// 				for _, E := range *result {
// 					if E.Key == "price" {

// 						TradePrice, errr := strconv.ParseFloat(E.Value.(string), 64)
// 						FatalOnError(errr, "Error parsing float64 from String")
// 						TradePrices1 = append(TradePrices1, TradePrice)
// 					}
// 				}

// 			}

// 			List1Channel <- true
// 		}()

// 		go func() {

// 			BinanceTradesCollectionPair2 := MongoDBclient.Database("poloniex_trades").Collection(Pairs[1])
// 			log.Println("Connected to db of binanace " + Pairs[1])

// 			Cursor2, err1 := BinanceTradesCollectionPair2.Find(context.TODO(), bson.M{"this.Date > " + FormattedTime})
// 			FatalOnError(err1, "oops")

// 			for Cursor2.Next(context.Background()) {
// 				result := &bson.D{}

// 				err := Cursor2.Decode(result)
// 				FatalOnError(err, "Error decoding result")

// 				for _, E := range *result {
// 					if E.Key == "price" {
// 						TradePrice, errr := strconv.ParseFloat(E.Value.(string), 64)
// 						FatalOnError(errr, "Error parsing float64 from String")
// 						TradePrices2 = append(TradePrices2, TradePrice)
// 					}
// 				}

// 			}
// 			List2Channel <- true

// 		}()

// 		<-List1Channel
// 		<-List2Channel

// 		correlation, err := stats.Correlation(TradePrices1, TradePrices2)
// 		FatalOnError(err, "Error finding correlation")

// 		Correlations := append(Correlations, correlation)

// 		if i == 7 {
// 			i = 10
// 		}
// 		if i == 10 {
// 			i = 15
// 		}
// 		if i == 15 {
// 			i = 30
// 		}
// 		if i == 30 {

// 		}

// 	}
// }
