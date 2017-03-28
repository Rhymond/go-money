package money

type Number struct {
	Amount int
	//mantissa int
	//characteristic int
}

func (n *Number) Negative() {
	if n.Amount > 0 {
		n.Amount = -n.Amount
	}
}

func (n *Number) Absolute() {
	if n.Amount < 0 {
		n.Amount = -n.Amount
	}
}