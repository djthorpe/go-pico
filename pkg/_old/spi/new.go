package spi

func New(cfg Config) *device {
	if d, err := cfg.New(); err != nil {
		panic(err)
	} else {
		return d
	}
}
