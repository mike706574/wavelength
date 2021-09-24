package main

import (
	"encoding/json"
	"fmt"
)

func PrintJson(label string, value interface{}) {
	json, err := json.MarshalIndent(value, "", "  ")

	if err != nil {
		panic("failed to marshal")
	}

	fmt.Printf("%s: %s\n", label, string(json))
}

func PrintStruct(label string, value interface{}) {
	fmt.Printf("%s: %+v\n", label, value)
}
