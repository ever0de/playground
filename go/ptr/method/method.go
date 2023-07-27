package ptr

type Origin struct{}

type Temp interface {
	Get() []byte
}

func (o *Origin) Get() []byte {
	return []byte("origin")
}

type Override struct {
	Temp
}

func NewOverride(o *Origin) *Override {
	return &Override{
		Temp: o,
	}
}

func (o Override) Get() []byte {
	return []byte("override")
}
