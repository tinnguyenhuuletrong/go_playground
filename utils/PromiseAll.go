package utils

import (
	"context"
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

type promiseAllResult[R any] struct {
	allRes []R
	err    error
}

func PromiseAllWithCtx[T any, R any](ctx context.Context, inputs []T, target func(T) R) ([]R, error) {

	finishChan := make(chan promiseAllResult[R])
	defer close(finishChan)

	go func() {
		allRes, err := PromiseAll(inputs, target)
		finishChan <- promiseAllResult[R]{
			allRes: allRes,
			err:    err,
		}
	}()

	select {
	case res := <-finishChan:
		{
			return res.allRes, res.err
		}
	case <-ctx.Done():
		{
			return nil, fmt.Errorf("ctx done was called")
		}
	}
}
