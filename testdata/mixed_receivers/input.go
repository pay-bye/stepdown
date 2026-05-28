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

func (f *Foo) SetID(id int) {
	f.id = id
}

func (f Foo) Value() int {
	return f.value()
}

func (f Foo) value() int {
	return f.id
}

func (f *Foo) Reset() {
	f.applyReset()
}

func (f *Foo) applyReset() {
	f.id = 0
}
