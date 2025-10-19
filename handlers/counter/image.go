package counter

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"time"

	"github.com/mn6/tinycounter/store"
	"github.com/rs/zerolog/log"
)

func GenerateCounter(s store.Store, key string, styles Styles, style string, count int64, cacheTime time.Duration) []byte {
	// Check if counter already was generated
	data, err := GetExistingCounter(s, key, style)
	if err == nil {
		return data
	}

	imgData := CreateCounterImage(styles[style], count)

	// Save the image async so it doesn't block response, with a TTL
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		err := s.SetImage(ctx, key, style, imgData, cacheTime)
		if err != nil {
			log.Error().Err(err).Msg("Failed to save counter image to store")
		}
		if err != nil {
			log.Error().Err(err).Msg("Failed to save counter image to store")
		}
	}()

	return imgData
}

func GetExistingCounter(s store.Store, key string, style string) ([]byte, error) {
	data, err := s.GetImage(context.Background(), key, style)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func CreateCounterImage(styleConfig StyleConfig, count int64) []byte {
	// Each number in count as a string
	// Prefixed with 0 to ensure at least 2 digits
	numStr := fmt.Sprintf("0%d", count)

	// Width = prefix + digits*digitWidth + suffix
	w := styleConfig.PrefixWidth + (len(numStr) * styleConfig.Width) + styleConfig.SuffixWidth
	h := styleConfig.Height

	img := image.NewRGBA(image.Rect(0, 0, w, h))

	x := 0

	// Draw prefix if present
	if styleConfig.PrefixWidth > 0 {
		prefixImg := GetDigitImage(styleConfig.Name, "P")
		drawImage(img, prefixImg, x, 0)
		x += styleConfig.PrefixWidth
	}

	// Draw digits
	for _, ch := range numStr {
		digitImg := GetDigitImage(styleConfig.Name, string(ch))
		drawImage(img, digitImg, x, 0)
		x += styleConfig.Width
	}

	// Draw suffix if present
	if styleConfig.SuffixWidth > 0 {
		suffixImg := GetDigitImage(styleConfig.Name, "S")
		drawImage(img, suffixImg, x, 0)
	}

	// Return img as PNG byte slice
	return EncodePNGToBytes(img)
}

func EncodePNGToBytes(img image.Image) []byte {
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil
	}
	return buf.Bytes()
}

func GetDigitImage(style string, digit string) image.Image {
	return StyleImages[style][digit]
}

func drawImage(dst *image.RGBA, src image.Image, xOffset, yOffset int) {
	if src == nil || dst == nil {
		return
	}

	sp := src.Bounds()
	dp := image.Rect(xOffset+sp.Min.X, yOffset+sp.Min.Y, xOffset+sp.Max.X, yOffset+sp.Max.Y)
	draw.Draw(dst, dp.Intersect(dst.Bounds()), src, sp.Min, draw.Over)
}
