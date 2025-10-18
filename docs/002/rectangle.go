package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Rectangle struct {
	Triangle [2]Polygon
	Normal   Vector3
	U        Vector3
	V        Vector3

	mp Matrix4
	mv Matrix4
}

// 01 ←左上の点から右回りに4つ
// 32
// vは頂点、tはテクスチャ座標
func NewRectangle(v0, v1, v2, v3 Vector3, t0, t1, t2, t3 Vector2, i, h *ebiten.Image) *Rectangle {
	x := v1.Sub(v0).Normalize()
	c := v2.Sub(v0).Normalize()
	y := v2.Sub(v1).Normalize()
	z := c.Cross(x).Normalize() // 法線

	p := [2]Polygon{
		NewPolygon(v0, v1, v2, t0, t1, t2, i, h),
		NewPolygon(v0, v2, v3, t0, t2, t3, i, h),
	}

	r := &Rectangle{
		Triangle: p,
		Normal:   z,
		U:        x,
		V:        y,
	}

	return r
}

// 射影変換
func (r *Rectangle) ProjTransform(m Matrix4) {
	r.Triangle[0].ProjTransform(m)
	r.Triangle[1].ProjTransform(m)
	r.mp = m
}

// ビューポート変換
func (r *Rectangle) ViewportTransform(m Matrix4) {
	r.Triangle[0].ViewportTransform(m)
	r.Triangle[1].ViewportTransform(m)
	r.mv = m
}

// カリング計算
func (r *Rectangle) CalcCulling() {
	r.Triangle[0].CalcCulling()
	r.Triangle[1].CalcCulling()
}

// ポリゴン接空間でのライトのベクトル算出。ほんとはポリゴン単位でやりたい
func (r *Rectangle) CalcLightVector(light Vector3) {
	z := r.Normal
	u := r.U
	v := r.V

	lu := u.Dot(light)
	lv := v.Dot(light)
	lz := z.Dot(light)
	l := Vector3{lu, lv, lz}.Normalize()

	r.Triangle[0].lightv = l
	r.Triangle[1].lightv = l
}

// デバッグ用描画
func (r *Rectangle) DrawNormal(screen *ebiten.Image) {
	z := r.Normal.Add((r.Triangle[0].base[0].Add(r.Triangle[0].base[2])).Mulf(0.5))
	u := r.U.Add((r.Triangle[0].base[0].Add(r.Triangle[0].base[2])).Mulf(0.5))
	v := r.V.Add((r.Triangle[0].base[0].Add(r.Triangle[0].base[2])).Mulf(0.5))

	zz := z.MulMatToVector4(r.mp)
	uu := u.MulMatToVector4(r.mp)
	vv := v.MulMatToVector4(r.mp)

	zz = zz.MulMat(r.mv)
	uu = uu.MulMat(r.mv)
	vv = vv.MulMat(r.mv)
	zz.X /= zz.W
	zz.Y /= zz.W
	uu.X /= uu.W
	uu.Y /= uu.W
	vv.X /= vv.W
	vv.Y /= vv.W

	x := (r.Triangle[0].Vec[0].X + r.Triangle[0].Vec[2].X) / 2
	y := (r.Triangle[0].Vec[0].Y + r.Triangle[0].Vec[2].Y) / 2
	vector.StrokeLine(screen, x, y, zz.X, zz.Y, 2, color.RGBA{255, 0, 0, 255}, true)
	vector.StrokeLine(screen, x, y, uu.X, uu.Y, 2, color.RGBA{255, 0, 0, 255}, true)
	vector.StrokeLine(screen, x, y, vv.X, vv.Y, 2, color.RGBA{255, 0, 0, 255}, true)
}

// 描画
func (r *Rectangle) Draw(screen *ebiten.Image, vertices []ebiten.Vertex, indices []uint16) ([]ebiten.Vertex, []uint16) {
	for i := 0; i < 2; i++ {
		if r.Triangle[i].Culling {
			continue
		}

		num := uint16(len(vertices))
		indices = append(indices, num, num+1, num+2)

		for j := 0; j < 3; j++ {
			var invW float32
			if PerspectiveCorrection {
				invW = 1 / r.Triangle[i].Vec[j].W
			} else {
				invW = 1
			}
			vertex := ebiten.Vertex{}
			vertex.DstX = r.Triangle[i].Vec[j].X
			vertex.DstY = r.Triangle[i].Vec[j].Y
			vertex.SrcX = r.Triangle[i].texV[j].X * invW // パースペクティブコレクト用にWで割る
			vertex.SrcY = r.Triangle[i].texV[j].Y * invW // パースペクティブコレクト用にWで割る
			vertex.ColorR = r.Triangle[i].lightv.X       // ライトのベクトル
			vertex.ColorG = r.Triangle[i].lightv.Y
			vertex.ColorB = r.Triangle[i].lightv.Z
			vertex.Custom3 = invW // パースペクティブコレクト用
			vertices = append(vertices, vertex)
		}
	}
	return vertices, indices
}
