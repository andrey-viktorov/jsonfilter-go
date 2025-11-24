package comparison

import (
	"unsafe"

	"github.com/tidwall/gjson"
)

// getJSONResult resolves a JSON path using gjson without the protective copies performed by GetBytes.
// gjson.GetBytes guarantees safety by copying substrings, but that creates allocations in hot paths. We
// convert the payload to a string via unsafe pointer casting, call gjson.Get, and only use the returned
// value within the scope of Evaluate, which keeps the access safe for read-only usage.
func getJSONResult(payload []byte, path string) gjson.Result {
	if len(payload) == 0 {
		return gjson.Result{}
	}
	jsonStr := *(*string)(unsafe.Pointer(&payload))
	return gjson.Get(jsonStr, path)
}
