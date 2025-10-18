package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Dice struct {
	Rotate    Matrix4
	diceImage *ebiten.Image
	heightMap *ebiten.Image
	shader    *ebiten.Shader
	rect      []*Rectangle
	box       []Vector3
	tindices  []Vector2
}

func NewDice() *Dice {
	d := &Dice{
		Rotate: NewMatrix4(),
	}

	// パースペクティブコレクト対応シェーダ
	shaderText := []byte(
		`//kage:unit pixels
		package main

		func Fragment(dstPos vec4, srcPos vec2, vc, custom vec4) vec4 {
			originPos := imageSrc0Origin()
			tx := srcPos - originPos
			tx *= 1 / custom.w // ←パースペクティブコレクトの計算
			tx += originPos

			color := imageSrc0UnsafeAt(tx)
			x1 := imageSrc1UnsafeAt(tx + vec2(-1,0)).r * 20 // ハイトマップの影響が少ないので無理やりおおきな数字にしている(ダメな例)
			x2 := imageSrc1UnsafeAt(tx + vec2(1,0)).r * 20
			y1 := imageSrc1UnsafeAt(tx + vec2(0,-1)).r * 20
			y2 := imageSrc1UnsafeAt(tx + vec2(0,1)).r * 20
			du := vec3(1, 0, x2 - x1)
			dv := vec3(0, 1, y2 - y1)
			normal := normalize(cross(du, dv))
			color.xyz *= dot(-vc.xyz, normal) * 0.5 + 0.5

			return color
		}
		`,
	)

	s, err := ebiten.NewShader(shaderText)
	if err != nil {
		panic(err)
	}
	d.shader = s

	// サイコロ画像
	d.diceImage = ebiten.NewImage(1202, 202)
	d.diceImage.Fill(color.White)
	vector.DrawFilledCircle(d.diceImage, 101, 101, 30, color.RGBA{0xff, 0x00, 0x00, 0xff}, true)
	vector.DrawFilledCircle(d.diceImage, 201+50, 51, 20, color.Black, true)
	vector.DrawFilledCircle(d.diceImage, 201+100, 101, 20, color.Black, true)
	vector.DrawFilledCircle(d.diceImage, 201+150, 151, 20, color.Black, true)
	vector.DrawFilledCircle(d.diceImage, 401+50, 51, 20, color.Black, true)
	vector.DrawFilledCircle(d.diceImage, 401+50, 101, 20, color.Black, true)
	vector.DrawFilledCircle(d.diceImage, 401+50, 151, 20, color.Black, true)
	vector.DrawFilledCircle(d.diceImage, 401+150, 51, 20, color.Black, true)
	vector.DrawFilledCircle(d.diceImage, 401+150, 101, 20, color.Black, true)
	vector.DrawFilledCircle(d.diceImage, 401+150, 151, 20, color.Black, true)
	vector.DrawFilledCircle(d.diceImage, 601+50, 51, 20, color.Black, true)
	vector.DrawFilledCircle(d.diceImage, 601+150, 51, 20, color.Black, true)
	vector.DrawFilledCircle(d.diceImage, 601+50, 151, 20, color.Black, true)
	vector.DrawFilledCircle(d.diceImage, 601+150, 151, 20, color.Black, true)
	vector.DrawFilledCircle(d.diceImage, 801+100, 101, 20, color.Black, true)
	vector.DrawFilledCircle(d.diceImage, 801+50, 51, 20, color.Black, true)
	vector.DrawFilledCircle(d.diceImage, 801+150, 51, 20, color.Black, true)
	vector.DrawFilledCircle(d.diceImage, 801+50, 151, 20, color.Black, true)
	vector.DrawFilledCircle(d.diceImage, 801+150, 151, 20, color.Black, true)
	vector.DrawFilledCircle(d.diceImage, 1001+50, 51, 20, color.Black, true)
	vector.DrawFilledCircle(d.diceImage, 1001+150, 151, 20, color.Black, true)

	// ハイトマップ
	d.heightMap = ebiten.NewImage(1202, 202)
	d.heightMap.Fill(color.RGBA{0x80, 0x00, 0x00, 0x00})

	for i := 30; i > 0; i-- {
		y := 128 - uint8(math.Sin(float64(30-i)*3*math.Pi/180.0)*128)
		vector.DrawFilledCircle(d.heightMap, 101, 101, float32(i), color.RGBA{y, 0x00, 0x00, 0xff}, true)
	}

	// 頂点座標
	d.box = []Vector3{
		{X: -1, Y: -1, Z: -1}, {X: 1, Y: -1, Z: -1}, {X: 1, Y: 1, Z: -1}, {X: -1, Y: 1, Z: -1}, // 手前側1
		{X: 1, Y: -1, Z: 1}, {X: -1, Y: -1, Z: 1}, {X: -1, Y: 1, Z: 1}, {X: 1, Y: 1, Z: 1}, // 反対側6
		{X: -1, Y: -1, Z: 1}, {X: -1, Y: -1, Z: -1}, {X: -1, Y: 1, Z: -1}, {X: -1, Y: 1, Z: 1}, // 左側2
		{X: 1, Y: -1, Z: -1}, {X: 1, Y: -1, Z: 1}, {X: 1, Y: 1, Z: 1}, {X: 1, Y: 1, Z: -1}, // 右側5
		{X: -1, Y: -1, Z: 1}, {X: 1, Y: -1, Z: 1}, {X: 1, Y: -1, Z: -1}, {X: -1, Y: -1, Z: -1}, // 上側4
		{X: -1, Y: 1, Z: -1}, {X: 1, Y: 1, Z: -1}, {X: 1, Y: 1, Z: 1}, {X: -1, Y: 1, Z: 1}, // 下側3
	}

	// テクスチャ座標(頂点座標と対応している)
	d.tindices = []Vector2{
		{1, 1}, {201, 1}, {201, 201}, {1, 201}, // 手前側1
		{401, 1}, {601, 1}, {601, 201}, {401, 201}, // 反対側6
		{1001, 1}, {1201, 1}, {1201, 201}, {1001, 201}, // 左側2
		{801, 1}, {1001, 1}, {1001, 201}, {801, 201}, // 右側5
		{601, 1}, {801, 1}, {801, 201}, {601, 201}, // 上側4
		{201, 1}, {401, 1}, {401, 201}, {201, 201}, // 下側3
	}

	// 四角形ポリゴン生成
	d.rect = []*Rectangle{
		NewRectangle(d.box[0], d.box[1], d.box[2], d.box[3], d.tindices[0], d.tindices[1], d.tindices[2], d.tindices[3], d.diceImage, d.heightMap),
		NewRectangle(d.box[4], d.box[5], d.box[6], d.box[7], d.tindices[4], d.tindices[5], d.tindices[6], d.tindices[7], d.diceImage, d.heightMap),
		NewRectangle(d.box[8], d.box[9], d.box[10], d.box[11], d.tindices[8], d.tindices[9], d.tindices[10], d.tindices[11], d.diceImage, d.heightMap),
		NewRectangle(d.box[12], d.box[13], d.box[14], d.box[15], d.tindices[12], d.tindices[13], d.tindices[14], d.tindices[15], d.diceImage, d.heightMap),
		NewRectangle(d.box[16], d.box[17], d.box[18], d.box[19], d.tindices[16], d.tindices[17], d.tindices[18], d.tindices[19], d.diceImage, d.heightMap),
		NewRectangle(d.box[20], d.box[21], d.box[22], d.box[23], d.tindices[20], d.tindices[21], d.tindices[22], d.tindices[23], d.diceImage, d.heightMap),
	}

	return d
}

