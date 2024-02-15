# trace

Minimal lib to wrap error with stacktrace. Works with std errors, slog packages.

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
		slog.Error(err.Error(), throw.SlogAttr(err))
	}

	slog.Info(fmt.Sprintf("errors.Is(err, os.ErrNotExist): %v", errors.Is(err, os.ErrNotExist)))
}
```

```json
{
  "time": "2024-02-15T20:23:37.393253477+09:00",
  "level": "ERROR",
  "msg": "f1: open non-existing-file: no such file or directory",
  "stack": [
    "main.f1:/home/rok/src/github.com/aca/throw/example/main.go:19",
    "main.main:/home/rok/src/github.com/aca/throw/example/main.go:26"
  ]
}
{
  "time": "2024-02-15T20:23:37.393368228+09:00",
  "level": "INFO",
  "msg": "errors.Is(err, os.ErrNotExist): true"
}
```
