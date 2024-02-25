package blockchain

type TxOutput struct {
	Value  int
	PubKey string
}

type TxInput struct {
	ID  []byte
	Out int
	Sig string
}

func (i *TxInput) CanUnlock(data string) bool {
	return i.Sig == data
}

func (o *TxOutput) CanBeUnlocked(data string) bool {
	return o.PubKey == data
}
