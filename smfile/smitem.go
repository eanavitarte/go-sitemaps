package smfile

import (
	"fmt"

	"github.com/eanavitarte/go-sitemaps/smutils"
	"github.com/eanavitarte/go-sitemaps/xml"
)

/*
Creates a new single sitemap file.
*/
func NewSitemap() *sitemapFile {
	sitemapFile := new(sitemapFile)
	return sitemapFile
}

type sitemapFile struct {
	urlSet xml.XMLURLSet
}

/*
Add individual urls.
*/
func (s *sitemapFile) Add(url xml.XMLURL) {
	s.urlSet.XMLURLs = append(s.urlSet.XMLURLs, url)
}

/*
Add an existing slice-pointer of urls.
*/
func (s *sitemapFile) Clone(urlList []xml.XMLURL) {
	s.urlSet.XMLURLs = urlList
}

func (s *sitemapFile) Save(compression CompressOption, filePath string) error {

	if len(s.urlSet.XMLURLs) > xml.MAXURLSETSIZE {
		return fmt.Errorf("more than 50.000 urls")
	}

	sitemapXML, err := s.urlSet.RenderXML()

	if err != nil {
		return err
	}

	switch compression {

	case NoCompress:
		err = smutils.WriteToFile(filePath, sitemapXML)
		if err != nil {
			return err
		}

	case Gzip:
		err = smutils.WriteToGzip(filePath+".gz", sitemapXML)
		if err != nil {
			return err
		}

	case Both:
		err = smutils.WriteToFile(filePath, sitemapXML)
		if err != nil {
			return err
		}

		err = smutils.WriteToGzip(filePath+".gz", sitemapXML)
		if err != nil {
			return err
		}

	default:
		return fmt.Errorf("invalid compression option")
	}

	return nil
}
