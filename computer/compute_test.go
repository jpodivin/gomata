package computer

import (
	"math"
	"testing"
)

func max(array []int8) int8 {
	max_e := array[0]

	for i := 1; i < len(array); i++ {
		if array[i] > max_e {
			max_e = array[i]
		}
	}
	return max_e
}

func min(array []int8) int8 {
	min_e := array[0]
	for i := 1; i < len(array); i++ {
		if array[i] < min_e {
			min_e = array[i]
		}
	}
	return min_e
}

func TestComputeState(t *testing.T) {

	test_world := World{
		OldState:     make([]int8, 10),
		CurrentState: make([]int8, 10),
	}
	test_rule := ComputeRule(0)
	err := ComputeState(test_world, test_rule)

	if max(test_world.CurrentState) > 1 || min(test_world.CurrentState) < 0 || err != nil {
		t.Fatalf("ComputeState(%v, %v) = %q", test_world.CurrentState, test_rule, err)
	}
}

func TestComputeStateEmpty(t *testing.T) {
	test_world := World{
		OldState:     make([]int8, 0),
		CurrentState: make([]int8, 0),
	}
	test_rule := ComputeRule(0)
	err := ComputeState(test_world, test_rule)

	if err == nil {
		t.Fatalf("ComputeState(%v, %v) = %q", test_world.CurrentState, test_rule, err)
	}
}

func TestIntPow(t *testing.T) {

	for exponent := 0; exponent < 10; exponent++ {
		for base := 1; base < 5; base++ {
			x := IntPow(base, exponent)

			if x != int(math.Pow(float64(base), float64(exponent))) {
				t.Fatalf("IntPow(%v, %v) = %v", base, exponent, x)
			}
		}
	}

}
