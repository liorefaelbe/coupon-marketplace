package models

type Coupon struct {
	ProductID        string  `json:"product_id"`
	CostPrice        float64 `json:"cost_price"`
	MarginPercentage float64 `json:"margin_percentage"`
	MinimumSellPrice float64 `json:"minimum_sell_price"`
	IsSold           bool    `json:"is_sold"`
	ValueType        string  `json:"value_type"`
	Value            string  `json:"value"`
}
