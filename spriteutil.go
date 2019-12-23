package spriteutil

import (
	"encoding/csv"
	"image"
	"io"
	"os"
	"strconv"
	"strings"

	_ "image/png" //decode png format images

	"github.com/faiface/pixel"
	"github.com/pkg/errors"
)

//image loading functions

//loadPicture decodes an image file at path and returns a pixel.Picture
func LoadPicture(path string) (pic pixel.Picture, err error) {
	//borrowed hack
	defer func() {
		if err != nil {
			err = errors.Wrap(err, "Error loading picture")
		}
	}()

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return pixel.PictureDataFromImage(img), nil
}

//LoadAnimationSheet decodes a sheet of animation frames saved in an image file at path, slicing
//the frames at frameWidth, it returns a pixel.Picture and array of pixel.Rect with the coordinates
//to each frame of the sheet of animation
func LoadAnimationSheet(sheetPath string, frameWidth float64) (sheet pixel.Picture, anims []pixel.Rect, err error) {
	// total hack, nicely format the error at the end, so I don't have to type it every time
	defer func() {
		if err != nil {
			err = errors.Wrap(err, "error loading animation sheet")
		}
	}()

	// open and load the spritesheet
	sheet, _ = LoadPicture(sheetPath)

	// create a slice of frames inside the spritesheet
	var frames []pixel.Rect
	for x := 0.0; x+frameWidth <= sheet.Bounds().Max.X; x += frameWidth {
		frames = append(frames, pixel.R(
			x,
			0,
			x+frameWidth,
			sheet.Bounds().H(),
		))
	}
	//skip lableing frames and return them immediately
	return sheet, frames, nil

}

//LoadAnimationSheetByCSV decodes a sheet of animation frames saved in an image file at path, slicing
//the frames at frameWidth and labeling the sequence of frames based on the contents of the description
//in descPath, a CSV formated file.  It returns a pixel.Picture and a map containing labeled sequences
//to the coordinates of each frame of the sheet of animation.
func LoadAnimationSheetByCSV(sheetPath, descPath string, frameWidth float64) (sheet pixel.Picture, anims map[string][]pixel.Rect, err error) {
	// total hack, nicely format the error at the end, so I don't have to type it every time
	defer func() {
		if err != nil {
			err = errors.Wrap(err, "error loading animation sheet")
		}
	}()

	// open and load the animation sheet
	sheet, frames, err := LoadAnimationSheet(sheetPath, frameWidth)
	if err != nil {
		return nil, nil, err
	}

	//load CSV description file
	descFile, err := os.Open(descPath)
	if err != nil {
		return nil, nil, err
	}
	defer descFile.Close()

	anims = make(map[string][]pixel.Rect)

	// load the animation information, name and interval inside the spritesheet
	desc := csv.NewReader(descFile)
	for {
		anim, err := desc.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}

		name := anim[0]
		start, _ := strconv.Atoi(anim[1])
		end, _ := strconv.Atoi(anim[2])

		anims[name] = frames[start : end+1]
	}

	return sheet, anims, nil
}

//LoadAnimationSheetByString decodes a sheet of animation frames saved in an image file at path, slicing
//the frames at frameWidth and labeling the sequence of frames based on the contents the comma separated
//string desc.  It returns a pixel.Picture and a map containing labeled sequences to the coordinates of
//each frame of the sheet of animation.
func LoadAnimationSheetByString(sheetPath, desc string, frameWidth float64) (sheet pixel.Picture, anims map[string][]pixel.Rect, err error) {
	// total hack, nicely format the error at the end, so I don't have to type it every time
	defer func() {
		if err != nil {
			err = errors.Wrap(err, "error loading animation sheet")
		}
	}()

	// open and load the animation sheet
	sheet, frames, err := LoadAnimationSheet(sheetPath, frameWidth)
	if err != nil {
		return nil, nil, err
	}

	anims = make(map[string][]pixel.Rect)

	// load the animation information, name and interval inside the spritesheet
	var anim []string
	for _, d := range strings.Split(desc, ",") {
		anim = append(anim, d)
		if len(anim) >= 3 {
			name := anim[0]
			start, _ := strconv.Atoi(anim[1])
			end, _ := strconv.Atoi(anim[2])

			anims[name] = frames[start : end+1]
			anim = anim[:0]
		}
	}

	return sheet, anims, nil
}
