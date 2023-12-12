package main

import (
	"context"
	"fmt"

	"github.com/evenboee/ams-generator/generator"
)

func guard(err error) {
	if err != nil {
		panic(fmt.Errorf("%T %w", err, err))
	}
}

func require[T any](d T, err error) T {
	guard(err)
	return d
}

func main() {
	config := require(generator.LoadConfig("config.yaml"))

	guard(generator.Generate(context.Background(), nil, config))
}
