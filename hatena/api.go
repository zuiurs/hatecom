package hatena

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	SleepTime = 3 // second base
)

var (
	BMNumRegexp    = regexp.MustCompile(`page-title.*?\(([0-9,]*)\)</h2>`)
	BMNumRegexpOld = regexp.MustCompile(`ブックマーク数</span>.*?([0-9,]*)</li>`)

	BMEntryRegexp = regexp.MustCompile(`\s<a href="(.*?)".*?(entry-link)`)

	BMCategoryRegexp = regexp.MustCompile(`class="category".*?/hotentry/(.*?)"`)

	urlCache = make(map[string]*LiteEntry)
)

type Hatena struct {
	//TODO: oauth 認証などのトークンを格納する構造体にする
}

type LiteEntry struct {
	Title      string     `json:"title"`
	Count      int        `json:"count"`
	URL        string     `json:"url"`
	EntryURL   string     `json:"entry_url"`
	Screenshot string     `json:"screenshot"`
	EID        int        `json:"eid"`
	Bookmarks  []Bookmark `json:"bookmarks"`
}

type Bookmark struct {
	User      string   `json:"user"`
	Tags      []string `json:"tags"`
	Timestamp string   `json:"timestamp"`
	Comment   string   `json:"comment"`
}

func (h Hatena) GetLiteEntry(url string) (*LiteEntry, error) {
	if entry, ok := urlCache[url]; ok {
		return entry, nil
	}

	resp, err := http.Get(fmt.Sprint("http://b.hatena.ne.jp/entry/jsonlite/?url=", url))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("StatusCode is %d", resp.StatusCode)
	}

	be, err := ioutil.ReadAll(resp.Body) //get byte based entry
	if err != nil {
		return nil, err
	}

	var entry LiteEntry
	if err = json.Unmarshal(be, &entry); err != nil {
		return nil, err
	}

	urlCache[url] = &entry

	//-----------------------------------------------------
	time.Sleep(SleepTime * time.Second)
	//-----------------------------------------------------

	return &entry, nil
}

func (h Hatena) GetBookmarkList(user string, limit int) ([]*LiteEntry, error) {
	var list []*LiteEntry
	var totalBookmarks int

	for offset := 0; offset < limit; offset += 20 {
		resp, err := http.Get(fmt.Sprintf("http://b.hatena.ne.jp/%s/?of=%d", user, offset))
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("StatusCode is %d", resp.StatusCode)
		}

		be, err := ioutil.ReadAll(resp.Body) //get byte based entry
		if err != nil {
			return nil, err
		}

		// First time only procedure
		// Extract total bookmark numbers
		if offset == 0 {
			var numberSubBytes [][]byte
			if numberSubBytes = BMNumRegexp.FindSubmatch(be); numberSubBytes == nil {
				fmt.Println("*** Detect old user page ***")
				numberSubBytes = BMNumRegexpOld.FindSubmatch(be)
			}

			numberStr := strings.Replace(string(numberSubBytes[1]), ",", "", -1) // e.g.) 1,043 -> 1043
			totalBookmarks, err = strconv.Atoi(numberStr)
			if err != nil {
				return nil, err
			}
			fmt.Printf("Total Bookmarks(%s): %d\n", user, totalBookmarks)

			if totalBookmarks < limit { // align with limit
				limit = totalBookmarks
			}
			list = make([]*LiteEntry, limit)
		}

		for i, v := range BMEntryRegexp.FindAllSubmatch(be, -1) {
			index := offset + i
			if index > limit-1 {
				break
			}

			entry, err := h.GetLiteEntry(string(v[1])) // v[1] is matched URL
			if err != nil {
				return nil, err
			}
			list[index] = entry

			fmt.Printf("%5d/%5d(%d): %s\n", index+1, limit, totalBookmarks, entry.Title)
		}
		//-----------------------------------------------------
		time.Sleep(SleepTime * time.Second)
		//-----------------------------------------------------
	}

	return list, nil
}

type CategoryCounter struct {
	General       int // 一般
	Social        int // 世の中
	Economics     int // 政治と経済
	Life          int // 暮らし
	Knowledge     int // 学び
	It            int // テクノロジー
	Fun           int // おもしろ
	Entertainment int // エンタメ
	Game          int // アニメとゲーム
}

func (h Hatena) GetUserCategoryCount(user string, limit int) (*CategoryCounter, error) {
	var categoryCounter CategoryCounter
	var totalBookmarks int

	for offset := 0; offset < limit; offset += 20 {
		resp, err := http.Get(fmt.Sprintf("http://b.hatena.ne.jp/%s/?of=%d", user, offset))
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("StatusCode is %d", resp.StatusCode)
		}

		be, err := ioutil.ReadAll(resp.Body) //get byte based entry
		if err != nil {
			return nil, err
		}

		// First time only procedure
		// Extract total bookmark numbers
		if offset == 0 {
			var numberSubBytes [][]byte
			if numberSubBytes = BMNumRegexp.FindSubmatch(be); numberSubBytes == nil {
				fmt.Println("*** Detect old user page ***")
				numberSubBytes = BMNumRegexpOld.FindSubmatch(be)
			}

			numberStr := strings.Replace(string(numberSubBytes[1]), ",", "", -1) // e.g.) 1,043 -> 1043
			totalBookmarks, err = strconv.Atoi(numberStr)
			if err != nil {
				return nil, err
			}
			fmt.Printf("Total Bookmarks(%s): %d\n", user, totalBookmarks)

			if totalBookmarks < limit { // align with limit
				limit = totalBookmarks
			}
		}

		for _, v := range BMCategoryRegexp.FindAllSubmatch(be, -1) {
			switch string(v[1]) {
			case "general":
				categoryCounter.General++
			case "social":
				categoryCounter.Social++
			case "economics":
				categoryCounter.Economics++
			case "life":
				categoryCounter.Life++
			case "knowledge":
				categoryCounter.Knowledge++
			case "it":
				categoryCounter.It++
			case "fun":
				categoryCounter.Fun++
			case "entertainment":
				categoryCounter.Entertainment++
			case "game":
				categoryCounter.Game++
			default:
				fmt.Printf("Category %s is not found\n", v[1])
			}
		}

		var index int
		if offset+20 > limit {
			index = limit
		} else {
			index = offset
		}
		fmt.Printf("%5d/%5d(%d): Category counted\n", index, limit, totalBookmarks)

		//-----------------------------------------------------
		time.Sleep(SleepTime * time.Second)
		//-----------------------------------------------------
	}

	return &categoryCounter, nil
}

func (c CategoryCounter) String() string {
	return fmt.Sprintf("General      : %d\n", c.General) +
		fmt.Sprintf("Social:      : %d\n", c.Social) +
		fmt.Sprintf("Economics    : %d\n", c.Economics) +
		fmt.Sprintf("Life         : %d\n", c.Life) +
		fmt.Sprintf("Knowledge    : %d\n", c.Knowledge) +
		fmt.Sprintf("It           : %d\n", c.It) +
		fmt.Sprintf("Fun          : %d\n", c.Fun) +
		fmt.Sprintf("Entertainment: %d\n", c.Entertainment) +
		fmt.Sprintf("Game         : %d\n", c.Game)
}
