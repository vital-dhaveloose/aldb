package ref

import "net/url"

type AttributeSetRef struct {
	ActivityRef
	AttributeSetId string
}

type ManifestRef struct {
	Id *url.URL
}
