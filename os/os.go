package os

import (
	"fmt"
	"os"
	"strconv"
)

/*
GetEnvAsInt takes the name of the environment variable as the first parameter. If 
the environment variable is found and the value is of type integer, the value is 
returned. Otherwise, 0 is returned with an error. If the environment variable is 
not found, the second parameter is used for a default value. If the second 
parameter is not set, an error is returned.
*/
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

/*
GetEnvAsInt takes the name of the environment variable as the first parameter. 
If the environment variable is found and the value is of type boolean, the value 
is returned. Otherwise, 0 is returned with an error. If the environment variable 
is not found, the second parameter is used for a default value. If the second 
parameter is not set, an error is returned.
*/
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
