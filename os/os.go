package os

import (
	"fmt"
	"os"
	"strconv"
)

func GetEnvAsInt(varName string, params ...int) (int, error) {
	val, ok := os.LookupEnv(varName)
	if !ok {
		if len(params) == 0 {
			return 0, fmt.Errorf("%s is not set", varName)
		}
		return params[0], nil
	}
	return strconv.Atoi(val)
}

func GetEnvAsBool(varName string, params ...bool) (bool, error) {
	val, ok := os.LookupEnv(varName)
	if !ok {
		if len(params) == 0 {
			return false, fmt.Errorf("%s is not set", varName)
		}
		return params[0], nil
	}
	return strconv.ParseBool(val)
}
