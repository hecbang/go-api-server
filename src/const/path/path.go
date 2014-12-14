/**
 * @author yorkershi
 * @create on December 13, 2014
 */
package path

import (
	"os"
)

const (
	DS          string = string(os.PathSeparator)
	ROOT_PATH   string = ".." + DS
	CONFIG_PATH string = ROOT_PATH + "configs" + DS
)
