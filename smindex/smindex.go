package smindex

import (
	"fmt"

	"github.com/eanavitarte/go-sitemaps/smfile"
	"github.com/eanavitarte/go-sitemaps/smutils"
	"github.com/eanavitarte/go-sitemaps/xml"
)

/*
Creates a new single sitemap file.
*/
func NewSitemapIndex() *sitemapIndex {
	sitemapFile := new(sitemapIndex)
	return sitemapFile
}

type sitemapIndex struct {
	index xml.XMLIndex
}

/*
Add individual urls.
*/
func (s *sitemapIndex) Add(url xml.XMLSitemap) {
	s.index.XMLSitemaps = append(s.index.XMLSitemaps, url)
}

func (s *sitemapIndex) Save(compression smfile.CompressOption, filePath string) error {

	if len(s.index.XMLSitemaps) > xml.MAXURLSETSIZE {
		return fmt.Errorf("more than 50.000 urls")
	}

	sitemapXML, err := s.index.RenderXML()

	if err != nil {
		return err
	}

	switch compression {

	case smfile.NoCompress:
		err = smutils.WriteToFile(filePath, sitemapXML)
		if err != nil {
			return err
		}

	case smfile.Gzip:
		err = smutils.WriteToGzip(filePath+".gz", sitemapXML)
		if err != nil {
			return err
		}

	case smfile.Both:
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
