// glue code for hatena and recommend

package recommend

import (
	"fmt"
	"github.com/zuiurs/hatecom/hatena"
)

func StoreCC(user string, cc *hatena.CategoryCounter) {
	m := make(map[string]float64)

	if cc.General != 0 {
		m["general"] = float64(cc.General)
	}
	if cc.Social != 0 {
		m["social"] = float64(cc.Social)
	}
	if cc.Economics != 0 {
		m["economics"] = float64(cc.Economics)
	}
	if cc.Life != 0 {
		m["life"] = float64(cc.Life)
	}
	if cc.Knowledge != 0 {
		m["knowledge"] = float64(cc.Knowledge)
	}
	if cc.It != 0 {
		m["it"] = float64(cc.It)
	}
	if cc.Fun != 0 {
		m["fun"] = float64(cc.Fun)
	}
	if cc.Entertainment != 0 {
		m["entertainment"] = float64(cc.Entertainment)
	}
	if cc.Game != 0 {
		m["game"] = float64(cc.Game)
	}

	Critics[user] = m
}

func OutputCategoryCode(uccs []hatena.UserCateCounter) {
	fmt.Printf("var Critics = map[string]map[string]float64{\n")
	for _, ucc := range uccs {
		fmt.Printf("\t\"%s\": map[string]float64{\n", ucc.User)
		if ucc.General != 0 {
			fmt.Printf("\t\t\"general\": %f,\n", float64(ucc.General))
		}
		if ucc.Social != 0 {
			fmt.Printf("\t\t\"social\": %f,\n", float64(ucc.Social))
		}
		if ucc.Economics != 0 {
			fmt.Printf("\t\t\"economics\": %f,\n", float64(ucc.Economics))
		}
		if ucc.Life != 0 {
			fmt.Printf("\t\t\"life\": %f,\n", float64(ucc.Life))
		}
		if ucc.Knowledge != 0 {
			fmt.Printf("\t\t\"knowledge\": %f,\n", float64(ucc.Knowledge))
		}
		if ucc.It != 0 {
			fmt.Printf("\t\t\"it\": %f,\n", float64(ucc.It))
		}
		if ucc.Fun != 0 {
			fmt.Printf("\t\t\"fun\": %f,\n", float64(ucc.Fun))
		}
		if ucc.Entertainment != 0 {
			fmt.Printf("\t\t\"entertainment\": %f,\n", float64(ucc.Entertainment))
		}
		if ucc.Game != 0 {
			fmt.Printf("\t\t\"game\": %f,\n", float64(ucc.Game))
		}
		fmt.Printf("\t},\n")
	}

	fmt.Printf("}\n")
}
