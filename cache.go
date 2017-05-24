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
	j := 0 // rand.Intn(longCache)
	minI := c[j].i
	minL := c[j].l
	minT := c[j].t
	k := 0
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

func (c *slCachType)newElement(i, n, t int) {
	(*c)[i] = cachType{1, 0, n, t}
	return
}

func (c *slCachType)increaceFirstElement(i, t int) {
	(*c)[i].i ++
	(*c)[i].t = t
	return
}

func (c *slCachType)increaceSecondElement(i, t int) {
	a1 := (*c)[i].i
	a2 := (*c)[i].l
	(*c)[i].l = math.Sqrt( (float64(a1*a1) + a2*a2)/2 )
	(*c)[i].t = t
	return
}

func (c *slCachType)increaceSecondElement1(i, t int) {
	a1 := (*c)[i].i
	a2 := (*c)[i].l
	(*c)[i].l = (float64(a1) + a2)/2
	(*c)[i].t = t
	return
}

func (c *slCachType)increaceSecondElement2(i, t int) {
	a1 := (*c)[i].i
	if a1 == 0 { a1 = 1 }
	a2 := (*c)[i].l
	if a2 == 0 { a2 = 1 }

	(*c)[i].l = 2/ (1/float64(a1) + 1/a2)
	(*c)[i].t = t
	return
}


func main() {

	flag.Parse()

	fi, err := os.Open(fiName)
	if err != nil { log.Fatalf("Error opening file %s - \n", fiName, err) }

	ch := make(slCachType, longCache, longCache)

	enc := csv.NewReader(fi)
	m := 1
	for {
		st, err := enc.Read()
		if err == io.EOF { break }
		if err != nil { log.Fatalf("Error parsing file %s - \n", fiName, err) }

		time, _	:= strconv.Atoi(st[0])
		blk, _	:= strconv.Atoi(st[1])
		//num, _	:= strconv.Atoi(st[2])


		found := false
		for i, n := range ch {

			if blk == n.num {

				ch.increaceFirstElement(i, time)
				ch.increaceSecondElement(i, time)
				found = true

				break
			}
		}
		if !found {
			k := ch.detectEmpty()
			ch.newElement(k, blk, time)
			ch.increaceSecondElement(k, time)
		}

		//if m %500 == 0 {
		//	//fmt.Println(ch)
		//	for i, _ := range ch {
		//		//ch.increaceSecondElement(i)
		//		ch[i].i = 0
		//	}
		//	//fmt.Println(m/100, ch)
		//}


		m++
	}

	fmt.Println(ch)

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