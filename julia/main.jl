using ZMQ

type Score
	myscore::Int
	yourscore::Int
end

const ROCK = "rock"
const PAPER = "paper"
const SCISSORS = "scissors"
const ME = "Me"
const YOU = "You"
const TIE = "Tie"

function getlocalip()
	ip = getipaddr()
	ip = dec((ip.host&(0xFF000000))>>24) * "." *
				dec((ip.host&(0xFF0000))>>16) * "." *
				dec((ip.host&(0xFF00))>>8) * "." *
				dec(ip.host&0xFF)
end

function computeresults(me, you)
	if me == ROCK && you == PAPER
		return YOU
	elseif me == SCISSORS && you == PAPER
		return ME
	elseif me == PAPER && you == ROCK
		return ME
	elseif me == SCISSORS && you == ROCK
		return YOU
	elseif me == PAPER && you == SCISSORS
		return YOU
	elseif me == ROCK && you == SCISSORS
		return ME
	else
		return TIE
	end
end

function updatescore!(score, winner)
	if winner == ME
		score.myscore += 1
	elseif winner == YOU
		score.yourscore += 1
	end
end

function computeoverall(score)
	if score.myscore > score.yourscore
		return ME
	elseif score.myscore < score.yourscore
		return YOU
	else
		return TIE
	end
end

function startserver()
	score = Score(0, 0)

	context = Context()
	server = Socket(context, PAIR)

	address = "tcp://" * getlocalip() * ":" * ENV["PORT"]
	games = ENV["GAMES"]

	ZMQ.bind(server, address)
	print("Address: " * address * "\n")
	print("Games: " * games * "\n")

	ZMQ.send(server, games)

	starttime = time() * 1000

	for game = 1:parse(games)
		yourmove = unsafe_string(ZMQ.recv(server))

		mymove = rand([ROCK, PAPER, SCISSORS])
		ZMQ.send(server, mymove)

		print("Game: " * dec(game) * "\n")
		print("Me: " * mymove * "\n")
		print("You: " * yourmove * "\n")

		winner = computeresults(mymove, yourmove)
		print("Winner: " * winner * "\n")

		updatescore!(score, winner)
		print("Score: " * dec(score.myscore) * "/" * dec(score.yourscore) * "\n")

		print("\n")
	end

	print("Overall: " * computeoverall(score) * "\n")

	endtime = time() * 1000
	print("Time Spent Playing: ")
	print(endtime - starttime)
	print("ms\n")

	ZMQ.send(server, "end")
end

function startclient()
	score = Score(0, 0)

	context = Context()
	client = Socket(context, PAIR)

	ZMQ.connect(client, ENV["ADDRESS"])

	games = unsafe_string(ZMQ.recv(client))
	print("Games: " * games * "\n")

	starttime = time() * 1000

	gamecounter = 0
	while true
		gamecounter = gamecounter + 1

		mymove = rand([ROCK, PAPER, SCISSORS])
		ZMQ.send(client, mymove)

		yourmove = unsafe_string(ZMQ.recv(client))
		if yourmove == "end"
			break
		end

		print("Game: " * dec(gamecounter) * "\n")
		print("Me: " * mymove * "\n")
		print("You: " * yourmove * "\n")

		winner = computeresults(mymove, yourmove)
		print("Winner: " * winner * "\n")

		updatescore!(score, winner)
		print("Score: " * dec(score.myscore) * "/" * dec(score.yourscore) * "\n")

		print("\n")
	end

	print("Overall: " * computeoverall(score) * "\n")

	endtime = time() * 1000
	print("Time Spent Playing: ")
	print(endtime - starttime)
	print("ms\n")
end

try
	address = ENV["ADDRESS"]
	startclient()
catch
end

try
	port = ENV["PORT"]
	games = ENV["GAMES"]
	startserver()
catch
end
