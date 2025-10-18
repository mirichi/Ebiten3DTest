package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sphere struct {
	Center    Vector3
	Radius    float32
	RadiusP   Vector3
	Vec       Vector4
	VecRadius Vector4
	shader    *ebiten.Shader
	image     *ebiten.Image
	Camera    *Camera
}

func NewSphere() *Sphere {
	shaderText := []byte(
		`//kage:unit pixels
		package main
		var Radius float
		var Light vec3
		var ScreenPos vec2

		func Fragment(dstPos vec4, srcPos vec2, vc, custom vec4) vec4 {
			px := srcPos - imageSrc0Origin()

			if distance(px, vec2(Radius,Radius)) >= Radius {
				discard()
			}

			// 法線を計算する
			ndc := (px - vec2(Radius, Radius)) / Radius
			z := -sqrt(1.0 - ndc.x * ndc.x - ndc.y * ndc.y)
			n := vec3(ndc.x, ndc.y, z)
			nv := n * vec3(-1.0, 1.0, 1.0) // 左右反転

			// 環境マッピング用テクスチャ座標計算
			a := (nv * (50 / -z)).xy + vec2(Radius, Radius)

			color := vec4(0.0,0.0,0.0,1)
			if a.x >= 0 && a.x < Radius * 2 && a.y >= 0 && a.y < Radius * 2 {
				color = imageSrc0UnsafeAt(imageSrc0Origin() + a.xy)
			}

			// 反射ベクトル
			r := reflect(Light,n)
			s := dot(r, vec3(0.0,0.0,-1.0))

			color.xyz += 0.3
			color.xyz *= dot(-Light, n) * 0.5 + 0.5
			color.xyz *= pow(clamp(s, 0, 1), 20.0) * 0.5 + 0.5 // スペキュラ
			return color
			// return vec4(1,1,1,1)
		}
		`,
	)

	sh, err := ebiten.NewShader(shaderText)
	if err != nil {
		panic(err)
	}

	i := ebiten.NewImage(500, 500)
	camera := NewCamera(
		Vector3{3000, 0, 14000},
		Vector3{0, 0, 0},
		Vector3{0, -1, 0},
		NewProjectionPerspective(50, 30000, 0, float32(i.Bounds().Dx()), 0, float32(i.Bounds().Dy())),
		NewViewport(float32(i.Bounds().Dx()), float32(i.Bounds().Dy())),
	)

	s := Sphere{
		Center:  Vector3{},
		Radius:  3000,
		RadiusP: Vector3{3000, 0, 0},
		shader:  sh,
		image:   i,
		Camera:  camera,
	}

	return &s
}

func (s *Sphere) GetModelMatrix() Matrix4 {
	m := NewMatrix4Translate(3000, 0, 14000)
	return m
}

func (s *Sphere) ProjTransform(m Matrix4) {
	s.Vec = s.Center.MulMatToVector4(m)
	s.VecRadius = s.RadiusP.MulMatToVector4(m)
}

func (s *Sphere) ViewportTransform(m Matrix4) {
	s.Vec = s.Vec.MulMat(m)
	s.Vec.X /= s.Vec.W
	s.Vec.Y /= s.Vec.W
	s.Vec.Z /= s.Vec.W
	s.VecRadius = s.VecRadius.MulMat(m)
	s.VecRadius.X /= s.VecRadius.W
	s.VecRadius.Y /= s.VecRadius.W
	s.VecRadius.Z /= s.VecRadius.W
}

func (s *Sphere) CalcCulling() {
}

func (s *Sphere) Draw(screen *ebiten.Image, light Vector3) {
	// 環境マッピング用テクスチャ描画
	s.DrawEnvironmentMapping(screen, light)

	// 半径算出
	radius := s.VecRadius.ToVec3().Sub(s.Vec.ToVec3()).Magnitude()

	// ポリゴンを2枚くっつけて円が入るサイズの矩形を描画する
	b := []Vector2{
		{s.Vec.X - radius, s.Vec.Y - radius},
		{s.Vec.X + radius, s.Vec.Y - radius},
		{s.Vec.X + radius, s.Vec.Y + radius},
		{s.Vec.X - radius, s.Vec.Y + radius},
	}
	t := []Vector2{
		{0, 0},
		{float32(s.image.Bounds().Dx()), 0},
		{float32(s.image.Bounds().Dx()), float32(s.image.Bounds().Dy())},
		{0, float32(s.image.Bounds().Dy())},
	}
	indices := []uint16{
		0, 1, 2,
		0, 2, 3,
	}

	vertices := []ebiten.Vertex{}
	for i := 0; i < 4; i++ {
		vertex := ebiten.Vertex{}
		vertex.DstX = b[i].X
		vertex.DstY = b[i].Y
		vertex.SrcX = t[i].X
		vertex.SrcY = t[i].Y
		vertices = append(vertices, vertex)
	}

	op := &ebiten.DrawTrianglesShaderOptions{}
	op.AntiAlias = true
	op.FillRule = ebiten.FillRuleFillAll
	op.Images[0] = s.image
	op.Uniforms = map[string]any{
		"Radius":    float32(s.image.Bounds().Dx() / 2),
		"Light":     [3]float32{light.X, light.Y, light.Z},
		"ScreenPos": [2]float32{s.Vec.X, s.Vec.Y},
	}
	screen.DrawTrianglesShader(vertices, indices, s.shader, op)
}

// 環境マップ描画
func (s *Sphere) DrawEnvironmentMapping(screen *ebiten.Image, light Vector3) {
	s.image.Fill(color.RGBA{0, 0, 0, 255})
	s.Camera.Draw(s.image, dice, light)
	// screen.DrawImage(s.image, nil)
}
