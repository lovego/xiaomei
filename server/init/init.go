// this package must be imported first,
// to ensure it is initialized before all the other package.
package init

import (
	"fmt"

	"github.com/lovego/xiaomei/config"
	"github.com/lovego/xiaomei/utils"
)

func init() {
	utils.Log(fmt.Sprintf(`starting.(%s)`, config.Env()))
}
