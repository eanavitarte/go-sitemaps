// Generates sitemaps and index files based on the sitemaps.org protocol.
// facilitates the creation of sitemaps for large amounts of urls.
// For a full guide visit https://github.com/StudioSol/Sitemap
package sitemaps

import (
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/eanavitarte/go-sitemaps/xml"
)

// Creates a new group of sitemaps that used a common name.
// If the sitemap exceed the limit of 50k urls, new sitemaps will have a numeric suffix to the name. Example:
// - blog_1.xml.gz
// - blog_2.xml.gz
// func NewSitemapGroup(name string, isMobile bool) *SitemapsSerie {
// 	s := new(SitemapsSerie)
// 	s.Configure(name, isMobile)
// 	return s
// }

// Creates a new group of sitemaps indice that used a common name.
// If the sitemap exceed the limit of 50k urls, new sitemaps will have a numeric suffix to the name. Example:
// - blog_1.xml.gz
// - blog_2.xml.gz
func NewIndexGroup(folder string, name string) (*IndexGroup, error) {
	s := new(IndexGroup)
	err := s.Configure(name, folder)
	if err != nil {
		return s, err
	}
	go s.Initialize()
	return s, nil
}

// Search all the xml.gz sitemaps_dir directory, uses the modified date of the file as lastModified
// path_index is included for the function does not include the url of the index in your own content, if it is present in the same directory.
func CreateIndexByScanDir(targetDir string, indexFileName string, public_url string) (index xml.Index) {
	index = xml.Index{Sitemaps: []xml.Sitemap{}}

	fs, err := ioutil.ReadDir(targetDir)
	if err != nil {
		return
	}

	for _, f := range fs {
		if strings.HasSuffix(f.Name(), ".xml.gz") && !strings.HasSuffix(indexFileName, f.Name()) {
			lastModified := f.ModTime()
			index.Sitemaps = append(index.Sitemaps, xml.Sitemap{Loc: public_url + f.Name(), LastMod: &lastModified})
		}
	}
	return
}

// Returns an index sitemap starting from a slice of urls
func CreateIndexBySlice(smFileNames []string, public_url string) (index xml.Index) {
	index = xml.Index{Sitemaps: []xml.Sitemap{}}
	if len(smFileNames) > 0 {
		for _, fileName := range smFileNames {
			lastModified := time.Now()
			index.Sitemaps = append(index.Sitemaps, xml.Sitemap{Loc: public_url + fileName, LastMod: &lastModified})
		}
	}
	return
}

// Creates and gzip the xml index
func CreateSitemapIndex(indexFilePath string, index xml.Index) (err error) {
	//create xml
	indexXml, err := xml.CreateSitemapIndexXml(index)
	if err != nil {
		return err
	}
	err = xml.SaveXml(indexXml, indexFilePath)
	log.Printf("Sitemap Index created on %s", indexFilePath)
	return err
}
