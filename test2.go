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
	// "bufio"
	// "strings"
	"time"
	. "github.com/hunterhug/go_image"
	// "math"
)

func getRandom(max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	time.Sleep(1*time.Millisecond)
	return r.Intn(max)
}

var (
	dpi      = flag.Float64("dpi", 72, "screen resolution in Dots Per Inch")
	fontfile = flag.String("fontfile", "tt1.TTF", "filename of the ttf font")
	hinting  = flag.String("hinting", "none", "none | full")
	size     = flag.Float64("size", 20, "font size in points")
	spacing  = flag.Float64("spacing", 1.5, "line spacing (e.g. 2 means double spaced)")
	wonb     = flag.Bool("whiteonblack", false, "white text on a black background")
	jindu    = flag.String("jindu", "117.350", "usage")
	weidu    = flag.String("weidu", "30.294", "usage")
	times    = flag.String("time", "2017-01-13 12:", "usage")
	addr     = flag.String("addr", "站点:AHCZ34/池州高坦吴角", "usage")
	IMSI     = flag.String("IMSI", "IMSI:460078157306249", "usage")
	phone = flag.String("phone", "13912649654", "usage")
	srcPath  = flag.String("src", "data/rawpic/", "usage")
	dstPath  = flag.String("dst", "data/repic/", "usage")
)

func main() {
	fileName:="data/config.conf"
	f, err := os.Open(fileName)

	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	arra := [8]string{}
	buf := bufio.NewReader(f)
	// defer buf.Close()
	i := 0
	for {

		line, err := buf.ReadString('\n')
		
		if err != nil {
			break
		}
		line = strings.TrimSpace(line) 
		// println(line)
		arra[i] = line
		i=i+1
	}
	jindu    = &arra[0]
	weidu    = &arra[1]
	times    = &arra[2]
	phone    = &arra[3]
	addr     = &arra[4]
	IMSI     = &arra[5]
	srcPath  = &arra[6]
	dstPath  = &arra[7]

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
		// println(dir )
		// println(path)
		// println(*dstPath)
		err = ThumbnailF2F(path, path, 600,800)
		if err != nil{
			fmt.Println(err)
		}
		// println(strings.Replace(path, dir, *dstPath, 1))
		genpic(path, strings.Replace(path, dir, *dstPath, 1))
		return nil
	})
}

func genpic(picPath string, dstPath string) {
	var text = []string{
		fmt.Sprintf("经度:%s%d%d%d", *jindu, getRandom(10), getRandom(10), getRandom(10)),
		fmt.Sprintf("维度:%s%d%d%d", *weidu, getRandom(10), getRandom(10), getRandom(10)),
		fmt.Sprintf("方位角:%d", getRandom(360)),
		fmt.Sprintf("%s%d%d", *times, getRandom(6), getRandom(10)),
		*phone,
		*addr,
		*IMSI,
	}
	// println(*jindu)
	// println(text[0])
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
	imgddd, err := jpeg.Decode(imgb)
	if err != nil {
		return
	}
	defer imgb.Close()
	fg, _ := image.NewUniform(color.RGBA{uint8(240), uint8(102), uint8(5), 255}), image.White
	if *wonb {
		fg, _ = image.White, image.NewUniform(color.RGBA{uint8(0), uint8(0), uint8(0), 255})
	}
	// fmt.Println(imgddd.Bounds().Dx())
	rgba := image.NewRGBA(image.Rect(0, 0,600, 800))
	// println(imgddd)
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

	pt := freetype.Pt(rgba.Bounds().Dx()-len(text[0])*20/2, 10+int(c.PointToFixed(*size)>>6))
	for _, s := range text {
		// de := 0
		// if strings.Contains(s, "经度") {
		//  de = 1
		// } else if strings.Contains(s, "维度") {
		//  de = 1
		// } else if strings.Contains(s, "方位角") {
		//  de = 2
		// }
		// fmt.Println(strings.Count(s, "")-1, len(s))
		wt := 0
		de := 0
		if strings.Count(s, "")-1 != len(s) {
			// println(int(math.Ceil(*size)))
			// println(c.PointToFixed(math.Ceil(*size))>>6)

			de = len(s) - (len(s)-strings.Count(s, "")+1)/2 
			if strings.Contains(s,"站点"){
				wt = (len(s)-strings.Count(s, "")+1)/2-4
			}
		} else {
			de = len(s) + 1
			if strings.Contains(s,"-"){
				de = de -1
			}else if !strings.Contains(s,":"){
				wt = -5
			}
			if strings.Contains(s,"IMSI"){
				de = de 
			}

		}
		pt.X = c.PointToFixed(float64(rgba.Bounds().Dx() - de*22/2 +wt))
		a, err := c.DrawString(s, pt)
		// if a != nil{
		fmt.Println(s)
		fmt.Println(a.X - a.Y)
		// }
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
	// fmt.Println("Wrote out.png OK.")

}
