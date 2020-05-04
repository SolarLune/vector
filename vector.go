package vector

import (
	"errors"
	"fmt"
	"math"
)

// Vector is the definition of a row vector that contains scalars as
// 64 bit floats
type Vector []float64

// Axis is an integer enum type that describes vector axis
type Axis int

const (
	// the consts below are used to represent vector axis, they are useful
	// to lookup values within the vector.
	X Axis = iota
	Y
	Z
)

var (
	// ErrNot3Dimensional is an error that is returned in functions that only
	// supports 3 dimensional vectors
	ErrNot3Dimensional   = errors.New("vector is not 3 dimensional")
	// ErrNotSameDimensions is an error that is returned when functions need both
	// Vectors provided to be the same dimensionally
	ErrNotSameDimensions = errors.New("the two vectors provided aren't the same dimensional size")
)

// Clone a vector
func Clone(v Vector) Vector {
	return v.Clone()
}

// Clone a vector
func (v Vector) Clone() Vector {
	clone := make(Vector, len(v))
	copy(clone, v)
	return clone
}

// Add a vector with a vector or a set of vectors
func Add(v1 Vector, vs ...Vector) Vector {
	return v1.Clone().Add(vs...)
}

// Add a vector with a vector or a set of vectors
func (v Vector) Add(vs ...Vector) Vector {
	dim := len(v)

	for i := range vs {
		if len(vs[i]) > dim {
			axpyUnitaryTo(v, 1, v, vs[i][:dim])
		} else {
			axpyUnitaryTo(v, 1, v, vs[i])
		}
	}

	return v
}

// Sub subtracts a vector with another vector or a set of vectors
func Sub(v1 Vector, vs ...Vector) Vector {
	return v1.Clone().Sub(vs...)
}

// Sub subtracts a vector with another vector or a set of vectors
func (v Vector) Sub(vs ...Vector) Vector {
	dim := len(v)

	for i := range vs {
		if len(vs[i]) > dim {
			axpyUnitaryTo(v, -1, vs[i][:dim], v)
		} else {
			axpyUnitaryTo(v, -1, vs[i], v)
		}
	}

	return v
}

// Scale vector with a given size
func Scale(v Vector, size float64) Vector {
	return v.Clone().Scale(size)
}

// Scale vector with a given size
func (v Vector) Scale(size float64) Vector {
	scalUnitaryTo(v, size, v)
	return v
}

// Equal compares that two vectors are equal to each other
func Equal(v1, v2 Vector) bool {
	return v1.Equal(v2)
}

// Equal compares that two vectors are equal to each other
func (v Vector) Equal(v2 Vector) bool {
	if len(v) != len(v2) {
		return false
	}

	for i := range v {
		if math.Abs(v[i]-v2[i]) > 1e-8 {
			return false
		}
	}

	return true
}

// Magnitude of a vector
func Magnitude(v Vector) float64 {
	return v.Magnitude()
}

// Magnitude of a vector
func (v Vector) Magnitude() float64 {
	var result float64

	for _, scalar := range v {
		result += scalar * scalar
	}

	return math.Sqrt(result)
}

// Unit returns a direction vector with the length of one.
func Unit(v Vector) Vector {
	return v.Clone().Unit()
}

// Unit returns a direction vector with the length of one.
func (v Vector) Unit() Vector {
	l := v.Magnitude()

	if l < 1e-8 {
		return v
	}

	for i := range v {
		v[i] = v[i] / l
	}

	return v
}

// Dot product of two vectors
func Dot(v1, v2 Vector) float64 {
	result, dim1, dim2 := 0., len(v1), len(v2)

	if dim1 > dim2 {
		v2 = append(v2, make(Vector, dim1-dim2)...)
	}

	if dim1 < dim2 {
		v1 = append(v1, make(Vector, dim2-dim1)...)
	}

	for i := range v1 {
		result += v1[i] * v2[i]
	}

	return result
}

// Dot product of two vectors
func (v Vector) Dot(v2 Vector) float64 {
	return Dot(v, v2)
}

// Cross product of two vectors
func Cross(v1, v2 Vector) (Vector, error) {
	return v1.Cross(v2)
}

// Cross product of two vectors
func (v Vector) Cross(v2 Vector) (Vector, error) {
	if len(v) != 3 || len(v2) != 3 {
		return nil, ErrNot3Dimensional
	}

	return Vector{
		v[Y]*v2[Z] - v2[Y]*v[Z],
		v[Z]*v2[X] - v2[Z]*v[X],
		v[X]*v2[Y] - v2[X]*v[Y],
	}, nil
}

