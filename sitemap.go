// Generates sitemaps and index files based on the sitemaps.org protocol.
// facilitates the creation of sitemaps for large amounts of urls.
// For a full guide visit https://github.com/StudioSol/Sitemap
package sitemaps

import (
	"io/ioutil"
	"strings"

	"github.com/eanavitarte/go-sitemaps/xml"
)

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
