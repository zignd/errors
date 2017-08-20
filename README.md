# errors [![GoDoc](https://godoc.org/github.com/zignd/errors?status.svg)](https://godoc.org/github.com/zignd/errors) [![Report card](https://goreportcard.com/badge/github.com/zignd/errors)](https://goreportcard.com/report/github.com/zignd/errors)     

An errors package that will help you handle them gracefully. It allows you to add contextual information to your errors, to wrap them and they even get a stack trace. Inspired by the [github.com/pgk/errors](https://www.github.com/pgk/errors) package and Node.js' [verror](https://github.com/joyent/node-verror) module.

# Features

* Add contextual information to error values preventing long and hard to read error messages
* Wrap existing error values into new ones
* Stack traces
* MultiError, wrap multiple errors values into a single one; great for concurrent workflows that may generate multiple errors
* Pretty print of the whole error value and support JSON marshalling to ease the serialization (check the ["Quick demo"](https://github.com/zignd/errors#quick-demo) section)

# Documentation

For a better understanding of the features provided by the package check the documentation at: [godoc.org/github.com/zignd/errors](https://godoc.org/github.com/zignd/errors)

# Quick demo

**Consider the following usages of the package in the functions defined below**

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/zignd/errors"
)

func foo() error {
	model := "iop-40392"

	if err := launch(model); err != nil {
		return errors.Wrapc(err, map[string]interface{}{
			"model": model,
		}, "failed to launch rocket")
	}

	return nil
}

func launch(model string) error {
	return errors.Errorc(map[string]interface{}{
		"rocket": map[string]interface{}{
			"ID":        "123",
			"Fuel":      10,
			"AutoPilot": true,
		},
	}, "something catastrofic just happened to rocket #123")
}
```

**JSON marshalling an error value**

```go
func main() {
	if err := foo(); err != nil {
		// Type assertions using the exposed Error type
		if err, ok := err.(*errors.Error); ok {
			b, _ := json.MarshalIndent(err, "", "\t")
			fmt.Printf("%s", b)
		}
	}
}
```

**Output**

```
{
        "Message": "failed to launch rocket",
        "Context": {
            "model": "iop-40392"
        },
        "Stack": "main.foo\n\t/home/zignd/go/src/github.com/zignd/test/main.go:38\nmain.main\n\t/home/zignd/go/src/github.com/zignd/test/main.go:11\nruntime.main\n\t/usr/local/go/src/runtime/proc.go:194\nruntime.goexit\n\t/usr/local/go/src/runtime/asm_amd64.s:2198",
        "Cause": {
            "Message": "something catastrofic just happened to rocket #123",
            "Context": {
                "rocket": {
                    "AutoPilot": true,
                    "Fuel": 10,
                    "ID": "123"
                }
            },
            "Stack": "main.launch\n\t/home/zignd/go/src/github.com/zignd/test/main.go:51\nmain.foo\n\t/home/zignd/go/src/github.com/zignd/test/main.go:35\nmain.main\n\t/home/zignd/go/src/github.com/zignd/test/main.go:11\nruntime.main\n\t/usr/local/go/src/runtime/proc.go:194\nruntime.goexit\n\t/usr/local/go/src/runtime/asm_amd64.s:2198",
            "Cause": null
        }
    }
```

**`fmt.Formatter` implementation supporting the `+v` format for recursive pretty print of the whole error value**

```go
func main() {
	if err := foo(); err != nil {
		if err, ok := err.(*errors.Error); ok {
			fmt.Printf("%+v", err)
		}
	}
}
```

**Output**

```
    Message:
        "failed to launch rocket"
    Context:
        model: iop-40392
    Stack:
        main.foo
            /home/zignd/go/src/github.com/zignd/test/main.go:38
        main.main
            /home/zignd/go/src/github.com/zignd/test/main.go:11
        runtime.main
            /usr/local/go/src/runtime/proc.go:194
        runtime.goexit
            /usr/local/go/src/runtime/asm_amd64.s:2198
    Cause:
        Message:
            "something catastrofic just happened to rocket #123"
        Context:
            rocket: map[ID:123 Fuel:10 AutoPilot:true]
        Stack:
            main.launch
                /home/zignd/go/src/github.com/zignd/test/main.go:51
            main.foo
                /home/zignd/go/src/github.com/zignd/test/main.go:35
            main.main
                /home/zignd/go/src/github.com/zignd/test/main.go:11
            runtime.main
                /usr/local/go/src/runtime/proc.go:194
            runtime.goexit
                /usr/local/go/src/runtime/asm_amd64.s:2198
```

**The usual `%s` format support**

```go
func main() {
	if err := foo(); err != nil {
		if err, ok := err.(*errors.Error); ok {
			fmt.Printf("%s", err)
		}
	}
}
```
**Output**

```
failed to launch rocket: something catastrofic just happened to rocket #123
```
