package finder

import (
	"context"
	"errors"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/css"
	"github.com/chromedp/chromedp"
)

func FindLogo(domain string) (*Logo, error) {

	if logo, err := findInBody(domain); err == nil {
		return logo, nil
	}

	if logo, err := findClearbit(domain); err == nil {
		return logo, nil
	}

	return nil, errors.New("Unable to parse logo from domain: " + domain)
}

func findInBody(domain string) (*Logo, error) {
	doc, css, err := newDocumentFromUrl("http://" + domain)
	if err != nil {
		return nil, err
	}

	if url := findImg(doc); url != "" {
		log.Printf("findImg match: %s", url)
		return NewLogoFromUrl(url, domain)
	}

	if raw := findSvg(doc); raw != "" {
		log.Printf("findSvg match")
		return NewLogoFromRaw(domain, raw, ".svg")
	}

	if url := findCss(css); url != "" {
		log.Printf("findCss match: %s", url)
		return NewLogoFromUrl(url, domain)
	}

	return nil, errors.New("Unable to find any image url in response body")
}

func findImg(doc *goquery.Document) string {
	var exists bool
	var attr string

	el := doc.Find("img[src*='logo'], img[data-src*='logo'], [id*='logo'], [class*='logo']").First()

	attr, exists = el.Attr("src")
	if exists {
		return attr
	}

	attr, exists = el.Attr("data-src")
	if exists {
		return attr
	}

	el = el.Find("img").First()

	attr, exists = el.Attr("data-src")
	if exists {
		return attr
	}

	attr, exists = el.Attr("src")
	if exists {
		return attr
	}

	return ""
}

func findSvg(doc *goquery.Document) string {
	el := doc.Find("[id*='logo'], [class*='logo']").First().Find("svg").First()
	if el.Find("use").Length() != 0 {
		return ""
	}
	html, _ := ParserSelection(*el).OuterHtml()
	return html
}

func findCss(css *[]*css.ComputedStyleProperty) string {

	style := *css
	for i := range style {
		if style[i].Name == "background" || style[i].Name == "background-image" {
			re := regexp.MustCompile(`url\((.+)\)`)
			match := re.FindStringSubmatch(style[i].Value)
			if len(match) >= 2 {
				return strings.Trim(match[1], "\"")
			}
		}
	}

	return ""
}

func findClearbit(domain string) (*Logo, error) {
	return NewLogoFromUrlWithExtension("https://logo.clearbit.com/"+domain+"?size=800", domain, ".png")
}

func newDocumentFromUrl(url string) (*goquery.Document, *[]*css.ComputedStyleProperty, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	html := "<html></html>"
	style := []*css.ComputedStyleProperty{}
	if err := chromedp.Run(
		ctx,
		chromedp.Tasks{
			chromedp.Navigate(url),
			chromedp.OuterHTML("html", &html),
			chromedp.ComputedStyle("img[src*='logo'], img[data-src*='logo'], [id*='logo'], [class*='logo']", &style),
		},
	); err != nil {
		log.Printf("Error fetching body: %s", err)
		return nil, nil, err
	}

	log.Printf("Received body: %s", url)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, nil, err
	}

	log.Printf("Parsed body: %s", url)

	return doc, &style, nil
}
