package jsonparser

type ValueType int

const (
	VString = iota
	VNumber
	VBoolean
	VNull
	VObject
	VArray
)

type Value struct {
	Type ValueType
	Data any
}
