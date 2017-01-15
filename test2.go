package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/golang/freetype"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func getRandom(max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(max)
}

var (
	dpi      = flag.Float64("dpi", 72, "screen resolution in Dots Per Inch")
	fontfile = flag.String("fontfile", "tt1.TTF", "filename of the ttf font")
	hinting  = flag.String("hinting", "none", "none | full")
	size     = flag.Float64("size", 16, "font size in points")
	spacing  = flag.Float64("spacing", 1.5, "line spacing (e.g. 2 means double spaced)")
	wonb     = flag.Bool("whiteonblack", false, "white text on a black background")
	jindu    = flag.String("jindu", "117.350", "usage")
	weidu    = flag.String("weidu", "30.294", "usage")
	times    = flag.String("time", "2017-01-13 12:", "usage")
	addr     = flag.String("addr", "站点:AHCZ34/池州高坦吴角", "usage")
	IMSI     = flag.String("IMSI", "IMSI:460078157306249", "usage")
	srcPath  = flag.String("src", "data/rawpic/", "usage")
	dstPath  = flag.String("dst", "data/repic/", "usage")
)

func main() {
	flag.Parse()
	path := *srcPath
	dir := ""
	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			dir = path
			return nil
		}
		println(dir + "/")
		println(path)
		genpic(path, strings.Replace(path, dir+"/", *dstPath, 1))
		return nil
	})
	// Read the font data.

}

func genpic(picPath string, dstPath string) {
	var text = []string{
		fmt.Sprintf("经度:%s%d%d%d", *jindu, getRandom(10), getRandom(10), getRandom(10)),
		fmt.Sprintf("维度:%s%d%d%d", *weidu, getRandom(10), getRandom(10), getRandom(10)),
		fmt.Sprintf("方位角:%d", getRandom(360)),
		fmt.Sprintf("%s%d%d", *times, getRandom(6), getRandom(10)),
		*addr,
		*IMSI,
	}
	fontBytes, err := ioutil.ReadFile(*fontfile)
	if err != nil {
		log.Println(err)
		return
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}

	// Initialize the context.
	imgb, _ := os.Open(picPath)
	imgddd, _ := jpeg.Decode(imgb)
	defer imgb.Close()

	fg, _ := image.NewUniform(color.RGBA{uint8(206), uint8(139), uint8(16), 255}), image.White
	if *wonb {
		fg, _ = image.White, image.NewUniform(color.RGBA{uint8(0), uint8(0), uint8(0), 255})
	}
	// fmt.Println(imgddd.Bounds().Dx())
	rgba := image.NewRGBA(image.Rect(0, 0, imgddd.Bounds().Dx(), imgddd.Bounds().Dy()))
	draw.Draw(rgba, rgba.Bounds(), imgddd, image.ZP, draw.Src)
	c := freetype.NewContext()

	c.SetDPI(*dpi)
	c.SetFont(f)
	c.SetFontSize(*size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)

	switch *hinting {
	default:
		c.SetHinting(font.HintingNone)
	case "full":
		c.SetHinting(font.HintingFull)
	}

	// Draw the text.
	// pt := freetype.Pt(10, 10+int(c.PointToFixed(*size)>>6))
	// fmt.Println(int(c.PointToFixed(*size)) >> 6)
	// fmt.Println(len(text[0]))

	pt := freetype.Pt(rgba.Bounds().Dx()-len(text[0])*16/2, 10+int(c.PointToFixed(*size)>>6))
	for _, s := range text {
		// de := 0
		// if strings.Contains(s, "经度") {
		//  de = 1
		// } else if strings.Contains(s, "维度") {
		//  de = 1
		// } else if strings.Contains(s, "方位角") {
		//  de = 2
		// }
		fmt.Println(strings.Count(s, "")-1, len(s))
		de := 0
		if strings.Count(s, "")-1 != len(s) {
			de = len(s) - (len(s)-strings.Count(s, "")+1)/2
		} else {
			de = len(s)
		}
		pt.X = c.PointToFixed(float64(rgba.Bounds().Dx() - de*16/2))
		_, err = c.DrawString(s, pt)
		if err != nil {
			log.Println(err)
			return
		}
		pt.Y += c.PointToFixed(*size * *spacing)
	}

	// Save that RGBA image to disk.
	outFile, err := os.Create(dstPath)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	err = png.Encode(b, rgba)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = b.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Println("Wrote out.png OK.")

}
