package main

import (
	"reflect"
)

func itemExists(arrayType interface{}, item interface{}) bool {
	arr := reflect.ValueOf(arrayType)

	if arr.Kind() != reflect.Slice {
		//fmt.Println(arr.Kind())
		//fmt.Println(reflect.Array)
		panic("Invalid data-type")
	}

	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true
		}
	}

	return false
}
