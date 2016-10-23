package rps

import (
	"math/rand"
	"time"
)

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

type LeastThrownStrategy struct {
	Games []*Game
}

func (strategy *LeastThrownStrategy) Throw() string {
	throws := map[string]int{
		ROCK: 0,
		PAPER: 0,
		SCISSORS: 0,
	}

	for _, game := range strategy.Games {
		if game != nil {
			throws[game.YourThrow] = throws[game.YourThrow] + 1
		}
	}

	if throws[ROCK] < throws[PAPER] && throws[ROCK] < throws[SCISSORS] {
		return PAPER
	} else if throws[PAPER] < throws[ROCK] && throws[PAPER] < throws[SCISSORS] {
		return SCISSORS
	} else if throws[SCISSORS] < throws[ROCK] && throws[SCISSORS] < throws[PAPER] {
		return ROCK
	}

	randomStrategy := &RandomStrategy{}
	return randomStrategy.Throw()
}

type MostThrownStrategy struct {
	Games []*Game
}

func (strategy *MostThrownStrategy) Throw() string {
	throws := map[string]int{
		ROCK: 0,
		PAPER: 0,
		SCISSORS: 0,
	}

	for _, game := range strategy.Games {
		if game != nil {
			throws[game.YourThrow] = throws[game.YourThrow] + 1
		}
	}

	if throws[ROCK] > throws[PAPER] && throws[ROCK] > throws[SCISSORS] {
		return PAPER
	} else if throws[PAPER] > throws[ROCK] && throws[PAPER] > throws[SCISSORS] {
		return SCISSORS
	} else if throws[SCISSORS] > throws[ROCK] && throws[SCISSORS] > throws[PAPER] {
		return ROCK
	}

	randomStrategy := &RandomStrategy{}
	return randomStrategy.Throw()
}

type LiteralSequenceStrategy struct {
	PreviousGame *Game
}

func (strategy *LiteralSequenceStrategy) Throw() string {
	if strategy.PreviousGame != nil {
		switch strategy.PreviousGame.MyThrow {
		case ROCK:
			return PAPER
		case PAPER:
			return SCISSORS
		case SCISSORS:
			return ROCK
		}
	}

	return ROCK
}
