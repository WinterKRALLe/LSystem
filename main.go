package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

type Turtle struct {
	x, y, heading float64
}

func DrawLSystem(axiom string, rules map[rune]string, iterations int, angle float64) {
	const (
		width  = 800
		height = 600
	)

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.White)
		}
	}
	angle = angle * math.Pi / 180

	turtle := Turtle{x: 400.0, y: 300.0, heading: 0.0}

	// Zde umístěte kreslení, které má být provedeno ihned
	newX := turtle.x + 10*math.Cos(turtle.heading)
	newY := turtle.y + 10*math.Sin(turtle.heading)
	drawLine(img, int(turtle.x), int(turtle.y), int(newX), int(newY), color.RGBA{B: 255, A: 255})
	turtle.x = newX
	turtle.y = newY
	// Konec umístění kreslení

	var stack []Turtle

	for i := 0; i < iterations; i++ {
		nextAxiom := ""
		for _, char := range axiom {
			if rule, ok := rules[char]; ok {
				nextAxiom += rule
			} else {
				nextAxiom += string(char)
			}
		}
		axiom = nextAxiom
	}

	for i, char := range axiom {
		if i == 0 {
			continue
		}
		switch char {
		case 'F':
			newX := turtle.x + 10*math.Cos(turtle.heading)
			newY := turtle.y + 10*math.Sin(turtle.heading)
			drawLine(img, int(turtle.x), int(turtle.y), int(newX), int(newY), color.RGBA{B: 255, A: 255})
			turtle.x = newX
			turtle.y = newY
		case '+':
			turtle.heading -= angle
		case '-':
			turtle.heading += angle
		case '|':
			turtle.heading += 180
		case '[':
			stack = append(stack, turtle)
		case ']':
			if len(stack) > 0 {
				turtle = stack[len(stack)-1]
				stack = stack[:len(stack)-1]
			}
		}
	}

	file, err := os.Create("lsystem.png")
	if err != nil {
		fmt.Println("Chyba při vytváření souboru:", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	if err := png.Encode(file, img); err != nil {
		fmt.Println("Chyba při ukládání obrázku:", err)
	}
	fmt.Println("L-systém byl vykreslen do souboru lsystem.png")
}

func drawLine(img *image.RGBA, x1, y1, x2, y2 int, col color.Color) {
	dx := x2 - x1
	dy := y2 - y1

	if dx == 0 && dy == 0 {
		img.Set(x1, y1, col)
		return
	}

	steps := math.Max(math.Abs(float64(dx)), math.Abs(float64(dy)))

	xIncrement := float64(dx) / steps
	yIncrement := float64(dy) / steps

	x, y := float64(x1), float64(y1)
	for i := 0; i <= int(steps); i++ {
		img.Set(int(x+0.5), int(y+0.5), col)
		x += xIncrement
		y += yIncrement
	}
}

func main() {
	axiom := "F"
	rules := map[rune]string{
		'F': "F+F-F-F+F",
	}
	iterations := 3
	angle := 90.0

	DrawLSystem(axiom, rules, iterations, angle)
}
