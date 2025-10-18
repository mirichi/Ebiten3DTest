package main

import (
	_ "image/png"
	"math"
	"myproject/control"
	"myproject/ui"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	touch                 ui.TouchInfo
	button                *control.Button
	label                 *control.Label
	PerspectiveCorrection bool = true
	dice                  *Dice
	shere                 *Sphere
	camera                *Camera
	light                 Vector3 = Vector3{1, 1, 0}.Normalize()
)

type Game struct {
}

func (g *Game) Update() error {
	ui.Input_Update()

	// スライドした方向に回転する
	if touch != nil {
		if !touch.IsJustReleased() {
			oldx, oldy := touch.OldPos()
			posx, posy := touch.Pos()
			if oldx != posx || oldy != posy {
				dx := posx - oldx
				dy := posy - oldy
				len := float32(math.Sqrt(float64(dx*dx + dy*dy)))
				dice.Rotate = dice.Rotate.Mul(NewMatrix4Rotate(float32(posy-oldy), float32(oldx-posx), 0, len*math.Pi/180))
			}
		} else {
			touch = nil
		}
	} else {
		touch = ui.FirstTouch()
	}

	// UI入力判定
	t := ui.FirstTouch()
	if t != nil && t.IsJustPressed() {
		button.ProcessTouch(t)
	}

	// ui処理
	button.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	camera.Draw(screen, shere, light)
	camera.Draw(screen, dice, light)
	button.Draw(screen)
	label.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func main() {
	// ボタン
	button = control.NewButton(20, 20, 100, 50, "ON", 28, ui.AdjustCenter, nil, func() {
		PerspectiveCorrection = !PerspectiveCorrection
		if PerspectiveCorrection {
			button.Label = "ON"
		} else {
			button.Label = "OFF"
		}
	})

	// ラベル
	label = control.NewLabel(130, 30, 0, 0, "PerspectiveCorrection", 20, ui.AdjustLeft, nil)

	// サイコロ作成
	dice = NewDice()

	// 球体作成
	shere = NewSphere()

	// カメラ作成
	camera = NewCamera(
		Vector3{0, 0, 0},
		Vector3{0, 0, 1},
		Vector3{0, -1, 0},
		NewProjectionPerspective(1000, 20000, 0, 640, 0, 480),
		NewViewport(640, 480),
	)

	ebiten.SetWindowSize(640, 480)
	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}
