package gocash
//
//import "testing"
//
//func TestNumber_New(t *testing.T) {
//	tc := map[int][2]int{
//		123456:  {56, 1234},
//		159:     {59, 1},
//		56:      {56, 0},
//		1:       {1, 0},
//		0:       {0, 0},
//		-1:      {1, 0},
//		-56:     {56, 0},
//		-159:    {59, 1},
//		-123456: {56, 1234},
//	}
//
//	for amount, expected := range tc {
//		n := new(Number).New(amount)
//
//		if n.mantissa != expected[0] {
//			t.Errorf("Expected %d mantissa to be %d got %d", n.Amount, expected[0], n.mantissa)
//		}
//
//		if n.characteristic != expected[1] {
//			t.Errorf("Expected %d characteristic to be %d got %d", n.Amount, expected[1], n.characteristic)
//		}
//	}
//}
