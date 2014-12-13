package path

import (
	"os"
)

const (
	DS          string = string(os.PathSeparator)
	ROOT_PATH   string = ".." + DS
	CONFIG_PATH string = ROOT_PATH + "configs" + DS
)
