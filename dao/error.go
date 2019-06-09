package dao

type DbErrorType string

const (
	CONFLICT DbErrorType = "conflict"
	UNKNOWN  DbErrorType = "unknown"
)

// DbError allow to give more context when an error happend
// in DAO layer
type DbError struct {
	Message string
	Type    DbErrorType
}
