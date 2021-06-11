package main

import (
	"reflect"
)

const UintSize = 32 << (^uint(0) >> 32 & 1) // 32 or 64


const (
	MaxInt  = 1<<(UintSize-1) - 1 // 1<<31 - 1 or 1<<63 - 1
	MaxUint = 1<<UintSize - 1     // 1<<32 - 1 or 1<<64 - 1
)


func itemExists(arrayType interface{}, item interface{}) bool {
	arr := reflect.ValueOf(arrayType)

	if arr.Kind() != reflect.Slice {
		panic("Invalid data-type")
	}

	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true
		}
	}

	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func normalizeFloatSlice(sl *[]float64){
	sum := float64(0)
	for _,elem := range *sl {
		sum+=elem
	}
	for i,_ := range *sl {
		(*sl)[i]/=sum
	}
}

func scaleFloat(x,a,b float64) float64{
	return a+((b-a)*x)
}
