package data

import (
	"encoding/xml"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"net/http"
	"strconv"
)

type ExchangeRates struct {
	log hclog.Logger
	rates map[string]float64
}

func NewRate( l hclog.Logger) (*ExchangeRates, error)  {
	er := &ExchangeRates{log: l, rates: map[string]float64{}}
	err := er.getRates()
	return er, err
}

func (e *ExchangeRates) GetRate(base, dest string) (float64,  error) {
	br, ok := e.rates[base]
	if !ok {
		return 0, fmt.Errorf("Rate not found for currency %s", base)
	}

	dr, ok := e.rates[dest]
	if !ok {
		return 0, fmt.Errorf("Rate not found for currency %s", dest)
	}
	return dr / br, nil
}

func (e *ExchangeRates) getRates() error {
	resp, err := http.DefaultClient.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("EXpected error code 200 got %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	md := &Cubes{}
	_ = xml.NewDecoder(resp.Body).Decode(&md)
	
	for _, c := range md.CubeData {
		rate, err := strconv.ParseFloat(c.Rate, 64)
		if err != nil {
			return err
		}
		e.rates[c.Currency] = rate
	}
	e.rates["EUR"] = 1 // EUR based currency
	return nil
}

type Cube struct {
	Currency string `xml:"currency,attr"`
	Rate string `xml:"rate,attr"`
}

type Cubes struct {
	CubeData []Cube `xml:"Cube>Cube>Cube"`
}