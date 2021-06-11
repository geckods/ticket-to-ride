package main

import (
	"reflect"
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
