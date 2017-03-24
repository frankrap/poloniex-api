package publicapi

import (
	"encoding/json"
	"fmt"
	"strings"
)

type LoanOrders struct {
	Offers  []LoanOrder `json:"offers"`
	Demands []LoanOrder `json:"demands"`
}

type LoanOrder struct {
	Rate     float64 `json:"rate,string"`
	Amount   float64 `json:"amount,string"`
	RangeMin int     `json:"rangeMin"`
	RangeMax int     `json:"rangeMax"`
}

func (client *PublicClient) GetLoanOrders(currency string) (*LoanOrders, error) {

	params := map[string]string{
		"command":  "returnLoanOrders",
		"currency": strings.ToUpper(currency),
	}

	url := buildUrl(params)

	resp, err := client.do("GET", url, "", false)
	if err != nil {
		return nil, fmt.Errorf("get: %v", err)
	}

	res := LoanOrders{}

	if err := json.Unmarshal(resp, &res); err != nil {
		return nil, fmt.Errorf("json unmarshal: %v", err)
	}

	return &res, nil
}
