package main

import "fmt"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func checkPlus(e error, extra string) {
	if e != nil {
		fmt.Printf("Extra: %s", extra)
		panic(e)
	}
}

// type ErrorPlus struct {
// 	err   error
// 	extra string
// }

// func NewErrorPlus(err error, extra string) *ErrorPlus {
// 	if err == nil {
// 		return nil
// 	}
// 	return &ErrorPlus{err, extra}
// }

// func (e *ErrorPlus) Error() string {
// 	return fmt.Sprintf("Info: %s\n%s", e.extra, e.err.Error())
// }
