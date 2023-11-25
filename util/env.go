package util

import "os"

func ByField(key string) string {
	return os.Getenv(key)
}
