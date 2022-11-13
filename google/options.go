package google

import (
	"strings"
	"time"
)

const flagUsageText = "Usage: %s google [OPTIONS] QUERY\n\n" +
	"*Output Options:\n" +
	"\t-v                  Set verbose.\n\n" +
	"*Request Options:\n" +
	"\t-t DURATION         Maximum time allowed for connection / e.g. 10s, 1m, ... . (default 20s)\n" +
	"\t-H HEADER           Pass custom header(s) to google.\n" +
	"\t                    Usage: ... -H 'KEY1: VALUE1' -H 'KEY2: VALUE2'\n" +
	"\t-C COOKIE           Send cookie(s) to google. (default request to google to get cookies)\n" +
	"\t                    Usage: ... -C 'KEY=VALUE' -C 'KEY2=VALUE2'\n" +
	"\t-U User-Agent       Pass custom User-Agent header.\n\n" +
	"*Search Options:\n" +
	"\t-n NUMBER           Number of results. (default 10)\n" +
	"\t-safe               Safe search. (defualt false)\n" +
	"\t-start NUMBER       Start of results. (defualt 0)\n" +
	"\t-tld TLD            Top level domain. (default '.com')\n" +
	"\t-lang LANGUAGE      Language.\n" +
	"\t-country COUNTRY    Country or region to focus the search on.\n\n" +
	"*Query Helpers:\n" +
	"\t-inurl TEXT         ... inurl:\"TEXT\"\n" +
	"\t-intext TEXT        ... intext:\"TEXT\"\n" +
	"\t-filetype TEXT      ... filetype:\"TEXT\"\n" +
	"\t-ext TEXT           ... ext:\"TEXT\"\n"

type options struct {
	// (Output options) Verbose level
	Verbose bool

	// (Request Options) Request Cookies
	Cookies *collector

	// (Request Options) Request Headers
	Header *collector

	// (Request Options) Request User-Agent
	UserAgent string

	// (Request Options) Request timeout
	Timeout time.Duration

	// (Search Options) Query to search
	Query string

	// (Search Options) Start of results
	Start int

	// (Search Options) Top level domain
	Tld string

	// (Search Options) Language
	Lang string

	// (Search Options) Number of results
	Num int

	// (Search Options) Safe search.
	Safe bool

	// (Search Options) Country or region to focus the search on. Similar to
	// changing the Tld, but does not yield exactly the same results
	Country string

	// (Query helper) ... inurl:"TEXT" ...
	Inurl string

	// (Query helper) ... intext:"TEXT" ...
	Intext string

	// (Query helper) ... filetype:"TEXT" ...
	Filetype string

	// (Query helper) ... ext:"TEXT" ...
	Ext string
}

type collector struct {
	Seperator string
	Collected [][]string
}

func (c *collector) Set(s string) error {
	values := strings.SplitN(s, c.Seperator, 2)

	for i, v := range values {
		values[i] = strings.TrimSpace(v)
	}

	c.Collected = append(c.Collected, values)

	return nil
}

func (c collector) String() string {
	return ""
}
