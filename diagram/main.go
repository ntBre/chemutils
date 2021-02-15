// diagram uses imagemagick to produce molecular diagrams
package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
	"os/exec"
)

// Globals
var (
	Viewer = "sxiv"
)

// Colors
var (
	Black = color.NRGBA{0, 0, 0, 255}
)

// Display encodes img to a temporary file and displays it using Viewer
func Display(img image.Image) error {
	tmp, err := ioutil.TempFile("", "img*.png")
	defer func() {
		tmp.Close()
		os.Remove(tmp.Name())
	}()
	if err != nil {
		return err
	}
	err = png.Encode(tmp, img)
	if err != nil {
		return err
	}
	cmd := exec.Command(Viewer, tmp.Name())
	return cmd.Run()
}

// NRGBA returns the image.NRGBA of img
func NRGBA(img image.Image) image.NRGBA {
	rect := img.Bounds()
	ret := image.NewNRGBA(rect)
	height, width := rect.Max.Y, rect.Max.X
	for h := 0; h <= height; h++ {
		for w := 0; w <= width; w++ {
			ret.Set(w, h, img.At(w, h))
		}
	}
	return *ret
}

// DrawGrid draws h horizontal and v vertical grid lines on img and
// returns the updated image
func DrawGrid(img image.Image, h, v int) image.Image {
	pic := NRGBA(img)
	rect := pic.Bounds()
	height, width := rect.Max.Y, rect.Max.X
	var (
		hsize, wsize int
		label        image.Image
	)
	if h > 0 {
		hsize = height / h
	} else {
		height = 0
	}
	if v > 0 {
		wsize = width / v
	} else {
		width = 0
	}
	for h := hsize; h < height; h += hsize {
		label = Label(fmt.Sprintf("%d", h), 36)
		lrect := label.Bounds()
		lw, lh := lrect.Max.X, lrect.Max.Y
		draw.Draw(&pic, image.Rect(0, h, lw, h+lh), label,
			image.Point{0, 0}, draw.Over)
		for w := 0; w <= width; w++ {
			pic.Set(w, h, Black)
		}
	}
	for w := wsize; w < width; w += wsize {
		label = Label(fmt.Sprintf("%d", w), 36)
		lrect := label.Bounds()
		lw, lh := lrect.Max.X, lrect.Max.Y
		draw.Draw(&pic, image.Rect(w, 0, w+lw, lh), label,
			image.Point{0, 0}, draw.Over)
		for h := 0; h <= height; h++ {
			pic.Set(w, h, Black)
		}
	}
	return &pic
}

// Label uses imagemagick with pango to generate a transparent PNG of
// text with size in points
func Label(text string, size int) image.Image {
	cmd := exec.Command("convert", "-background", "transparent",
		fmt.Sprintf("pango:<span face=\"sans\" size=\"%d\">%s</span>",
			1024*size, text), "png:-")
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Run()
	pic, err := png.Decode(&buf)
	if err != nil {
		panic(err)
	}
	return pic
}

// want to have a -grid flag to show grid
// - use Go image stuff for that

// DONE use pango to label the grid lines
// TODO make size of labels scale with size of image

// then parse input file and place everything where it says

// I guess it should have a command line interface with flags too, not
// just infile

// should take Viewer from config file I guess or use xdg-open

// Label should allow you to select a font, see
// https://developer.gnome.org/pygtk/stable/pango-markup-language.html
// for information

func main() {
	infile, _ := os.Open("tests/c2h4.png")
	img, err := png.Decode(infile)
	if err != nil {
		panic(err)
	}
	// label := Label("C<sub>3</sub>H<sub>2</sub>", 48)
	img = DrawGrid(img, 4, 8)
	err = Display(img)
	if err != nil {
		panic(err)
	}
	// Display(label)
}
