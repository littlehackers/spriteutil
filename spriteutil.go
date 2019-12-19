package spriteutil

import (
    "encoding/csv"
    "image"
    "io"
    "os"
    "strconv"
    "strings"

    _ "image/png"

    "github.com/faiface/pixel"
    "github.com/pkg/errors"
)

func LoadAnimationSheet(sheetPath, descPath string, frameWidth float64) (sheet pixel.Picture, anims map[string][]pixel.Rect, err error) {
    // total hack, nicely format the error at the end, so I don't have to type it every time
    defer func() {
        if err != nil {
            err = errors.Wrap(err, "error loading animation sheet")
        }
    }()

    // open and load the spritesheet
    sheetFile, err := os.Open(sheetPath)
    if err != nil {
        return nil, nil, err
    }
    defer sheetFile.Close()
    sheetImg, _, err := image.Decode(sheetFile)
    if err != nil {
        return nil, nil, err
    }
    sheet = pixel.PictureDataFromImage(sheetImg)

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

func LoadAnimationSheetByString(sheetPath, desc string, frameWidth float64) (sheet pixel.Picture, anims map[string][]pixel.Rect, err error) {
    // total hack, nicely format the error at the end, so I don't have to type it every time
    defer func() {
        if err != nil {
            err = errors.Wrap(err, "error loading animation sheet with string")
        }
    }()

    // open and load the spritesheet
    sheetFile, err := os.Open(sheetPath)
    if err != nil {
        return nil, nil, err
    }
    defer sheetFile.Close()
    sheetImg, _, err := image.Decode(sheetFile)
    if err != nil {
        return nil, nil, err
    }
    sheet = pixel.PictureDataFromImage(sheetImg)

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
