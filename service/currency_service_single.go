// // CurrencyService is used to get the currency data from adapters (PriceProviderRepository)
// // and save it to the database (CurrencyRepository)
package service

// import (
// 	"fmt"
// 	"math"

// 	"github.com/razagr/pensionera/domain"
// )

// var (

// 	// channels is a map of channels for each currency
// 	// each country has it's own channel
// 	channels = map[string]chan domain.Currency{}

// 	// channelOpened keeps track of which channels are open
// 	channelOpened = map[string]bool{}
// )

// // service
// type service struct {

// 	// window size
// 	Window int

// 	// currency symbol
// 	Symbols map[string]float32

// 	// prices slice
// 	Prices map[string][]float64

// 	// priceSliceIndex tracks the index of the slice
// 	PriceSliceIndex map[string]int

// 	// storage is the repository for saving currency data
// 	storage CurrencyRepository
// }

// // NewService creates a new service for a currency
// func NewService(Window int, Symbols map[string]float32, storage CurrencyRepository) CurrencyService {
// 	Prices := make(map[string][]float64)
// 	PriceSliceIndex := make(map[string]int)

// 	for s := range Symbols {
// 		Prices[s] = make([]float64, Window)
// 		PriceSliceIndex[s] = 0
// 	}
// 	return &service{Window, Symbols, Prices, PriceSliceIndex, storage}
// }

// // AddPrice adds the price to the prices slice and calls GetAverage()
// // It also resets the priceSliceIndex to 0 if it reaches the window size
// func (s *service) AddPrice(currency domain.Currency) error {
// 	if math.IsNaN(currency.Price) && math.IsInf(currency.Price, 0) {
// 		return fmt.Errorf("Invalid price: %f", currency.Price)
// 	}

// 	// get current currncy coming from price provider and assign its prices to new slice
// 	prices := s.Prices[currency.Symbol]
// 	priceSliceIndex := s.PriceSliceIndex[currency.Symbol]

// 	// add the new price to the prices slice
// 	prices[priceSliceIndex] = currency.Price

// 	// our slice is full, lets calculate the moving average and reset the slice
// 	// we have to subtract 1 from the window because we are using a circular buffer
// 	if priceSliceIndex >= (s.Window - 1) {
// 		avg := s.GetAverage(currency)

// 		// reset the slice index for circular buffer
// 		priceSliceIndex = 0
// 		// reset the prices slice values to 0
// 		prices = make([]float64, s.Window)

// 		// starts a goroutine to save the average to the database
// 		go func() {
// 			fmt.Printf("-> Currency: %s: Window: %d:  Timestamp: %d: Average: %v\n", currency.Symbol, s.Window, currency.Time, avg)
// 			s.storage.Save(currency, avg)
// 		}()
// 	} else {
// 		// increase the index of the slice
// 		priceSliceIndex++
// 	}

// 	s.Prices[currency.Symbol] = prices
// 	s.PriceSliceIndex[currency.Symbol] = priceSliceIndex
// 	fmt.Println(priceSliceIndex, "-> ", currency.Symbol, "-> ", currency.Price)
// 	return nil
// }

// // Returns the average of the prices slice divided by the window size
// func (s *service) GetAverage(currency domain.Currency) float64 {
// 	var sum float64 = 0

// 	// count the sum of the prices slice
// 	for _, price := range s.Prices[currency.Symbol] {
// 		sum += price
// 	}
// 	return sum / float64(s.Window)
// }

// // add retrived currency detail from websocket to the map using go channel
// func (s *service) AddToChannel(currency domain.Currency) error {

// 	// check if the channel is open for this currency
// 	if channelOpened[currency.Symbol] == false {
// 		// create a new channel for this currency
// 		channels[currency.Symbol] = make(chan domain.Currency)
// 		// set the channel as open
// 		channelOpened[currency.Symbol] = true
// 		fmt.Println("Channel opened for: ", currency.Symbol)
// 		// start a goroutine to listen to the channel
// 		go func(c <-chan domain.Currency) {

// 			// loop through channel until the channel is closed
// 			for cur := range c {
// 				// when channel sends a value then send it to the service for processing
// 				err := s.AddPrice(cur)
// 				if err != nil {
// 					panic(err)
// 				}
// 			}
// 			fmt.Println("Channel closed for: ", currency.Symbol)
// 		}(channels[currency.Symbol])
// 	}
// 	// send the currency to the channel
// 	channels[currency.Symbol] <- currency
// 	return nil
// }