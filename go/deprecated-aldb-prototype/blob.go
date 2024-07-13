package aldb

import "io"

type Blob struct {
	Renditions []BlobRendition
}

type BlobRendition interface {
	io.ReadWriteCloser
	Type() BlobRenditionFunction
	MediaType() MediaType
	//TODO PhysicalLocation? LegacyLocation?
}

type BlobRenditionFunction string

const (
	BlobRenditionFunctionMain      = BlobRenditionFunction("main")
	BlobRenditionFunctionThumbnail = BlobRenditionFunction("thumbnail")
	BlobRenditionFunctionOcrText   = BlobRenditionFunction("ocr-text")
)

type MediaType struct {
	Type, SubType string
	Params        map[string]string
}
