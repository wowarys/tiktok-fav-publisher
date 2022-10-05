package internal

type Cache interface {
	Add(value string) error
	IsExist(value string) bool
}
