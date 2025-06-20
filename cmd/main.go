package main

import (
	"ton-utils-go/internal/app"
	scan "ton-utils-go/internal/scanner"
	"ton-utils-go/internal/storage"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}

	// if err := wallet.StartTracking(); err != nil {
	// 	panic(err)
	// }

	scanner, err := scan.NewScanner()
	if err != nil {
		panic(err)
	}
	scanner.Listen()


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

