package algorithm

import (
	"errors"
	"math"
)

var (
	ErrInvalidLength       = errors.New("counts and rewards must be of equal length")
	ErrInvalidArms         = errors.New("arms must be greater than zero")
	ErrArmsIndexOutOfRange = errors.New("arms index is out of range")
	ErrInvalidReward       = errors.New("reward must be greater than zero")
)

// Algorithm interface
type Algorithm interface {
	Reset(int) error
	SelectArm() int
	Update(int, float64)
}

// UCB1 algorithm
type UCB1 struct {
	Counts  []int
	Rewards []float64
}

// NewUCB1 returns a pointer to the UCB1 struct
func NewUCB1(counts []int, rewards []float64) (*UCB1, error) {
	if len(counts) != len(rewards) {
		return nil, ErrInvalidLength
	}

	return &UCB1{
		Counts:  counts,
		Rewards: rewards,
	}, nil
}

// Reset will set the counts and rewards with the provided number of arms
func (u *UCB1) Reset(nArms int) error {
	if nArms < 1 {
		return ErrInvalidArms
	}

	u.Counts = make([]int, nArms)
	u.Rewards = make([]float64, nArms)

	return nil
}

// SelectArm chooses an arm that exploits if the value is more than the epsilon
// threshold, and explore if the value is less than epsilon
func (b *UCB1) SelectArm() (index int) {
	for i, v := range b.Counts {
		if v == 0 {
			return i
		}
	}

	totalCounts := 0
	for _, v := range b.Counts {
		totalCounts = totalCounts + v
	}

	ucbValues := []float64{}
	for i := range b.Counts {
		bonus := math.Sqrt(2.0 * math.Log(float64(totalCounts)) / float64(b.Counts[i]))
		ucbValues = append(ucbValues, b.Rewards[i]+bonus)
	}

	maxIndex, _ := max(ucbValues)

	return maxIndex
}

// Update will update an arm with some reward value
func (b *UCB1) Update(chosenArm int, reward float64) error {
	if chosenArm < 0 || chosenArm >= len(b.Rewards) {
		return ErrArmsIndexOutOfRange
	}
	if reward < 0 {
		return ErrInvalidReward
	}

	b.Counts[chosenArm] = b.Counts[chosenArm] + 1
	n := b.Counts[chosenArm]

	value := b.Rewards[chosenArm]
	if n == 0 {
		b.Rewards[chosenArm] = reward
	} else {
		b.Rewards[chosenArm] = (float64(n-1)/float64(n))*value + (1.0/float64(n))*reward
	}

	return nil
}

// Max will return the max index
func max(d []float64) (index int, value float64) {
	if len(d) > 0 {
		value = d[0]
	}

	index = 0

	for i, v := range d {
		if v > value {
			index = i
			value = v
		}
	}

	return
}
