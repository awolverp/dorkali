package google

import (
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/awolverp/dorkali"
	"github.com/awolverp/dorkali/html"
)

var (
	_ dorkali.Engine = (*GoogleEngine)(nil)
	_ dorkali.Result = (*GoogleResult)(nil)
)

const (
	Version = "v1.1.4"

	stdUserAgent = "Mozilla/5.0 (Windows NT 10.0; rv:91.0) Gecko/20100101 Firefox/91.0"

	URL = "https://www.google%s/search"
)

func init() {
	dorkali.RegisterEngine("google", NewGoogleEngine)
}

type GoogleEngine struct {
	Opt options
}

func NewGoogleEngine() dorkali.Engine {
	return &GoogleEngine{
		Opt: options{
			Cookies: &collector{Seperator: "="},
			Header:  &collector{Seperator: ":"},
		},
	}
}

func (engine *GoogleEngine) Start() error {
	parser := flag.NewFlagSet("google", flag.ExitOnError)

	parser.Usage = func() { fmt.Printf("Use '%s help google' to see help information.\n", os.Args[0]) }

	parser.BoolVar(&engine.Opt.Verbose, "v", false, "")              // verbose
	parser.Var(engine.Opt.Cookies, "C", "")                          // cookies
	parser.Var(engine.Opt.Header, "H", "")                           // headers
	parser.StringVar(&engine.Opt.UserAgent, "U", stdUserAgent, "")   // user agent
	parser.DurationVar(&engine.Opt.Timeout, "t", time.Second*20, "") // timeout
	parser.IntVar(&engine.Opt.Start, "start", 0, "")                 // start
	parser.StringVar(&engine.Opt.Tld, "tld", ".com", "")             // tld
	parser.StringVar(&engine.Opt.Lang, "lang", "", "")               // lang
	parser.BoolVar(&engine.Opt.Safe, "safe", false, "")              // verbose
	parser.IntVar(&engine.Opt.Num, "n", 10, "")                      // num
	parser.StringVar(&engine.Opt.Country, "country", "", "")         // country
	parser.StringVar(&engine.Opt.Inurl, "inurl", "", "")             // inurl
	parser.StringVar(&engine.Opt.Intext, "intext", "", "")           // intext
	parser.StringVar(&engine.Opt.Filetype, "filetype", "", "")       // filetype
	parser.StringVar(&engine.Opt.Ext, "ext", "", "")                 // ext

	parser.Parse(os.Args[2:])

	engine.Opt.Query = parser.Arg(0)

	if engine.Opt.Query == "" {
		return fmt.Errorf("error: query is required. use '%s help google' to see information", os.Args[0])
	}
	return nil
}

func (engine *GoogleEngine) Version() string {
	return Version
}

func (engine *GoogleEngine) Description() string {
	return "Searches in google search engine"
}

func (engine *GoogleEngine) Usage() {
	fmt.Printf(flagUsageText, os.Args[0])
}

func (engine *GoogleEngine) Search(interface{}) (*http.Response, error) {
	if engine.Opt.Tld == "" {
		engine.Opt.Tld = ".com"
	} else if engine.Opt.Tld[0] != '.' {
		engine.Opt.Tld = "." + engine.Opt.Tld
	}

	cli := http.Client{Timeout: engine.Opt.Timeout}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.google%s/", engine.Opt.Tld), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", engine.Opt.UserAgent)

	if len(engine.Opt.Cookies.Collected) == 0 {
		resp, err := cli.Do(req)
		if err == nil {
			for _, c := range resp.Cookies() {
				req.AddCookie(c)
			}
			req.Header.Add("Referer", resp.Request.URL.String())
		}
	}

	uri := generate_url(
		engine.Opt.Query, engine.Opt.Tld, engine.Opt.Lang,
		engine.Opt.Country, engine.Opt.Inurl, engine.Opt.Intext,
		engine.Opt.Filetype, engine.Opt.Ext, engine.Opt.Num, engine.Opt.Start,
		engine.Opt.Safe,
	)

	req.URL, _ = url.Parse(uri)

	req.Header["DNT"] = []string{"1"}
	req.Header.Add("Accept", "text/html")
	req.Header.Add("Alt-Used", "www.google.com")
	req.Header.Add("Host", fmt.Sprintf("www.google%s", engine.Opt.Tld))

	for _, v := range engine.Opt.Header.Collected {
		if len(v) == 2 {
			req.Header.Add(v[0], v[1])
		}
	}

	for _, v := range engine.Opt.Cookies.Collected {
		if len(v) == 2 {
			req.AddCookie(&http.Cookie{Name: v[0], Value: v[1]})
		}
	}

	if engine.Opt.Verbose {

		print("|  " + uri + "\n\n")

		for k, values := range req.Header {
			for _, v := range values {
				println("|> " + k + ": " + v)
			}
		}
	}

	resp, err := cli.Do(req)

	if err == nil && resp.StatusCode == 403 {
		resp.Body.Close()
		return nil, fmt.Errorf("google blocked. ( returns status code 403 )")
	}

	if engine.Opt.Verbose && err == nil {

		println()

		for k, values := range resp.Header {
			for _, v := range values {
				println("|< " + k + ": " + v)
			}
		}

		println()
	}

	return resp, err
}

