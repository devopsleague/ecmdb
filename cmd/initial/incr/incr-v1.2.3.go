package incr

import (
	"github.com/Duke1616/ecmdb/cmd/initial/ioc"
)

type incrV123 struct {
	App *ioc.App
}

func NewIncrV123(app *ioc.App) InitialIncr {
	return &incrV123{
		App: app,
	}
}

func (i *incrV123) Version() string {
	return "v1.2.3"
}

func (i *incrV123) Commit() error {
	return nil
}

func (i *incrV123) Rollback() error {
	//TODO implement me
	panic("implement me")
}