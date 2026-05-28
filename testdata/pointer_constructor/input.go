package alpha

type Foo struct {
	id int
}

func NewFoo(id int) *Foo {
	return &Foo{id: id}
}

func (f Foo) ID() int {
	return f.id
}
