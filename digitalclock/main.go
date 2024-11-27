//go:build !solution

package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	charWidth    = 8
	charHeight   = 12
	colonWidth   = 4
	defaultScale = 10
)

var digitMap = map[rune]string{
	'0': Zero,
	'1': One,
	'2': Two,
	'3': Three,
	'4': Four,
	'5': Five,
	'6': Six,
	'7': Seven,
	'8': Eight,
	'9': Nine,
	':': Colon,
}

func renderTimeToImage(timeString string, scale int) (image.Image, error) {
	if len(timeString) != 8 || timeString[2] != ':' || timeString[5] != ':' {
		return nil, errors.New("invalid time format")
	}

	imgWidth := (6*charWidth + 2*colonWidth) * scale
	imgHeight := charHeight * scale

	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)

	xOffset := 0
	for _, ch := range timeString {
		symbol, ok := digitMap[ch]
		if !ok {
			return nil, errors.New("invalid character in time string")
		}

		drawSymbol(img, symbol, xOffset, scale)
		if ch == ':' {
			xOffset += colonWidth * scale
		} else {
			xOffset += charWidth * scale
		}
	}

	return img, nil
}

func drawSymbol(img *image.RGBA, symbol string, xOffset, scale int) {
	yOffset := 0
	for _, line := range strings.Split(symbol, "\n") {
		for x, pixel := range line {
			if pixel == '1' {
				fillRect(img, xOffset+x*scale, yOffset, scale, scale, Cyan)
			}
		}
		yOffset += scale
	}
}

func fillRect(img *image.RGBA, x, y, w, h int, col color.Color) {
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			img.Set(x+i, y+j, col)
		}
	}
}

func isValidTime(t string) bool {
	parts := strings.Split(t, ":")
	if len(parts) != 3 {
		return false
	}

	hours, err1 := strconv.Atoi(parts[0])
	minutes, err2 := strconv.Atoi(parts[1])
	seconds, err3 := strconv.Atoi(parts[2])

	return err1 == nil && err2 == nil && err3 == nil &&
		hours >= 0 && hours <= 23 &&
		minutes >= 0 && minutes <= 59 &&
		seconds >= 0 && seconds <= 59
}

func clockHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	timeParam := query.Get("time")
	if timeParam == "" {
		timeParam = time.Now().Format("15:04:05")
	} else if !isValidTime(timeParam) {
		http.Error(w, "invalid time format", http.StatusBadRequest)
		return
	}

	scaleParam := query.Get("k")
	scale := defaultScale
	if scaleParam != "" {
		s, err := strconv.Atoi(scaleParam)
		if err != nil || s < 1 || s > 30 {
			http.Error(w, "invalid k", http.StatusBadRequest)
			return
		}
		scale = s
	}

	img, err := renderTimeToImage(timeParam, scale)
	if err != nil {
		http.Error(w, "invalid time format", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	if err := png.Encode(w, img); err != nil {
		http.Error(w, "failed to encode image", http.StatusInternalServerError)
		return
	}
}

func main() {
	port := flag.Int("port", 6029, "port to listen on")
	flag.Parse()
	http.HandleFunc("/", clockHandler)
	addr := fmt.Sprintf(":%d", *port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Println("Server error:", err)
		os.Exit(1)
	}
}
