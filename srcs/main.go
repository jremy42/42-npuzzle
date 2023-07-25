package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func areFlagsOk(workers int, seenNodesSplit int, heuristic string, mapSize int) (fx eval, ok bool) {
	if workers < 1 || workers > 16 {
		fmt.Println("Invalid number of workers")
		os.Exit(1)
	}
	if seenNodesSplit < 1 || seenNodesSplit > 64 {
		fmt.Println("Invalid number of splits")
		os.Exit(1)
	}
	if mapSize < 3 {
		fmt.Println("Invalid map size")
		os.Exit(1)
	}
	for _, current := range evals {
		if current.name == heuristic {
			return current, true
		}
	}
	fmt.Println("Invalid heuristic")
	return eval{}, false
}

func handleSignals() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGKILL)
	go func() {
		<-sigc
		fmt.Println("\b\bExiting after receiving a signal")
		os.Exit(1)
	}()
}

func getAvailableRAM() (uint64, error) {
	var info syscall.Sysinfo_t
	err := syscall.Sysinfo(&info)
	if err != nil {
		return 0, fmt.Errorf("Error while getting info about memory: %v", err)
	}
	availableRAM := info.Freeram*uint64(info.Unit) + info.Bufferram*uint64(info.Unit)
	return availableRAM, nil
}

func iterateAlgo(board [][]int, workers int, seenNodesSplit int, evalfx evalFx, data *safeData) (result Result) {
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(board [][]int, evalfx evalFx, data *safeData, i int, workers int, seenNodesSplit int) {

			algo(board, evalfx, data, i, workers, seenNodesSplit)
			wg.Done()
		}(board, evalfx, data, i, workers, seenNodesSplit)
	}
	wg.Wait()
	if data.path != nil {
		for _, value := range data.seenNodes {
			data.closedSetComplexity += len(value)
		}
	}
	return Result{data.path, data.closedSetComplexity, data.tries, false}
}

func main() {
	var (
		file                string
		mapSize             int
		heuristic           string
		workers             int
		seenNodesSplit      int
		speedDisplay        int
		noIterativeDepth    bool
		debug               bool
		disableUI           bool
		closedSetComplexity int
		board               [][]int
	)
	handleSignals()

	flag.StringVar(&file, "f", "", "usage : -f [filename]")
	flag.IntVar(&mapSize, "s", 3, "usage : -s [size]")
	flag.StringVar(&heuristic, "h", "astar_manhattan", "usage : -h [heuristic]")
	flag.IntVar(&workers, "w", 1, "usage : -w [workers] between 1 and 16")
	flag.IntVar(&seenNodesSplit, "split", 1, "usage : -split [setNodesSplit] between 1 and 32")
	flag.IntVar(&speedDisplay, "speed", 100, "usage : -speed [speedDisplay] between 1 and 1000")
	flag.BoolVar(&noIterativeDepth, "no-i", false, "usage : -no-i")
	flag.BoolVar(&debug, "d", false, "usage : -d")
	flag.BoolVar(&disableUI, "no-ui", false, "usage : -no-ui")
	flag.Parse()

	eval, ok := areFlagsOk(workers, seenNodesSplit, heuristic, mapSize)
	if !ok {
		os.Exit(1)
	}
	if !debug {
		newstderr, _ := os.Open("/dev/null")
		defer newstderr.Close()
		os.Stderr = newstderr
	}
	if file != "" {
		file := OpenFile(file)
		fmt.Println("Opening user provided map in file", file.Name())
		_, board = ParseInput(file)
	} else if mapSize > 0 {
		fmt.Println("Generating a map with size", mapSize)
		board = gridGenerator(mapSize)
	} else {
		fmt.Println("No valid map size or filename option missing")
		os.Exit(1)
	}
	if !isSolvable(board) {
		fmt.Println("Board is not solvable")
		if !disableUI {
			displayBoard(board, []byte{}, eval.name, "", 0, 0, workers, seenNodesSplit, speedDisplay)
		}
		os.Exit(0)
	}
	fmt.Println("Board is :", board)
	fmt.Println("Now starting with :", eval.name)
	start := time.Now()
	algoResult := Result{}
	if !noIterativeDepth {
		fmt.Println("Search Method : IDA*")
		data := idaData{}
		data.maxScore = eval.fx(board, board, goal(len(board)), []byte{})
		data.states = append(data.states, Deep2DSliceCopy(board))
		hash, _, _ := matrixToStringSelector(board, 1, 1)
		data.hashes = append(data.hashes, hash)
		data.fx = eval.fx
		data.goal = goal(len(board))
		algoResult = iterateIDA(&data)
	} else {
		fmt.Println("Search Method : A*")
		data := initData(board, workers, seenNodesSplit)
		algoResult = iterateAlgo(board, workers, seenNodesSplit, eval.fx, &data)
	}
	end := time.Now()
	elapsed := end.Sub(start)
	if algoResult.path != nil {
		fmt.Println("Succes with :", eval.name, "in ", elapsed.String(), "!")
		fmt.Printf("len of solution : %v, time complexity / tries : %d, space complexity : %d\n", len(algoResult.path), algoResult.tries, closedSetComplexity)
		if !disableUI {
			displayBoard(board, algoResult.path, eval.name, elapsed.String(), algoResult.tries, closedSetComplexity, workers, seenNodesSplit, speedDisplay)
		}
		fmt.Println(string(algoResult.path))
	} else {
		fmt.Println("No solution !")
	}
}
