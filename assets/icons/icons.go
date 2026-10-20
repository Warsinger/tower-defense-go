// Package icons embeds the application icons without loading the source photo.
package icons

import (
	"bytes"
	"embed"
	"image"
	"image/png"
)

//go:embed icon-32.png icon-48.png icon-256.png
var files embed.FS

// Images returns multiple resolutions for the native window manager.
func Images() ([]image.Image, error) {
	var images []image.Image
	for _, name := range []string{"icon-32.png", "icon-48.png", "icon-256.png"} {
		data, err := files.ReadFile(name)
		if err != nil {
			return nil, err
		}
		img, err := png.Decode(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return images, nil
}
