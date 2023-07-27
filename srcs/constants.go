package main

var evals = []eval{
	{"dijkstra", dijkstra},
	{"greedy_hamming", greedy_hamming},
	{"greedy_manhattan", greedy_manhattan},
	{"greedy_max_manhattan", greedy_max_manhattan},
	{"astar_hamming", astar_hamming},
	{"astar_max_manhattan", astar_max_manhattan},
	{"astar_manhattan", astar_manhattan_generator(1)},
	{"astar_manhattan2", astar_manhattan_generator(2)},
}

var directions = []struct {
	name byte
	fx   moveFx
}{
	{'U', moveUp},
	{'D', moveDown},
	{'L', moveLeft},
	{'R', moveRight},
}

var minRAMAvailableMB uint64 = 512
