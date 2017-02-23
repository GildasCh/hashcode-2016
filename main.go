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

type PredictionByImportance []Prediction

func (p PredictionByImportance) Len() int {
	return len(p)
}

func (p PredictionByImportance) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p PredictionByImportance) Less(i, j int) bool {
	scorei := p[i].n * Endpoints[p[i].e].Ld
	scorej := p[j].n * Endpoints[p[j].e].Ld
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

	solve()
}

type Cache struct {
	size int
	v    []int
}

func createCaches() {
	Caches = make([]Cache, C)
	for i := 0; i < C; i++ {
		Caches[i] = Cache{X, nil}
	}
}

func sizeLeft() int {
	ret := 0
	for _, c := range Caches {
		ret += X - c.size
	}
	return ret
}

func inCache(ci int, vi int) bool {
	ret := false
	for _, cvi := range Caches[ci].v {
		if cvi == vi {
			ret = true
		}
	}
	return ret
}

func solve() interface{} {
	//

	// fmt.Printf("Videos: %+v\n", Videos)
	// fmt.Printf("Predictions: %+v\n", Predictions)
	createCaches()

	totalLeft := sizeLeft()
	previousLeft := totalLeft + 1
	counter := 0

	for {
		sort.Sort(PredictionByImportance(Predictions))

		pChosen := Predictions[0]
		// fmt.Printf("pChosen: %+v\n", pChosen)
		// fmt.Println(len(Videos))
		// fmt.Println(Endpoints[pChosen.e])

		minCacheId := -1
		minCacheLat := -1
		for cid, clat := range Endpoints[pChosen.e].Lc {
			// fmt.Printf("Caches: %+v\n", Caches)
			// return 0
			// fmt.Println("cid:", cid, "; Caches:", Caches)
			if Caches[cid].size < Videos[pChosen.v] {
				continue
			}
			if inCache(cid, pChosen.v) {
				continue
			}

			if minCacheLat == -1 || minCacheLat > clat {
				minCacheId = cid
				minCacheLat = clat
			}
		}

		if minCacheId != -1 {
			Caches[minCacheId].v = append(Caches[minCacheId].v, pChosen.v)
			Caches[minCacheId].size -= Videos[pChosen.v]

			Predictions = Predictions[1:]
		}

		previousLeft = totalLeft
		totalLeft = sizeLeft()

		if totalLeft == previousLeft {
			counter++
			if counter > 5 {
				break
			}
		} else {
			counter = 0
		}
	}

	fmt.Fprintf(output, "%d\n", C)
	for ci, c := range Caches {
		fmt.Fprintf(output, "%d", ci)
		for _, vi := range c.v {
			fmt.Fprintf(output, " %d", vi)
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
