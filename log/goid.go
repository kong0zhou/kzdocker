package log

import (
	"kzdocker/log/g"
	"unsafe"
)

// var goidOffset uintptr

// func init() {
// 	// f, ok := reflect.TypeOf(g.G()).FieldByName(`goid`)
// 	p := (*interface{})(unsafe.Pointer(g.G()))
// 	g := reflect.TypeOf(p)
// 	if g.Kind() == reflect.Ptr {
// 		fmt.Println(g.Kind())
// 		g = g.Elem()
// 	}
// 	fmt.Println(g.Kind())
// 	if g.Kind() != reflect.Struct {
// 		panic(`g is not struct`)
// 	}
// 	f, ok := g.FieldByName(`goid`)
// 	if !ok {
// 		panic(`can not get goid offset`)
// 	}
// 	goidOffset = f.Offset
// }

var goidOffset uintptr = 152 // Go1.8/Go1.9/Go1.10/Go1.11/Go1.12/Go1.13

func getGroutineID() int64 {
	g := g.G()
	p := (*int64)(unsafe.Pointer(uintptr(g) + goidOffset))
	return *p
}
