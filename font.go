package eui

import (
	"bytes"
	"fmt"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Font struct {
	name   string
	source *text.GoTextFaceSource
	font   map[int]*text.GoTextFace
}

func NewFont(name string, data []byte, size int) (*Font, error) {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	fnt := &text.GoTextFace{Source: s, Size: float64(size)}
	fonts := make(map[int]*text.GoTextFace, 0)
	fonts[size] = fnt
	f := &Font{
		name:   name,
		source: s,
		font:   fonts,
	}
	return f, nil
}

func (f Font) add(size int) Font {
	font := &text.GoTextFace{Source: f.source, Size: float64(size)}
	f.font[size] = font
	return f
}

func (f Font) Get(size int) *text.GoTextFace {
	if font, exists := f.font[size]; exists {
		return font
	}
	f.add(size)
	return f.font[size]
}

func (f Font) CalcFontSize(txt string, rect Rect[int]) int {
	if rect.IsEmpty() {
		panic(fmt.Sprintln("Empty rect", txt, rect))
	}
	percent := 0.9
	w0, h0 := rect.Size()
	sz := min(w0, h0)
	var fontSize float64
	for {
		fontSize = percent * float64(sz)
		fnt := f.Get(int(fontSize))
		w, h := text.Measure(txt, fnt, fnt.Size*1.2)
		if w0 > int(w) && h0 > int(h) {
			break
		}
		percent -= 0.01
	}
	return int(fontSize)
}

func (f Font) WordWrapText(txt string, fontSize float64, width int) (string, Point[int]) {
	if len(txt) == 0 {
		return txt, NewPoint(0, 0)
	}
	var (
		fnt    *text.GoTextFace
		result strings.Builder
		maxW   float64
		lines  int
	)
	fnt = f.Get(int(fontSize))
	origLines := strings.Split(txt, "\n")
	for li, origLine := range origLines {
		words := strings.Fields(origLine)
		line := ""
		for i, str := range words {
			testLine := line
			if testLine != "" {
				testLine += " "
			}
			testLine += str
			w, _ := text.Measure(testLine, fnt, fnt.Size*1.2)
			if w > float64(width) && line != "" {
				result.WriteString(line)
				result.WriteString("\n")
				lines++
				line = str
			} else {
				line = testLine
			}
			if w > maxW {
				maxW = w
			}
			// Последнее слово в строке
			if i == len(words)-1 {
				result.WriteString(line)
				lines++
			}
		}
		// Если строка была пустой (например, двойной \n)
		if len(words) == 0 {
			result.WriteString("\n")
			lines++
		}
		// Не добавлять лишний перенос после последней строки
		if li < len(origLines)-1 && len(words) > 0 {
			result.WriteString("\n")
		}
	}
	return result.String(), NewPoint(int(maxW), int(fnt.Size*1.2*float64(lines)))
}

func (f Font) DrawString(surface *ebiten.Image, txt string, fontSize int, rect Rect[int], hAlign text.Align, vAlign text.Align, fg color.Color, wordWrap bool) {
	var (
		x, y, w, h float64
	)
	x, y = float64(rect.X), float64(rect.Y)
	w, h = float64(rect.Width()), float64(rect.Height())
	switch hAlign {
	case text.AlignStart:
		x += 0
	case text.AlignCenter:
		x += w / 2
	case text.AlignEnd:
		x += w
	}
	switch vAlign {
	case text.AlignStart:
		y += 0
	case text.AlignCenter:
		y += h / 2
	case text.AlignEnd:
		y += h
	}
	if fontSize == 0 && !wordWrap {
		fontSize = f.CalcFontSize(txt, rect)
	}
	if wordWrap {
		txt, _ = f.WordWrapText(txt, float64(fontSize), rect.W)
	}
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(fg)
	op.LineSpacing = float64(fontSize) * 1.2
	op.PrimaryAlign = hAlign
	op.SecondaryAlign = vAlign
	text.Draw(surface, txt, f.Get(fontSize), op)
}
