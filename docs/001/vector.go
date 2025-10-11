package main

import "math"

type Vector3 struct {
	X, Y, Z float32
}

type Vector4 struct {
	X, Y, Z, W float32
}

func NewVecor3(x, y, z float32) Vector3 {
	return Vector3{X: x, Y: y, Z: z}
}

func (v Vector3) Add(o Vector3) Vector3 {
	v.X += o.X
	v.Y += o.Y
	v.Z += o.Z
	return v
}

func (v Vector3) Sub(o Vector3) Vector3 {
	v.X -= o.X
	v.Y -= o.Y
	v.Z -= o.Z
	return v
}

func (v Vector3) Magnitude() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
}

func (v Vector3) Normalize() Vector3 {
	l := v.Magnitude()
	if l < 1e-8 || l == 1 {
		return v
	}
	v.X, v.Y, v.Z = v.X/l, v.Y/l, v.Z/l
	return v
}

func (v Vector3) Dot(o Vector3) float32 {
	return v.X*o.X + v.Y*o.Y + v.Z*o.Z
}

func (v Vector3) Cross(o Vector3) Vector3 {

	x := v.X
	y := v.Y
	z := v.Z

	v.X = y*o.Z - z*o.Y
	v.Y = z*o.X - x*o.Z
	v.Z = x*o.Y - y*o.X

	return v
}

func (v Vector3) RotateX(angle float32) Vector3 {
	a := float64(angle)
	y := float64(v.Y)
	z := float64(v.Z)

	v.Y = float32(y*math.Cos(a) - z*math.Sin(a))
	v.Z = float32(y*math.Sin(a) + z*math.Cos(a))
	return v
}

func (v Vector3) RotateY(angle float32) Vector3 {
	a := float64(angle)
	x := float64(v.X)
	z := float64(v.Z)

	v.X = float32(x*math.Cos(a) + z*math.Sin(a))
	v.Z = float32(-x*math.Sin(a) + z*math.Cos(a))
	return v
}

func (v Vector3) RotateZ(angle float32) Vector3 {
	a := float64(angle)
	x := float64(v.X)
	y := float64(v.Y)

	v.X = float32(x*math.Cos(a) - y*math.Sin(a))
	v.Y = float32(x*math.Sin(a) + y*math.Cos(a))
	return v
}

func (v Vector3) MulMat(m Matrix4) Vector3 {

	return Vector3{
		X: m[0][0]*v.X + m[1][0]*v.Y + m[2][0]*v.Z + m[3][0],
		Y: m[0][1]*v.X + m[1][1]*v.Y + m[2][1]*v.Z + m[3][1],
		Z: m[0][2]*v.X + m[1][2]*v.Y + m[2][2]*v.Z + m[3][2],
	}
}

func (v Vector3) MulMatToVector4(m Matrix4) Vector4 {

	return Vector4{
		X: m[0][0]*v.X + m[1][0]*v.Y + m[2][0]*v.Z + m[3][0],
		Y: m[0][1]*v.X + m[1][1]*v.Y + m[2][1]*v.Z + m[3][1],
		Z: m[0][2]*v.X + m[1][2]*v.Y + m[2][2]*v.Z + m[3][2],
		W: m[0][3]*v.X + m[1][3]*v.Y + m[2][3]*v.Z + m[3][3],
	}
}

func (v Vector4) MulMat(m Matrix4) Vector4 {

	return Vector4{
		X: m[0][0]*v.X + m[1][0]*v.Y + m[2][0]*v.Z + m[3][0]*v.W,
		Y: m[0][1]*v.X + m[1][1]*v.Y + m[2][1]*v.Z + m[3][1]*v.W,
		Z: m[0][2]*v.X + m[1][2]*v.Y + m[2][2]*v.Z + m[3][2]*v.W,
		W: m[0][3]*v.X + m[1][3]*v.Y + m[2][3]*v.Z + m[3][3]*v.W,
	}
}

func (v Vector4) ToVec3() Vector3 {
	return Vector3{v.X, v.Y, v.Z}
}

func (v Vector4) Divf(f float32) Vector4 {
	return Vector4{v.X / f, v.Y / f, v.Z / f, v.W / f}
}
