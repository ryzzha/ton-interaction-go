package main

import (
	"ton-utils-go/internal/app"
	"ton-utils-go/internal/stonfi_dedust"

	"ton-utils-go/internal/storage"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}

	// if err := wallet.StartTracking(); err != nil {
	// 	panic(err)
	// }

	// if err := nft.NftActions(); err != nil {
	// 	panic(err)
	// }

	if err := stonfi_dedust.DedustSwap(); err != nil {
		panic(err)
	}

	// if err := stonfi_dedust.StonfiSwap(); err != nil {
	// 	panic(err)
	// }

	// scanner, err := scan.NewScanner()
	// if err != nil {
	// 	panic(err)
	// }
	// scanner.Listen()


}

func run() error {

	if err := app.InitApp(); err != nil {
		return err
	}

	app.DB.AutoMigrate(
		&storage.Block{},
		&storage.NftCollection{},
		&storage.NftItem{},
	)

	return nil
}

