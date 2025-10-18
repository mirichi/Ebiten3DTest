package main

import "github.com/hajimehoshi/ebiten/v2"

type Polygon struct {
	base        [3]Vector3 // 元の頂点情報(右回り)
	texV        [3]Vector2
	tex, height *ebiten.Image
	Vec         [3]Vector4 // 計算後頂点
	Culling     bool       // カリングする
	lightv      Vector3
}

// 右回りに3つ。vは頂点、tはテクスチャ座標
func NewPolygon(v0, v1, v2 Vector3, t0, t1, t2 Vector2, i, h *ebiten.Image) Polygon {
	return Polygon{
		base:   [3]Vector3{v0, v1, v2},
		texV:   [3]Vector2{t0, t1, t2},
		tex:    i,
		height: h,
	}
}

// 射影変換
func (p *Polygon) ProjTransform(m Matrix4) {
	p.Vec[0] = p.base[0].MulMatToVector4(m)
	p.Vec[1] = p.base[1].MulMatToVector4(m)
	p.Vec[2] = p.base[2].MulMatToVector4(m)
}

// ビューポート変換
func (p *Polygon) ViewportTransform(m Matrix4) {
	p.Vec[0] = p.Vec[0].MulMat(m)
	p.Vec[1] = p.Vec[1].MulMat(m)
	p.Vec[2] = p.Vec[2].MulMat(m)

	for i := 0; i < len(p.Vec); i++ {
		p.Vec[i].X /= p.Vec[i].W
		p.Vec[i].Y /= p.Vec[i].W
		p.Vec[i].Z /= p.Vec[i].W
	}
}

// カリング計算
func (p *Polygon) CalcCulling() {
	x1 := p.Vec[1].X - p.Vec[0].X
	y1 := p.Vec[1].Y - p.Vec[0].Y
	x2 := p.Vec[2].X - p.Vec[1].X
	y2 := p.Vec[2].Y - p.Vec[1].Y
	v := Vector2{x2, y2}.Normalize()

	p.Culling = x1*v.Y-v.X*y1 <= 0
}
