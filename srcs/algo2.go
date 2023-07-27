package main

import (
	"fmt"
	"os"
)

func iterateIDA(data *idaData) (result Result) {
	fmt.Fprintln(os.Stderr, "Selected ALGO : IDA*")
	for data.maxScore < 1<<30 {
		fmt.Fprintln(os.Stderr, "Cut off is now :", data.maxScore)
		newMaxScore, found := IDA(data)
		if found {
			return Result{data.path, data.closedSetComplexity, data.tries, data.ramFailure}
		}
		data.maxScore = newMaxScore
	}
	data.path = nil
	return
}

func IDA(data *idaData) (newMaxScore int, found bool) {
	currentState := data.states[len(data.states)-1]
	score := data.fx(currentState, data.states[0], data.goal, data.path)
	data.tries++
	if data.tries > 0 && data.tries%100000 == 0 {
		fmt.Fprintf(os.Stderr, "%d * 100k tries\n", data.tries/100000)
	}
	if currentComplexity := len(data.states); currentComplexity > data.closedSetComplexity {
		data.closedSetComplexity = currentComplexity
	}
	if score > data.maxScore {
		return score, false
	}
	if isEqual(currentState, data.goal) {
		return -1, true
	}
	minScoreAboveCutOff := 1 << 30
	for _, dir := range directions {
		ok, nextPos := dir.fx(currentState)
		if !ok {
			continue
		}
		nextHash, _, _ := matrixToStringSelector(nextPos, 1, 1)
		if index := Index(data.hashes, nextHash); index != -1 {
			continue
		}
		data.path = append(data.path, dir.name)
		data.states = append(data.states, nextPos)
		data.hashes = append(data.hashes, nextHash)

		newMaxScore, found := IDA(data)
		if found {
			return newMaxScore, true
		}
		if newMaxScore < minScoreAboveCutOff {
			minScoreAboveCutOff = newMaxScore
		}
		data.path = data.path[:len(data.path)-1]
		data.states = data.states[:len(data.states)-1]
		data.hashes = data.hashes[:len(data.hashes)-1]
	}
	return minScoreAboveCutOff, false
}
