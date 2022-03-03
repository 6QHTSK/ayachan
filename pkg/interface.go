package pkg

type InputMap interface {
	MapCheck() (bool, error)
	Decode() Chart
	Encode(chart Chart) error
}
