package convert

func NullIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func StringPtr(s string) *string {
	return &s
}
