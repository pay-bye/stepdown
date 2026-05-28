package alpha

type Foo struct {
	id int
}

func NewFoo(id int) Foo {
	return Foo{id: id}
}

func (f Foo) ID() int {
	return f.id
}

func (f Foo) Process() int {
	return f.prepare()
}

func (f Foo) prepare() int {
	return f.finish()
}

func (f Foo) finish() int {
	return f.id
}
