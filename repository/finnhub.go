//Finnhub adapatar implements PriceProviderRepository interface
// it is a source adapter for the Finnhub Websocket API
package repository

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/razagr/pensionera/domain"
	"github.com/razagr/pensionera/service"
)

// CurrencyJSON saves slice of Currency and its type
type CurrencyJSON struct {
	Data []domain.Currency `json:"data"`
	Type string            `json:"type"`
}

// repo implements the PriceProviderRepository interface
type repo struct {

	// currency symbols
	symbols map[string]float32

	// APIKey is used to authenticate the websocket connection
	APIKey string
}

// NewFinnHubRepository is a constructor for the Finnhub adapter
func NewFinnHubRepository(

	symbols map[string]float32,
	APIKey string) service.PriceProviderRepository {
	return &repo{symbols, APIKey}
}

// Run starts the websocket connection and calls the subscribe and start functions
func (r *repo) Run(channels map[string]chan domain.Currency) error {
	w, _, err := websocket.DefaultDialer.Dial("wss://ws.finnhub.io?token="+r.APIKey, nil)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	// subscribe to the websocket
	r.subscribe(w)

	// start listening to the websocket
	r.startListening(w, channels)
	return nil
}

// subscribe to the websocket
func (r *repo) subscribe(w *websocket.Conn) {
	for s := range r.symbols {
		fmt.Println("Subscribe to ", s)
		msgReceived, _ := json.Marshal(map[string]interface{}{"type": "subscribe", "symbol": s})
		w.WriteMessage(websocket.TextMessage, msgReceived)
	}
}

// start listening to the websocket and passes the pricing data to the CurrencyServices
func (r *repo) startListening(w *websocket.Conn, channels map[string]chan domain.Currency) {
	for {

		// respone will save the websocket JSON data
		var respone CurrencyJSON
		err := w.ReadJSON(&respone)
		if err != nil {
			panic(err)
		}
		// as per challenge requirement, we only need to process the data if it is a type of trade
		if respone.Type == "trade" {

			// TODO: there are duplicate entries with same symbol, Only difference
			// in the map is the Volume value,  we should discuess if we need to process all of them or not

			// duplicateSymbol helps to avoid duplicate entries in the response
			duplicateSymbol := []string{}

			for _, curr := range respone.Data {

				// check if symbol already exist in the duplicateSymbol slice
				duplicate, _ := r.duplicate(curr.Symbol, duplicateSymbol)

				// if  symbol does not exist in the slice then add it to
				//  duplicateSymbol slice and also perform the add operation
				if !duplicate {
					duplicateSymbol = append(duplicateSymbol, curr.Symbol)

					// start a new goroutine
					go func() {
						// send the pricing data to its channel
						channels[curr.Symbol] <- curr
						if err != nil {
							fmt.Printf("Error: %v\n", err)
						}
					}()
				}
			}

		}
	}

}

// duplicate is required to avoid duplicate entries I was receivng in the websocket JSON response
func (r *repo) duplicate(val string, array []string) (ok bool, i int) {
	for i = range array {
		if ok = array[i] == val; ok {
			return
		}
	}
	return
}
