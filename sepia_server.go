// Web Server for transforming images to a sepia format
package main

import (
	"fmt"
	"net/http"
	"html/template"
	"os"
	"image"
	"image/color"
	"image/jpeg"

	_ "image/gif"
	_ "image/png"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Rendering Index")
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, r.URL)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	formFile, _, err := r.FormFile("image")
	defer formFile.Close()

	if err != nil {
		fmt.Println("Unable to open file")
		return
	}

	src, _, _ := image.Decode(formFile)

	bounds := src.Bounds()
	sepia  := image.NewRGBA(bounds)

	width, height := bounds.Max.X, bounds.Max.Y
	for x:= 0; x < width; x++ {
		for y:= 0; y < height; y++ {
			sepia.Set(x, y, colorToSepia(src.At(x, y)))
		}
	}

	out, _ := os.Create("assets/sepia_image.jpg")
	defer out.Close()
	jpeg.Encode(out, sepia, &jpeg.Options{jpeg.DefaultQuality})

	fmt.Println("Rendering Show")
	t, _ := template.ParseFiles("templates/show.html")
	t.Execute(w, out)
}

func colorToSepia(src color.Color) color.Color {
	r, g, b, _ := src.RGBA()

	fr := float64(r)
	fg := float64(g)
	fb := float64(b)

	sr := fr * .393 + fg * .769 + fb * .189
	sg := fr * .349 + fg * .686 + fb * .168
	sb := fr * .272 + fg * .534 + fb * .131

	if sr > 65535.0 { sr = 65535.0 }
	if sg > 65535.0 { sg = 65535.0 }
	if sb > 65535.0 { sb = 65535.0 }

	return color.RGBA64{uint16(sr), uint16(sg), uint16(sb), ^uint16(0)}
}

func main() {
	fmt.Println("Starting server...")

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/images", createHandler)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	err := http.ListenAndServe(":8080", nil)
	if err != nil { panic(err) }
}
