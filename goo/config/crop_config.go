package config

const (
	// Crop config keys
	CC_TOP_CAM   = "crop_top-cam"
	CC_FRONT_CAM = "crop_front-cam"
	CC_DSLR      = "crop_dslr"
)

type CropConfig struct {
	BottomAbs int `yaml:"bottom_abs"`
	BottomRel int `yaml:"bottom_rel"`
	LeftAbs   int `yaml:"left_abs"`
	LeftRel   int `yaml:"left_rel"`
	RightAbs  int `yaml:"right_abs"`
	RightRel  int `yaml:"right_rel"`
	TopAbs    int `yaml:"top_abs"`
	TopRel    int `yaml:"top_rel"`
}
