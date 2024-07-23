package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 300
	screenHeight = 600
	blockSize    = 30
	boardWidth   = 10
	boardHeight  = 20
)

var (
	backgroundColor = color.RGBA{0x0, 0x0, 0x0, 0xff}
	blockColors     = []color.RGBA{
		{0xff, 0x00, 0x00, 0xff}, // Red
		{0x00, 0xff, 0x00, 0xff}, // Green
		{0x00, 0x00, 0xff, 0xff}, // Blue
		{0xff, 0xff, 0x00, 0xff}, // Yellow
		{0xff, 0x00, 0xff, 0xff}, // Magenta
		{0x00, 0xff, 0xff, 0xff}, // Cian
		{0xff, 0x7f, 0x00, 0xff}, // Naranja
	}
)

type Block struct {
	x, y int
}
type Piece struct {
	blocks [4]Block
	color  color.RGBA
}

var pieces = []Piece{
	{[4]Block{{0, 0}, {1, 0}, {2, 0}, {3, 0}}, blockColors[0]}, // I
	{[4]Block{{0, 0}, {1, 0}, {0, 1}, {1, 1}}, blockColors[1]}, // O
	{[4]Block{{0, 0}, {1, 0}, {2, 0}, {1, 1}}, blockColors[2]}, // T
	{[4]Block{{0, 0}, {1, 0}, {2, 0}, {2, 1}}, blockColors[3]}, // L
	{[4]Block{{0, 0}, {1, 0}, {2, 0}, {0, 1}}, blockColors[4]}, // J
	{[4]Block{{0, 0}, {1, 0}, {1, 1}, {2, 1}}, blockColors[5]}, // S
	{[4]Block{{1, 0}, {2, 0}, {0, 1}, {1, 1}}, blockColors[6]}, // Z
}

type Game struct {
	board        [boardHeight][boardWidth]color.RGBA
	currentPiece Piece
	currentX     int
	currentY     int
	score        int
	gameOver     bool
	lastMoveTime time.Time
	moveInternal time.Duration
}

func NewGame() *Game {
	g := &Game{
		moveInternal: time.Second,
		lastMoveTime: time.Now(),
	}
	// Inicializar el tablero con el color de fondo
	for y := 0; y < boardHeight; y++ {
		for x := 0; x < boardWidth; x++ {
			g.board[y][x] = backgroundColor
		}
	}
	g.spawnPiece()
	return g
}
func (g *Game) spawnPiece() {
	g.currentPiece = pieces[rand.Intn(len(pieces))]
	g.currentX = boardWidth/2 - 2
	g.currentY = -1 // Empezar justo encima del tablero visible
	if g.collision() {
		g.gameOver = true
	}
}
func (g *Game) collision() bool {
	for _, block := range g.currentPiece.blocks {
		x, y := g.currentX+block.x, g.currentY+block.y
		if x < 0 || x >= boardWidth || y >= boardHeight {
			return true
		}
		if y >= 0 && g.board[y][x] != backgroundColor {
			return true
		}
	}
	return false
}
func (g *Game) moveDown() bool {
	g.currentY++
	if g.collision() {
		g.currentY--
		g.lockPiece()
		return false
	}
	return true
}
func (g *Game) moveLeft() {
	g.currentX--
	if g.collision() {
		g.currentX++
	}
}

func (g *Game) moveRight() {
	g.currentX++
	if g.collision() {
		g.currentX--
	}
}
func (g *Game) rotate() {
	original := g.currentPiece
	for i := range g.currentPiece.blocks {
		x := g.currentPiece.blocks[i].y
		y := g.currentPiece.blocks[i].x
		g.currentPiece.blocks[i].x = x
		g.currentPiece.blocks[i].y = y
	}
	if g.collision() {
		g.currentPiece = original
	}
}
func (g *Game) lockPiece() {
	for _, block := range g.currentPiece.blocks {
		x, y := g.currentX+block.x, g.currentY+block.y
		if y >= 0 && y < boardHeight && x >= 0 && x < boardWidth {
			g.board[y][x] = g.currentPiece.color
		}
	}
	g.clearLines()
	g.spawnPiece()
}
func (g *Game) clearLines() {
	for y := boardHeight - 1; y >= 0; y-- {
		full := true
		for x := 0; x < boardWidth; x++ {
			if g.board[y][x] == backgroundColor {
				full = false
				break
			}
		}
		if full {
			g.score += 100
			for yy := y; yy > 0; yy-- {
				g.board[yy] = g.board[yy-1]
			}
			g.board[0] = [boardWidth]color.RGBA{}
			y++
		}
	}
}
func (g *Game) Update() error {
	if g.gameOver {
		return nil
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.moveLeft()
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.moveRight()
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.moveDown()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		g.rotate()
	}

	if time.Since(g.lastMoveTime) >= g.moveInternal {
		if !g.moveDown() {
			g.spawnPiece()
		}
		g.lastMoveTime = time.Now()
	}

	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(backgroundColor)

	// Draw the board
	for y := 0; y < boardHeight; y++ {
		for x := 0; x < boardWidth; x++ {
			if g.board[y][x] != backgroundColor {
				drawRect(screen, float64(x*blockSize), float64(y*blockSize), blockSize, blockSize, g.board[y][x])
			}
		}
	}

	// Draw the current piece
	for _, block := range g.currentPiece.blocks {
		x, y := g.currentX+block.x, g.currentY+block.y
		if y >= 0 && y < boardHeight && x >= 0 && x < boardWidth {
			drawRect(screen, float64(x*blockSize), float64(y*blockSize), blockSize, blockSize, g.currentPiece.color)
		}
	}

	// Draw score
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Score: %d", g.score))

	if g.gameOver {
		ebitenutil.DebugPrint(screen, "Game Over!")
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	rand.Seed(time.Now().UnixNano())
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Tetris en Go")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
func drawRect(screen *ebiten.Image, x, y, width, height float64, clr color.Color) {
	ebitenutil.DrawRect(screen, x, y, width, height, clr)
}
