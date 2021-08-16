package face

type Face interface {
	Recognize(image []byte) (string, error)
}
