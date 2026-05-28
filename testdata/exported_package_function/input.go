package alpha

type Foo struct {
	name string
}

func NewFoo(name string) Foo {
	return Foo{name: name}
}

func BuildFoo(name string) Foo {
	return NewFoo(name)
}

func normalize(value string) string {
	return value
}
