package main

import (
	"flag"
	"fmt"
	//"github.com/k0kubun/pp"
	"github.com/zuiurs/grawl/hatena"
	"log"
	"os"
)

var (
	h   hatena.Hatena
	url string
)

func main() {
	flag.StringVar(&url, "-u", "http://b.hatena.ne.jp/", "Target URL")
	flag.Parse()

	user := flag.Arg(0)
	fmt.Printf("---->  %s start...\n", user)
	//userBMList, err := h.GetBookmarkList(user, 5)
	//if err != nil {
	//	log.Println(err)
	//	os.Exit(1)
	//}
	//fmt.Println(userBMList)
	userCategory, err := h.GetUserCategoryCount(user, 100)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Println(userCategory)

	entry, err := h.GetLiteEntry(url)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	targetUsers := make([]string, len(entry.Bookmarks))
	for i, v := range entry.Bookmarks {
		targetUsers[i] = v.User
	}

	for i, v := range targetUsers {
		fmt.Printf("----> (%4d/%4d) %s start...\n", i+1, len(entry.Bookmarks), v)
		category, err := h.GetUserCategoryCount(v, 100)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		fmt.Println(category)
	}

}
