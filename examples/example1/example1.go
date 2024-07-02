package main

import (
	"encoding/json"
	"fmt"

	"github.com/zignd/errors"
)

func createTransaction(id string) error {
	bank := "bank_123456"
	if err := updateDatabase(); err != nil {
		return errors.Wrapdf(err, errors.Data{
			"transactionId": id,
			"userId":        "67890",
		}, "failed to complete the transaction on %s", bank)
	}

	return nil
}

func updateDatabase() error {
	if err := createConnection(); err != nil {
		return errors.Wrapd(err, errors.Data{
			"tableName": "transactions",
			"operation": "update",
		}, "failed to update the database")
	}

	return nil
}

func createConnection() error {
	if err := open(); err != nil {
		return errors.Wrapd(err, errors.Data{
			"server":         "db-server-01",
			"timeoutSeconds": 30,
		}, "connection timeout")
	}

	return nil
}

func open() error {
	return errors.Errord(errors.Data{
		"network":  "internal",
		"severity": "high",
	}, "network instability detected")
}

func main() {
	if err := createTransaction("tx_123456"); err != nil {
		b, _ := json.MarshalIndent(err, "", "  ")
		fmt.Println("Error logged as a JSON structure using the JSON.MarshalIndent:")
		fmt.Printf("%s\n", b)

		b, _ = json.Marshal(err)
		fmt.Println("\nError logged as a JSON structure using the JSON.Marshal:")
		fmt.Printf("%s\n", b)

		fmt.Println("\nError logged using the s format specifier:")
		fmt.Printf("%s\n", err)

		fmt.Println("\nError logged using the +v format specifier:")
		fmt.Printf("%+v\n", err)
	}
}
