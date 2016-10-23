package rps

type Game struct {
	MyThrow string
	YourThrow string
	Winner string
}

func NewGame(me, you string) *Game {
	var game = &Game{MyThrow: me, YourThrow: you, Winner: ""}

	switch {
	case me == ROCK && you == PAPER:
		game.Winner = YOU
	case me == SCISSORS && you == PAPER:
		game.Winner = ME
	case me == PAPER && you == ROCK:
		game.Winner = ME
	case me == SCISSORS && you == ROCK:
		game.Winner = YOU
	case me == PAPER && you == SCISSORS:
		game.Winner = YOU
	case me == ROCK && you == SCISSORS:
		game.Winner = ME
	default:
		game.Winner = TIE
	}

	return game
}
