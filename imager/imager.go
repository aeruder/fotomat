package imager

import (
	"errors"
)

var (
	UnknownFormat = errors.New("Unknown image format")
	TooBig        = errors.New("Image is too wide or tall")
)

const (
	minDimension = 2             // Avoid off-by-one divide-by-zero errors.
	maxDimension = (2 << 14) - 2 // Avoid signed int16 overflows.
)

type Imager struct {
	blob         []byte
	Width        uint
	Height       uint
	InputFormat  string
	OutputFormat string
	JpegQuality  uint
	Sharpen      bool
}

func New(blob []byte) (*Imager, error) {
	// Security: Guess at formats.  Limit formats we pass to ImageMagick
	// to just JPEG, PNG, GIF.
	inputFormat, outputFormat := detectFormats(blob)
	if inputFormat == "" {
		return nil, UnknownFormat
	}

	// Ask ImageMagick to parse metadata.
	width, height, format, err := imageMetaData(blob)
	if err != nil {
		return nil, UnknownFormat
	}

	// Security: Confirm that detectFormat() and imageMagick agreed on
	// format and that image sizes are sane.
	if format != inputFormat {
		return nil, UnknownFormat
	} else if width < minDimension || height < minDimension {
		return nil, UnknownFormat
	} else if width > maxDimension || height > maxDimension {
		return nil, TooBig
	}

	img := &Imager{
		blob:         blob,
		Width:        width,
		Height:       height,
		InputFormat:  inputFormat,
		OutputFormat: outputFormat,
		JpegQuality:  85,
		Sharpen:      true,
	}

	return img, nil
}

func (img *Imager) Thumbnail(width, height uint, within bool) ([]byte, error) {
	width, height = scaleAspect(img.Width, img.Height, width, height, within)

	result, err := img.NewResult(width, height)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	if result.Width > width || result.Height > height {
		if err := result.Resize(width, height); err != nil {
			return nil, err
		}
	}

	return result.Get()
}

func (img *Imager) Close() {
	*img = Imager{}
}
