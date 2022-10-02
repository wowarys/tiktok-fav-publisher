package internal

type Database interface {
	Add(value string) error
	IsExist(value string) bool
}
