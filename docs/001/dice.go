package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	// 頂点座標
	box = []Vector3{
		{X: -1, Y: -1, Z: -1}, // 手前側
		{X: 1, Y: -1, Z: -1},
		{X: 1, Y: 1, Z: -1},
		{X: -1, Y: 1, Z: -1},
		{X: 1, Y: -1, Z: 1}, // 反対側
		{X: -1, Y: -1, Z: 1},
		{X: -1, Y: 1, Z: 1},
		{X: 1, Y: 1, Z: 1},
		{X: -1, Y: -1, Z: 1}, // 左側
		{X: -1, Y: -1, Z: -1},
		{X: -1, Y: 1, Z: -1},
		{X: -1, Y: 1, Z: 1},
		{X: 1, Y: -1, Z: -1}, // 右側
		{X: 1, Y: -1, Z: 1},
		{X: 1, Y: 1, Z: 1},
		{X: 1, Y: 1, Z: -1},
		{X: -1, Y: -1, Z: 1}, // 上側
		{X: 1, Y: -1, Z: 1},
		{X: 1, Y: -1, Z: -1},
		{X: -1, Y: -1, Z: -1},
		{X: -1, Y: 1, Z: -1}, // 下側
		{X: 1, Y: 1, Z: -1},
		{X: 1, Y: 1, Z: 1},
		{X: -1, Y: 1, Z: 1},
	}

	// テクスチャ座標(頂点座標と対応している)
	tindices [][]float32 = [][]float32{
		{0, 0}, {200, 0}, {200, 200}, {0, 200},
		{400, 0}, {600, 0}, {600, 200}, {400, 200},
		{1000, 0}, {1200, 0}, {1200, 200}, {1000, 200},
		{800, 0}, {1000, 0}, {1000, 200}, {800, 200},
		{600, 0}, {800, 0}, {800, 200}, {600, 200},
		{200, 0}, {400, 0}, {400, 200}, {200, 200},
	}

	// 頂点インデックス(3個で三角形1個を表現、上の座標の番号)
	indices []uint16 = []uint16{
		0, 1, 2, 0, 2, 3, 4, 5, 6, 4, 6, 7, 8, 9, 10, 8, 10, 11, 12, 13, 14, 12, 14, 15, 16, 17, 18, 16, 18, 19, 20, 21, 22, 20, 22, 23,
	}

	Rotate    Matrix4 = NewMatrix4Rotate(1, 1, 0, 1)
	diceImage *ebiten.Image
	shader    *ebiten.Shader
)

