// diagram uses imagemagick to add captions to images
package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/template"
)

// Colors
var (
	Black    = color.NRGBA{0, 0, 0, 255}
	ARGS     []string
	lastTemp string
)

const (
	help = `diagram <captions> <image.png>
where captions is a caption file with lines of the format
    text size xpos,ypos 
text is a string for the caption (without spaces), size is an integer
font size in points, and xpos and ypos are coordinates for the caption
in pixels.
Flags:`
)

// Flags
var (
	grid = flag.String("grid", "",
		"h,v: draw a grid of h horizontal and v "+
			"vertical lines on the image")
	outfile = flag.String("o", "",
		"save the resulting image to file")
	web = flag.Bool("web", false,
		"run the program interactively in the browser")
	debug = flag.Bool("debug", false, "toggle debug printing")
)

// Display encodes img to a temporary file and displays it using the
// system default image viewer
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
	cmd := exec.Command("xdg-open", tmp.Name())
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
		hsize = height - 1
	}
	if v > 0 {
		wsize = width / v
	} else {
		wsize = width - 1
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

func parseCaption(fields []string) (Caption, error) {
	size, err := strconv.Atoi(fields[1])
	if err != nil {
		log.Printf("error parsing caption size %q as an integer, skipping", fields[1])
		return Caption{}, errors.New("malformed caption")
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
	return Caption{
		Text:     fields[0],
		Size:     size,
		Position: image.Point{ptx, pty},
	}, nil
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
			cap, err := parseCaption(fields)
			if err == nil {
				ret = append(ret, cap)
			}
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

type Index struct {
	Img string
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if *debug {
		fmt.Printf("indexHandler url: %q\n", r.URL)
	}
	index := template.Must(template.ParseFiles("index.template"))
	err := index.ExecuteTemplate(w, "index", &Index{Img: ARGS[0]})
	if err != nil {
		panic(err)
	}
}

// TODO I'd like a crop function as well, should be very easy to
// implement in the Go part. have to make sure it puts the captions
// first though or it'll throw off their locations

// TODO load caption file so you can resume working in here

// TODO list of captions with del/edit buttons next to each

// probably change this name to requestHandler or similar and take all
// kinds of these requests, then use /req as the url for it

// how am I going to handle captions AND grids? if captionHandler is
// separate, I probably have to redo the grid each time too. how to
// maintain state? do i need a global img to use everywhere?
func reqHandler(w http.ResponseWriter, r *http.Request) {
	reqs := r.URL.Query()
	grid := reqs["grid"][0]
	caps := reqs["cap"]
	if *debug {
		fmt.Printf("gridHandler url: %q\n", r.URL)
		fmt.Printf("gridHandler GET grid: %q\n", grid)
		fmt.Printf("gridHandler GET  cap: %q\n", caps)
	}
	img := loadPic(ARGS[0])
	f, err := os.CreateTemp("", "diagram*.png")
	if err != nil {
		panic(err)
	}
	out := img
	if grid != "" && grid != "," {
		h, v := ParseGrid(grid)
		out = DrawGrid(img, h, v)
	}
	for _, c := range caps {
		if c != "" {
			fields := strings.Split(c, ",")
			fmt.Printf("%q\n", fields)
			if len(fields) == 4 {
				fields[2] = strings.Join(
					[]string{fields[2], fields[3]}, ",",
				)
				fields = fields[0:3]
				cap, err := parseCaption(fields)
				if err == nil {
					drawCaption(&out, cap)
				}
			}
			fmt.Printf("%q\n", fields)
			fmt.Printf("%q\n", c)
		}
	}
	err = png.Encode(f, &out)
	if err != nil {
		panic(err)
	}
	if *debug {
		fmt.Printf("gridHandler generated file: %q\n", f.Name())
	}
	if lastTemp != "" {
		os.Remove(lastTemp)
	}
	lastTemp = f.Name()
	io.WriteString(w, lastTemp)
}

func fileHandler(filename string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filename)
	}
}

func initialize() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Usage: %s\n", help)
		flag.PrintDefaults()
	}
	log.SetFlags(0)
	log.SetPrefix("diagram: ")
	flag.Parse()
	ARGS = flag.Args()
}

func miscHandler(w http.ResponseWriter, r *http.Request) {
	if *debug {
		fmt.Printf("miscHandler requested url: %q\n", r.URL.Path)
	}
	http.ServeFile(w, r, r.URL.Path)
}

func webInterface() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/req", reqHandler)
	http.HandleFunc("/main.css", fileHandler("main.css"))
	http.HandleFunc("/main.js", fileHandler("main.js"))
	http.HandleFunc("/"+ARGS[0], fileHandler(ARGS[0]))
	http.HandleFunc("/tmp/", miscHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func loadPic(filename string) image.NRGBA {
	infile, _ := os.Open(filename)
	img, err := png.Decode(infile)
	if err != nil {
		panic(err)
	}
	return NRGBA(img)
}

func dumpPic(pic image.NRGBA, filename string) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	err = png.Encode(f, &pic)
	if err != nil {
		panic(err)
	}
}

func drawCaption(pic draw.Image, caption Caption) {
	label := Label(caption.Text, caption.Size)
	lrect := label.Bounds()
	lw, lh := lrect.Max.X, lrect.Max.Y
	draw.Draw(pic, image.Rect(
		caption.Position.X-lw/2,
		caption.Position.Y-lh/2,
		caption.Position.X+lw/2,
		caption.Position.Y+lh/2,
	), label, image.Point{0, 0}, draw.Over)
}

func main() {
	initialize()
	switch {
	case *web && len(ARGS) < 1:
		log.Fatal("missing image for -web")
	case *web:
		webInterface()
	case len(ARGS) < 2:
		log.Fatal("not enough input arguments")
	}
	capfile := ARGS[0]
	captions := ParseCaptions(capfile)
	pic := loadPic(ARGS[1])
	if *grid != "" {
		h, v := ParseGrid(*grid)
		pic = DrawGrid(pic, h, v)
	}
	for _, caption := range captions {
		drawCaption(&pic, caption)
	}
	if *outfile != "" {
		dumpPic(pic, *outfile)
	} else {
		err := Display(&pic)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
