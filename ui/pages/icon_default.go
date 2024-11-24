package pages

type DefaultIcon struct{}

func (DefaultIcon) Name() string {
	return "default"
}
func (DefaultIcon) Content() []byte {
	return []byte{}
}
