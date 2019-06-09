// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pdqsort

import (
	"math"
	"math/rand"
	stdsort "sort"
	"testing"
	"github.com/stretchr/testify/assert"
)

var intsTriggerMaxInsertion = [...]int{74, 59, 238, -784, 9845, 959, 905, 0, 0, 42, 7586, -5467984, 7586}
var intsLessThanOneBlock = [...]int{74, 59, 238, -784, 9845, 959, 905, 0, 0, 42, 7586, -5467984, 7586, 19, 13, 27, 103, 42, 23, 29, 31, 51, 10}
var intsSlightMoreThanOneBlock = [...]int{90337,88478,26565,58029,97724,21943,74586,69778,38123,96767,62447,28961,29183,35336,54815,25791,67154,12981,37441,61888,50918,20343,28477,44445,18070,48789,21573,40409,66596,96451,6231,83381,77153,46297,4258,74942,91850,34493,54317,49711,22500,21353,93441,21251,67018,69925,68045,10752,10335,7649,60327,74202,23945,75596,64937,1485,6246,48012,52374,52426,29995,22322,92172,59959,5321,32484,5256,7126,71306,27399,19513,75199,83464,19069,12954,85341,21235,19425,52653,4306,82569,34531,64719,48199,858,1951,85560,45017,427,96208,62045,20635,44607,85578,17793,97095,14269,44351,68591,80670,68872,14928,57226,75167,23886,82035,96756,90760,13192,37199,73043,48940,96878,71939,49360,99496,77306,6813,92131,20280,14474,64530,54156,93176,76113,52879,64552,52065,23053,439,51936,13551,91191,58436,71272,17751,92107,67330,71522,97552,98338,2472,51371,65428,22605,20904,51144,4647,41523,65802}
var intsMultipleBlocks = [...]int{60811,7470,28654,90705,91164,89132,51441,26090,96814,85314,87768,58739,73401,72817,65401,89469,62030,57712,43551,11153,29916,51198,54122,29950,37423,79969,16741,9422,24222,16987,67416,23228,33271,5519,94740,67764,47394,65829,38585,89684,55002,72608,70302,70354,58349,20718,22803,27193,40576,4088,32243,11516,70001,15809,15157,58343,9164,76093,52485,573,86657,91328,71370,10744,34931,60423,96875,37247,38507,31919,22209,17197,13855,96321,26489,29366,22499,89886,96160,49518,79941,84110,5905,67203,55089,20610,65712,16071,95308,22189,17626,17223,97462,28450,51547,83597,81889,25597,89217,52865,842,32987,93924,59943,17622,75920,5098,4212,44510,36220,37098,34100,25604,79553,88567,3938,73767,36700,95881,49223,16077,70556,71164,22412,92794,43835,8840,52554,72966,68710,56175,1671,878,33390,35918,96576,49816,76419,25377,18224,29990,1076,89228,1124,22284,23903,24043,42351,67533,68298,51462,98415,73892,49241,33277,57387,15932,57004,45463,43021,3339,83412,22896,68741,72164,55648,42311,47132,72037,10454,65515,57793,70310,93740,18156,9651,53849,24721,82253,82558,28346,85573,94265,55266,14323,85976,47756,64322,90892,16403,28132,62798,19609,81904,8556,49184,72161,20960,74452,38893,17246,53885,14747,27985,96565,2717,24412,31158,81927,72517,71675,70205,69861,5999,53863,24384,60310,81712,45849,45444,52977,99562,57672,80297,83039,42911,47627,86939,26654,12703,71230,86359,75555,17159,27825,87537,52362,43001,70075,79468,6752,30736,13288,40688,99294,91188,40299,11755,43459,83191,82990,81152,68676,88515,1079,17396,95476,432,59574,57579,48822,30396,62617,80819,31792,64667,38804,87491,6553,70349,39342,54321,90236,69549,21912,94588,41531,21051,29438,64312,42953,18086,63318,82556,5408,46549,58065,32139,13487,50340,18078,22233,43670,14636,23920,79765,90071,13516,12295,20731,25628,95452,82124,27800,50627,42468,2865,4278,86272,64158,71053,32891,44127,88053,25307,97138,57291,65922,92467,19167,31361,37635,11965,1849,74851,24147,20896,54473,36127,92631,39290,58540,43924,81862,65987,76819,98874,49951,72040,3879,24773,93927,21350,90350,64190,27004,93550,16071,44584,41645,20440,74304,53460,297,74068,43704,67572,57224,10902,29946,6686,20071,30949,35582,22721,9277,67712,4565,65027,4701,79458,80356,93988,27594,61590,54529,84380,65390,82115,33353,68176,9017,63480,78417,93691,28889,89048,43223,22826,89285,37020,55082,58325,43049,56879,47063,21402,10908,44105,264,49943,13266,75857,98763,19011,57105,80822,41764,64463,36181,66668,19753,42213,8002,81003,48780,22814,25326,41055,88826,27490,72653,59203,57312,10373,19633,5317,22502,16641,28171,87777,66910,96763,65189,4444,90680,87024,58964,61461,62172,6613,65466,21284,25617,87154,32330,75266,32757,74962,91197,50113,69348,85681,46486,61833,93855,56459,46410,33433,58749,32373,93569,33376,16024,34665,46665,91667,14021,48239,43260,67544,83320,70256,76493,70067,92575,24121,9412,86794,13732,90824,71944,57659,69140,10776,92446,8436,93109,55240,87597,34159,61058,6549,49913,15210,45190,15704,71093,28375,42260}

var float64s = [...]float64{74.3, 59.0, math.Inf(1), 238.2, -784.0, 2.3, math.NaN(), math.NaN(), math.Inf(-1), 9845.768, -959.7485, 905, 7.8, 7.8}
var strings = [...]string{"", "Hello", "foo", "bar", "foo", "f00", "%*&^*&^&", "***"}

