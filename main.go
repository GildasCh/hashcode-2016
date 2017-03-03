package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gildasch/go-algos/priorityqueue"
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

type Endpoint struct {
	Ld int         // Latency to datacenter
	Lc map[int]int // Caches: id -> Lc
	P  map[int]int // Predictions: vi -> n
}

type Prediction struct {
	v, e, n int
}

type Cache struct {
	size int
	v    []int
	e    map[int]int
}

var Caches []Cache

func createCaches() {
	Caches = make([]Cache, C)
	for i := 0; i < C; i++ {
		Caches[i] = Cache{X, nil, make(map[int]int)}
	}

	// Adding endpoints
	for ei, e := range Endpoints {
		for ci, lc := range e.Lc {
			Caches[ci].e[ei] = lc
		}
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

type Gain struct {
	value        int
	video, cache int
}

var Gains map[string]Gain
var HighestGains *priorityqueue.PriorityQueue

// type GainByImportance []Gain

// func (g GainByImportance) Len() int { return len(g) }

// func (g GainByImportance) Swap(i, j int) { g[i], g[j] = g[j], g[i] }

// func (g GainByImportance) Less(i, j int) bool { return g[i].value < g[j].value }

func calculateGains(ciu, viu int) {
	newGains := make(map[string]Gain)
	if HighestGains == nil {
		HighestGains = &priorityqueue.PriorityQueue{}
	}

	for ci, c := range Caches {
		es := c.e
		for vi, v := range Videos {
			if v > c.size || inCache(ci, vi) {
				continue
			}

			key := strconv.Itoa(ci) + "," + strconv.Itoa(vi)
			g := 0

			if ciu != -1 && ci != ciu && vi != viu {
				newGains[key] = Gains[key]
				continue
			}

			for e, elc := range es {
				if n, ok := Endpoints[e].P[vi]; ok {
					diff := n * (Endpoints[e].Ld - elc)
					if diff < 0 {
						continue
					}
					g += diff
				}
			}
			gain := Gain{g, vi, ci}
			newGains[key] = gain
			HighestGains.Add(key, gain.value)
		}
	}

	Gains = newGains
}

// func updateGains(ciu, viu int) {
// 	HighestGain = Gain{-1, 0, 0}

// 	for ci, c := range Caches {
// 		es := c.e
// 		for vi, v := range Videos {
// 			if ci != ciu && vi != viu {
// 				continue
// 			}
// 			if v > c.size || inCache(ci, vi) {
// 				continue
// 			}

// 			g := 0
// 			for e, elc := range es {
// 				if n, ok := Endpoints[e].P[vi]; ok {
// 					diff := n * (Endpoints[e].Ld - elc)
// 					if diff < 0 {
// 						continue
// 					}
// 					g += diff
// 				}
// 			}
// 			gain := Gain{g, vi, ci}
// 			Gains = append(Gains, gain)

// 			if g > HighestGain.value {
// 				HighestGain = gain
// 			}
// 		}
// 	}
// }

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
		Endpoints[i] = Endpoint{Ld, C, make(map[int]int)}
	}

	// Predictions
	Predictions = make([]Prediction, R)
	for i := 0; i < R; i++ {
		v := readInt()
		e := readInt()
		n := readInt()
		Predictions[i] = Prediction{v, e, n}
		Endpoints[e].P[v] = n
	}

	solve()
}
func solve() interface{} {
	// fmt.Printf("Videos: %+v\n", Videos)
	// fmt.Printf("Endpoints: %+v\n", Endpoints)
	// fmt.Printf("Predictions: %+v\n", Predictions)

	createCaches()

	totalLeft := sizeLeft()
	previousLeft := totalLeft + 1
	counter := 0

	calculateGains(-1, -1)
	// fmt.Println(HighestGains)
	// fmt.Println(Gains)
	for {
		keyi, _ := HighestGains.Pop()
		fmt.Println("keyi:", keyi)
		if keyi == -1 {
			break
		}
		key, _ := keyi.(string)
		HighestGain := Gains[key]
		// fmt.Println(HighestGain)

		if Caches[HighestGain.cache].size-Videos[HighestGain.video] < 0 {
			continue
		}

		// fmt.Println("HighestGain:", HighestGain)
		// fmt.Println("Gains:", Gains)
		Caches[HighestGain.cache].v = append(Caches[HighestGain.cache].v, HighestGain.video)
		Caches[HighestGain.cache].size -= Videos[HighestGain.video]
		if Caches[HighestGain.cache].size < 0 {
			fmt.Println()
			fmt.Println()
			fmt.Println(Caches)
			fmt.Println(HighestGain)
			return nil
		}

		previousLeft = totalLeft
		totalLeft = sizeLeft()

		fmt.Println(totalLeft, "/", X*C)
		// fmt.Println(Caches)
		if totalLeft == previousLeft {
			// fmt.Println("Caches:", Caches)
			// fmt.Println("HighestGain:", HighestGain)
			counter++
			if counter > 5 {
				break
			}
		} else {
			counter = 0
		}

		calculateGains(HighestGain.cache, HighestGain.video)
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
