package utils

import (
	"context"
	"fmt"
	"sync/atomic"
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

type promiseRes[R any] struct {
	res R
	err error
}

func PromiseRace[T any, R any](inputs []T, target func(T) R) (R, error) {
	var hasDone atomic.Bool
	oneRes := make(chan promiseRes[R], 1)
	hasDone.Store(false)

	defer func() {
		hasDone.Store(true)
		close(oneRes)
	}()

	for _, v := range inputs {
		activeValue := v
		go func() {
			inp := activeValue

			defer func() {
				if r := recover(); r != nil {
					err := fmt.Errorf("%v", r)
					if !hasDone.Load() {
						oneRes <- promiseRes[R]{
							err: err,
						}
					}
				}
			}()

			res := target(inp)
			if !hasDone.Load() {
				oneRes <- promiseRes[R]{
					res: res,
					err: nil,
				}

			}
		}()
	}

	finalRes := <-oneRes
	return finalRes.res, finalRes.err
}

func PromiseRaceWithCtx[T any, R any](ctx context.Context, inputs []T, target func(T) R) (R, error) {
	finishChan := make(chan promiseRes[R])
	defer close(finishChan)

	go func() {
		res, err := PromiseRace(inputs, target)
		finishChan <- promiseRes[R]{
			res: res,
			err: err,
		}
	}()

	select {
	case res := <-finishChan:
		{
			return res.res, res.err
		}
	case <-ctx.Done():
		{
			var noop R
			return noop, fmt.Errorf("ctx done was called")
		}
	}
}
