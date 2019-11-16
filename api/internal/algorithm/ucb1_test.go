package algorithm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewUCB1(t *testing.T) {
	testCases := map[string]struct {
		counts  []int
		rewards []float64
		err     error
	}{
		"nil value":           {nil, nil, nil},
		"empty slice":         {make([]int, 0), make([]float64, 0), nil},
		"correct value":       {make([]int, 3), make([]float64, 3), nil},
		"counts is nil":       {nil, make([]float64, 3), ErrInvalidLength},
		"rewards is nil":      {make([]int, 3), nil, ErrInvalidLength},
		"in counts len less":  {make([]int, 3), make([]float64, 5), ErrInvalidLength},
		"in rewards len less": {make([]int, 5), make([]float64, 3), ErrInvalidLength},
	}

	for _, testCase := range testCases {
		ucb1, err := NewUCB1(testCase.counts, testCase.rewards)

		if testCase.err != nil {
			assert.Equal(t, testCase.err, err, "errors should match")
		} else {
			assert.Nil(t, err, "errors should be nil")
			assert.Equal(t, testCase.counts, ucb1.Counts, "counts should be equal")
			assert.Equal(t, testCase.rewards, ucb1.Rewards, "rewards should be equal")
		}
	}
}

func TestUCB1_Reset(t *testing.T) {
	testCases := []struct {
		arms int
		err  error
	}{
		{-1, ErrInvalidArms},
		{0, ErrInvalidArms},
		{1, nil},
		{3, nil},
	}

	for _, testCase := range testCases {
		ucb1, _ := NewUCB1(nil, nil)
		err := ucb1.Reset(testCase.arms)

		if testCase.err != nil {
			assert.Equal(t, err, testCase.err, "should throw error for invalid arms length")
		} else {
			assert.Nil(t, err)
			assert.Equal(t, testCase.arms, len(ucb1.Counts), "counts should be of equal length with arm")
			assert.Equal(t, testCase.arms, len(ucb1.Rewards), "rewards should be of equal length with arm")
		}
	}
}

func TestUCB1_SelectArm(t *testing.T) {
	ucb1, err := NewUCB1(nil, nil)
	assert.Nil(t, err)
	ucb1.Reset(4)

	arm := ucb1.SelectArm()
	assert.Equal(t, 0, arm)
	ucb1.Update(arm, 0.0)

	arm = ucb1.SelectArm()
	assert.Equal(t, 1, arm)
	ucb1.Update(arm, 0.0)

	arm = ucb1.SelectArm()
	assert.Equal(t, 2, arm)
	ucb1.Update(arm, 1.0)

	arm = ucb1.SelectArm()
	assert.Equal(t, 3, arm)
	ucb1.Update(arm, 1.0)

	arm = ucb1.SelectArm()
	assert.Equal(t, 2, arm)
	ucb1.Update(arm, 1.0)

	arm = ucb1.SelectArm()
	assert.Equal(t, 3, arm)
	ucb1.Update(arm, 0.0)

	arm = ucb1.SelectArm()
	assert.Equal(t, 2, arm)
	ucb1.Update(arm, 1.0)

	assert.Equal(t, []int{1, 1, 3, 2}, ucb1.Counts)
	assert.Equal(t, []float64{0, 0, 1, 0.5}, ucb1.Rewards)
}

func TestUCB1_Update(t *testing.T) {
	ucb1, err := NewUCB1(nil, nil)
	assert.Nil(t, err)

	testCases := []struct {
		arms      int
		chosenArm int
		reward    float64
		err       error
	}{
		{1, 0, 0.0, nil},
		{1, -1, 0.0, ErrArmsIndexOutOfRange},
		{1, 1, 0.0, ErrArmsIndexOutOfRange},
		{1, 2, 0.0, ErrArmsIndexOutOfRange},
		{3, 1, 1.0, nil},
		{3, 1, -1.0, ErrInvalidReward},
		{3, 5, -1.0, ErrArmsIndexOutOfRange},
	}

	for _, testCase := range testCases {
		ucb1.Reset(testCase.arms)
		err := ucb1.Update(testCase.chosenArm, testCase.reward)
		if testCase.err != nil {
			assert.Equal(t, testCase.err, err, "should throw error for invalid params")
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestMax(t *testing.T) {
	testCases := []struct {
		params        []float64
		expectedIndex int
		expectedValue float64
	}{
		{[]float64{1.1, 2.1, 3.1, 4.1, 5.1}, 4, 5.1},
		{[]float64{}, 0, 0},
		{[]float64{-1, 1}, 1, 1},
		{[]float64{10.5, 30.5}, 1, 30.5},
	}

	for _, testCase := range testCases {
		index, value := max(testCase.params)

		assert.Equal(t, testCase.expectedIndex, index, "should return the max index")
		assert.Equal(t, testCase.expectedValue, value, "should return the max value")
	}
}

func TestUCB1_SelectArm2(t *testing.T) {
	selected := []int{4, 1, 4, 3, 4}
	reward := []float64{0, 1, 3, 10, 5}
	expectedArm := 3

	ucb1, _ := NewUCB1(selected, reward)

	chosenArm := ucb1.SelectArm()
	assert.Equal(t, expectedArm, chosenArm)
}
