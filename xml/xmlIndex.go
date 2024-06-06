package xml

import (
	"encoding/xml"
	"time"
)

type XMLIndex struct {
	XMLName     xml.Name     `xml:"sitemapindex"`
	XMLNS       string       `xml:"xmlns,attr"`
	XMLSitemaps []XMLSitemap `xml:"sitemap"`
}

type XMLSitemap struct {
	Loc     string     `xml:"loc"`
	LastMod *time.Time `xml:"lastmod,omitempty"`
}

func (index *XMLIndex) RenderXML() (rendered []byte, errors error) {

	if len(index.XMLSitemaps) > MAXURLSETSIZE {
		errors = ErrMaxUrlSetSize
		return
	}

	index.XMLNS = XMLNS

	var smIndexXMLContent []byte
	smIndexXMLContent, err := xml.Marshal(index)
	if err != nil {
		errors = err
		return
	}

	smIndexXMLFile := append(
		[]byte(PREAMBLE), smIndexXMLContent...)

	if len(smIndexXMLFile) > MAXFILESIZE {
		errors = ErrMaxFileSize
		return
	}

	rendered = smIndexXMLFile

	return
}
