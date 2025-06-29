package migrate

// SeqInfo migrate seq info
type SeqInfo struct {
	Seq           string
	Description   string
	Ext           string
	Instance      interface{}
	Ignore        bool
	AlterOpertion bool
}
