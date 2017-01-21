# hatecom

hatecom is a scraper to recommend entry of Hatena Bookmark.

## Usage

Simple usage.

```
hatecom -u http://b.hatena.ne.jp/ ZuiUrs
```

You can use some options.

```
hatecom -c it -u http://b.hatena.ne.jp/ -b ZuiUrs > bookmark_list.html
```

Other option information is given by help.

```
Usage of hatecom:
  -b    output HTML and not display log
  -c string
        set target Category (default "it")
  -nbb int
        (advanced) set border number of user's bookmark
        The more severe users are selected. (default 100)
  -nbbc int
        (advanced) set border number of bookmark
        The better entries are recommended. (default 100)
  -ncc int
        (advanced) set number of entry to count
        The more accurate the feature value will be. (default 100)
  -nentry int
        (advanced) set number of entry
        The more entries are recommended. (default 20)
  -ntop int
        (advanced) set number of ranking (default 5)
  -nuser int
        (advanced) set number of user to check
        The more similar users are recommended. (default 30)
  -t int
        (advanced) set sleep time of HTTP request interval (default 1)
  -u string
        set target URL (source URL) (default "http://b.hatena.ne.jp/")
```

## Installation

You can use `go get`.

```
go get -u github.com/zuiurs/hatecom
```

## Caution

In scraiping web contents, you should follow dosage and administration.  
I shall not be responsible for any loss, damages and troubles.

## TODO

- judge category for recommendation automatically

## License

This software is released under the MIT License, see LICENSE.txt.
