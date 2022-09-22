package trace

import (
	"bytes"
	"runtime"
	"strconv"
)

const EnableGetGoroutineID = true

// https://blog.sgmansfield.com/2015/12/goroutine-ids/
func __caution__GetGoroutineID() uint64 {
	if EnableGetGoroutineID {
		b := make([]byte, 64)
		b = b[:runtime.Stack(b, false)]
		b = bytes.TrimPrefix(b, []byte("goroutine "))
		b = b[:bytes.IndexByte(b, ' ')]
		n, _ := strconv.ParseUint(string(b), 10, 64)
		return n
	}
	return 0
}
