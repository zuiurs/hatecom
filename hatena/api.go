package hatena

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	// SleepTime is sleep time of GET request.
	// The unit is 'second'.
	SleepTime = 1
)

var (
	// BMNumRegexp is defined to get the number of bookmarks
	BMNumRegexp = regexp.MustCompile(`page-title.*?\(([0-9,]*)\)</h2>`)
	// BMNumRegexp is defined to get the number of bookmarks for old styled page
	BMNumRegexpOld = regexp.MustCompile(`ブックマーク数</span>.*?([0-9,]*)</li>`)

	// BMEntryRegexp is defined to get URL of each bookmark entry
	BMEntryRegexp = regexp.MustCompile(`\s<a href="(.*?)".*?(entry-link)`)

	// BMCategoryRegexp is defined to get category of each bookmark entry
	BMCategoryRegexp = regexp.MustCompile(`class="category".*?/hotentry/(.*?)"`)

	// localCache is a hashmap of URL and bookmark entry information
	localCache = make(map[string]*LiteEntry)
)

// Hatena includes user's information and etc.
type Hatena struct {
	//TODO: add user's information, OAuth token parameter and etc. in the future.
}

// LiteEntry is a json struct.
// Endpoint: http://b.hatena.ne.jp/entry/jsonlite
// Query: url, callback
type LiteEntry struct {
	Title      string     `json:"title"`
	Count      int        `json:"count"`
	URL        string     `json:"url"`
	EntryURL   string     `json:"entry_url"`
	Screenshot string     `json:"screenshot"`
	EID        int        `json:"eid"`
	Bookmarks  []Bookmark `json:"bookmarks"`
	Category   string
}

// Bookmark is a json struct depended on LiteEntry.
type Bookmark struct {
	User      string   `json:"user"`
	Tags      []string `json:"tags"`
	Timestamp string   `json:"timestamp"`
	Comment   string   `json:"comment"`
}

// GetLiteEntry returns the bookmark entry information of url,
// and caches the information to localCache.
func (h Hatena) GetLiteEntry(url string) (*LiteEntry, error) {
	return h.GetLiteEntryC(url, localCache)
}

func (h Hatena) GetLiteEntryC(url string, cache map[string]*LiteEntry) (*LiteEntry, error) {
	// check the cache
	if entry, ok := cache[url]; ok {
		return entry, nil
	}

	bs, err := getPageBytes(fmt.Sprintf("http://b.hatena.ne.jp/entry/jsonlite/?url=%s", url))
	if err != nil {
		return nil, err
	}

	var entry LiteEntry
	if err = json.Unmarshal(bs, &entry); err != nil {
		return nil, err
	}

	cache[url] = &entry

	return &entry, nil
}

// GetBookmarkList returns array of LiteEntry.
func (h Hatena) GetBookmarkList(user string, limit int) ([]*LiteEntry, error) {
	return h.GetBookmarkListC(user, limit, localCache)
}

func (h Hatena) GetBookmarkListC(user string, limit int, cache map[string]*LiteEntry) ([]*LiteEntry, error) {
	var list []*LiteEntry
	var totalBookmarks int

	// offset must be separated 20.
	for offset := 0; offset < limit; offset += 20 {
		be, err := getPageBytes(fmt.Sprintf("http://b.hatena.ne.jp/%s/?of=%d", user, offset))
		if err != nil {
			return nil, err
		}

		// First time only procedure
		// Extract total bookmark numbers
		if offset == 0 {
			var numberSubBytes [][]byte
			if numberSubBytes = BMNumRegexp.FindSubmatch(be); numberSubBytes == nil {
				fmt.Fprintln(os.Stderr, "*** Warning: Old user page detected ***")
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

			entry, err := h.GetLiteEntryC(string(v[1]), cache) // v[1] is matched URL
			if err != nil {
				return nil, err
			}
			list[index] = entry

			fmt.Printf("%5d/%5d(%d): %s\n", index+1, limit, totalBookmarks, entry.Title)
		}

		// add category field
		for i, v := range BMCategoryRegexp.FindAllSubmatch(be, -1) {
			index := offset + i
			if i+1 == offset+20 {
				break
			}
			list[index].Category = string(v[1])

		}

	}

	return list, nil
}

// CategoryCounter is a counter of category of bookmarks
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

// UserCateCounter is a wrapper for glue code in recommend package.
type UserCateCounter struct {
	User string
	*CategoryCounter
}

// GetUserCategoryCount returns CategoryCounter.
// This counts up the number of user's each category.
// If set limit -1, then limitless.
func (h Hatena) GetUserCategoryCount(user string, limit int) (*CategoryCounter, error) {
	var categoryCounter CategoryCounter
	var totalBookmarks int

	for offset := 0; offset < limit; offset += 20 {
		be, err := getPageBytes(fmt.Sprintf("http://b.hatena.ne.jp/%s/?of=%d", user, offset))

		// First time only procedure
		// Extract total bookmark numbers
		if offset == 0 {
			var numberSubBytes [][]byte
			if numberSubBytes = BMNumRegexp.FindSubmatch(be); numberSubBytes == nil {
				fmt.Fprintln(os.Stderr, "*** Warning: Old user page detected ***")
				numberSubBytes = BMNumRegexpOld.FindSubmatch(be)
			}

			numberStr := strings.Replace(string(numberSubBytes[1]), ",", "", -1) // e.g.) 1,043 -> 1043
			totalBookmarks, err = strconv.Atoi(numberStr)
			if err != nil {
				return nil, err
			}
			fmt.Printf("Total Bookmarks(%s): %d\n", user, totalBookmarks)

			if limit < 0 {
				limit = totalBookmarks
			} else if totalBookmarks < limit { // align with limit
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
			index = offset + 20
		}
		fmt.Printf("%5d/%5d(%d): Category counted\n", index, limit, totalBookmarks)
	}

	return &categoryCounter, nil
}

// GetUserBMCount returns the number of user's bookmarks.
// If it's failed to get the count, then returns -1.
func (h Hatena) GetUserBMCount(user string) (int, error) {
	bs, err := getPageBytes(fmt.Sprintf("http://b.hatena.ne.jp/%s/?of=0", user))
	if err != nil {
		return -1, err
	}

	var numberSubBytes [][]byte
	if numberSubBytes = BMNumRegexp.FindSubmatch(bs); numberSubBytes == nil {
		numberSubBytes = BMNumRegexpOld.FindSubmatch(bs)
	}

	numberStr := strings.Replace(string(numberSubBytes[1]), ",", "", -1) // e.g.) 1,043 -> 1043
	totalBookmarks, err := strconv.Atoi(numberStr)
	if err != nil {
		return -1, err
	}

	return totalBookmarks, nil
}

// getPageBytes returns bytes array of responce indicated URL.
// The request method is GET.
// The URL should include query parameter.
func getPageBytes(url string) ([]byte, error) {
	resp, err := http.Get(url)
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

	//-----------------------------------------------------
	time.Sleep(SleepTime * time.Second)
	//-----------------------------------------------------

	return be, nil
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
