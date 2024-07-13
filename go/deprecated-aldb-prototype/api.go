package aldb

import (
	"context"
)

type RawApi interface {
	Create(ctx context.Context, req CreateRequest) (CreateResponse, error)
	Delete(ctx context.Context, req DeleteRequest) (DeleteResponse, error)
	Read(ctx context.Context, req ReadRequest) (ReadResponse, error)
}

type CreateRequest struct {
	//TODO: what attributes/blobdata/relations are copied from the previous to
	// the new version? Or must the caller "repeat" al of that?
	ToCreate ActivityVersion
}

type CreateResponse struct {
	Created ActivityVersion
}

type DeleteRequest struct {
	// ToDelete VersionSelection
}

type DeleteResponse struct {
	//TODO fields
}

type ReadRequest struct {
	//TODO fields
}

type ReadResponse struct {
	//TODO fields
}
