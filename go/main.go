package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	zmq "github.com/pebbe/zmq4"

	"github.com/blainsmith/zeromq-rock-paper-scissors/go/rps"
)

type score struct {
	myScore   int
	yourScore int
}

func getLocalIP() (net.IP, error) {
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, address := range addresses {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP, nil
			}
		}
	}

	return nil, nil
}

func updateScore(score *score, winner string) {
	if winner == rps.ME {
		score.myScore += 1
	} else if winner == rps.YOU {
		score.yourScore += 1
	}
}

func computeOverall(score *score) string {
	if score.myScore > score.yourScore {
		return rps.ME
	} else if score.myScore < score.yourScore {
		return rps.YOU
	}
	return rps.TIE
}

func startServer() {
	score := &score{}

	// Discover the local IP address and construct the TCP address with the given PORT
	ip, _ := getLocalIP()
	address := "tcp://" + ip.String() + ":" + os.Getenv("PORT")

	// Start a new server PAIR socket and start listening for connections
	server, _ := zmq.NewSocket(zmq.PAIR)
	defer server.Close()
	server.Bind(address)

	// Set the game count to be play with the given GAMES
	numGames, _ := strconv.Atoi(os.Getenv("GAMES"))

	strat := os.Getenv("STRATEGY")
	if strat == "" {
		strat = "random"
	}
	var strategy rps.Strategy = &rps.RandomStrategy{}

	games := make([]*rps.Game, numGames)

	// Display the TCP address and number of games to be played so a client can connect
	fmt.Println("Address:", address)
	fmt.Println("Games:", numGames)
	fmt.Println("")

	// Block untol a client connect and on connection send the number of games to be played
	server.Send(os.Getenv("GAMES"), 0)

	startTime := time.Now()

	for game := 0; game < numGames; game++ {
		// Block and wait for the client to send their move
		yourThrow, _ := server.Recv(0)

		// Throw your move to the client
		myThrow := strategy.Throw()

		_, err := server.Send(myThrow, 0)
		if err != nil {
			break //  Interrupted
		}

		games[game] = rps.NewGame(myThrow, yourThrow)
		switch strat {
		case "wsls":
			strategy = &rps.WinStayLoseShiftStrategy{PreviousGame: games[game]}
		case "lt":
			strategy = &rps.LeastThrownStrategy{Games: games}
		case "random":
		default:
			strategy = &rps.RandomStrategy{}
		}

		// Display the game and the moves
		fmt.Println("Game:", (game+1))
		fmt.Println("Me:", myThrow)
		fmt.Println("You:", yourThrow)

		// Compute and display the winner
		winner := games[game].Winner
		fmt.Println("Winner:", winner)

		// Update and display the running score
		updateScore(score, winner)
		fmt.Println("Score:", score.myScore, "/", score.yourScore)
		fmt.Println("")
	}

	// Send "end" to the client to tell them the games are over and to disconnect
	server.Send("end", 0)

	// After quiting th loop the games are over so display the final winner
	fmt.Println("Overall:", computeOverall(score))
	fmt.Println("Time Spent Playing:", time.Now().Sub(startTime))
}

func startClient() {
	score := &score{}

	// Start a new client PAIR socket and connect to the given ADDRESS
	client, _ := zmq.NewSocket(zmq.PAIR)
	defer client.Close()
	client.Connect(os.Getenv("ADDRESS"))

	// After connecting receive the first message from the server as the number of games to be played
	gamesFromServer, _ := client.Recv(0)
	numGames, _ := strconv.Atoi(gamesFromServer)
	fmt.Println("Games:", numGames)
	fmt.Println("")

	startTime := time.Now()

	strat := os.Getenv("STRATEGY")
	if strat == "" {
		strat = "random"
	}
	var strategy rps.Strategy = &rps.RandomStrategy{}

	games := make([]*rps.Game, numGames)

	// Start and increment a simple display counter
	gameCounter := 0
	for {
		// Start playing by sending a move
		myThrow := strategy.Throw()
		_, err := client.Send(myThrow, 0)
		if err != nil {
			break //  Interrupted
		}

		// Block until you receive a move back, if the move is "end" then games are done and break the loop
		yourThrow, _ := client.Recv(0)
		if yourThrow == "end" {
			break
		}

		games[gameCounter] = rps.NewGame(myThrow, yourThrow)
		switch strat {
		case "wsls":
			strategy = &rps.WinStayLoseShiftStrategy{PreviousGame: games[gameCounter]}
		case "random":
		default:
			strategy = &rps.RandomStrategy{}
		}

		// Display the game and the moves
		fmt.Println("Game:", (gameCounter+1))
		fmt.Println("Me:", myThrow)
		fmt.Println("You:", yourThrow)

		// Compute and display the winner
		winner := games[gameCounter].Winner
		fmt.Println("Winner:", winner)

		// Update and display the running score
		updateScore(score, winner)
		fmt.Println("Score:", score.myScore, "/", score.yourScore)
		fmt.Println("")

		gameCounter++
	}

	// After quiting th loop the games are over so display the final winner
	fmt.Println("Overall:", computeOverall(score))
	fmt.Println("Time Spent Playing:", time.Now().Sub(startTime))
}

func main() {
	if os.Getenv("ADDRESS") != "" {
		startClient()
	} else {
		startServer()
	}
}
