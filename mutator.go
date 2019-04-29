package money

type mutator struct {
	calc *calculator
}

// initialize our default mutator here.
var mutate = mutator{calc: &calculator{}}
