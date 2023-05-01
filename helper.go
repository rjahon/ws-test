package main

import "strconv"

func Parse(key string, v any) (Order, error) {
	data := v.(map[string]interface{})[key].([]interface{})

	amount, err := strconv.ParseFloat(data[0].(string), 64)
	if err != nil {
		return Order{}, err
	}
	price, err := strconv.ParseFloat(data[1].(string), 64)
	if err != nil {
		return Order{}, err
	}

	return Order{
		amount,
		price,
	}, nil
}
