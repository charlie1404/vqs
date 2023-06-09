package utils

import (
	"bytes"
	"net/url"
)

func ParseUrlEncodedBodyParamKV(kv []byte, sep byte) (key, value string) {
	pos := bytes.IndexByte(kv, sep)
	if pos < 0 {
		key, _ = url.QueryUnescape(string(kv))
		value = ""
	} else {
		key, _ = url.QueryUnescape(string(kv[:pos]))
		value, _ = url.QueryUnescape(string(kv[pos+1:]))
	}

	return
}
