package apps

func JsonData(status int, message string, data any) map[string]any {
	return map[string]any{
		"status":  status,
		"message": message,
		"data":    data,
	}
}