func TestShiftTailIntSlice(t *testing.T) {
	data := [8]int{2,3,4,5,6,7,8,1}
	answer := [8]int{1,2,3,4,5,6,7,8}

	a := IntSlice(data[0:])
	a.ShiftTail(0, 8)
	b := IntSlice(answer[0:])

	assert.Equal(t, b, a, "not the same")
}

func TestShiftHeadIntSlice(t *testing.T) {
	data := [8]int{8,1,2,3,4,5,6,7}
	answer := [8]int{1,2,3,4,5,6,7,8}

	a := IntSlice(data[0:])
	a.ShiftHead(0, 8)
	b := IntSlice(answer[0:])

	assert.Equal(t, b, a, "not the same")
}

func TestCyclicSwapIntSlice(t *testing.T) {
	data := [8]int{1,2,3,4,5,6,7,8}
	answer := [8]int{5,6,8,4,2,3,7,1}
	
	a := IntSlice(data[0:])
	a.CyclicSwaps([]int{0,1,2}, []int{4,5,7})
	b := IntSlice(answer[0:])

	assert.Equal(t, b, a, "not the same")
}

func TestInsertionSortIntSlice(t *testing.T) {
	data := intsLessThanOneBlock
	a := IntSlice(data[0:])
	insertionSort(a, 0, len(a))
	if !IsSorted(a) {
		t.Errorf("sorted %v", intsLessThanOneBlock)
		t.Errorf("   got %v", data)
	}
}

func TestHeapSortIntSlice(t *testing.T) {
	data := intsLessThanOneBlock
	a := IntSlice(data[0:])
	heapSort(a, 0, len(a))
	if !IsSorted(a) {
		t.Errorf("sorted %v", intsLessThanOneBlock)
		t.Errorf("   got %v", data)
	}
}

func TestSortIntSliceByTriggeringInsertionSort(t *testing.T) {
	data := intsTriggerMaxInsertion
	a := IntSlice(data[0:])
	Sort(a)
	if !IsSorted(a) {
		t.Errorf("sorted %v", intsLessThanOneBlock)
		t.Errorf("   got %v", data)
	}
}

func TestSortIntSliceLessThanOneBlock(t *testing.T) {
	data := intsLessThanOneBlock
	a := IntSlice(data[0:])
	Sort(a)
	if !IsSorted(a) {
		t.Errorf("sorted %v", intsLessThanOneBlock)
		t.Errorf("   got %v", data)
	}
}


func TestSortIntSliceSlightMoreThanOneBlock(t *testing.T) {
	data := intsSlightMoreThanOneBlock
	a := IntSlice(data[0:])
	Sort(a)
	if !IsSorted(a) {
		t.Errorf("sorted %v", intsSlightMoreThanOneBlock)
		t.Errorf("   got %v", data)
	}
}

func TestSortIntSliceMultipleBlocks(t *testing.T) {
	data := intsMultipleBlocks
	a := IntSlice(data[0:])
	Sort(a)
	if !IsSorted(a) {
		t.Errorf("sorted %v", intsMultipleBlocks)
		t.Errorf("   got %v", data)
	}
}

func TestSortFloat64Slice(t *testing.T) {
	data := float64s
	a := Float64Slice(data[0:])
	Sort(a)
	if !IsSorted(a) {
		t.Errorf("sorted %v", float64s)
		t.Errorf("   got %v", data)
	}
}

func TestSortStringSlice(t *testing.T) {
	data := strings
	a := StringSlice(data[0:])
	Sort(a)
	if !IsSorted(a) {
		t.Errorf("sorted %v", strings)
		t.Errorf("   got %v", data)
	}
}

func TestIntsLessThanOneBlock(t *testing.T) {
	data := intsLessThanOneBlock
	Ints(data[0:])
	if !IntsAreSorted(data[0:]) {
		t.Errorf("sorted %v", intsLessThanOneBlock)
		t.Errorf("   got %v", data)
	}
}

func TestFloat64s(t *testing.T) {
	data := float64s
	Float64s(data[0:])
	if !Float64sAreSorted(data[0:]) {
		t.Errorf("sorted %v", float64s)
		t.Errorf("   got %v", data)
	}
}

func TestStrings(t *testing.T) {
	data := strings
	Strings(data[0:])
	if !StringsAreSorted(data[0:]) {
		t.Errorf("sorted %v", strings)
		t.Errorf("   got %v", data)
	}
}

func TestSortLarge_Random(t *testing.T) {
	n := 1000000
	if testing.Short() {
		n /= 100
	}
	data := make([]int, n)
	for i := 0; i < len(data); i++ {
		data[i] = rand.Intn(100)
	}
	if IntsAreSorted(data) {
		t.Fatalf("terrible rand.rand")
	}
	Ints(data)
	if !IntsAreSorted(data) {
		t.Errorf("sort didn't sort - 1M ints")
	}
}

func BenchmarkPDQSortInt1K(b *testing.B) {
  	b.StopTimer()
  	for i := 0; i < b.N; i++ {
    	data := make([]int, 1<<10)
    	for i := 0; i < len(data); i++ {
      		data[i] = i ^ 0x2cc
    	}
    	b.StartTimer()
    	Ints(data)
    	b.StopTimer()
  	}
}

func BenchmarkStdSortInt1K(b *testing.B) {
	b.StopTimer()
	for i := 0; i < b.N; i++ {
	  data := make([]int, 1<<10)
	  for i := 0; i < len(data); i++ {
			data[i] = i ^ 0x2cc
	  }
	  b.StartTimer()
	  stdsort.Ints(data)
	  b.StopTimer()
	}
}
