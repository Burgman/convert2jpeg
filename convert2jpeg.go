package main

import (
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func convertToJPEG(w io.Writer, r io.Reader) (error, bool) {
	img, err := png.Decode(r)
	if err != nil {
		fmt.Println("Cannot decode png file")
		return err, false
	}

	width := img.Bounds().Size().X
	height := img.Bounds().Size().Y

	var hasTransparentBg bool

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a == 0 {
				hasTransparentBg = true
			}
		}
	}

	if hasTransparentBg {
		fmt.Println("Contains transparent bg")
		return nil, false
	} else {
		return jpeg.Encode(w, img, nil), true
	}

}

func visit(path string, f os.FileInfo, err error) error {
	fmt.Printf("Visited: %s\n", path)

	match, _ := regexp.MatchString(`.*(\.png)`, path)
	if !match {
		return nil
	}

	pngImage, err := os.Open(path)
	if err != nil {
		fmt.Println("Cannot open file at path %s", path)
		return err
	}
	defer pngImage.Close()

	newPath := renameFile(path)
	jpegImage, e := os.Create(newPath)
	if e != nil {
		fmt.Println("Cannot create new JPEG file")
		return e
	}
	defer jpegImage.Close()

	err, didConvert := convertToJPEG(jpegImage, pngImage)

	if didConvert {
		deleteFile(path)
	} else {
		deleteFile(newPath)
	}

	return nil
}

func deleteFile(path string) error {
	err := os.Remove(path)
	return err
}

func renameFile(path string) string {
	return strings.Replace(path, ".png", ".jpeg", -1)
}

func main() {
	root := "Resources"
	err := filepath.Walk(root, visit)
	fmt.Printf("filepath.Walk() returned %v\n", err)
}
