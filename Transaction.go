package assignment02IBC_master

type Transaction struct {
	Amount   int
	Sender   string
	Receiver string
}

func (T Transaction) IsEmpty() bool{
	if T.Amount == 0 && T.Receiver == "" && T.Sender == "" {
		return true
	}
	return false
}