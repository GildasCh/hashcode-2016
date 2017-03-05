package main

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"time"
)

var input *os.File
var output *os.File

var V int
var E int
var R int
var C int
var X int
var Videos []int
var Endpoints []*Endpoint
var Predictions []Prediction
var Caches []Cache

type Endpoint struct {
	Ld int         // Latency to datacenter
	Lc map[int]int // Caches: id -> Lc
	P  map[int]int // Predictions: id video -> number of view
	Pl map[int]int // Predictions: id video -> current latency
}

type Prediction struct {
	v, e, n int
}

type Cache struct {
	Endpoints []int
}

func CacheInts() []int {
	ret := make([]int, C)
	for i := 0; i < C; i++ {
		ret[i] = i
	}
	return ret
}

type CacheByEndpoint []int

func (c CacheByEndpoint) Len() int      { return len(c) }
func (c CacheByEndpoint) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c CacheByEndpoint) Less(i, j int) bool {
	return rand.Intn(1)%1 == 0
}

func CacheFromEndpoints(C int, E []*Endpoint) []Cache {
	ret := make([]Cache, C)

	for i, e := range E {
		for c, _ := range e.Lc {
			ret[c].Endpoints = append(ret[c].Endpoints, i)
		}
	}
	return ret
}

// Add predictions to E
func AddPredictions() {
	for _, e := range Endpoints {
		e.P = make(map[int]int)
		e.Pl = make(map[int]int)
	}
	for _, p := range Predictions {
		Endpoints[p.e].P[p.v] = p.n
		Endpoints[p.e].Pl[p.v] = Endpoints[p.e].Ld
	}
}

type weightvideo struct {
	idvideo int
	weight  int
}

type ByWeight []weightvideo

func (a ByWeight) Len() int           { return len(a) }
func (a ByWeight) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByWeight) Less(i, j int) bool { return a[i].weight > a[j].weight }

func removeDuplicates(a []int) []int {
	result := []int{}
	seen := map[int]int{}
	for _, val := range a {
		if _, ok := seen[val]; !ok {
			result = append(result, val)
			seen[val] = val
		}
	}
	return result
}

func interestingVids(idcache int) (idvids []int) {
	var videos []weightvideo

	for _, iEndpoint := range Caches[idcache].Endpoints {
		// from Predictions, extract the videos for a given endpoint
		e := Endpoints[iEndpoint]
		for idvideo, n := range e.P {
			videos = append(videos, weightvideo{idvideo, n * (e.Pl[idvideo] - e.Lc[idcache])})
		}
	}

	sort.Sort(ByWeight(videos))

	for _, iv := range videos {
		idvids = append(idvids, iv.idvideo)
	}

	idvids = removeDuplicates(idvids)

	return idvids
}

var fileOut string

func main() {
	sample := os.Args[1]
	fileIn := sample + ".in"
	fileOut = sample + ".out"

	rand.Seed(time.Now().Unix())

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
	Endpoints = make([]*Endpoint, E)
	for i := 0; i < E; i++ {
		Ld := readInt()
		K := readInt()
		C := make(map[int]int)
		for j := 0; j < K; j++ {
			cid := readInt()
			C[cid] = readInt()
		}
		Endpoints[i] = &Endpoint{Ld, C, nil, nil}
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

func removeVidFromEnpoints(ci int, vi int) {
	for _, ei := range Caches[ci].Endpoints {
		Endpoints[ei].Pl[vi] = Endpoints[ei].Lc[ci]
	}
}

func calculateTotalLat() int {
	tot := 0
	for _, p := range Predictions {
		e := Endpoints[p.e]
		bestLat := e.Ld
		for ci, lc := range e.Lc {
			if lc < bestLat && inOutputCache(ci, p.v) {
				bestLat = lc
			}
		}
		tot += bestLat * p.n
	}
	return tot
}

var bestLat = -1

func solve() interface{} {
	//

	for {
		AddPredictions()

		// fmt.Printf("V %d E %d R %d C %d X %d\n\n", V, E, R, C, X)
		// fmt.Printf("Videos: %v\n\nEndpoints: %+v\n\nPredictions: %+v\n\n", Videos, Endpoints, Predictions)

		Caches = CacheFromEndpoints(C, Endpoints)
		// fmt.Printf("Caches: %v\n\n", Caches)

		cacheSorted := CacheInts()
		sort.Sort(CacheByEndpoint(cacheSorted))
		// fmt.Printf("Caches order: %v\n\n", cacheSorted)

		createOutputCaches()
		for _, ci := range cacheSorted {
			iVids := interestingVids(ci)
			for _, iv := range iVids {
				if Videos[iv] > OutputCaches[ci].size {
					continue
				}
				OutputCaches[ci].v = append(OutputCaches[ci].v, iv)
				OutputCaches[ci].size -= Videos[iv]
				// Remove from endpoints
				removeVidFromEnpoints(ci, iv)
			}
		}

		totalLatency := calculateTotalLat()
		fmt.Printf("Total Lat: %d\r", totalLatency)
		if bestLat == -1 || totalLatency < bestLat {
			fmt.Println("\nNew best total Lat:", totalLatency)
			bestLat = totalLatency
			output, _ = os.Create(fileOut)
			fmt.Fprintf(output, "%d\n", C)
			for ci, c := range OutputCaches {
				fmt.Fprintf(output, "%d", ci)
				for _, vi := range c.v {
					fmt.Fprintf(output, " %d", vi)
				}
				fmt.Fprintf(output, "\n")
			}
		}
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

type OutputCache struct {
	size int
	v    []int
}

var OutputCaches []OutputCache

func createOutputCaches() {
	OutputCaches = make([]OutputCache, C)
	for i := 0; i < C; i++ {
		OutputCaches[i] = OutputCache{X, nil}
	}
}

func sizeLeft() int {
	ret := 0
	for _, c := range OutputCaches {
		ret += X - c.size
	}
	return ret
}

func inOutputCache(ci int, vi int) bool {
	ret := false
	for _, cvi := range OutputCaches[ci].v {
		if cvi == vi {
			ret = true
		}
	}
	return ret
}
