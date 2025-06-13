package main

import (
	// "context"
	// "ton-utils-go/internal/app"
	"ton-utils-go/internal/wallet"
	// "ton-utils-go/internal/scanner"
)

func main() {
	// if err := app.InitApp(); err != nil {
	// 	panic(err)
	// }

	if err := wallet.StartTracking(); err != nil {
		panic(err)
	}

	// scanner, err := scan.NewScanner()
	// if err != nil {
	// 	panic(err)
	// }
	// scanner.Listen()

}

