package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
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
	P  map[int]int // Predictions: id video -> number of view
}

type Prediction struct {
	v, e, n int
}

type Cache struct {
	Endpoints         []int
	EndpointsPNumbers int
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
	return Caches[c[i]].EndpointsPNumbers > Caches[c[j]].EndpointsPNumbers
}

func CacheFromEndpoints(C int, E []Endpoint) []Cache {
	ret := make([]Cache, C)

	for i, e := range E {
		for c, _ := range e.Lc {
			ret[c].Endpoints = append(ret[c].Endpoints, i)
			ret[c].EndpointsPNumbers += len(e.P)
		}
	}
	return ret
}

// Add predictions to E
func AddPredictions() {
	for _, p := range Predictions {
		Endpoints[p.e].P[p.v] = p.n
	}
}

type weightvideo struct {
	idvideo int
	weight  int
}

type ByWeight struct {
	weights map[int]int
	idvids  []int
}

func (a *ByWeight) Len() int           { return len(a.idvids) }
func (a *ByWeight) Swap(i, j int)      { a.idvids[i], a.idvids[j] = a.idvids[j], a.idvids[i] }
func (a *ByWeight) Less(i, j int) bool { return a.weights[a.idvids[i]] > a.weights[a.idvids[j]] }

func interestingVids(idcache int) (idvids []int) {
	byweight := &ByWeight{make(map[int]int), nil}

	for _, iEndpoint := range Caches[idcache].Endpoints {
		// from Predictions, extract the videos for a given endpoint
		e := &Endpoints[iEndpoint]
		for idvideo, n := range e.P {
			byweight.weights[idvideo] += n * (bestRoute(iEndpoint, idvideo) - e.Lc[idcache]) / Videos[idvideo]
		}
	}

	for iv, _ := range byweight.weights {
		byweight.idvids = append(byweight.idvids, iv)
	}

	sort.Sort(byweight)

	// fmt.Println(byweight)
	return byweight.idvids
}

var BestRoutes = make(map[string]int)

func addRoute(ei, vi, lat int) {
	key := strconv.Itoa(ei) + "," + strconv.Itoa(vi)
	if current, ok := BestRoutes[key]; !ok || (ok && current > lat) {
		BestRoutes[key] = lat
	}
}

func bestRoute(ei, vi int) int {
	key := strconv.Itoa(ei) + "," + strconv.Itoa(vi)
	if lat, ok := BestRoutes[key]; ok {
		return lat
	}
	return Endpoints[ei].Ld
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
		Endpoints[i] = Endpoint{Ld, C, make(map[int]int)}
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
		addRoute(ei, vi, Endpoints[ei].Lc[ci])
	}
}

func solve() interface{} {
	//

	AddPredictions()

	// fmt.Printf("V %d E %d R %d C %d X %d\n\n", V, E, R, C, X)
	// fmt.Printf("Videos: %v\n\nEndpoints: %+v\n\nPredictions: %+v\n\n", Videos, Endpoints, Predictions)

	Caches = CacheFromEndpoints(C, Endpoints)
	// fmt.Printf("Caches: %v\n\n", Caches)

	cacheSorted := CacheInts()
	sort.Sort(CacheByEndpoint(cacheSorted))
	// fmt.Printf("Caches order: %v\n\n", cacheSorted)

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
