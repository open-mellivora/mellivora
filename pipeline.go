package mellivora

type Pipeline interface {
	ProcessItems(items ...interface{})
}
