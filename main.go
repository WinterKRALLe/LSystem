package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
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

	angle = angle * math.Pi / 180

	turtle := Turtle{x: startPosX, y: startPosY, heading: 0.0}

	var frames []*image.Paletted

	for iter := 0; iter <= iterations; iter++ {
		img := image.NewPaletted(image.Rect(0, 0, width, height), color.Palette{color.White})

		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				img.Set(x, y, color.White)
			}
		}

		turtleCopy := turtle

		for _, char := range axiom {
			switch char {
			case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				newX := turtleCopy.x + 10*math.Cos(turtleCopy.heading)
				newY := turtleCopy.y + 10*math.Sin(turtleCopy.heading)
				drawLine(img, int(turtleCopy.x), int(turtleCopy.y), int(newX), int(newY))
				turtleCopy.x = newX
				turtleCopy.y = newY
			case '+':
				turtleCopy.heading -= angle
			case '-':
				turtleCopy.heading += angle
			case '|':
				turtleCopy.heading += 180
			case '[':
				// Not supported in this example
			case ']':
				// Not supported in this example
			}
		}

		frames = append(frames, img)

		axiom = applyRules(axiom, rules)
	}

	file, err := os.Create(name + ".gif")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	delays := make([]int, len(frames))
	for i := range delays {
		delays[i] = 100
	}
	gif.EncodeAll(file, &gif.GIF{
		Image: frames,
		Delay: delays,
	})
	fmt.Println("L-system animation was saved to", name+".gif")
}

func applyRules(axiom string, rules map[rune]string) string {
	result := ""
	for _, char := range axiom {
		if rule, ok := rules[char]; ok {
			result += rule
		} else {
			result += string(char)
		}
	}
	return result
}

func drawLine(img *image.Paletted, x1, y1, x2, y2 int) {
	dx := x2 - x1
	dy := y2 - y1

	if dx == 0 && dy == 0 {
		img.SetColorIndex(x1, y1, 1)
		return
	}

	steps := math.Max(math.Abs(float64(dx)), math.Abs(float64(dy)))

	xIncrement := float64(dx) / steps
	yIncrement := float64(dy) / steps

	x, y := float64(x1), float64(y1)
	for i := 0; i <= int(steps); i++ {
		img.SetColorIndex(int(x+0.5), int(y+0.5), 1)
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
		{"Kdo ví co", "F", map[rune]string{'F': "F++F+F--"}, 5, 20.0, 200.0, 200.0},
	}

	for _, lsystem := range lsystems {
		DrawLSystem(lsystem.name, lsystem.axiom, lsystem.rules, lsystem.iterations, lsystem.angle, lsystem.startPosX, lsystem.startPosY)
	}
}
