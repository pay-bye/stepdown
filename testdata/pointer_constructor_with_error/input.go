package alpha

import "errors"

var errInvalid = errors.New("invalid")

type Foo struct {
	id int
}

func NewFoo(id int) (*Foo, error) {
	if id < 0 {
		return nil, errInvalid
	}
	return &Foo{id: id}, nil
}

func (f Foo) ID() int {
	return f.id
}
