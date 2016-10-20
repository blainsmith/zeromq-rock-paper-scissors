package main

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	zmq "github.com/pebbe/zmq4"
)

type games struct {
	count     int
	myScore   int
	yourScore int
}

const (
	ROCK     = "rock"
	PAPER    = "paper"
	SCISSORS = "scissors"

	ME  = "Me"
	YOU = "You"
	TIE = "Tie"
)

var moves []string = []string{ROCK, PAPER, SCISSORS}

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

func throwMove() string {
	rand.Seed(time.Now().UTC().UnixNano())
	return moves[rand.Intn(3)]
}

func computeResult(me, you string) string {
	switch {
	case me == ROCK && you == PAPER:
		return YOU
	case me == SCISSORS && you == PAPER:
		return ME
	case me == PAPER && you == ROCK:
		return ME
	case me == SCISSORS && you == ROCK:
		return YOU
	case me == PAPER && you == SCISSORS:
		return YOU
	case me == ROCK && you == SCISSORS:
		return ME
	}
	return TIE
}

func updateScore(games *games, winner string) {
	if winner == ME {
		games.myScore += 1
	} else if winner == YOU {
		games.yourScore += 1
	}
}

func computeOverall(games *games) string {
	if games.myScore > games.yourScore {
		return ME
	} else if games.myScore < games.yourScore {
		return YOU
	}
	return TIE
}

func startServer() {
	serverGames := &games{}

	// Discover the local IP address and construct the TCP address with the given PORT
	ip, _ := getLocalIP()
	address := "tcp://" + ip.String() + ":" + os.Getenv("PORT")

	// Start a new server PAIR socket and start listening for connections
	server, _ := zmq.NewSocket(zmq.PAIR)
	defer server.Close()
	server.Bind(address)

	// Set the game count to be play with the given GAMES
	serverGames.count, _ = strconv.Atoi(os.Getenv("GAMES"))

	// Display the TCP address and number of games to be played so a client can connect
	fmt.Println("Address:", address)
	fmt.Println("Games:", serverGames.count)
	fmt.Println("")

	// Block untol a client connect and on connection send the number of games to be played
	server.Send(os.Getenv("GAMES"), 0)

	startTime := time.Now()

	gameCounter := 0
	for {
		gameCounter++

		// Block and wait for the client to send their move
		yourMove, _ := server.Recv(0)

		// Throw your move to the client
		myMove := throwMove()
		_, err := server.Send(myMove, 0)
		if err != nil {
			break //  Interrupted
		}

		// Display the game and the moves
		fmt.Println("Game:", gameCounter)
		fmt.Println("Me:", myMove)
		fmt.Println("You:", yourMove)

		// Compute and display the winner
		winner := computeResult(myMove, yourMove)
		fmt.Println("Winner:", winner)

		// Update and display the running score
		updateScore(serverGames, winner)
		fmt.Println("Score:", serverGames.myScore, "/", serverGames.yourScore)
		fmt.Println("")

		// If the current game counter is the number of games to be played then break the loop to stop playing
		if gameCounter == serverGames.count {
			// Send "end" to the client to tell them the games are over and to disconnect
			server.Send("end", 0)
			break
		}
	}

	// After quiting th loop the games are over so display the final winner
	fmt.Println("Overall:", computeOverall(serverGames))
	fmt.Println("Time Spent Playing:", time.Now().Sub(startTime))
}

func startClient() {
	clientGames := &games{}

	// Start a new client PAIR socket and connect to the given ADDRESS
	client, _ := zmq.NewSocket(zmq.PAIR)
	defer client.Close()
	client.Connect(os.Getenv("ADDRESS"))

	// After connecting receive the first message from the server as the number of games to be played
	games, _ := client.Recv(0)
	clientGames.count, _ = strconv.Atoi(games)
	fmt.Println("Games:", clientGames.count)
	fmt.Println("")

	startTime := time.Now()

	// Start and increment a simple display counter
	gameCounter := 0
	for {
		gameCounter++

		// Start playing by sending a move
		myMove := throwMove()
		_, err := client.Send(myMove, 0)
		if err != nil {
			break //  Interrupted
		}

		// Block until you receive a move back, if the move is "end" then games are done and break the loop
		yourMove, _ := client.Recv(0)
		if yourMove == "end" {
			break
		}

		// Display the game and the moves
		fmt.Println("Game:", gameCounter)
		fmt.Println("Me:", myMove)
		fmt.Println("You:", yourMove)

		// Compute and display the winner
		winner := computeResult(myMove, yourMove)
		fmt.Println("Winner:", winner)

		// Update and display the running score
		updateScore(clientGames, winner)
		fmt.Println("Score:", clientGames.myScore, "/", clientGames.yourScore)
		fmt.Println("")
	}

	// After quiting th loop the games are over so display the final winner
	fmt.Println("Overall:", computeOverall(clientGames))
	fmt.Println("Time Spent Playing:", time.Now().Sub(startTime))
}

func main() {
	if os.Getenv("ADDRESS") != "" {
		startClient()
	} else {
		startServer()
	}
}