func (d *Dice) GetModelMatrix() Matrix4 {
	m := d.Rotate.Mul(NewMatrix4Scale(800, 800, 800))
	m = m.Mul(NewMatrix4Translate(-1500, 0, 10000))
	return m
}

func (d *Dice) ProjTransform(m Matrix4) {
	for _, r := range d.rect {
		r.ProjTransform(m)
	}
}

func (d *Dice) ViewportTransform(m Matrix4) {
	for _, r := range d.rect {
		r.ViewportTransform(m)
	}
}

func (d *Dice) CalcCulling() {
	for _, r := range d.rect {
		r.CalcCulling()
	}
}

func (d *Dice) Draw(screen *ebiten.Image, light Vector3) {
	// ポリゴン接空間でのライトのベクトル算出
	lightv := light.MulMat(d.Rotate.Inverted()).Normalize()
	for _, r := range d.rect {
		r.CalcLightVector(lightv)
	}

	// 頂点情報作成
	vertices := []ebiten.Vertex{}
	indices := []uint16{}
	for _, r := range d.rect {
		vertices, indices = r.Draw(screen, vertices, indices)
	}

	// 描画
	op := &ebiten.DrawTrianglesShaderOptions{}
	op.AntiAlias = true
	op.FillRule = ebiten.FillRuleFillAll
	op.Images[0] = d.diceImage
	op.Images[1] = d.heightMap
	screen.DrawTrianglesShader(vertices, indices, d.shader, op)

	// // デバッグ用描画
	// for _, r := range d.rect {
	// 	r.DrawNormal(screen)
	// }
}