// Rotate is rotating a vector around a specified axis.
// If no axis are specified, it will default to the Z axis.
//
// If a vector with more than 3-dimensions is rotated, it will cut the extra
// dimensions and return a 3-dimensional vector.
//
// NOTE: the ...Axis is just syntactic sugar that allows the axis to not be
// specified and default to Z, if multiple axis is passed the first will be
// set as the rotation axis
func Rotate(v Vector, angle float64, as ...Axis) Vector {
	return v.Clone().Rotate(angle, as...)
}

// Rotate is rotating a vector around a specified axis.
// If no axis are specified, it will default to the Z axis.
//
// If a vector with more than 3-dimensions is rotated, it will cut the extra
// dimensions and return a 3-dimensional vector.
//
// NOTE: the ...Axis is just syntactic sugar that allows the axis to not be
// specified and default to Z, if multiple axis is passed the first will be
// set as the rotation axis
func (v Vector) Rotate(angle float64, as ...Axis) Vector {
	axis, dim := Z, len(v)

	if dim == 0 {
		return v
	}

	if len(as) > 0 {
		axis = as[0]
	}

	if dim == 1 && axis != Z {
		v = append(v, 0, 0)
	}

	if (dim < 2 && axis == Z) || (dim == 2 && axis != Z) {
		v = append(v, 0)
	}

	x, y := v[X], v[Y]

	cos, sin := math.Cos(angle), math.Sin(angle)

	switch axis {
	case X:
		z := v[Z]
		v[Y] = y*cos - z*sin
		v[Z] = y*sin + z*cos
	case Y:
		z := v[Z]
		v[X] = x*cos + z*sin
		v[Z] = -x*sin + z*cos
	case Z:
		v[X] = x*cos - y*sin
		v[Y] = x*sin + y*cos
	}

	if dim > 3 {
		return v[:3]
	}

	return v
}

// Angle returns the angle in radians from the first Vector to the second, the Vector of rotation, and an error if the
// two Vectors aren't of equal dimensions (length). For 0-dimension Vectors, the angle is 0, and the rotation Vector
// is empty. For 1-dimension Vectors, the angle is 0 if they both have the same sign, and pi if they don't.
// The Vector of rotation is, again, empty. For 2-dimension Vectors, the Vector of rotation is a Unit Vector in the Z
// direction (0, 0, 1).
func Angle(v1, v2 Vector) (float64, Vector, error) {
	return v1.Angle(v2)
}

// Angle returns the angle in radians from the first Vector to the second, the Vector of rotation, and an error if the
// two Vectors aren't of equal dimensions (length). For 0-dimension Vectors, the angle is 0, and the rotation Vector
// is empty. For 1-dimension Vectors, the angle is 0 if they both have the same sign, and pi if they don't.
// The Vector of rotation is, again, empty. For 2-dimension Vectors, the Vector of rotation is a Unit Vector in the Z
// direction (0, 0, 1).
func (v Vector) Angle(v2 Vector) (float64, Vector, error) {

	dim := len(v)
	dim2 := len(v2)
	zeroVec := make(Vector, 0)

	if dim != dim2 {
		return 0, zeroVec, ErrNotSameDimensions
	}

	if dim == 0 {
		return 0, zeroVec, nil
	}
	if dim == 1 {
		if (v[0] > 0 && v2[0] < 0) || (v[0] < 0 && v2[0] > 0) {
			return math.Pi, zeroVec, nil
		}
		return 0, zeroVec, nil
	} else if dim == 2 {
		return (math.Atan2(v2.Y(), v2.X()) - math.Atan2(v.Y(), v.X())), Vector{0, 0, 1}, nil
	}

	// 3 or more dimensions
	angle := math.Acos(Dot(v.Clone().Unit(), v2.Clone().Unit()))
	axis, _ := Cross(v, v2)
	axis.Unit()
	return angle, axis, nil

}

// String returns the string representation of a vector
func (v Vector) String() (str string) {
	if v == nil {
		return "[]"
	}

	for i := range v {
		if v[i] < 1e-8 && v[i] > 0 {
			str += "0 "
		} else {
			str += fmt.Sprint(v[i]) + " "
		}

	}

	return "[" + str[:len(str)-1] + "]"
}

// X is corresponding to doing a v[0] lookup, if index 0 does not exist yet, a
// 0 will be returned instead
func (v Vector) X() float64 {
	if len(v) < 1 {
		return 0.
	}

	return v[X]
}

// Y is corresponding to doing a v[1] lookup, if index 1 does not exist yet, a
// 0 will be returned instead
func (v Vector) Y() float64 {
	if len(v) < 2 {
		return 0.
	}

	return v[Y]
}

// Z is corresponding to doing a v[2] lookup, if index 2 does not exist yet, a
// 0 will be returned instead
func (v Vector) Z() float64 {
	if len(v) < 3 {
		return 0.
	}

	return v[Z]
}
