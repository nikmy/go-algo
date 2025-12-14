package xmath

import (
	"fmt"

	"golang.org/x/exp/constraints"

	"github.com/nikmy/algo/container/slices"
)

type ErrSizeMismatch [2][2]int

func (e ErrSizeMismatch) Error() string {
	l, r := e[0], e[1]
	return fmt.Sprintf("incompatible sizes (%d, %d) and (%d, %d)", l[0], l[1], r[0], r[1])
}

func NewDiagonalMatrix[T Complex](values ...T) *Matrix[T] {
	m := make(Matrix[T], len(values))
	for i := range m {
		m[i] = make([]T, len(values))
		m[i][i] = values[i]
	}
	return &m
}

func NewIdMatrix[T Complex](n int) *Matrix[T] {
	return NewDiagonalMatrix(slices.Filled[T](1, n)...)
}

func NewZeroMatrix[T Complex](m, n int) *Matrix[T] {
	return NewDiagonalMatrix(slices.Filled[T](0, n)...)
}

type Matrix[T Complex] [][]T

func (m *Matrix[T]) Clone() *Matrix[T] {
	c := make(Matrix[T], len(*m))
	for i, row := range *m {
		c[i] = make([]T, len(row))
		copy(c[i], row)
	}
	return &c
}

func (m *Matrix[T]) Rows() int {
	return len(*m)
}

func (m *Matrix[T]) Cols() int {
	if len(*m) == 0 {
		return 0
	}
	return len((*m)[0])
}

func (m *Matrix[T]) Shape() [2]int {
	return [2]int{m.Rows(), m.Cols()}
}

func (m *Matrix[T]) Elem(i, j int) T {
	return m.Row(i)[j]
}

func (m *Matrix[T]) Row(i int) []T {
	return (*m)[i]
}

func (m *Matrix[T]) Column(j int) []T {
	c := make([]T, m.Rows())
	for i := range c {
		c[i] = m.Elem(i, j)
	}
	return c
}

func (m *Matrix[T]) Diag() []T {
	n := min(m.Rows(), m.Cols())
	d := make([]T, n)
	for i := 0; i < n; i++ {
		d[i] = m.Elem(i, i)
	}
	return d
}

func (m *Matrix[T]) Trace() T {
	if m.Rows() != m.Cols() {
		panic("trace is undefined for non-square matrices")
	}
	return slices.Sum(m.Diag())
}

func (m *Matrix[T]) Det() T {
	if m.Rows() != m.Cols() {
		panic("determinant is undefined for non-square matrices")
	}
	return det(m)
}

func det[T Complex](m *Matrix[T]) T {
	n := m.Rows()
	switch n {
	case 1:
		return m.Elem(0, 0)
	case 2:
		return m.Elem(0, 0)*m.Elem(1, 1) - m.Elem(1, 0)*m.Elem(0, 1)
	}

	d, sign := T(0), T(1)
	cofactor := NewZeroMatrix[T](n-1, n-1)
	for k := 0; k < n; k++ {
		subI := 0
		for i := 0; i < n; i++ {
			subJ := 0
			for j := 0; j < n; j++ {
				if j == k {
					continue
				}
				(*cofactor)[subI][subJ] = m.Elem(i, j)
				subJ++
			}
			subI++
		}
		d += sign * m.Elem(0, k) * det(cofactor)
		sign = -sign
	}
	return d
}

func (m *Matrix[T]) Transposed() *Matrix[T] {
	if m.Rows() == m.Cols() {
		return m.Clone()
	}

	t := make(Matrix[T], m.Cols())
	for i := range t {
		t[i] = make([]T, m.Rows())
		for j := range t[i] {
			t[i][j] = m.Elem(j, i)
		}
	}

	return &t
}

func Hermit[T constraints.Complex](m *Matrix[T]) *Matrix[T] {
	t := m.Transposed()
	for _, row := range *t {
		for j := range row {
			row[j] = Conjugate(row[j])
		}
	}
	return t
}

func (m *Matrix[T]) Add(rhs *Matrix[T]) error {
	if m.Shape() != rhs.Shape() {
		return ErrSizeMismatch{m.Shape(), rhs.Shape()}
	}

	for i, row := range *m {
		for j := range row {
			(*m)[i][j] += rhs.Elem(i, j)
		}
	}

	return nil
}

func (m *Matrix[T]) Multiplied(rhs *Matrix[T]) (*Matrix[T], error) {
	if m.Cols() != rhs.Rows() {
		return nil, ErrSizeMismatch{m.Shape(), rhs.Shape()}
	}

	result := make(Matrix[T], m.Cols())
	for i := range result {
		result[i] = make([]T, rhs.Cols())
		for j := 0; j < rhs.Cols(); j++ {
			for k := 0; k < m.Cols(); k++ {
				result[i][j] += m.Elem(i, k) * rhs.Elem(k, j)
			}
		}
	}

	return &result, nil
}

func (m *Matrix[T]) Pow(exp int) error {
	if m.Rows() != m.Cols() {
		return ErrSizeMismatch{m.Shape(), m.Shape()}
	}

	p, a := NewIdMatrix[T](m.Rows()), m
	for exp > 0 {
		if exp%2 == 0 {
			a, _ = a.Multiplied(a)
			exp /= 2
			continue
		}

		p, _ = p.Multiplied(a)
		exp--
	}

	*m = *p
	return nil
}
