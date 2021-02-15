// diagram uses imagemagick to add captions to images
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Globals
var (
	Viewer = "sxiv"
)

// Colors
var (
	Black = color.NRGBA{0, 0, 0, 255}
)

// Flags
var (
	grid = flag.String("grid", "",
		"h,v: draw a grid of h horizontal and v "+
			"vertical lines on the image")
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
func DrawGrid(img image.NRGBA, h, v int) image.NRGBA {
	rect := img.Bounds()
	height, width := rect.Max.Y, rect.Max.X
	font := int(math.Sqrt(float64(height*width))) / 100
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
		label = Label(fmt.Sprintf("%d", h), font)
		lrect := label.Bounds()
		lw, lh := lrect.Max.X, lrect.Max.Y
		draw.Draw(&img, image.Rect(0, h, lw, h+lh), label,
			image.Point{0, 0}, draw.Over)
		for w := 0; w <= width; w++ {
			img.Set(w, h, Black)
		}
	}
	for w := wsize; w < width; w += wsize {
		label = Label(fmt.Sprintf("%d", w), font)
		lrect := label.Bounds()
		lw, lh := lrect.Max.X, lrect.Max.Y
		draw.Draw(&img, image.Rect(w, 0, w+lw, lh), label,
			image.Point{0, 0}, draw.Over)
		for h := 0; h <= height; h++ {
			img.Set(w, h, Black)
		}
	}
	return img
}

// Label uses imagemagick with pango to generate a transparent PNG of
// text with size in points
func Label(text string, size int) image.Image {
	cmd := exec.Command("convert", "-background", "transparent",
		fmt.Sprintf("pango:<span face=\"sans\" "+
			"size=\"%d\">%s</span>",
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

// Caption holds the information for a caption
type Caption struct {
	Text     string
	Size     int
	Position image.Point
}

// ParseCaptions parses caption input from filename and returns a
// slice of Captions
func ParseCaptions(filename string) (ret []Caption) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 3 {
			size, err := strconv.Atoi(fields[1])
			if err != nil {
				log.Fatalf("error parsing caption size %q", fields[1])
			}
			strpt := strings.Split(fields[2], ",")
			ptx, err := strconv.Atoi(strpt[0])
			if err != nil {
				panic(err)
			}
			pty, err := strconv.Atoi(strpt[1])
			if err != nil {
				panic(err)
			}
			ret = append(ret,
				Caption{
					Text:     fields[0],
					Size:     size,
					Position: image.Point{ptx, pty},
				})
		}
	}
	return
}

// ParseGrid parses the string from the -grid flag and returns its
// components as ints
func ParseGrid(str string) (h, v int) {
	split := strings.Split(str, ",")
	if len(split) != 2 {
		log.Fatal("diagram: malformed -grid argument")
	}
	h, _ = strconv.Atoi(split[0])
	v, _ = strconv.Atoi(split[1])
	return
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("diagram: ")
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		log.Fatal("not enough input arguments")
	}
	capfile := args[0]
	captions := ParseCaptions(capfile)
	inpic := args[1]
	infile, _ := os.Open(inpic)
	img, err := png.Decode(infile)
	if err != nil {
		panic(err)
	}
	pic := NRGBA(img)
	if *grid != "" {
		h, v := ParseGrid(*grid)
		pic = DrawGrid(pic, h, v)
	}
	for _, caption := range captions {
		label := Label(caption.Text, caption.Size)
		lrect := label.Bounds()
		lw, lh := lrect.Max.X, lrect.Max.Y
		draw.Draw(&pic, image.Rect(
			caption.Position.X-lw/2,
			caption.Position.Y-lh/2,
			caption.Position.X+lw/2,
			caption.Position.Y+lh/2,
		), label, image.Point{0, 0}, draw.Over)
	}
	err = Display(&pic)
	if err != nil {
		panic(err)
	}
}
