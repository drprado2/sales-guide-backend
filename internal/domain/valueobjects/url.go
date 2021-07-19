package valueobjects

import "net/url"

type (
	Uri string
)

func (u *Uri) AsString() string {
	return (string)(*u)
}

func NewUri(value string) (*Uri, error) {
	_, err := url.ParseRequestURI(value)
	if err != nil {
		return nil, err
	}
	return (*Uri)(&value), nil
}
