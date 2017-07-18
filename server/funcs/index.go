package funcs

import (
	"strings"

	"github.com/lovego/xiaomei/utils/slice/union"
	"github.com/lovego/xiaomei/utils/strnum"
)

func Index() map[string]interface{} {
	return map[string]interface{}{
		`asset`:        AssetFunc(),
		`html_safe`:    HtmlSafe,
		`dict`:         MakeDict,
		`keys`:         MapKeys,
		`values`:       MapValues,
		`IF`:           IF,
		`field`:        StructOrMapField,
		`union`:        union.Union,
		`keys_union`:   MapKeysUnion,
		`thousand_sep`: strnum.ThousandSep,
		`contains`:     strings.Contains,
	}
}
