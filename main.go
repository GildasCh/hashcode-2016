package main

import (
	"fmt"
	"os"
	"sort"
)

var input *os.File
var output *os.File

var V int
var E int
var R int
var C int
var X int
var Videos []int
var Endpoints []Endpoint
var Predictions []Prediction
var Caches []Cache

type Endpoint struct {
	Ld int         // Latency to datacenter
	Lc map[int]int // Caches: id -> Lc
}

type Prediction struct {
	v, e, n int
}

type Cache struct {
	vi int
}

type PredictionByImportance []Prediction

func (p PredictionByImportance) Len() int {
	return len(p)
}

func (p PredictionByImportance) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p PredictionByImportance) Less(i, j int) bool {
	scorei := p[i].n * Endpoints[p[i].e]
	scorej := p[j].n * Endpoints[p[j].e]
	return scorei < scorej
}

func main() {
	sample := os.Args[1]
	fileIn := sample + ".in"
	fileOut := sample + ".out"

	var err error
	input, err = os.Open(fileIn)
	if err != nil {
		panic(fmt.Sprintf("open %s: %v", fileIn, err))
	}
	output, err = os.Create(fileOut)
	if err != nil {
		panic(fmt.Sprintf("creating %s: %v", fileOut, err))
	}
	defer input.Close()
	defer output.Close()

	// Global
	V = readInt()
	E = readInt()
	R = readInt()
	C = readInt()
	X = readInt()

	// Videos
	Videos = make([]int, V)
	for i := 0; i < V; i++ {
		Videos[i] = readInt()
	}

	// Endpoints
	Endpoints = make([]Endpoint, E)
	for i := 0; i < E; i++ {
		Ld := readInt()
		K := readInt()
		C := make(map[int]int)
		for j := 0; j < K; j++ {
			cid := readInt()
			C[cid] = readInt()
		}
		Endpoints[i] = Endpoint{Ld, C}
	}

	// Predictions
	Predictions = make([]Prediction, R)
	for i := 0; i < R; i++ {
		v := readInt()
		e := readInt()
		n := readInt()
		Predictions[i] = Prediction{v, e, n}
	}

	solve(V, E, R, C, X, Videos, Endpoints, Predictions)
}

func removeVidFromEnpoints(ci int, vi int) {
	for _, ei := range Caches[ci].Endpoints {
		delete(Endpoints[ei].P, vi)
	}
}

func solve() interface{} {
	//

	sort.Sort(PredictionByImportance(Predictions))

	fmt.Fprintf(output, "%d\n", C)
	for _, ci := range cacheSorted {
		fmt.Fprintf(output, "%d", ci)
		iVids := interestingVids(ci)
		// fmt.Printf("Interesting vids: %v\n\n", iVids)
		sizeCache := X
		for _, iv := range iVids {
			if Videos[iv] > sizeCache {
				continue
			}
			fmt.Fprintf(output, " %d", iv)
			sizeCache -= Videos[iv]
			// Remove from endpoints
			removeVidFromEnpoints(ci, iv)
		}
		fmt.Fprintf(output, "\n")
	}

	return 0
}

func readInt() int {
	var i int
	fmt.Fscanf(input, "%d", &i)
	return i
}

func readString() string {
	var str string
	fmt.Fscanf(input, "%s", &str)
	return str
}

func readFloat() float64 {
	var x float64
	fmt.Fscanf(input, "%f", &x)
	return x
}
