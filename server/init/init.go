// this package must be imported first,
// to ensure it is initialized before all the other package.
package init

import (
	"log"

	"github.com/lovego/xiaomei/config"
)

func init() {
	log.Printf(`starting.(%s)`, config.Env())
}
