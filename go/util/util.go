package util

func GetEntryBool(m map[string]any, key string, defaultValue bool) bool {
	itf, found := m[key]
	if !found {
		return defaultValue
	}
	b, castOk := itf.(bool)
	if !castOk {
		return defaultValue
	}
	return b
}

func GetEntryString(m map[string]any, key string, defaultValue string) string {
	itf, found := m[key]
	if !found {
		return defaultValue
	}
	s, castOk := itf.(string)
	if !castOk {
		return defaultValue
	}
	return s
}
