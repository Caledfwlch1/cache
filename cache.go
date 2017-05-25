package main

import (
	"fmt"
	"flag"
	"encoding/csv"
	"os"
	"log"
	"io"
	"strconv"
	"math"
	//"math/rand"
	"sort"
)

var longCache	int
var fiName		string

func init() {
	flag.StringVar(&fiName, "f", "input.csv", "input file name in csv-format (default: input.csv)")
	flag.IntVar(&longCache, "l", 10, "length of cache (default: 10)")
}

type cachType struct {
	i	int
	l	float64
	num	int
	t	int
}

type slCachType []cachType

func (c slCachType)detectEmpty() int {
	j := 0 // len(c) - 1 // rand.Intn(longCache) //
	minI := c[j].i
	minL := c[j].l
	minT := c[j].t
	k := j
	for i, n := range c {
		if n.i == 0 && n.l == 0 && n.t == 0 {
			return i
		}

		if minL > n.l {
			minI = n.i
			minL = n.l
			minT = n.t
			k = i
			continue
		}

		if minI > n.i {
			minI = n.i
			minL = n.l
			minT = n.t
			k = i
			continue
		}

		if minT > n.t {
			minI = n.i
			minL = n.l
			minT = n.t
			k = i
			continue
		}
	}

	return k
}

func (c slCachType)detectEmpty1() int {
	j := 0 // len(c) - 1 // rand.Intn(longCache) //
	minI := c[j].i
	minL := c[j].l
	minT := c[j].t
	k := j
	for i, n := range c {
		if n.i == 0 && n.l == 0 && n.t == 0 {
			return i
		}

		if minL > n.l {
			minI = n.i
			minL = n.l
			minT = n.t
			k = i
			continue
		}

		if minL == n.l {
			if minI > n.i {
				minI = n.i
				minL = n.l
				minT = n.t
				k = i
				continue
			}
			if minI == n.i && minT > n.t {
				minI = n.i
				minL = n.l
				minT = n.t
				k = i
				continue
			}
		}

		if minI > n.i {
			minI = n.i
			minL = n.l
			minT = n.t
			k = i
			continue
		}

		if minI == n.i && minT > n.t {
			minI = n.i
			minL = n.l
			minT = n.t
			k = i
			continue

		}

		if minT >= n.t {
			minI = n.i
			minL = n.l
			minT = n.t
			k = i
			continue
		}
	}

	return k
}

func (c slCachType)newElement(i, k, n, t int) {
	c[i] = cachType{k, 0, n, t}
	return
}

func (c slCachType)increaceFirstElement(i, t int) {
	c[i].i ++
	c[i].t = t
	return
}

func (c slCachType)increaceSecondElement(i, t int) {
	a1 := c[i].i
	a2 := c[i].l
	c[i].l = math.Sqrt( (float64(a1*a1) + a2*a2)/2 )
	c[i].t = t
	// c[i].i = 0
	return
}

func (c slCachType)increaceSecondElement1(i, t int) {
	a1 := c[i].i
	a2 := c[i].l
	c[i].l = (float64(a1) + a2)/2
	c[i].t = t
	return
}

func (c slCachType)increaceSecondElement2(i, t int) {
	a1 := c[i].i
	if a1 == 0 { a1 = 1 }
	a2 := c[i].l
	if a2 == 0 { a2 = 1 }

	c[i].l = 2/ (1/float64(a1) + 1/a2)
	c[i].t = t
	return
}

func (c slCachType)increaceSecondElement3(i, t int) {
	a1 := c[i].i
	a2 := c[i].l
	c[i].l = math.Sqrt(float64(a1) + a2)
	c[i].t = t
	return
}

func (c slCachType)Len() int { return len(c) }

func (c slCachType)Swap(i, j int) { c[i], c[j] = c[j], c[i] }

func (c slCachType)Less(i, j int) bool {
	return c[i].num < c[j].num
}

func main() {

	flag.Parse()

	fi, err := os.Open(fiName)
	if err != nil { log.Fatalf("Error opening file %s - \n", fiName, err) }

	ch := make(slCachType, longCache, longCache)

	enc := csv.NewReader(fi)
	m := 1
	miss := 0
	goal := 0

	hist := make(map[int]int)
	for {
		st, err := enc.Read()
		if err == io.EOF { break }
		if err != nil { log.Fatalf("Error parsing file %s - \n", fiName, err) }

		time, _	:= strconv.Atoi(st[0])
		blk, _	:= strconv.Atoi(st[1])
		//num, _	:= strconv.Atoi(st[2])

		tmp := hist[blk]
		hist[blk] = tmp +1

		found := false
		for i, n := range ch {

			if blk == n.num {

				ch.increaceFirstElement(i, time)
				ch.increaceSecondElement(i, time)
				found = true
				goal ++

				break
			}
		}
		if !found {

			miss++
			k := ch.detectEmpty()
			ch.newElement(k, hist[blk], blk, time)
			ch.increaceSecondElement(k, time)

		}

		//if m % 1000 == 0 {
		//	for i, n := range ch {
		//		ch.increaceSecondElement(i, n.t)
		//	}
		//}

		if m % 10000 == 0 {
			//fmt.Println(ch)
			for i, j := range hist {
				if j < 2 { delete(hist, i) }
			}
			//fmt.Println(m/100, ch)
		}

		m++
	}

	sort.Sort(ch)
	fmt.Println(ch)
	fmt.Println("miss=", miss, "goal=", goal, "len(hist)=", len(hist))
	fmt.Println("effect=", float64(goal)/float64(miss))

	return
}

func (c slCachType)String() string {
	s := ""
	for _, l := range c {
		m := fmt.Sprint(l.l)
		if len(m) > 4 { m = m[:4] }
		s += fmt.Sprintf("{ %d, %s, %6d, %3d }\n", l.i, m, l.num, l.t)
	}
	return s
}