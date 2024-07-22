package group

type GroupInfo struct {
	Desc string
	Obj  interface{}
}
type Group map[string]GroupInfo
type Command struct {
	Header  string
	Options Group
	Handler func(args ...interface{}) interface{}
}
