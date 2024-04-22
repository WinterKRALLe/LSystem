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

func DrawLSystem(name, axiom string, rules map[rune]string, iterations int, angle float64, startPosX, startPosY float64) {
	const (
		width  = 800
		height = 800
	)

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.White)
		}
	}

	angle = angle * math.Pi / 180

	turtle := Turtle{x: startPosX, y: startPosY, heading: 0.0}

	newX := turtle.x + 10*math.Cos(turtle.heading)
	newY := turtle.y + 10*math.Sin(turtle.heading)
	drawLine(img, int(turtle.x), int(turtle.y), int(newX), int(newY), color.RGBA{B: 255, A: 255})
	turtle.x = newX
	turtle.y = newY

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
		case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
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

	file, err := os.Create(name + ".png")
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
	fmt.Println("L-systém byl vykreslen do souboru " + name + ".png")
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
	lsystems := []struct {
		name       string
		axiom      string
		rules      map[rune]string
		iterations int
		angle      float64
		startPosX  float64
		startPosY  float64
	}{
		{"Koch curve", "F", map[rune]string{'F': "F+F-F-F+F"}, 3, 90.0, 200.0, 400.0},
		{"Dragon curve", "F", map[rune]string{'F': "F+G", 'G': "F-G"}, 10, 90.0, 300.0, 500.0},
		{"Sierpinski triangle", "F-G-G", map[rune]string{'F': "F-G+F+G-F", 'G': "GG"}, 4, 120.0, 200.0, 400.0},
		{"Snowflake", "F++F++F++F++F", map[rune]string{'F': "F++F++F+++++F-F++F"}, 4, 36.0, 150.0, 750.0},
	}

	for _, lsystem := range lsystems {
		DrawLSystem(lsystem.name, lsystem.axiom, lsystem.rules, lsystem.iterations, lsystem.angle, lsystem.startPosX, lsystem.startPosY)
	}
}
