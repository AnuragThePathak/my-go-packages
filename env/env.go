package env

import (
	"fmt"
	"os"
	"strconv"
)

/*
GetEnv takes the name of the environment variable as the first parameter. If 
the environment variable is found, the value is returned. If the environment 
variable is not found, the second parameter is used for a default value. If the 
second parameter is not set, an error is returned. You may choose to provide 
default value depending on your needs.
*/
func GetEnv(varName string, params ...string) (string, error) {
	val, ok := os.LookupEnv(varName)
	if !ok {
		if len(params) == 0 {
			return val, fmt.Errorf("%s is not set", varName)
		}
		return params[0], nil
	}
	return val, nil
}

/*
GetEnvAsInt takes the name of the environment variable as the first parameter. If 
the environment variable is found and the value is of type integer, the value is 
returned. If the environment variable is not found, the second parameter is used 
for a default value. If the second parameter is not set, an error is returned. You 
may choose to provide default value depending on your needs.
*/
func GetEnvAsInt(varName string, params ...int) (int, error) {
	val, ok := os.LookupEnv(varName)
	if !ok {
		if len(params) == 0 {
			return 0, fmt.Errorf("%s is not set", varName)
		}
		return params[0], nil
	}
	num, err := strconv.Atoi(val)
	if err != nil {
		return num, fmt.Errorf("%s can't be parsed as an integer", varName)
	}
	return num, nil
}

/*
GetEnvAsBool takes the name of the environment variable as the first parameter. 
If the environment variable is found and the value is of type boolean, the value 
is returned. If the environment variable is not found, the second parameter is 
used for a default value. If the second parameter is not set, an error is 
returned. You may choose to provide default value depending on your needs.
*/
func GetEnvAsBool(varName string, params ...bool) (bool, error) {
	val, ok := os.LookupEnv(varName)
	if !ok {
		if len(params) == 0 {
			return false, fmt.Errorf("%s is not set", varName)
		}
		return params[0], nil
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		return b, fmt.Errorf("%s can't be parsed as a boolean", varName)
	}
	return b, nil
}
