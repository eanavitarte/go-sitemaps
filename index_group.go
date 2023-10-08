package sitemaps

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/eanavitarte/go-sitemaps/xml"
)

type IndexGroup struct {
	name            string
	folder          string
	group_count     int
	sitemaps        []xml.Sitemap
	sitemap_channel chan xml.Sitemap
	done            chan bool
}

// Add a sitemap.Sitemap to the group
func (s *IndexGroup) Add(entry xml.Sitemap) {
	s.sitemap_channel <- entry
}

// Clean Urls not yet added to the group
func (s *IndexGroup) Clear() {
	s.sitemaps = []xml.Sitemap{}
}

// Returns one sitemap.Index of Urls not yet added to the group
func (s *IndexGroup) getSitemapSet() xml.Index {
	return xml.Index{Sitemaps: s.sitemaps}
}

func (s *IndexGroup) getSitemapName() string {
	return s.name + "_" + strconv.Itoa(s.group_count) + ".xml.gz"
}

// Saves the sitemap from the sitemap.URLSet
func (s *IndexGroup) Create(index xml.Index) {
	var path string
	var remnant []xml.Sitemap
	xmlF, err := xml.CreateSitemapIndexXml(index)
	if err == xml.ErrMaxFileSize {
		//splits into two sitemaps recursively
		newlimit := xml.MAXURLSETSIZE / 2
		s.Create(xml.Index{Sitemaps: index.Sitemaps[newlimit:]})
		s.Create(xml.Index{Sitemaps: index.Sitemaps[:newlimit]})
		return
	} else if err == xml.ErrMaxUrlSetSize {
		remnant = index.Sitemaps[xml.MAXURLSETSIZE:]
		index.Sitemaps = index.Sitemaps[:xml.MAXURLSETSIZE]
		xmlF, err = xml.CreateSitemapIndexXml(index)
	}

	if err != nil {
		log.Fatal("File not saved:", err)
	}

	sitemap_name := s.getSitemapName()
	path = filepath.Join(s.folder, sitemap_name)

	err = xml.SaveXml(xmlF, path)
	if err != nil {
		log.Fatal("File not saved:", err)
	}
	s.group_count++
	s.Clear()
	//append remnant urls if exists
	if len(remnant) > 0 {
		s.sitemaps = append(s.sitemaps, remnant...)
	}
	log.Printf("Sitemap created on %s", path)

}

// Starts to run the given list of Sitemap Groups concurrently.
func CloseIndexGroups(groups ...*IndexGroup) (done <-chan bool) {
	var wg sync.WaitGroup
	wg.Add(len(groups))

	ch := make(chan bool, 1)
	for _, group := range groups {
		go func(g *IndexGroup) {
			<-g.Close()
			wg.Done()
		}(group)
	}
	go func() {
		wg.Wait()
		ch <- true
	}()
	return ch
}

// Mandatory operation, handle the rest of the url that has not been added to any sitemap and add.
// Furthermore performs cleaning of variables and closes the channel group
func (s *IndexGroup) Close() <-chan bool {
	var closeDone = make(chan bool, 1)
	close(s.sitemap_channel)

	go func() {
		<-s.done
		closeDone <- true
	}()

	return closeDone
}

// Initialize channel
func (s *IndexGroup) Initialize() {
	s.done = make(chan bool, 1)
	s.sitemap_channel = make(chan xml.Sitemap)

	for entry := range s.sitemap_channel {
		s.sitemaps = append(s.sitemaps, entry)
		if len(s.sitemaps) == xml.MAXURLSETSIZE {
			s.Create(s.getSitemapSet())
		}
	}

	//remnant urls
	s.Create(s.getSitemapSet())
	s.Clear()

	s.done <- true
}

// Configure name and folder of group
func (s *IndexGroup) Configure(name string, folder string) error {
	s.name = strings.Replace(name, ".xml.gz", "", 1)
	s.group_count = 1
	s.folder = folder
	_, err := ioutil.ReadDir(folder)
	if err != nil {
		err = os.MkdirAll(folder, 0655)
	}
	return err
}
