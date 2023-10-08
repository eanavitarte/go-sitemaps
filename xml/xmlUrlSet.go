package xml

import (
	"encoding/xml"
	"time"
)

type XMLURLSet struct {
	XMLName     xml.Name `xml:"urlset"`
	XMLNS       string   `xml:"xmlns,attr"`
	XMLNSMOBILE string   `xml:"xmlns:mobile,attr,omitempty"`
	XMLURLs     []XMLURL `xml:"url"`

	IsMobile bool `xml:"-"`
}

type XMLURL struct {
	Loc        string     `xml:"loc"`
	LastMod    *time.Time `xml:"lastmod,omitempty"`
	ChangeFreq ChangeFreq `xml:"changefreq,omitempty"`
	Priority   float64    `xml:"priority,omitempty"`
	Mobile     *struct{}  `xml:"mobile:mobile,omitempty"`
}

func (Set *XMLURLSet) RenderXML() (rendered []byte, errors error) {

	if len(Set.XMLURLs) > MAXURLSETSIZE {
		errors = ErrMaxUrlSetSize
		return
	}

	Set.XMLNS = XMLNS

	if Set.IsMobile {
		Set.XMLNSMOBILE = XMLNSMOBILE
	}

	var urlSetXMLContent []byte
	urlSetXMLContent, err := xml.Marshal(Set)

	if err != nil {
		errors = err
		return
	}

	sitemapXMLContent := append(
		[]byte(PREAMBLE), urlSetXMLContent...)

	if len(sitemapXMLContent) > MAXFILESIZE {
		errors = ErrMaxFileSize
		return
	}

	rendered = sitemapXMLContent
	return
}
