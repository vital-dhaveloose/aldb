package ref

import "net/url"

type AttributeSetRef struct {
	NodeRef
	AttributeSetId string
}

type ManifestRef struct {
	Id *url.URL
}
