package util

import (
	"image"
	"image/color"

	stddraw "image/draw"

	"golang.org/x/image/draw"
)

func CreateThumbnail(src image.Image, maxWidth, maxHeight int) image.Image {
	srcBounds := src.Bounds()
	srcW, srcH := srcBounds.Dx(), srcBounds.Dy()

	// Compute scaling factor
	ratioW := float64(maxWidth) / float64(srcW)
	ratioH := float64(maxHeight) / float64(srcH)
	ratio := ratioW
	if ratioH < ratioW {
		ratio = ratioH
	}

	newW := int(float64(srcW) * ratio)
	newH := int(float64(srcH) * ratio)

	// Resize into a temporary image
	resized := image.NewRGBA(image.Rect(0, 0, newW, newH))
	draw.CatmullRom.Scale(resized, resized.Bounds(), src, src.Bounds(), draw.Over, nil)

	// Create final canvas with padding (white background)
	dst := image.NewRGBA(image.Rect(0, 0, maxWidth, maxHeight))
	stddraw.Draw(dst, dst.Bounds(), &image.Uniform{color.White}, image.Point{}, stddraw.Src)

	// Center the resized image
	offsetX := (maxWidth - newW) / 2
	offsetY := (maxHeight - newH) / 2
	stddraw.Draw(dst, image.Rect(offsetX, offsetY, offsetX+newW, offsetY+newH), resized, image.Point{}, stddraw.Over)

	return dst
}
