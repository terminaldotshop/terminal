package assert

import (
	"log/slog"
	"os"
)

var assertContext map[string]interface{} = map[string]interface{}{}

func AddContext(key string, value interface{}) {
    assertContext[key] = value
}

var lastKey string
func LastKeyPressed(key string) {
    lastKey = key
}

func Assert(truth bool, msg string) {
    if !truth {
        slog.Error(msg, "key", lastKey, "context", assertContext)
        os.Exit(1)
    }
}

