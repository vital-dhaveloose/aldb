package ref

type Ref interface {
	ToName() string
	FromName(name string) error
	IsComplete() bool
}