func init() {
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
			color *= custom.x

			return color
		}
		`,
	)

	s, err := ebiten.NewShader(shaderText)
	if err != nil {
		panic(err)
	}
	shader = s

	// サイコロ画像
	diceImage = ebiten.NewImage(1202, 202)
	diceImage.Fill(color.White)
	vector.DrawFilledCircle(diceImage, 100, 100, 30, color.RGBA{0xff, 0x00, 0x00, 0xff}, true)
	vector.DrawFilledCircle(diceImage, 200+50, 50, 20, color.Black, true)
	vector.DrawFilledCircle(diceImage, 200+100, 100, 20, color.Black, true)
	vector.DrawFilledCircle(diceImage, 200+150, 150, 20, color.Black, true)
	vector.DrawFilledCircle(diceImage, 400+50, 50, 20, color.Black, true)
	vector.DrawFilledCircle(diceImage, 400+50, 100, 20, color.Black, true)
	vector.DrawFilledCircle(diceImage, 400+50, 150, 20, color.Black, true)
	vector.DrawFilledCircle(diceImage, 400+150, 50, 20, color.Black, true)
	vector.DrawFilledCircle(diceImage, 400+150, 100, 20, color.Black, true)
	vector.DrawFilledCircle(diceImage, 400+150, 150, 20, color.Black, true)
	vector.DrawFilledCircle(diceImage, 600+50, 50, 20, color.Black, true)
	vector.DrawFilledCircle(diceImage, 600+150, 50, 20, color.Black, true)
	vector.DrawFilledCircle(diceImage, 600+50, 150, 20, color.Black, true)
	vector.DrawFilledCircle(diceImage, 600+150, 150, 20, color.Black, true)
	vector.DrawFilledCircle(diceImage, 800+100, 100, 20, color.Black, true)
	vector.DrawFilledCircle(diceImage, 800+50, 50, 20, color.Black, true)
	vector.DrawFilledCircle(diceImage, 800+150, 50, 20, color.Black, true)
	vector.DrawFilledCircle(diceImage, 800+50, 150, 20, color.Black, true)
	vector.DrawFilledCircle(diceImage, 800+150, 150, 20, color.Black, true)
	vector.DrawFilledCircle(diceImage, 1000+50, 50, 20, color.Black, true)
	vector.DrawFilledCircle(diceImage, 1000+150, 150, 20, color.Black, true)

	// パースペクティブコレクト確認用ライン
	// for x := 0; x < 1200; x += 50 {
	// 	vector.DrawFilledRect(diceImage, float32(x), 0, 10, 200, color.Black, true)
	// }
	// for y := 0; y < 200; y += 50 {
	// 	vector.DrawFilledRect(diceImage, 0, float32(y), 1200, 10, color.Black, true)
	// }

}

func normalize(x, y float32) (float32, float32) {
	len := float32(math.Sqrt(float64(x*x + y*y)))
	if len != 0 {
		x /= len
		y /= len
	}
	return x, y
}

func DrawDice(screen *ebiten.Image) {
	// ワールド変換
	m := NewMatrix4Scale(200, 200, 200)
	m = m.Mul(Rotate)
	m = m.Mul(NewMatrix4Translate(0, 0, 600))

	// ビュー変換は省略
	// (サイコロの座標をあらかじめカメラの前に設定しているので)

	// 射影変換
	m = m.Mul(NewProjection(200, 2000, 0, 640, 0, 480))

	// 変換実行
	b := []Vector4{}
	for _, v := range box {
		b = append(b, v.MulMatToVector4(m))
	}

	// このタイミングで視錐台カリングする(今は面倒なのでやらない)

	// 法線を頂点に設定するために法線スライスを作成する
	var normal []float32 = make([]float32, len(b))
	for i := 0; i < len(indices); i += 3 {
		v1 := b[indices[i+1]].ToVec3().Sub(b[indices[i]].ToVec3())
		v2 := b[indices[i+2]].ToVec3().Sub(b[indices[i]].ToVec3())
		n := v1.Cross(v2).Normalize() // ポリゴンの法線

		// と思ったが、法線の表現には項目が3つ(xyz)必要になり、Ebitenではシェーダに渡せる項目が少ないので、
		// 適当な光源を元に暗さを計算して渡すようにした
		s := n.Dot(NewVecor3(1, 1, 1).Normalize())
		normal[indices[i]] = s*0.5 + 0.5
		normal[indices[i+1]] = s*0.5 + 0.5
		normal[indices[i+2]] = s*0.5 + 0.5
	}

	// ビューポート変換
	m = NewViewport(640, 480)

	// 変換実行
	for i, v := range b {
		b[i] = v.MulMat(m)
	}

	// 座標をwで割る
	for i, v := range b {
		v.X /= v.W
		v.Y /= v.W
		v.Z /= v.W
		b[i] = v
	}

	// 背面カリングする
	dindices := []uint16{}
	for i := 0; i < len(indices); i += 3 {
		x1 := b[indices[i+1]].X - b[indices[i]].X
		y1 := b[indices[i+1]].Y - b[indices[i]].Y
		x2 := b[indices[i+2]].X - b[indices[i+1]].X
		y2 := b[indices[i+2]].Y - b[indices[i+1]].Y
		x, y := normalize(x2, y2)

		if x1*y-x*y1 > 0 {
			// 表を向いたポリゴンの頂点だけをインデックスバッファに入れる
			dindices = append(dindices, indices[i], indices[i+1], indices[i+2])
		}
	}

	// vector用頂点バッファ作成
	var vertices []ebiten.Vertex = []ebiten.Vertex{}
	for i, v := range b {
		var invW float32
		if PerspectiveCorrection {
			invW = 1 / v.W
		} else {
			invW = 1
		}
		vertex := ebiten.Vertex{}
		vertex.DstX = v.X
		vertex.DstY = v.Y
		vertex.SrcX = tindices[i][0] * invW // パースペクティブコレクト用にWで割る
		vertex.SrcY = tindices[i][1] * invW // パースペクティブコレクト用にWで割る
		vertex.Custom0 = normal[i]          // 暗さ
		vertex.Custom3 = invW               // パースペクティブコレクト用
		vertices = append(vertices, vertex)
	}
	op := &ebiten.DrawTrianglesShaderOptions{}
	op.AntiAlias = true
	op.FillRule = ebiten.FillRuleFillAll
	op.Images[0] = diceImage
	screen.DrawTrianglesShader(vertices, dindices, shader, op)
}
