# errors [![GoDoc](https://godoc.org/github.com/zignd/errors?status.svg)](https://godoc.org/github.com/zignd/errors) [![Report card](https://goreportcard.com/badge/github.com/zignd/errors)](https://goreportcard.com/report/github.com/zignd/errors)     

An errors package that will help you handle them gracefully. It allows you to add additional data to your errors, to wrap it and you even get a stack trace. Inspired by the [github.com/pkg/errors](https://www.github.com/pkg/errors) package and Node.js' [verror](https://github.com/joyent/node-verror) module.

# Features

* Add additional data to error values preventing long and hard to read error messages
* Wrap existing error values into new ones
* Stack traces for each error value
* MultiError, wrap multiple errors values into a single one; great for concurrent workflows that may generate multiple errors
* Pretty print of the whole error value and support JSON marshalling to ease the serialization (check the ["Quick demo"](https://github.com/zignd/errors#quick-demo) section)

# Installation

```bash
go get -u github.com/zignd/errors
```

# Documentation

For a better understanding of the features provided by the package check the documentation at: [pkg.go.dev/github.com/zignd/errors](https://pkg.go.dev/github.com/zignd/errors)

# Quick demo

There's an example at `examples/example1/example1.go` that shows how to use the package.

<details><summary>Here's the code for the example:</summary>

```go
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
		fmt.Println("Error logged as a JSON structure using the json.MarshalIndent:")
		fmt.Printf("%s\n", b)

		b, _ = json.Marshal(err)
		fmt.Println("\nError logged as a JSON structure using the json.Marshal:")
		fmt.Printf("%s\n", b)

		fmt.Println("\nError logged using the s format specifier:")
		fmt.Printf("%s\n", err)

		fmt.Println("\nError logged using the +v format specifier:")
		fmt.Printf("%+v\n", err)
	}
}
```

</details>

<details><summary>Here's the execution of the example:</summary>

```
$ go run examples/example1/example1.go
Error logged as a JSON structure using the json.MarshalIndent:
[
  {
    "data": {
      "transactionId": "tx_123456",
      "userId": "67890"
    },
    "message": "failed to complete the transaction on bank_123456",
    "stack": [
      "main.createTransaction @ /root/hack/errors/examples/example1/example1.go:13",
      "main.main @ /root/hack/errors/examples/example1/example1.go:52",
      "runtime/internal/atomic.(*Uint32).Load @ /root/go/version/go1.21.0/src/runtime/internal/atomic/types.go:194",
      "runtime.goexit @ /root/go/version/go1.21.0/src/runtime/asm_amd64.s:1651"
    ]
  },
  {
    "data": {
      "operation": "update",
      "tableName": "transactions"
    },
    "message": "failed to update the database",
    "stack": [
      "main.updateDatabase @ /root/hack/errors/examples/example1/example1.go:24",
      "main.createTransaction @ /root/hack/errors/examples/example1/example1.go:12",
      "main.main @ /root/hack/errors/examples/example1/example1.go:52",
      "runtime/internal/atomic.(*Uint32).Load @ /root/go/version/go1.21.0/src/runtime/internal/atomic/types.go:194",
      "runtime.goexit @ /root/go/version/go1.21.0/src/runtime/asm_amd64.s:1651"
    ]
  },
  {
    "data": {
      "server": "db-server-01",
      "timeoutSeconds": 30
    },
    "message": "connection timeout",
    "stack": [
      "main.createConnection @ /root/hack/errors/examples/example1/example1.go:35",
      "main.updateDatabase @ /root/hack/errors/examples/example1/example1.go:23",
      "main.createTransaction @ /root/hack/errors/examples/example1/example1.go:12",
      "main.main @ /root/hack/errors/examples/example1/example1.go:52",
      "runtime/internal/atomic.(*Uint32).Load @ /root/go/version/go1.21.0/src/runtime/internal/atomic/types.go:194",
      "runtime.goexit @ /root/go/version/go1.21.0/src/runtime/asm_amd64.s:1651"
    ]
  },
  {
    "data": {
      "network": "internal",
      "severity": "high"
    },
    "message": "network instability detected",
    "stack": [
      "main.open @ /root/hack/errors/examples/example1/example1.go:45",
      "main.createConnection @ /root/hack/errors/examples/example1/example1.go:34",
      "main.updateDatabase @ /root/hack/errors/examples/example1/example1.go:23",
      "main.createTransaction @ /root/hack/errors/examples/example1/example1.go:12",
      "main.main @ /root/hack/errors/examples/example1/example1.go:52",
      "runtime/internal/atomic.(*Uint32).Load @ /root/go/version/go1.21.0/src/runtime/internal/atomic/types.go:194",
      "runtime.goexit @ /root/go/version/go1.21.0/src/runtime/asm_amd64.s:1651"
    ]
  }
]

Error logged as a JSON structure using the json.Marshal:
[{"data":{"transactionId":"tx_123456","userId":"67890"},"message":"failed to complete the transaction on bank_123456","stack":["main.createTransaction @ /root/hack/errors/examples/example1/example1.go:13","main.main @ /root/hack/errors/examples/example1/example1.go:52","runtime/internal/atomic.(*Uint32).Load @ /root/go/version/go1.21.0/src/runtime/internal/atomic/types.go:194","runtime.goexit @ /root/go/version/go1.21.0/src/runtime/asm_amd64.s:1651"]},{"data":{"operation":"update","tableName":"transactions"},"message":"failed to update the database","stack":["main.updateDatabase @ /root/hack/errors/examples/example1/example1.go:24","main.createTransaction @ /root/hack/errors/examples/example1/example1.go:12","main.main @ /root/hack/errors/examples/example1/example1.go:52","runtime/internal/atomic.(*Uint32).Load @ /root/go/version/go1.21.0/src/runtime/internal/atomic/types.go:194","runtime.goexit @ /root/go/version/go1.21.0/src/runtime/asm_amd64.s:1651"]},{"data":{"server":"db-server-01","timeoutSeconds":30},"message":"connection timeout","stack":["main.createConnection @ /root/hack/errors/examples/example1/example1.go:35","main.updateDatabase @ /root/hack/errors/examples/example1/example1.go:23","main.createTransaction @ /root/hack/errors/examples/example1/example1.go:12","main.main @ /root/hack/errors/examples/example1/example1.go:52","runtime/internal/atomic.(*Uint32).Load @ /root/go/version/go1.21.0/src/runtime/internal/atomic/types.go:194","runtime.goexit @ /root/go/version/go1.21.0/src/runtime/asm_amd64.s:1651"]},{"data":{"network":"internal","severity":"high"},"message":"network instability detected","stack":["main.open @ /root/hack/errors/examples/example1/example1.go:45","main.createConnection @ /root/hack/errors/examples/example1/example1.go:34","main.updateDatabase @ /root/hack/errors/examples/example1/example1.go:23","main.createTransaction @ /root/hack/errors/examples/example1/example1.go:12","main.main @ /root/hack/errors/examples/example1/example1.go:52","runtime/internal/atomic.(*Uint32).Load @ /root/go/version/go1.21.0/src/runtime/internal/atomic/types.go:194","runtime.goexit @ /root/go/version/go1.21.0/src/runtime/asm_amd64.s:1651"]}]

Error logged using the s format specifier:
failed to complete the transaction on bank_123456: failed to update the database: connection timeout: network instability detected

Error logged using the +v format specifier:
message:
        "failed to complete the transaction on bank_123456"
data:
        userId: 67890
        transactionId: tx_123456
stack:
        main.createTransaction @ /root/hack/errors/examples/example1/example1.go:13
        main.main @ /root/hack/errors/examples/example1/example1.go:52
        runtime/internal/atomic.(*Uint32).Load @ /root/go/version/go1.21.0/src/runtime/internal/atomic/types.go:194
        runtime.goexit @ /root/go/version/go1.21.0/src/runtime/asm_amd64.s:1651
cause:
        message:
                "failed to update the database"
        data:
                tableName: transactions
                operation: update
        stack:
                main.updateDatabase @ /root/hack/errors/examples/example1/example1.go:24
                main.createTransaction @ /root/hack/errors/examples/example1/example1.go:12
                main.main @ /root/hack/errors/examples/example1/example1.go:52
                runtime/internal/atomic.(*Uint32).Load @ /root/go/version/go1.21.0/src/runtime/internal/atomic/types.go:194
                runtime.goexit @ /root/go/version/go1.21.0/src/runtime/asm_amd64.s:1651
        cause:
                message:
                        "connection timeout"
                data:
                        server: db-server-01
                        timeoutSeconds: 30
                stack:
                        main.createConnection @ /root/hack/errors/examples/example1/example1.go:35
                        main.updateDatabase @ /root/hack/errors/examples/example1/example1.go:23
                        main.createTransaction @ /root/hack/errors/examples/example1/example1.go:12
                        main.main @ /root/hack/errors/examples/example1/example1.go:52
                        runtime/internal/atomic.(*Uint32).Load @ /root/go/version/go1.21.0/src/runtime/internal/atomic/types.go:194
                        runtime.goexit @ /root/go/version/go1.21.0/src/runtime/asm_amd64.s:1651
                cause:
                        message:
                                "network instability detected"
                        data:
                                severity: high
                                network: internal
                        stack:
                                main.open @ /root/hack/errors/examples/example1/example1.go:45
                                main.createConnection @ /root/hack/errors/examples/example1/example1.go:34
                                main.updateDatabase @ /root/hack/errors/examples/example1/example1.go:23
                                main.createTransaction @ /root/hack/errors/examples/example1/example1.go:12
                                main.main @ /root/hack/errors/examples/example1/example1.go:52
                                runtime/internal/atomic.(*Uint32).Load @ /root/go/version/go1.21.0/src/runtime/internal/atomic/types.go:194
                                runtime.goexit @ /root/go/version/go1.21.0/src/runtime/asm_amd64.s:1651
```

</details>
