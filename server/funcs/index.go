package funcs

import (
	"strings"

	"github.com/lovego/xiaomei/utils/slice"
	"github.com/lovego/xiaomei/utils/strnum"
)

func Map() map[string]interface{} {
	return map[string]interface{}{
		`asset`:        AssetFunc(),
		`html_safe`:    HtmlSafe,
		`dict`:         MakeDict,
		`keys`:         MapKeys,
		`values`:       MapValues,
		`IF`:           IF,
		`field`:        StructOrMapField,
		`union`:        slice.Union,
		`keys_union`:   MapKeysUnion,
		`thousand_sep`: strnum.ThousandSep,
		`contains`:     strings.Contains,
	}
}
