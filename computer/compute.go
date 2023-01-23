package computer

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

type World struct {
	OldState, CurrentState []int8
	right                  bool
	networkCell            int
}

func init() {
	rand.Seed(time.Now().UnixMicro())
	log.SetFlags(1)
}

func InitWorld(worldSize int, density float32, randomInit bool, right bool, remoteConn ...string) (World, error) {

	initializedWorld := World{
		OldState:     make([]int8, worldSize),
		CurrentState: make([]int8, worldSize),
	}

	if randomInit {
		for i := 0; i < int(float32(worldSize)*density); i++ {
			initializedWorld.CurrentState[rand.Intn(worldSize)] = 1
		}
	} else {
		initializedWorld.CurrentState[int(worldSize/2)] = int8(1)
	}

	if len(remoteConn) == 2 {
		initializedWorld.right = right
		if right {
			initializedWorld.networkCell = worldSize - 1
		} else {
			initializedWorld.networkCell = 0
		}
		StartServer(initializedWorld, remoteConn[0])
	} else if len(remoteConn) > 2 {
		return initializedWorld, fmt.Errorf("invalid connection parameters %v", remoteConn)
	}

	return initializedWorld, nil
}

func IntPow(base int, exponent int) int {

	if exponent == 0 {
		return 1
	}
	x := base
	for i := 1; i < exponent; i++ {
		x *= base
	}

	return x
}

func ComputeRule(ruleCode int) [8]int8 {
	rule := [8]int8{}

	for i := 0; i < 8; i++ {
		ruleCode, rule[i] = ruleCode/2, int8(ruleCode%2)
	}
	return rule
}

func GetMinBound(index int, neighborhoodSize int, worldsize int) int {

	minBound := index - neighborhoodSize
	if minBound < 0 {
		minBound = worldsize + minBound
	}

	return minBound
}

func (t *World) GetRemoteState(right bool, state *int8) error {

	if right {
		*state = t.CurrentState[len(t.CurrentState)-1]
	} else {
		*state = t.CurrentState[0]
	}
	return nil
}

func GetNeighborhood(index int, neighborhoodSize int, worldSize int) []int {
	indices := make([]int, (neighborhoodSize*2)+1)
	minBound := GetMinBound(index, neighborhoodSize, worldSize)
	for i := 0; i < (neighborhoodSize*2)+1; i++ {
		indices[i] = minBound % worldSize
		minBound++
	}
	return indices
}

func UpdateState(index int, world World, rule [8]int8, networked bool, remoteConn ...string) error {
	var err error
	var newValue int8
	if index == world.networkCell && networked {
		newValue, err = RetrieveRemoteState(world, remoteConn[0])

	} else {
		newValue, err = NewCellState(
			GetNeighborhood(index, 1, len(world.OldState)),
			world.OldState,
			rule,
		)
	}
	if err != nil {
		return fmt.Errorf("error during computation of state of cell %v", index)
	}
	world.CurrentState[index] = newValue

	return nil
}

func ComputeState(world World, rule [8]int8, remoteConn ...string) error {

	networked := false

	if len(world.OldState) == 0 || len(world.CurrentState) == 0 {
		return fmt.Errorf("World of size '0' given")
	}
	if len(remoteConn) == 1 {
		networked = true
	}
	var cellGroup sync.WaitGroup

	copy(world.OldState, world.CurrentState)
	for i := 0; i < len(world.CurrentState); i++ {
		cellGroup.Add(1)
		go func(index int) {
			defer cellGroup.Done()
			if networked {
				UpdateState(index, world, rule, networked, remoteConn[0])
			} else {
				UpdateState(index, world, rule, networked)
			}

		}(i)
	}
	cellGroup.Wait()
	return nil
}

func NewCellState(neighborhood []int, worldState []int8, rule [8]int8) (int8, error) {

	var newState int8 = 0

	if len(neighborhood) <= 0 {
		return 0, fmt.Errorf("invalid neighborhood %v", neighborhood)
	}

	for i := 0; i < len(neighborhood); i++ {
		newState += worldState[neighborhood[i]] * int8(IntPow(2, i))
	}
	return rule[newState], nil
}
