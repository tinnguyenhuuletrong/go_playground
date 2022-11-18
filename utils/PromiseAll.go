package utils

import (
	"fmt"
	"sync"
)

/**
Understand and rewrite ( not sure better than author )
	https://medium.com/@kikochaudhuri/handling-concurrent-goroutines-f222967f12f4

Leason learned
	Fixed slice CC access is safe - https://klotzandrew.com/blog/concurrent_writing_to_slices_in_go

TODO
	- early stop when error occur
	- accept context -> can stop by callee
	- rate limit
**/

func PromiseAll[T any, R any](inputs []T, target func(T) R) ([]R, error) {
	var wg sync.WaitGroup
	len := len(inputs)
	allRes := make([]R, len)
	var err error = nil

	wg.Add(len)
	for index, v := range inputs {
		activeIndex := index
		activeValue := v
		go func() {
			inp := activeValue

			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("%v", r)
				}
				wg.Done()
			}()

			res := target(inp)
			allRes[activeIndex] = res
		}()
	}

	wg.Wait()

	return allRes, err
}
