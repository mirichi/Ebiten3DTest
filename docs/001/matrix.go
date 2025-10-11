package main

import "math"

type Matrix4 [4][4]float32

func NewMatrix4() Matrix4 {
	return Matrix4{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}
}

func NewEmptyMatrix4() Matrix4 {
	return Matrix4{
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
		{0, 0, 0, 0},
	}
}

func NewProjection(zn, zf, left, right, top, bottom float32) Matrix4 {
	return Matrix4{
		{(2 * zn) / (right - left), 0, 0, 0},
		{0, (2 * zn) / (bottom - top), 0, 0},
		{0, 0, zf / (zf - zn), 1},
		{0, 0, -zn * zf / (zf - zn), 0},
	}
}

func NewViewport(width, height float32) Matrix4 {
	return Matrix4{
		{width / 2, 0, 0, 0},
		{0, height / 2, 0, 0},
		{0, 0, 1, 0},
		{width / 2, height / 2, 0, 1},
	}
}

func NewMatrix4Translate(x, y, z float32) Matrix4 {
	mat := NewMatrix4()
	mat[3][0] = x
	mat[3][1] = y
	mat[3][2] = z
	return mat
}

func NewMatrix4Scale(x, y, z float32) Matrix4 {
	mat := NewMatrix4()
	mat[0][0] = x
	mat[1][1] = y
	mat[2][2] = z
	mat[3][3] = 1
	return mat
}

func NewMatrix4Rotate(x, y, z, angle float32) Matrix4 {

	if x == 0 && y == 0 && z == 0 {
		y = 1
	}

	mat := NewMatrix4()
	v := Vector3{X: x, Y: y, Z: z}.Normalize()
	s := float32(math.Sin(float64(angle)))
	c := float32(math.Cos(float64(angle)))
	m := 1 - c

	mat[0][0] = m*v.X*v.X + c
	mat[0][1] = m*v.X*v.Y + v.Z*s
	mat[0][2] = m*v.Z*v.X - v.Y*s

	mat[1][0] = m*v.X*v.Y - v.Z*s
	mat[1][1] = m*v.Y*v.Y + c
	mat[1][2] = m*v.Y*v.Z + v.X*s

	mat[2][0] = m*v.Z*v.X + v.Y*s
	mat[2][1] = m*v.Y*v.Z - v.X*s
	mat[2][2] = m*v.Z*v.Z + c

	return mat
}

func (m Matrix4) Mul(o Matrix4) Matrix4 {
	newM := NewMatrix4()

	newM[0][0] = m[0][0]*o[0][0] + m[0][1]*o[1][0] + m[0][2]*o[2][0] + m[0][3]*o[3][0]
	newM[1][0] = m[1][0]*o[0][0] + m[1][1]*o[1][0] + m[1][2]*o[2][0] + m[1][3]*o[3][0]
	newM[2][0] = m[2][0]*o[0][0] + m[2][1]*o[1][0] + m[2][2]*o[2][0] + m[2][3]*o[3][0]
	newM[3][0] = m[3][0]*o[0][0] + m[3][1]*o[1][0] + m[3][2]*o[2][0] + m[3][3]*o[3][0]

	newM[0][1] = m[0][0]*o[0][1] + m[0][1]*o[1][1] + m[0][2]*o[2][1] + m[0][3]*o[3][1]
	newM[1][1] = m[1][0]*o[0][1] + m[1][1]*o[1][1] + m[1][2]*o[2][1] + m[1][3]*o[3][1]
	newM[2][1] = m[2][0]*o[0][1] + m[2][1]*o[1][1] + m[2][2]*o[2][1] + m[2][3]*o[3][1]
	newM[3][1] = m[3][0]*o[0][1] + m[3][1]*o[1][1] + m[3][2]*o[2][1] + m[3][3]*o[3][1]

	newM[0][2] = m[0][0]*o[0][2] + m[0][1]*o[1][2] + m[0][2]*o[2][2] + m[0][3]*o[3][2]
	newM[1][2] = m[1][0]*o[0][2] + m[1][1]*o[1][2] + m[1][2]*o[2][2] + m[1][3]*o[3][2]
	newM[2][2] = m[2][0]*o[0][2] + m[2][1]*o[1][2] + m[2][2]*o[2][2] + m[2][3]*o[3][2]
	newM[3][2] = m[3][0]*o[0][2] + m[3][1]*o[1][2] + m[3][2]*o[2][2] + m[3][3]*o[3][2]

	newM[0][3] = m[0][0]*o[0][3] + m[0][1]*o[1][3] + m[0][2]*o[2][3] + m[0][3]*o[3][3]
	newM[1][3] = m[1][0]*o[0][3] + m[1][1]*o[1][3] + m[1][2]*o[2][3] + m[1][3]*o[3][3]
	newM[2][3] = m[2][0]*o[0][3] + m[2][1]*o[1][3] + m[2][2]*o[2][3] + m[2][3]*o[3][3]
	newM[3][3] = m[3][0]*o[0][3] + m[3][1]*o[1][3] + m[3][2]*o[2][3] + m[3][3]*o[3][3]

	return newM
}
