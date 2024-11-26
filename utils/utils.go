package utils

import (
	"reflect"
	"songLibrary/customLog"
	"strings"

	"github.com/joho/godotenv"
)

// GetConfFromEnvFile receives data for the database from the environment file. If successful, returns a non-empty map.
func GetConfFromEnvFile() map[string]string {
	resp := make(map[string]string)
	envFile, err := godotenv.Read(".env")
	if err == nil {
		resp = envFile
	} else {
		customLog.Logging(err)
	}
	return resp
}

// ConcatSlice returns a string from the elements of the passed slice with strings. Separator - space.
func ConcatSlice(strSlice []string) string {
	resp := ""
	if len(strSlice) > 0 {
		var strBuilder strings.Builder
		for _, val := range strSlice {
			strBuilder.WriteString(val)
		}
		resp = strBuilder.String()
		strBuilder.Reset()
	}
	return resp
}

func CompareMapsByStringKeys(map1, map2 map[string]string) bool {
	var resp bool
	len1 := len(map1)
	len2 := len(map2)
	if len1 == len2 {
		keysSlice1 := GetMapKeys(map1)
		keysSlice2 := GetMapKeys(map2)
		resp = reflect.DeepEqual(keysSlice1, keysSlice2)
	}
	return resp
}

func GetMapKeys(mapArg map[string]string) []string {
	var resp []string
	if len(mapArg) > 0 {
		for i := range mapArg {
			resp = append(resp, i)
		}
	}
	return resp
}

func GetMapValues(mapArg map[string]string) []string {
	var resp []string
	if len(mapArg) > 0 {
		for _, value := range mapArg {
			if value != "" {
				resp = append(resp, value)
			}
		}
	}
	return resp
}
