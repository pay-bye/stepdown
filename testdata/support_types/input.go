package alpha

type Config struct {
	Host string
	Port int
}

type Settings struct {
	Debug bool
	Limit int
}

type Service struct {
	config Config
}

func NewService(config Config) Service {
	return Service{config: config}
}

func (s Service) Run() error {
	return nil
}
