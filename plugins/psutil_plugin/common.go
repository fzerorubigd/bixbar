package main

func getString(cfg map[string]interface{}, key string, def string) string {
	v, ok := cfg[key].(string)
	if ok {
		return v
	}

	return def
}
