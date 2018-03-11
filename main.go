package main

import (
	"fmt"

	"gitlab.com/etherlabs/pkg/errors"
)

func main() {
	err := errors.New(errors.IO)
	fmt.Println(err)
}
