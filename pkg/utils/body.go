package utils

import (
	"net/url"
	"strconv"
)

func GetFormValueString(form url.Values, key string) string {
	if vals := form[key]; len(vals) > 0 {
		return vals[0]
	}
	return ""
}

func GetFormValueUint(form url.Values, key string, defaultVal uint) uint {
	val := ""
	if vals := form[key]; len(vals) > 0 {
		val = vals[0]
	}

	if intVal, err := strconv.Atoi(val); err != nil {
		return defaultVal
	} else {
		return uint(intVal)
	}
}
