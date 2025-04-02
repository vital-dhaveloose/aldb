package mediatype

import "mime"

type MediaType struct {
	Type       string
	Parameters map[string]string
}

func MediaTypeMustParse(raw string) MediaType {
	mediatype, params, err := mime.ParseMediaType(raw)
	if err != nil {
		panic(err)
	}
	//TODO not correct!
	return MediaType{
		Type:       mediatype,
		Parameters: params,
	}
}
