package main

import (
	"flag"
	"fmt"
	//	"github.com/k0kubun/pp"
	"github.com/zuiurs/hatecom/hatena"
	"github.com/zuiurs/hatecom/recommend"
	"log"
	"os"
)

const (
	DefaultURL     = "http://b.hatena.ne.jp/"
	TargetCategory = "it"

	LimitCategoryCount = 100
	LimitUser          = 30
	LimitTopUser       = 5
	LimitEntry         = 20

	// user's bookmark number
	BorderBookmark = 100
	// entry's bookmark count
	BorderBookmarkCount = 100
)

var (
	h   hatena.Hatena
	url string
)

func main() {
	flag.StringVar(&url, "u", DefaultURL, "Target URL")
	flag.Parse()
	user := flag.Arg(0)

	// get user's preferences
	fmt.Printf("---->  %s Checking Bookmarks\n", user)
	cc, err := h.GetUserCategoryCount(user, LimitCategoryCount) // cc is Category Counter
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	recommend.StoreCC(user, cc)

	uccs := make([]hatena.UserCateCounter, 0)
	uccs = append(uccs, hatena.UserCateCounter{User: user, CategoryCounter: cc})

	fmt.Println(cc)

	// get targetURL entry
	entry, err := h.GetLiteEntry(url)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// get targetUsers list
	targetUsers := make([]string, len(entry.Bookmarks))
	for i, v := range entry.Bookmarks {
		targetUsers[i] = v.User
	}

	// get each targetUsers's preferences
	for i, v := range targetUsers {
		fmt.Printf("----> %s Checking Bookmarks (%4d/%4d)\n", v, i+1, len(entry.Bookmarks))

		// filter user who bookmarks few entry
		if count, _ := h.GetUserBMCount(v); count < BorderBookmark {
			fmt.Printf("----> %s is skipped\n", v)
			continue
		}

		cc, err := h.GetUserCategoryCount(v, BorderBookmark)
		uccs = append(uccs, hatena.UserCateCounter{User: v, CategoryCounter: cc})
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		recommend.StoreCC(v, cc)

		if i >= LimitUser {
			break
		}
	}

	// Code Generator
	//recommend.OutputCategoryCode(uccs)

	//------------transform category struct to recomend data type-------------------
	//pp.Print(recommend.Critics)

	// user と似ている人を順位付けして表示
	// euclid
	scoresEuclid := recommend.TopMatches(recommend.Critics, user, LimitTopUser, recommend.SimDistance)
	fmt.Println("#### Similar User Ranking (Euclid Distance)####")
	for i, v := range scoresEuclid {
		fmt.Printf("%2d: %.4f (%s)\n", i+1, v.Score, v.Key)
	}
	fmt.Println()

	// pearson
	scoresPearson := recommend.TopMatches(recommend.Critics, user, LimitTopUser, recommend.SimPearson)
	fmt.Println("#### Similar User Ranking (Pearson Correlation Coefficient)####")
	for i, v := range scoresPearson {
		fmt.Printf("%2d: %.4f (%s)\n", i+1, v.Score, v.Key)
	}
	fmt.Println()

	// 自分のエントリを取得
	fmt.Printf("---> %s Get Bookmark Entry\n", user)
	bmlist, err := h.GetBookmarkList(user, LimitEntry)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	// for judging exist
	userEntryMap := make(map[string]bool)
	for _, v := range bmlist {
		userEntryMap[v.URL] = true
	}

	// use pearson result
	selectedEntryMap := make(map[string]bool)
	var result []string
	for i, v := range scoresPearson {
		fmt.Printf("---> %s Get Bookmark Entry (%2d/%2d)\n", v.Key, i+1, len(scoresPearson))
		bl, err := h.GetBookmarkList(v.Key, LimitEntry)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		for _, e := range bl {
			// どちらにも含まれていない且つ目的のカテゴリ
			if !(userEntryMap[e.URL] || selectedEntryMap[e.URL]) && e.Category == TargetCategory {
				if e.Count > BorderBookmarkCount {
					result = append(result, fmt.Sprintf("%s ( %s )", e.Title, e.URL))
				}
				selectedEntryMap[e.URL] = true
			}
		}
	}

	fmt.Println()
	fmt.Println("----------- Recomend Link -----------")
	for _, v := range result {
		fmt.Println(v)
	}
}
