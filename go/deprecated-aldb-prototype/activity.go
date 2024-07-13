package aldb

type ActivityId string
type VersionId string

type ActivityRef struct {
	ActivityId string
}

func (r *ActivityRef) String() string {
	return r.ActivityId
}

func (r *ActivityRef) FromString(s string) error {
	r.ActivityId = s
	return nil
}

type Activity struct {
	Ref      ActivityRef
	Versions []ActivityVersion
}

type VersionRef struct {
	ActivityRef
	// VersionId identifies the version of the Activity amongst the other versions. It should be strictly
	// increasing, i.e. it must come after the previous version when sorted lexicographically. Actual ActivityVersions
	// should not have a versionId from the set SpecialVersionIds, as these are used for specific semantics in
	// requests.
	// TODO define collation algorithm
	VersionId string
}

func (r *VersionRef) String() string {
	return r.ActivityId + "|" + r.VersionId //TODO definitive separator + versionId optional?
}

func (r *VersionRef) FromString(s string) error {
	r.ActivityId = "foo"
	r.VersionId = "bar"
	return nil //TODO
}

type ActivityVersion struct {
	Ref        VersionRef
	Attributes Value
	Blob       Blob
}
