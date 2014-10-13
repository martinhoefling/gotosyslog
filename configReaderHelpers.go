package main

func readBool(j map[string]interface{}, key string, def bool) bool {
	if val, ok := j[key]; ok {
		return val.(bool)
	} else {
		return def
	}
}

func readString(j map[string]interface{}, key string, def string) string {
	if val, ok := j[key]; ok {
		return val.(string)
	} else {
		return def
	}
}

func readStringSlice(j map[string]interface{}, key string, def []string) []string {
	if val, ok := j[key]; ok {
		strings := make([]string, 0)
		for _, v := range val.([]interface{}) {
			strings = append(strings, v.(string))
		}
		return strings
	} else {
		return def
	}
}

type readconfig func(map[string]interface{}) interface{}
type defaultconfig func() interface{}

func readMap(j map[string]interface{}, key string, read readconfig, def defaultconfig) interface{} {
	if val, ok := j[key]; ok {
		return read(val.(map[string]interface{}))
	} else {
		return def()
	}
}