func (engine *GoogleEngine) ParseResponse(response *http.Response) ([]dorkali.Result, error) {
	var (
		b   []byte
		err error
	)

	if response.Header.Get("Content-Encoding") == "gzip" {
		b, err = gzipDecode(response.Body)
		fmt.Println("gzip decoded")
	} else {
		b, err = io.ReadAll(response.Body)
	}

	response.Body.Close()
	if err != nil {
		return nil, err
	}

	doc, err := html.Parse(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	var res []dorkali.Result

	doc.FindAllFunc(&html.Match{Name: "div", Attributes: map[string]string{"class": "g"}}, func(e *html.Element) {
		res = append(res, &GoogleResult{e})
	})

	return res, nil
}

func (engine *GoogleEngine) ParseHTML(h string) ([]dorkali.Result, error) {
	doc, err := html.Parse(strings.NewReader(h))
	if err != nil {
		return nil, err
	}

	var res []dorkali.Result

	doc.FindAllFunc(&html.Match{Name: "div", Attributes: map[string]string{"class": "g"}}, func(e *html.Element) {
		res = append(res, &GoogleResult{e})
	})

	return res, nil
}

func generate_url(query, tld, lang, country, inurl, intext, filetype, ext string, num, start int, safe bool) string {
	u, _ := url.Parse(fmt.Sprintf(URL, tld))
	q := u.Query()

	if lang != "" {
		q.Set("lr", "lang_"+url.QueryEscape(lang))
	}

	if safe {
		q.Set("safe", "on")
	} else {
		q.Set("safe", "off")
	}

	if country != "" {
		q.Set("cr", url.QueryEscape(country))
	}

	if start != 0 {
		q.Set("start", strconv.Itoa(start))
	}

	q.Set("num", strconv.Itoa(num+3))

	if inurl != "" {
		query += " inurl:" + inurl
	}
	if intext != "" {
		query += " intext:" + intext
	}
	if filetype != "" {
		query += " filetype:" + filetype
	}
	if ext != "" {
		query += " ext:" + ext
	}

	q.Set("q", url.QueryEscape(query))

	u.RawQuery = q.Encode()

	return u.String()
}

type GoogleResult struct {
	Doc *html.Element
}

func (r *GoogleResult) Title() string {
	title := r.Doc.Find(&html.Match{Name: "h3", Parent: &html.Match{Name: "a"}})
	if title == nil {
		return ""
	}

	return title.Text()
}

func (r *GoogleResult) Description() string {
	d := r.Doc.Find(&html.Match{Name: "span", Parent: &html.Match{Name: "div"}})

	if d == nil {
		return ""
	}

	return d.Text()
}

func (r *GoogleResult) Url() string {
	h := r.Doc.Find(&html.Match{Name: "a"})
	if h == nil {
		return ""
	}

	return filter_url(h.Attr("href"))
}

func (r *GoogleResult) String() string {
	return fmt.Sprintf("> %s\n%s\n%s\n", r.Url(), r.Title(), r.Description())
}

func filter_url(u string) string {
	if u == "" {
		return ""
	}

	parsed, err := url.Parse(u)
	if err != nil {
		return u
	}

	if parsed.Path == "/search" {
		return ""
	}

	if strings.Contains(parsed.Host, "translate.google.com") {
		return parsed.Query().Get("u")
	}

	return u
}

func gzipDecode(r io.Reader) ([]byte, error) {
	decoder, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer decoder.Close()

	b, err := io.ReadAll(decoder)
	if err != nil {
		return nil, err
	}

	return b, nil
}
