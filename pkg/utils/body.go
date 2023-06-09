package utils

import (
	"bytes"
)

func ParseUrlEncodedBodyParamKV(kv []byte, sep byte) (key, value string) {
	pos := bytes.IndexByte(kv, sep)
	if pos < 0 {
		key = string(kv)
		value = ""
	} else {
		key, value = string(kv[:pos]), string(kv[pos+1:])
	}

	return
}
