package alpha

import "errors"

const MaxLorem = 100

var ErrInvalid = errors.New("invalid")

type Foo struct {
	id  int
	bar string
}

func NewFoo(id int, bar string) (Foo, error) {
	if err := requireNonBlank(bar); err != nil {
		return Foo{}, err
	}
	return Foo{id: id, bar: bar}, nil
}

func (f Foo) ID() int {
	return f.id
}

func (f Foo) Bar() string {
	return f.bar
}

func (f *Foo) SetBar(bar string) {
	f.bar = bar
}

func (f Foo) Lorem() (Foo, error) {
	if err := f.requireValid(); err != nil {
		return Foo{}, err
	}
	return f.applyLorem(), nil
}

func (f Foo) requireValid() error {
	if f.bar == "" {
		return ErrInvalid
	}
	return nil
}

func (f Foo) applyLorem() Foo {
	f.bar = f.bar + "_lorem"
	return f
}

func requireNonBlank(value string) error {
	if value == "" {
		return ErrInvalid
	}
	return nil
}
