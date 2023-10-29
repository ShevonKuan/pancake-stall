package pancake

import (
	"fmt"

	"github.com/zztroot/color"
)

func (c *Canvas) StdOut() {
	xLength, yLength := len(*c), len((*c)[0])
	for i := 0; i < xLength; i++ {
		for j := 0; j < yLength; j++ {
			if (*c)[i][j] != '\x00' {
				p := string((*c)[i][j])
				switch (*colorsContainer[c])[i][j] {
				case 1:
					fmt.Print(color.Coat(p, color.Green))
				case 2:
					fmt.Print(color.Coat(p, color.Blue))
				case 3:
					fmt.Print(color.Coat(p, color.Yellow))
				}
				fmt.Print()
			} else {
				fmt.Print("　")
			}

		}
		fmt.Print("|\n")
	}
}

func (c *Canvas) GenerateLatexCode() string {
	latexCode := "\\documentclass{standalone}\n"
	latexCode += "\\usepackage{tikz,ctex}\n"
	latexCode += "\\begin{document}\n"
	latexCode += "\\begin{tikzpicture}\n"

	height := len(*c)
	width := len((*c)[0])
	latexCode += "\\definecolor{color1}{HTML}{edae49}\n"
	latexCode += "\\definecolor{color2}{HTML}{d1495b}\n"
	latexCode += "\\definecolor{color3}{HTML}{00798c}\n"
	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			if (*c)[row][col] != '\x00' {
				x := float32(col) / 2
				y := float32(height-row-1) / 2
				latexCode += fmt.Sprintf("\\node[text=color%s] at (%f, %f) {%c};\n", fmt.Sprint((*colorsContainer[c])[row][col]), x, y, (*c)[row][col])
			}
		}
	}

	latexCode += "\\end{tikzpicture}\n"
	latexCode += "\\end{document}"

	return latexCode
}
