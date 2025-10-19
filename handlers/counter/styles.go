package counter

import (
	"fmt"
	"image"
	"os"

	"go.yaml.in/yaml/v3"

	"github.com/mn6/tinycounter/utils"
	"github.com/rs/zerolog/log"
)

const StylesConfigPath = "configs/styles.yaml"
const StylesResourcesPath = "resources/styles"

type StyleConfig struct {
	Width       int `yaml:"width"`
	Height      int `yaml:"height"`
	SuffixWidth int `yaml:"suffixWidth"`
	PrefixWidth int `yaml:"prefixWidth"`
	Name        string
}

type Styles map[string]StyleConfig

var StyleSet Styles
var StyleImages map[string]map[string]image.Image

func init() {
	// StyleSet initialization
	var err error
	StyleSet, err = LoadStyles()

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load styles")
	}

	for k, v := range StyleSet {
		v.Name = k
		StyleSet[k] = v
	}

	// Preload style images into StyleImages map for quick access during request handling
	StyleImages = make(map[string]map[string]image.Image)
	for styleName := range StyleSet {
		StyleImages[styleName] = make(map[string]image.Image)
		stylePath := fmt.Sprintf("%s/%s", StylesResourcesPath, styleName)

		// Load digits 0-9
		for i := 0; i <= 9; i++ {
			filename := fmt.Sprintf("%d.png", i)
			imgPath := fmt.Sprintf("%s/%s", stylePath, filename)
			img, err := utils.LoadImage(imgPath)
			if err != nil {
				log.Error().Err(err).Msgf("Failed to load image: %s", imgPath)
				continue
			}
			StyleImages[styleName][fmt.Sprintf("%d", i)] = img
		}

		// Load prefix and suffix if they exist
		prefixPath := fmt.Sprintf("%s/pre.png", stylePath)
		if img, err := utils.LoadImage(prefixPath); err == nil {
			StyleImages[styleName]["P"] = img
		}

		suffixPath := fmt.Sprintf("%s/suf.png", stylePath)
		if img, err := utils.LoadImage(suffixPath); err == nil {
			StyleImages[styleName]["S"] = img
		}
	}
}

// ValidateStyle returns true if the provided style string is one of the allowed styles.
func ValidateStyle(style string) bool {
	_, ok := StyleSet[style]
	return ok
}

func LoadStyles() (Styles, error) {
	var styles Styles
	data, err := os.ReadFile(StylesConfigPath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &styles)
	return styles, err
}
