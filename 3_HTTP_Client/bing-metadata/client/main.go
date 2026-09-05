package main

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"bing-metadata/metadata"

	"github.com/PuerkitoBio/goquery"
)

// extractRealURL decodes Bing's tracking redirect link to get the actual
// document URL. Bing wraps real links as bing.com/ck/a?...&u=a1<base64>&...
func extractRealURL(bingLink string) (string, error) {
	parsed, err := url.Parse(bingLink)
	if err != nil {
		return "", err
	}
	u := parsed.Query().Get("u")
	if len(u) < 2 {
		return "", fmt.Errorf("no encoded URL found")
	}
	// Bing prefixes the base64 payload with "a1" — strip it before decoding.
	encoded := strings.TrimPrefix(u, "a1")
	decoded, err := base64.RawURLEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

func handler(domain string, i int, s *goquery.Selection) {
	rawURL, ok := s.Find("a").Attr("href")
	if !ok {
		return
	}

	docURL, err := extractRealURL(rawURL)
	if err != nil {
		return
	}

	// Safety check: skip results that don't actually belong to the domain
	// we searched for. Bing sometimes falls back to unrelated "related
	// searches" results when there are zero real matches.
	if !strings.Contains(docURL, domain) {
		return
	}

	fmt.Printf("%d: %s\n", i, docURL)

	res, err := http.Get(docURL)
	if err != nil {
		return
	}
	defer res.Body.Close()

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}

	r, err := zip.NewReader(bytes.NewReader(buf), int64(len(buf)))
	if err != nil {
		return
	}

	cp, ap, err := metadata.NewProperties(r)
	if err != nil {
		return
	}

	log.Printf(
		"%25s %25s - %s %s\n",
		cp.Creator,
		cp.LastModifiedBy,
		ap.Application,
		ap.GetMajorVersion())
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Usage: main.go domain")
	}
	domain := os.Args[1]

	filetypes := []string{"docx", "xlsx", "pptx"}

	for _, filetype := range filetypes {
		fmt.Printf("=== Searching for .%s files on %s ===\n", filetype, domain)

		q := fmt.Sprintf("site:%s filetype:%s", domain, filetype)
		search := fmt.Sprintf("http://www.bing.com/search?q=%s", url.QueryEscape(q))

		res, err := http.Get(search)
		if err != nil {
			log.Println(err)
			continue
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		res.Body.Close()
		if err != nil {
			log.Println(err)
			continue
		}

		s := "li.b_algo h2"
		doc.Find(s).Each(func(i int, sel *goquery.Selection) {
			handler(domain, i, sel)
		})
	}
}
