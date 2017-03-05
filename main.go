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

type Endpoint struct {
	Ld int         // Latency to datacenter
	Lc map[int]int // Caches: id -> Lc
}

type Prediction struct {
	v, e, n int
}

type Cache struct {
	size int
	v    []int
}

var Caches []Cache

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

type Gain struct {
	value        int
	video, cache int
}

var Gains []Gain

type GainByImportance []Gain

func (g GainByImportance) Len() int { return len(g) }

func (g GainByImportance) Swap(i, j int) { g[i], g[j] = g[j], g[i] }

func (g GainByImportance) Less(i, j int) bool { return g[i].value > g[j].value }

func calculateGains() {
	for _, p := range Predictions {
		e := Endpoints[p.e]
		cls := CacheLatency{}
		for ci, lc := range e.Lc {
			cls.id = append(cls.id, ci)
			cls.lat = append(cls.lat, lc)
		}
		sort.Sort(cls)
		currentLat := e.Ld
		for i := 0; i < cls.Len(); i++ {
			diff := currentLat - cls.lat[i]
			if diff < 0 {
				continue
			}
			Gains = append(Gains, Gain{p.n * (currentLat - cls.lat[i]), p.v, cls.id[i]})
			currentLat = cls.lat[i]
		}
	}

	// fmt.Println(Gains)
	sort.Sort(GainByImportance(Gains))
	// fmt.Println(Gains)
}

type CacheLatency struct {
	id, lat []int
}

func (cl CacheLatency) Len() int { return len(cl.id) }

func (cl CacheLatency) Swap(i, j int) {
	cl.id[i], cl.id[j] = cl.id[j], cl.id[i]
	cl.lat[i], cl.lat[j] = cl.lat[j], cl.lat[i]
}

func (cl CacheLatency) Less(i, j int) bool { return cl.lat[i] < cl.lat[j] }

func calculateTotalLat() int {
	tot := 0
	for _, p := range Predictions {
		e := Endpoints[p.e]
		bestLat := e.Ld
		for ci, lc := range e.Lc {
			if lc < bestLat && inCache(ci, p.v) {
				bestLat = lc
			}
		}
		tot += bestLat * p.n
	}
	return tot
}

var fileOut string
var bFrom, bTo int

func main() {
	sample := os.Args[1]
	fileIn := sample + ".in"
	fileOut = sample + ".out"

	bFroms := os.Args[2]
	bTos := os.Args[3]
	bFrom, _ = strconv.Atoi(bFroms)
	bTo, _ = strconv.Atoi(bTos)

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

var totalGain = 0

var bestLat = -1
var bestBiais = -1

func solve() interface{} {
	for ei, e := range Endpoints {
		fmt.Println(ei, e.Lc)
	}
	return nil

	// fmt.Printf("Videos: %+v\n", Videos)
	// fmt.Printf("Endpoints: %+v\n", Endpoints)
	// fmt.Printf("Predictions: %+v\n", Predictions)

	// totalLeft := sizeLeft()
	// previousLeft := totalLeft + 1
	// counter := 0

	calculateGains()
	// return nil
	for obiais := bFrom; obiais < bTo; obiais++ {
		biais := obiais
		totalGain = 0
		createCaches()

		for _, g := range Gains {
			if Caches[g.cache].size-Videos[g.video] < 0 {
				continue
			}
			if inCache(g.cache, g.video) {
				continue
			}

			if biais%2 == 1 {
				biais /= 2
				continue
			}
			biais /= 2

			Caches[g.cache].v = append(Caches[g.cache].v, g.video)
			Caches[g.cache].size -= Videos[g.video]
			totalGain += g.value
		}

		totalLatency := calculateTotalLat()
		if bestLat == -1 || totalLatency < bestLat {
			fmt.Println("New best Biais:", obiais, "Total Lat:", totalLatency)
			bestLat = totalLatency
			bestBiais = obiais
			output, _ = os.Create(fileOut)
			fmt.Fprintf(output, "%d\n", C)
			for ci, c := range Caches {
				fmt.Fprintf(output, "%d", ci)
				for _, vi := range c.v {
					fmt.Fprintf(output, " %d", vi)
				}
				fmt.Fprintf(output, "\n")
			}
		}
	}

	biais := bestBiais
	totalGain = 0
	createCaches()

	for _, g := range Gains {
		if Caches[g.cache].size-Videos[g.video] < 0 {
			continue
		}
		if inCache(g.cache, g.video) {
			continue
		}

		if biais%2 == 1 {
			biais /= 2
			continue
		}

		Caches[g.cache].v = append(Caches[g.cache].v, g.video)
		Caches[g.cache].size -= Videos[g.video]
		totalGain += g.value
	}

	output, _ = os.Create(fileOut)
	fmt.Fprintf(output, "%d\n", C)
	for ci, c := range Caches {
		fmt.Fprintf(output, "%d", ci)
		for _, vi := range c.v {
			fmt.Fprintf(output, " %d", vi)
		}
		fmt.Fprintf(output, "\n")
	}

	fmt.Println("Best biais:", bestBiais, bestLat)
	fmt.Println("Total Gain:", totalGain)

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
