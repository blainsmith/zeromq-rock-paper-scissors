using ZMQ

function getlocalip()
	ip = getipaddr()
	ip = dec((ip.host&(0xFF000000))>>24) * "." *
				dec((ip.host&(0xFF0000))>>16) * "." *
				dec((ip.host&(0xFF00))>>8) * "." *
				dec(ip.host&0xFF)
end

context = Context()
server = Socket(context, PAIR)

address = "tcp://" * getlocalip() * ":" * ENV["PORT"]
games = ENV["GAMES"]

ZMQ.bind(server, address)
print("Address: " * address * "\n")
print("Games: " * games * "\n")

ZMQ.send(server, games)

for game = 1:parse(games)
	move = unsafe_string(ZMQ.recv(server))
	ZMQ.send(server, "rock")
end

ZMQ.send(server, "end")
