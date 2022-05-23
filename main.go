package main

import (
	"fmt"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
)

func loadImg(filePath string) (*[][][3]uint8, error) {
	imgFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer imgFile.Close()

	image, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}

	bounds := image.Bounds()

	RGBMatrix := [][][3]uint8{}

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		newX := [][3]uint8{}

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, _ := image.At(x, y).RGBA()

			newY := [3]uint8{uint8(r / 257), uint8(g / 257), uint8(b / 257)}

			newX = append(newX, newY)
		}

		RGBMatrix = append(RGBMatrix, newX)
	}

	return &RGBMatrix, nil
}

func createImgFromMatrix(matrix *[][][3]uint8) *image.NRGBA {
	width := len(*matrix)
	height := len((*matrix)[0])

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewNRGBA(image.Rectangle{upLeft, lowRight})

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.RGBA{(*matrix)[x][y][0], (*matrix)[x][y][1], (*matrix)[x][y][2], 255})
		}
	}

	return img
}

func saveImg(img *image.NRGBA) {
	thisPath, errPath := filepath.Abs("./")
	if errPath != nil {
		fmt.Println(errPath)
		return
	}

	imgPath := thisPath + "\\img-flip-result.png"

	f, _ := os.Create(imgPath)

	png.Encode(f, img)

	err := f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Saved new image to " + imgPath)

	cmd := exec.Command("C:\\Windows\\explorer.exe", imgPath)

	cmd.Output()
}

func reverseArr[T any](arr []T) []T {
	newArr := []T{}

	for i := len(arr) - 1; i >= 0; i-- {
		newArr = append(newArr, arr[i])
	}

	return newArr
}

func main() {
	fmt.Println("Running...")

	if len(os.Args) < 2 {
		fmt.Println("Not enough args!")
		return
	}

	imgPath := os.Args[1]

	matrix, err := loadImg(imgPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	reversedMatrix := reverseArr(*matrix)

	combinedMatrix := [][][3]uint8{}

	combinedMatrix = append(combinedMatrix, (*matrix)...)

	combinedMatrix = append(combinedMatrix, reversedMatrix...)

	newImg := createImgFromMatrix(&combinedMatrix)

	fmt.Println("Completed!")

	saveImg(newImg)
}
