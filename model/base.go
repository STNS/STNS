package model

type Base struct {
	ID   int    `toml:"id" json:"id" yaml:"id"`
	Name string `toml:"name" json:"name" yaml:"name"`
}

func (b *Base) id() int {
	return b.ID
}

func (b *Base) name() string {
	return b.Name
}

func (b *Base) setName(n string) {
	b.Name = n
}
