package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/aca/throw"
)

func f1() ( err error ) {
	_, err = os.Open("non-existing-file")

	// wrap with std fmt.Errorf
	err = fmt.Errorf("f1: %w", err)

	// wrap with throw
	err = throw.Wrap(err)
	return
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	err := f1()
	if err != nil {
		slog.Error(err.Error(), throw.SlogAttr(err))
	}

	slog.Info(fmt.Sprintf("errors.Is(err, os.ErrNotExist): %v", errors.Is(err, os.ErrNotExist)))
}
