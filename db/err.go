package db

import "fmt"

type Err struct {
	Name   string
	Engine string
	Type   string
	Msg    string
}

func (e *Err) Error() string {
	return fmt.Sprintf("dbName=%s, dbEngine=%s, dbType=%s, err: %s", e.Name, e.Engine, e.Type, e.Msg)
}

func newErr(name, engine, dbType, msg string) *Err {
	return &Err{
		Name:   name,
		Engine: engine,
		Type:   dbType,
		Msg:    msg,
	}
}
