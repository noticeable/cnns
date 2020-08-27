package tensor

import (
	"math"
)

/*
	Simple math operations for Tensor type
*/

// Add - Element-wise summation.
func (t1 *Tensor) Add(t2 *Tensor) (*Tensor, error) {
	var ret = NewTensor(t1.Size.X, t1.Size.Y, t1.Size.Z)
	ok := t1.IsEqualDims(t2)
	if !ok {
		return ret, ErrDimensionsAreNotEqual
	}
	for i := 0; i < t2.Size.Total(); i++ {
		ret.Data[i] = t1.Data[i] + t2.Data[i]
	}
	return ret, nil
}

// Sub - Element-wise subtraction.
func (t1 *Tensor) Sub(t2 *Tensor) (*Tensor, error) {
	var ret = NewTensor(t1.Size.X, t1.Size.Y, t1.Size.Z)
	ok := t1.IsEqualDims(t2)
	if !ok {
		return ret, ErrDimensionsAreNotEqual
	}
	for i := 0; i < t2.Size.Total(); i++ {
		ret.Data[i] = t1.Data[i] - t2.Data[i]
	}
	return ret, nil
}

// MSE - Mean square error
func (t1 *Tensor) MSE(t2 *Tensor) float64 {
	sum := 0.0
	var ret = NewTensorCopy(t1)
	num := t2.Size.Total()
	for i := 0; i < num; i++ {
		ret.Data[i] = math.Pow((ret.Data[i] - t2.Data[i]), 2.0)
		sum += ret.Data[i]
	}
	return sum / float64(num)
}

// Transpose Transponse tensor by X and Y axis (2D). See ref. https://en.wikipedia.org/wiki/Transpose
func (t1 *Tensor) Transpose() *Tensor {
	ret := NewTensor(t1.Size.Y, t1.Size.X, t1.Size.Z)
	for z := 0; z < t1.Size.Z; z++ {
		for y := 0; y < t1.Size.Y; y++ {
			for x := 0; x < t1.Size.X; x++ {
				ret.Set(y, x, z, t1.Get(x, y, z))
			}
		}
	}
	return ret
}

// Product Product Tensor by Tensor (by X and Y axis, 2D). See ref. https://en.wikipedia.org/wiki/Matrix_multiplication
func (t1 *Tensor) Product(t2 *Tensor) (*Tensor, error) {
	ret := NewTensor(t2.Size.X, t1.Size.Y, t1.Size.Z)
	if t1.Size.Z != t2.Size.Z || t1.Size.X != t2.Size.Y {
		return ret, ErrDimensionsNotFit
	}
	for z := 0; z < t1.Size.Z; z++ {
		for y := 0; y < t1.Size.Y; y++ {
			for x := 0; x < t2.Size.X; x++ {
				var e float64
				for i := 0; i < t1.Size.X; i++ {
					e += t1.Get(i, y, z) * t2.Get(x, i, z)
				}
				ret.Set(x, y, z, e)
			}
		}
	}
	return ret, nil
}

// HadamardProduct Element-wise product. See ref. https://en.wikipedia.org/wiki/Hadamard_product_(matrices)
func HadamardProduct(t1, t2 *Tensor) (*Tensor, error) {
	ret := NewTensor(t1.Size.X, t1.Size.Y, t1.Size.Z)
	ok := t1.IsEqualDims(t2)
	if !ok {
		return ret, ErrDimensionsAreNotEqual
	}
	for i := range ret.Data {
		ret.Data[i] = t1.Data[i] * t2.Data[i]
	}
	return ret, nil
}

/* @experimental
func Reshape(a []float64, newshape ... int) (Tensor, error) {
	indexes := TDsize { X: 0, Y: 1, Z: 2, }
	sizes := TDsize { X: 1, Y: 1, Z: 1, }
	switch len(newshape) {
		case 0:
		case 1:
			if newshape[indexes.X] > 0 {
				sizes.X = newshape[indexes.X]
			} else {
				return NewTensor(sizes.X, sizes.Y, sizes.Z), errors.New("Invalid argument")
			}
		case 2:
			if newshape[indexes.X] > 0 && newshape[indexes.Y] > 0 {
				sizes.X = newshape[indexes.X]
				sizes.Y = newshape[indexes.Y]
			} else {
				return NewTensor(sizes.X, sizes.Y, sizes.Z), errors.New("Invalid argument")
			}
		case 3:
			if newshape[indexes.X] > 0 && newshape[indexes.Y] > 0 && newshape[indexes.Z] > 0 {
				sizes.X = newshape[indexes.X]
				sizes.Y = newshape[indexes.Y]
				sizes.Z = newshape[indexes.Z]
			} else {
				return NewTensor(sizes.X, sizes.Y, sizes.Z), errors.New("Invalid argument")
			}
		default:
			return NewTensor(sizes.X, sizes.Y, sizes.Z), errors.New("Invalid argument")
	}
	ret := NewTensor(sizes.X, sizes.Y, sizes.Z)
	for z := 0; z < ret.Size.Z; z++ {
		for y := 0; y < ret.Size.Y; y++ {
			for x := 0; x < ret.Size.X; x++ {
				if x + y * ret.Size.X + z * ret.Size.X * ret.Size.Y < len(a) {
					ret.Set(x, y, z, a[x + y * ret.Size.X + z * ret.Size.X * ret.Size.Y])
				} else {
					ret.Set(x, y, z, 0)
				}
			}
		}
	}
	return ret, nil
}
*/
