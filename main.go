package main

import (
	"flag"
	"fmt"
	"github.com/zuiurs/hatecom/hatena"
	"github.com/zuiurs/hatecom/log"
	"github.com/zuiurs/hatecom/recommend"
	"html"
	"os"
)

const (
	DefaultURL                 = "http://b.hatena.ne.jp/"
	DefaultLimitCategoryCount  = 100
	DefaultLimitUser           = 30
	DefaultLimitTopUser        = 5
	DefaultLimitEntry          = 20
	DefaultBorderBookmark      = 100
	DefaultBorderBookmarkCount = 100
	DefaultTargetCategory      = "it"
)

var (
	h hatena.Hatena

	url                string
	limitCategoryCount int
	limitUser          int
	limitTopUser       int
	limitEntry         int
	// user's bookmark number
	borderBookmark int
	// entry's bookmark count
	borderBookmarkCount int
	targetCategory      string
	browsingMode        bool

	categoryList = []string{
		"general",
		"social",
		"economics",
		"life",
		"knowledge",
		"it",
		"fun",
		"entertainment",
		"game",
	}
)

func main() {
	flag.StringVar(&url, "u", DefaultURL,
		"set target URL (source URL)")

	flag.StringVar(&targetCategory, "c", DefaultTargetCategory,
		"set target Category")

	flag.IntVar(&limitCategoryCount, "ncc", DefaultLimitCategoryCount,
		"(advanced) set number of entry to count\n\tThe more accurate the feature value will be.")

	flag.IntVar(&limitUser, "nuser", DefaultLimitUser,
		"(advanced) set number of user to check\n\tThe more similar users are recommended.")

	flag.IntVar(&limitTopUser, "ntop", DefaultLimitTopUser,
		"(advanced) set number of ranking")

	flag.IntVar(&limitEntry, "nentry", DefaultLimitEntry,
		"(advanced) set number of entry\n\tThe more entries are recommended.")

	flag.IntVar(&borderBookmark, "nbb", DefaultBorderBookmark,
		"(advanced) set border number of user's bookmark\n\tThe more severe users are selected.")

	flag.IntVar(&borderBookmarkCount, "nbbc", DefaultBorderBookmarkCount,
		"(advanced) set border number of bookmark\n\tThe better entries are recommended.")

	flag.IntVar(&hatena.SleepTime, "t", hatena.SleepTime,
		"(advanced) set sleep time of HTTP request interval")

	flag.BoolVar(&browsingMode, "b", false,
		"output HTML and not display log")

	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Fprintln(os.Stderr, "required an user name")
		os.Exit(1)
	}
	user := flag.Arg(0)

	// log control
	l := log.New(os.Stdout, "", 0, !browsingMode)
	hatena.DebugMode = !browsingMode

	if !categoryCheck(targetCategory) {
		fmt.Fprintln(os.Stderr, "You can use these category.")
		for _, v := range categoryList {
			fmt.Fprintln(os.Stderr, v)
		}
		os.Exit(1)
	}

	// get user's preferences
	l.Printf("---->  %s Checking Bookmarks\n", user)
	cc, err := h.GetUserCategoryCount(user, limitCategoryCount) // cc is Category Counter
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	recommend.StoreCC(user, cc)

	uccs := make([]hatena.UserCateCounter, 0)
	uccs = append(uccs, hatena.UserCateCounter{User: user, CategoryCounter: cc})

	l.Println(cc)

	// get targetURL entry
	entry, err := h.GetLiteEntry(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// get targetUsers list
	targetUsers := make([]string, len(entry.Bookmarks))
	for i, v := range entry.Bookmarks {
		targetUsers[i] = v.User
	}

	// get each targetUsers's preferences
	for i, v := range targetUsers {
		l.Printf("----> %s Checking Bookmarks (%4d/%4d)\n", v, i+1, len(entry.Bookmarks))

		// filter user who bookmarks few entry
		if count, _ := h.GetUserBMCount(v); count < borderBookmark {
			l.Printf("----> %s is skipped\n", v)
			continue
		}

		cc, err := h.GetUserCategoryCount(v, limitCategoryCount)
		uccs = append(uccs, hatena.UserCateCounter{User: v, CategoryCounter: cc})
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		recommend.StoreCC(v, cc)

		if i >= limitUser {
			break
		}
	}

	// show ranking
	// euclid
	scoresEuclid := recommend.TopMatches(recommend.Critics, user, limitTopUser, recommend.SimDistance)
	l.Println("#### Similar User Ranking (Euclid Distance)####")
	for i, v := range scoresEuclid {
		l.Printf("%2d: %.4f (%s)\n", i+1, v.Score, v.Key)
	}
	l.Println()

	// pearson
	scoresPearson := recommend.TopMatches(recommend.Critics, user, limitTopUser, recommend.SimPearson)
	l.Println("#### Similar User Ranking (Pearson Correlation Coefficient)####")
	for i, v := range scoresPearson {
		l.Printf("%2d: %.4f (%s)\n", i+1, v.Score, v.Key)
	}
	l.Println()

	// get my bookmark entry
	l.Printf("---> %s Get Bookmark Entry\n", user)
	bmlist, err := h.GetBookmarkList(user, limitEntry)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	// for checking existence
	userEntryMap := make(map[string]bool)
	for _, v := range bmlist {
		userEntryMap[v.URL] = true
	}

	// use pearson result
	selectedEntryMap := make(map[string]bool)
	result := make(map[string]string)
	for i, v := range scoresPearson {
		l.Printf("---> %s Get Bookmark Entry (%2d/%2d)\n", v.Key, i+1, len(scoresPearson))
		bl, err := h.GetBookmarkList(v.Key, limitEntry)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, e := range bl {
			// not included user's entry and deprecated
			if !(userEntryMap[e.URL] || selectedEntryMap[e.URL]) && e.Category == targetCategory {
				if e.Count > borderBookmarkCount {
					result[e.URL] = e.Title
				}
				selectedEntryMap[e.URL] = true
			}
		}
	}

	l.Println()
	l.Println("----------- Recomend Link -----------")
	for k, v := range result {
		if browsingMode {
			fmt.Printf("<a href=\"%s\">%s</a>\n<br />\n", html.EscapeString(k), html.EscapeString(v))
		} else {
			fmt.Printf("%s (%s)\n", v, k)
		}
	}
}

func categoryCheck(c string) bool {
	for _, v := range categoryList {
		if c == v {
			return true
		}
	}
	return false
}
