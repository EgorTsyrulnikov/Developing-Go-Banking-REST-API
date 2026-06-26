package main

import (
	"fmt"
	"bankapi/pkg/cbr"
)

func main() {
	rate, err := cbr.GetCentralBankRate()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Success! Rate with margin: %.2f%%\n", rate)
	}
}
