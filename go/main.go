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

type score struct {
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

func updateScore(score *score, winner string) {
	if winner == ME {
		score.myScore += 1
	} else if winner == YOU {
		score.yourScore += 1
	}
}

func computeOverall(score *score) string {
	if score.myScore > score.yourScore {
		return ME
	} else if score.myScore < score.yourScore {
		return YOU
	}
	return TIE
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
	games, _ := strconv.Atoi(os.Getenv("GAMES"))

	// Display the TCP address and number of games to be played so a client can connect
	fmt.Println("Address:", address)
	fmt.Println("Games:", games)
	fmt.Println("")

	// Block untol a client connect and on connection send the number of games to be played
	server.Send(os.Getenv("GAMES"), 0)

	startTime := time.Now()

	for game := 1; game <= games; game++ {
		// Block and wait for the client to send their move
		yourMove, _ := server.Recv(0)

		// Throw your move to the client
		myMove := throwMove()
		_, err := server.Send(myMove, 0)
		if err != nil {
			break //  Interrupted
		}

		// Display the game and the moves
		fmt.Println("Game:", game)
		fmt.Println("Me:", myMove)
		fmt.Println("You:", yourMove)

		// Compute and display the winner
		winner := computeResult(myMove, yourMove)
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
	games, _ := client.Recv(0)
	fmt.Println("Games:", games)
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
		updateScore(score, winner)
		fmt.Println("Score:", score.myScore, "/", score.yourScore)
		fmt.Println("")
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
