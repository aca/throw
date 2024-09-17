package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/xtdlib/trace"
)

func f1() (err error) {
	_, err = os.Open("non-existing-file")

	// wrap with std fmt.Errorf
	err = fmt.Errorf("f1: %w", err)

	// wrap with trace
	// err = trace.Wrap(err)
	err = trace.Errorf("wrapped %w", err)
	return
}

func f2() (error) {
	return trace.Wrap1(os.Open("non-existing-file"))
}

func f3() (error) {
	return nil
}


func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	err := f2()
	if err != nil {
		slog.Error("something is wrong", trace.SlogAttr(err))
	}

	slog.Info(fmt.Sprintf("this should print true: %v", errors.Is(err, os.ErrNotExist)))
}
