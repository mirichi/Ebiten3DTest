package main

import "github.com/hajimehoshi/ebiten/v2"

type Camera struct {
	Position              Vector3
	Target                Vector3
	Up                    Vector3
	ProjectionPerspective Matrix4
	Vewport               Matrix4
}

type Object interface {
	GetModelMatrix() Matrix4
	ProjTransform(m Matrix4)
	ViewportTransform(m Matrix4)
	Draw(screen *ebiten.Image, light Vector3)
	CalcCulling()
}

func NewCamera(pos, target, up Vector3, mp, mvp Matrix4) *Camera {
	return &Camera{
		Position:              pos,
		Target:                target,
		Up:                    up,
		ProjectionPerspective: mp,
		Vewport:               mvp,
	}
}

func (c *Camera) Draw(screen *ebiten.Image, o Object, light Vector3) {
	m := o.GetModelMatrix()
	// ビュー変換
	m = m.Mul(NewMatrix4LookAt(c.Position, c.Target, c.Up))
	// 射影変換
	m = m.Mul(c.ProjectionPerspective)
	o.ProjTransform(m)

	// ビューポート変換
	m = c.Vewport
	o.ViewportTransform(m)
	// 背面カリング
	o.CalcCulling()
	// 描画
	o.Draw(screen, light)
}
