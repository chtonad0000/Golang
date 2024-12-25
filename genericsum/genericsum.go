//go:build !solution

package genericsum

import (
	"golang.org/x/exp/constraints"
	"math"
	"sort"
)

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func SortSlice[T constraints.Ordered](a []T) {
	sort.Slice(a, func(i, j int) bool {
		return a[i] < a[j]
	})
}

func MapsEqual[T, V comparable](a, b map[T]V) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if w, ok := b[k]; !ok || v != w {
			return false
		}
	}
	return true
}

func SliceContains[T comparable](s []T, v T) bool {
	for _, x := range s {
		if x == v {
			return true
		}
	}
	return false
}
func MergeChans[T any](chs ...<-chan T) <-chan T {
	out := make(chan T)
	go func() {
		done := make(chan struct{}, len(chs))
		for _, ch := range chs {
			go func(ch <-chan T) {
				for val := range ch {
					out <- val
				}
				done <- struct{}{}
			}(ch)
		}
		for i := 0; i < len(chs); i++ {
			<-done
		}
		close(out)
	}()

	return out
}
func IsHermitianMatrix[T comparable](matrix [][]T) bool {
	if len(matrix) == 0 || len(matrix) != len(matrix[0]) {
		return false
	}
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix); j++ {
			switch v := any(matrix[i][j]).(type) {
			case complex64:
				b := any(matrix[j][i]).(complex64)
				if i == j && imag(v) != 0 {
					return false
				}
				if float64(real(v)) != float64(real(b)) || math.Abs(float64(imag(v))) != math.Abs(float64(imag(b))) {
					return false
				}
				continue
			case complex128:
				b := any(matrix[j][i]).(complex128)
				if i == j && imag(v) != 0 {
					return false
				}
				if real(v) != real(b) || math.Abs(imag(v)) != math.Abs(imag(b)) {
					return false
				}
				continue
			}
			if matrix[i][j] != matrix[j][i] {
				return false
			}
		}
	}
	return true
}
