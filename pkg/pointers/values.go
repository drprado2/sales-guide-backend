package pointers

func Bool(value bool) *bool {
	return &value
}

func String(value string) *string {
	return &value
}

func SafeString(value string) *string {
	newS := value
	return &newS
}
