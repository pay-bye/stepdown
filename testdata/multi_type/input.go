package alpha

import "errors"

var errInvalid = errors.New("invalid")

type Foo struct {
	id int
}

func NewFoo(id int) Foo {
	return Foo{id: id}
}

func (f Foo) ID() int {
	return f.id
}

func (f Foo) Process() error {
	return f.validate()
}

func (f Foo) validate() error {
	if f.id < 0 {
		return errInvalid
	}
	return nil
}

type Widget struct {
	name string
}

func NewWidget(name string) Widget {
	return Widget{name: name}
}

func (w Widget) Name() string {
	return w.name
}

func (w Widget) Render() string {
	return w.format()
}

func (w Widget) format() string {
	return w.name
}
