package smserie

import (
	"fmt"
	"time"

	"github.com/eanavitarte/go-sitemaps/smfile"
	"github.com/eanavitarte/go-sitemaps/smindex"
	"github.com/eanavitarte/go-sitemaps/xml"
)

type Options struct {
	Directory     string // "/path/to"
	Label         string
	Compression   smfile.CompressOption
	StartingSerie int
	Domain        string // "http://example.com"
}

/*
Creates a new single sitemap serie.
*/
func NewSitemapSerie(opts Options) (*sitemapSerie, error) {
	sitemapFile := new(sitemapSerie)

	if opts.Directory == "" {
		return nil, fmt.Errorf("you need to set a directory")
	}

	sitemapFile.dir = opts.Directory

	if opts.Label == "" {
		sitemapFile.label = "sitemap"
	} else {
		sitemapFile.label = opts.Label
	}

	if opts.Compression != 0 {
		sitemapFile.compress = opts.Compression
	}

	if opts.Domain != "" {
		sitemapFile.domain = opts.Domain
	}

	if opts.StartingSerie == 0 {
		sitemapFile.cycle = 1
	} else {
		sitemapFile.cycle = opts.StartingSerie
	}

	return sitemapFile, nil
}

type sitemapSerie struct {
	dir      string
	label    string
	compress smfile.CompressOption
	cycle    int
	domain   string
}

func (s *sitemapSerie) Process(urlList []xml.XMLURL) error {

	if len(urlList) > xml.MAXURLSETSIZE {
		return fmt.Errorf("more than 50.000 urls")
	}

	sitemapFile := smfile.NewSitemap()

	sitemapFile.Clone(urlList)

	filePath := fmt.Sprintf("%s/%s-%d.xml", s.dir, s.label, s.cycle)

	err := sitemapFile.Save(s.compress, filePath)

	if err != nil {
		return err
	}

	s.cycle++

	return nil
}

func (s *sitemapSerie) IndexFromCount(seriesDone int) error {

	if s.domain == "" {
		return fmt.Errorf("you need a domain for indexing based on countings")
	}

	sitemapIndex := smindex.NewSitemapIndex()

	for v := range seriesDone {

		index := v + 1
		loc := fmt.Sprintf("%s/%s-%d.xml", s.domain, s.label, index)

		now := time.Now()
		url := xml.XMLSitemap{
			Loc:     loc,
			LastMod: &now,
		}

		sitemapIndex.Add(url)
	}

	filePath := fmt.Sprintf("%s/%s-index.xml", s.dir, s.label)

	err := sitemapIndex.Save(s.compress, filePath)

	if err != nil {
		return err
	}

	return nil
}
