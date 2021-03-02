package controller

import (
	"buxiong/car/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStayDirection(t *testing.T) {
	leftStat := model.NewDriverStatus(model.DirectionBrake, 0)
	rightStat := model.NewDriverStatus(model.DirectionBrake, 0)
	assert.Equal(t, determineCarDirection(leftStat, rightStat), stay)
	leftStat = model.NewDriverStatus(model.DirectionGlide, 0)
	rightStat = model.NewDriverStatus(model.DirectionGlide, 0)
	assert.Equal(t, determineCarDirection(leftStat, rightStat), stay)
}

func TestForwardDirection(t *testing.T) {
	leftStat := model.NewDriverStatus(model.DirectionForward, 1)
	rightStat := model.NewDriverStatus(model.DirectionForward, 1)
	assert.Equal(t, determineCarDirection(leftStat, rightStat), forward)
}

func TestForwardLeftDirection(t *testing.T) {
	leftStat := model.NewDriverStatus(model.DirectionForward, 1)
	rightStat := model.NewDriverStatus(model.DirectionForward, 2)
	assert.Equal(t, determineCarDirection(leftStat, rightStat), forwardLeft)
	leftStat = model.NewDriverStatus(model.DirectionBackward, 1)
	rightStat = model.NewDriverStatus(model.DirectionForward, 2)
	assert.Equal(t, determineCarDirection(leftStat, rightStat), forwardLeft)
	leftStat = model.NewDriverStatus(model.DirectionBrake, 0)
	rightStat = model.NewDriverStatus(model.DirectionForward, 1)
	assert.Equal(t, determineCarDirection(leftStat, rightStat), forwardLeft)
	leftStat = model.NewDriverStatus(model.DirectionGlide, 0)
	rightStat = model.NewDriverStatus(model.DirectionForward, 1)
	assert.Equal(t, determineCarDirection(leftStat, rightStat), forwardLeft)
}

func TestForwardRightDirection(t *testing.T) {
	leftStat := model.NewDriverStatus(model.DirectionForward, 2)
	rightStat := model.NewDriverStatus(model.DirectionForward, 1)
	assert.Equal(t, determineCarDirection(leftStat, rightStat), forwardRight)
	leftStat = model.NewDriverStatus(model.DirectionForward, 2)
	rightStat = model.NewDriverStatus(model.DirectionBackward, 1)
	assert.Equal(t, determineCarDirection(leftStat, rightStat), forwardRight)
	leftStat = model.NewDriverStatus(model.DirectionForward, 1)
	rightStat = model.NewDriverStatus(model.DirectionBrake, 0)
	assert.Equal(t, determineCarDirection(leftStat, rightStat), forwardRight)
	leftStat = model.NewDriverStatus(model.DirectionForward, 1)
	rightStat = model.NewDriverStatus(model.DirectionGlide, 0)
	assert.Equal(t, determineCarDirection(leftStat, rightStat), forwardRight)
}

func TestBackwardDirection(t *testing.T) {
	leftStat := model.NewDriverStatus(model.DirectionBackward, 1000)
	rightStat := model.NewDriverStatus(model.DirectionBackward, 1000)
	assert.Equal(t, determineCarDirection(leftStat, rightStat), backward)
}

func TestBackwardLeftDirection(t *testing.T) {
	leftStat := model.NewDriverStatus(model.DirectionBackward, 1)
	rightStat := model.NewDriverStatus(model.DirectionBackward, 2)
	assert.Equal(t, determineCarDirection(leftStat, rightStat), backwardLeft)
	leftStat = model.NewDriverStatus(model.DirectionForward, 1)
	rightStat = model.NewDriverStatus(model.DirectionBackward, 2)
	assert.Equal(t, determineCarDirection(leftStat, rightStat), backwardLeft)
	leftStat = model.NewDriverStatus(model.DirectionBrake, 0)
	rightStat = model.NewDriverStatus(model.DirectionBackward, 1)
	assert.Equal(t, determineCarDirection(leftStat, rightStat), backwardLeft)
	leftStat = model.NewDriverStatus(model.DirectionGlide, 0)
	rightStat = model.NewDriverStatus(model.DirectionBackward, 1)
	assert.Equal(t, determineCarDirection(leftStat, rightStat), backwardLeft)
}

func TestBackwardRightDirection(t *testing.T) {
	leftStat := model.NewDriverStatus(model.DirectionBackward, 2)
	rightStat := model.NewDriverStatus(model.DirectionBackward, 1)
	assert.Equal(t, determineCarDirection(leftStat, rightStat), backwardRight)
	leftStat = model.NewDriverStatus(model.DirectionBackward, 2)
	rightStat = model.NewDriverStatus(model.DirectionForward, 1)
	assert.Equal(t, determineCarDirection(leftStat, rightStat), backwardRight)
	leftStat = model.NewDriverStatus(model.DirectionBackward, 1)
	rightStat = model.NewDriverStatus(model.DirectionBrake, 0)
	assert.Equal(t, determineCarDirection(leftStat, rightStat), backwardRight)
	leftStat = model.NewDriverStatus(model.DirectionBackward, 1)
	rightStat = model.NewDriverStatus(model.DirectionBrake, 0)
	assert.Equal(t, determineCarDirection(leftStat, rightStat), backwardRight)
}

func TestLeftDirection(t *testing.T) {
	leftStat := model.NewDriverStatus(model.DirectionBackward, 2)
	rightStat := model.NewDriverStatus(model.DirectionForward, 2)
	assert.Equal(t, determineCarDirection(leftStat, rightStat), left)
}

func TestRightDirection(t *testing.T) {
	leftStat := model.NewDriverStatus(model.DirectionForward, 2)
	rightStat := model.NewDriverStatus(model.DirectionBackward, 2)
	assert.Equal(t, determineCarDirection(leftStat, rightStat), right)
}
