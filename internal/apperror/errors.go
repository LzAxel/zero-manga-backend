package apperror

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound = errors.New("not found")
)

type AppErorr struct {
	Err     error
	Service string
	Func    string
	Data    map[string]interface{}
}

func (e AppErorr) Error() string {
	return fmt.Sprintf(
		"app error: %s.%s: %s: %v",
		e.Service, e.Func,
		e.Err.Error(),
		e.Data,
	)
}

type DBErorr struct {
	Err     error
	Service string
	Func    string

	Query string
	Args  []interface{}
}

func (e DBErorr) Error() string {
	return fmt.Sprintf(
		"db error: %s.%s: %s: query: %s args: %v",
		e.Service, e.Func,
		e.Err.Error(),
		e.Query,
		e.Args,
	)
}

func NewDBError(err error, service, funcName string, query string, args []interface{}) DBErorr {
	return DBErorr{
		Err:     err,
		Service: service,
		Func:    funcName,
		Query:   query,
		Args:    args,
	}
}
