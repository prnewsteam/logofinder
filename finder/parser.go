package finder

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type ParserSelection goquery.Selection

// Implemented from equivalent for InnerHtml (see .Html() in goquery/property.go)
func (s ParserSelection) OuterHtml() (ret string, e error) {
	// Since there is no .outerHtml, the HTML content must be re-created from
	// the node using html.Render
	var buf bytes.Buffer
	if len(s.Nodes) > 0 {
		c := s.Nodes[0]
		e = html.Render(&buf, c)
		if e != nil {
			return
		}
		ret = buf.String()
	}
	return
}
