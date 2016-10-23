package rps

import (
	"math/rand"
	"time"
)

var Strategies map[string]Strategy = map[string]Strategy{
	"random": &RandomStrategy{},
	"wsls": &WinStayLoseShiftStrategy{},
}

type Strategy interface {
	Throw() string
}

type RandomStrategy struct {}

func (strategy *RandomStrategy) Throw() string {
	rand.Seed(time.Now().UTC().UnixNano())
	return ThrowOptions[rand.Intn(3)]
}

type WinStayLoseShiftStrategy struct {
	PreviousGame *Game
}

func (strategy *WinStayLoseShiftStrategy) Throw() string {
	if strategy.PreviousGame.Winner == ME {
		switch strategy.PreviousGame.MyThrow {
		case ROCK:
			return PAPER
		case PAPER:
			return SCISSORS
		case SCISSORS:
			return ROCK
		}
	} else if strategy.PreviousGame.Winner == YOU {
		switch strategy.PreviousGame.YourThrow {
		case ROCK:
			return PAPER
		case PAPER:
			return SCISSORS
		case SCISSORS:
			return ROCK
		}
	}

	randomStrategy := &RandomStrategy{}
	return randomStrategy.Throw()
}
