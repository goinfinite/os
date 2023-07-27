package cliHelper

import (
	"encoding/json"
	"fmt"
)

type formattedResponse struct {
	Status bool        `json:"status"`
	Body   interface{} `json:"body"`
}

func ResponseWrapper(
	responseStatus bool,
	responseBody interface{},
) {
	formattedResponse := formattedResponse{
		Status: responseStatus,
		Body:   responseBody,
	}

	jsonResponse, err := json.MarshalIndent(formattedResponse, "", "  ")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(jsonResponse))
}
