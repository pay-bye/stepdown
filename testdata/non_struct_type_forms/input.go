package alpha

type ID string

func NewID(value string) ID {
	return ID(value)
}

func (i ID) String() string {
	return string(i)
}

type Handler func(ID) string

func NewHandler() Handler {
	return func(id ID) string {
		return id.String()
	}
}

type Items []ID

func NewItems(id ID) Items {
	return Items{id}
}

func (i Items) Count() int {
	return len(i)
}
