package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/lxn/win"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var ball *ebiten.Image
var paddle *ebiten.Image
var scr_width int = int(win.GetSystemMetrics(win.SM_CXSCREEN))
var scr_height int = int(win.GetSystemMetrics(win.SM_CYSCREEN))

var step = 5
var paddle_padding = 50
var upkey = ebiten.KeyW
var downkey = ebiten.KeyS

var FontFace font.Face

func randint(min, max int) int {
	return min + rand.Intn(max-min)
}

func init() {
	var err error
	ball, _, err = ebitenutil.NewImageFromFile("ball.png")
	if err != nil {
		log.Fatal(err)
	}

	paddle, _, err = ebitenutil.NewImageFromFile("paddle.png")
	if err != nil {
		log.Fatal(err)
	}

	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	//**convert** to font.Face
	FontFace, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    48,
		DPI:     72,
		Hinting: font.HintingFull,
	})

}

func isBallTouchingWall(bx, by, bw, bh int, g *Game) int {
	/* 0 is left, 1 is right, 2 is top, 3 is bottom, -1 is not touching, 4 is touching left but die */
	if bx <= 0 {
		if g.elapsedFrames > 3 {
			g.rightpoints++
		}
		return 0
	} else if bx <= paddle_padding && (by <= g.leftPaddleY+paddle.Bounds().Dy() && by >= g.leftPaddleY-paddle.Bounds().Dy()) {
		return 0
	} else if bx >= scr_width-bw {
		g.leftpoints++
		return 1
	} else if bx >= scr_width-bw-paddle_padding && (by <= g.rightPaddleY+paddle.Bounds().Dy() && by >= g.rightPaddleY-paddle.Bounds().Dy()) {
		return 1
	} else if by <= 0 {
		return 2
	} else if by >= scr_height-bh {
		return 3
	}

	return -1
}

type Game struct {
	ballPosX        int
	ballPosY        int
	ballIncreasingX bool
	ballIncreasingY bool
	rightPaddleX    int
	rightPaddleY    int
	leftPaddleX     int
	leftPaddleY     int
	leftpoints      int
	rightpoints     int
	elapsedFrames   int
}

func continueBallMomentum(g *Game) {
	if g.ballIncreasingX {
		g.ballPosX += step
	} else {
		g.ballPosX -= step
	}

	if g.ballIncreasingY {
		g.ballPosY += step
	} else {
		g.ballPosY -= step
	}
}

func (g *Game) Update() error {
	rand.Seed(time.Now().UnixNano())

	continueBallMomentum(g)

	if g.elapsedFrames < 3 {
		g.ballPosX = scr_width / 2
		g.ballPosY = scr_height / 2
	}

	walltouch := isBallTouchingWall(g.ballPosX, g.ballPosY, ball.Bounds().Dx(), ball.Bounds().Dy(), g)

	if walltouch == 0 {
		g.ballIncreasingX = true
	} else if walltouch == 1 {
		g.ballIncreasingX = false
	} else if walltouch == 2 {
		g.ballIncreasingY = true
	} else if walltouch == 3 {
		g.ballIncreasingY = false
	}

	g.rightPaddleY = g.ballPosY
	g.rightPaddleX = scr_width - paddle_padding

	if ebiten.IsKeyPressed(upkey) {
		g.leftPaddleY -= step
	} else if ebiten.IsKeyPressed(downkey) {
		g.leftPaddleY += step
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.elapsedFrames++
	if g.elapsedFrames < 3 {
		g.leftPaddleX = paddle_padding - paddle.Bounds().Dx()
		g.leftPaddleY = scr_height / 2
	}
	// screen.DrawImage(ball, nil)
	//fill transparent
	// screen.Fill(color.RGBA{1, 1, 1, 128})
	screen.Fill(color.NRGBA{1, 1, 1, 2})

	ballop := &ebiten.DrawImageOptions{}
	ballop.GeoM.Translate(float64(g.ballPosX), float64(g.ballPosY))
	screen.DrawImage(ball, ballop)

	txttodraw := fmt.Sprintf("Player Points: %d\nCPU Points: %d", g.leftpoints, g.rightpoints)
	txtwidth := text.BoundString(FontFace, txttodraw).Dx()
	text.Draw(screen, txttodraw, FontFace, scr_width/2-txtwidth/2, 50, color.White)

	paddleop := &ebiten.DrawImageOptions{}
	paddleop.GeoM.Translate(float64(g.leftPaddleX), float64(g.leftPaddleY))
	screen.DrawImage(paddle, paddleop)

	paddleop2 := &ebiten.DrawImageOptions{}
	paddleop2.GeoM.Translate(float64(g.rightPaddleX), float64(g.rightPaddleY))
	screen.DrawImage(paddle, paddleop2)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\nFPS: %0.2f", ebiten.ActualFPS()))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\nBall X: %d", g.ballPosX))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\n\nBall Y: %d", g.ballPosY))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\n\n\nRight Paddle X: %d", g.rightPaddleX))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\n\n\n\nRight Paddle Y: %d", g.rightPaddleY))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\n\n\n\n\nLeft Paddle X: %d", g.leftPaddleX))
	ebitenutil.DebugPrint(screen, fmt.Sprintf("\n\n\n\n\n\n\nLeft Paddle Y: %d", g.leftPaddleY))

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return scr_width, scr_height
}

func main() {
	ebiten.SetWindowSize(scr_width, scr_height)
	ebiten.SetWindowTitle("overlaypong")
	//transparency
	ebiten.SetWindowResizable(false)
	ebiten.SetWindowDecorated(false)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowPosition(0, 0)

	if err := ebiten.RunGameWithOptions(&Game{}, &ebiten.RunGameOptions{ScreenTransparent: true}); err != nil {
		log.Fatal(err)
	}
}
