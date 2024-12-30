package utils

import (
	"database/sql"
	"slices"
	"songLibrary/customLog"
	"strings"
	"time"

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

// CompareMapsByStringKeys for map-arguments, checks the keys of the first argument that contain non-empty values
// ​​to see if they are present in the second argument.
func CompareMapsByStringKeys(map1, map2 map[string]string) bool {
	var resp bool
	len1 := len(map1)
	len2 := len(map2)
	if len1 == len2 {
		keysSlice1 := GetMapKeysWithValue(map1)
		keysSlice2 := GetMapKeysWithValue(map2)
		check := true
		for _, val := range keysSlice1 {
			if !slices.Contains(keysSlice2, val) {
				check = false
				break
			}
		}
		resp = check
	}
	return resp
}

// GetMapKeysWithValue returns from the argument map, a map with keys with non-empty values.
func GetMapKeysWithValue(mapArg map[string]string) []string {
	var resp []string
	if len(mapArg) > 0 {
		for key, val := range mapArg {
			if val != "" {
				resp = append(resp, key)
			}
		}
	}
	return resp
}

// GetMapValues from the passed map returns a slice with its non-empty values.
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

// GetIndexByStrValue returns the integer index of the passed value in the passed slice; if the value is missing, then -1.
func GetIndexByStrValue(data []string, value string) int {
	resp := -1
	for i, val := range data {
		if val == value {
			resp = i
			break
		}
	}
	return resp
}

// SqlToMap the values ​​of the passed response structure are returned by the database as a map.
func SqlToMap(rows *sql.Rows) []map[string]interface{} {
	resp := make([]map[string]interface{}, 0)
	columns, err := rows.Columns()
	if err != nil {
		customLog.Logging(err)
	} else {
		scanArgs := make([]interface{}, len(columns))
		values := make([]interface{}, len(columns))
		for i := range values {
			scanArgs[i] = &values[i]
		}
		for rows.Next() {
			err = rows.Scan(scanArgs...)
			if err != nil {
				customLog.Logging(err)
			}
			record := make(map[string]interface{})
			for i, col := range values {
				if col != nil {
					switch col.(type) {
					case bool:
						record[columns[i]] = col.(bool)
					case int:
						record[columns[i]] = col.(int)
					case int64:
						record[columns[i]] = col.(int64)
					case float64:
						record[columns[i]] = col.(float64)
					case string:
						record[columns[i]] = col.(string)
					case time.Time:
						record[columns[i]] = col.(time.Time)
					case []byte:
						record[columns[i]] = string(col.([]byte))
					default:
						record[columns[i]] = col
					}
				}
			}
			resp = append(resp, record)
		}
	}
	return resp
}

// GetMapKeys returns a slice of the keys of the passed map.
func GetMapKeys(argMap map[string]string) []string {
	resp := make([]string, len(argMap))
	var i int
	for k := range argMap {
		resp[i] = k
		i++
	}
	return resp
}

// PresenceMapKeysInOtherMap returns a Boolean answer whether the keys of the first passed card are contained in the second.
func PresenceMapKeysInOtherMap(map1, map2 map[string]string) bool {
	var resp bool
	keys1 := GetMapKeys(map1)
	keys2 := GetMapKeys(map2)
	check := true
	for _, val := range keys1 {
		if !slices.Contains(keys2, val) {
			check = false
			break
		}
	}
	resp = check
	return resp
}

// GetMapWithoutKeys returns the transferred map without the transferred key.
func GetMapWithoutKeys(map1 map[string]string, exceptKeys []string) map[string]string {
	resp := make(map[string]string, len(map1)-len(exceptKeys))
	for k, v := range map1 {
		if !slices.Contains(exceptKeys, k) {
			resp[k] = v
		}
	}
	return resp
}
