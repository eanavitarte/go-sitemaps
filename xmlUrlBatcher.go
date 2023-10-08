package sitemaps

import (
	"fmt"
	"log"

	"github.com/eanavitarte/go-sitemaps/xml"
)

type xmlUrlBatcher struct {
	label    string
	dirPath  string
	filePath string
	serie    int
	stock    []xml.XMLURL
	errors   error
}

func (B *xmlUrlBatcher) tag(dirPath, label string, serie int) {
	B.label = label
	B.dirPath = dirPath
	B.filePath = dirPath + "/" + label
	B.serie = serie
}

func (B *xmlUrlBatcher) stack(xmlUrl xml.XMLURL) {
	B.stock = append(B.stock, xmlUrl)
}

func (B *xmlUrlBatcher) full() (isFull bool) {
	return len(B.stock) >= xml.MAXURLSETSIZE
}

func (B *xmlUrlBatcher) pump(pipeline seriePipes) {
	pack := B.pack()

	if B.errors != nil {
		log.Println("[Error]:", B.errors)
		// TODO: process error
	} else {
		pipeline.files <- pack
		pipeline.names <- pack.name
		B.stock = nil
		B.serie++
	}
}

func (B *xmlUrlBatcher) pack() (pack sitemapFile) {
	pack = sitemapFile{
		name:       B.stamp(),
		filePath:   B.sign(),
		xmlContent: B.process()}

	return
}

func (B *xmlUrlBatcher) process() (sitemapXml []byte) {
	xmlUrlSet := new(xml.XMLURLSet)

	xmlUrlSet.XMLURLs = B.stock

	sitemapXml, err := xmlUrlSet.RenderXML()

	if err != nil {
		B.errors = err
	}

	return
}

func (B *xmlUrlBatcher) stamp() (stamp string) {
	stamp = fmt.Sprintf("%s-%d.xml", B.label, B.serie)
	return
}

func (B *xmlUrlBatcher) sign() (sign string) {
	sign = fmt.Sprintf("%s-%d.xml", B.filePath, B.serie)
	return
}
