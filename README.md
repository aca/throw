# trace

Minimal lib to wrap error with stacktrace. Works with std errors, slog packages.

example
```go
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
		slog.Error("mail fail", throw.SlogAttr(err))
	}

	slog.Info(fmt.Sprintf("errors.Is(err, os.ErrNotExist): %v", errors.Is(err, os.ErrNotExist)))
}
```

output
```json
{
  "time": "2024-02-22T16:35:30.029289212+09:00",
  "level": "ERROR",
  "msg": "mail fail",
  "throw": {
    "error": "f1: open non-existing-file: no such file or directory",
    "stack": [
      "main.f1:/home/rok/src/github.com/aca/throw/example/main.go:19",
      "main.main:/home/rok/src/github.com/aca/throw/example/main.go:26"
    ]
  }
}
{
  "time": "2024-02-22T16:35:30.029464254+09:00",
  "level": "INFO",
  "msg": "errors.Is(err, os.ErrNotExist): true"
}
```
