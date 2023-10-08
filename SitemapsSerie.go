package sitemaps

import (
	"github.com/eanavitarte/go-sitemaps/xml"
)

type SitemapsSerieOpts struct {
	Compress CompressOption
	Label    string
	DirPath  string
	IsMobile bool
}

func (Serie *SitemapsSerie) Configure(opts SitemapsSerieOpts) {
	Serie.label = opts.Label
	Serie.serie = 1
	Serie.isMobile = opts.IsMobile
	Serie.compress = opts.Compress
	if opts.DirPath == "" {
		Serie.dirPath = "."
	} else {
		Serie.dirPath = opts.DirPath
	}
}

type SitemapsSerie struct {
	label        string
	dirPath      string
	styleSheet   string
	serie        int
	xmlUrlList   []xml.XMLURL
	isMobile     bool
	sitemapsList []string
	pipeline     seriePipes
	compress     CompressOption
}

// add a xml-url to the processing group
func (Serie *SitemapsSerie) Add(xmlUrl xml.XMLURL) {
	Serie.xmlUrlList = append(Serie.xmlUrlList, xmlUrl)
}

func (Serie *SitemapsSerie) Process() {
	Serie.pipeline.init()

	Serie.listenSave()
	Serie.listenArchive()
	Serie.batch()

	Serie.pipeline.wait()
}

// returns the url of already generated sitemaps
func (Serie *SitemapsSerie) List() []string {
	return Serie.sitemapsList
}

func (Serie *SitemapsSerie) listenSave() {
	go func() {
		for sitemapFile := range Serie.pipeline.files {
			sitemapFile.save(Serie.compress)
		}
		Serie.pipeline.doneSave <- true
	}()
}

func (Serie *SitemapsSerie) listenArchive() {
	go func() {
		for sitemapUrl := range Serie.pipeline.names {
			Serie.sitemapsList = append(Serie.sitemapsList, sitemapUrl)
		}
		Serie.pipeline.doneArchive <- true
	}()
}

func (Serie *SitemapsSerie) batch() {
	go func() {
		defer close(Serie.pipeline.files)
		defer close(Serie.pipeline.names)

		xmlUrlBatcher := new(xmlUrlBatcher)
		xmlUrlBatcher.tag(Serie.dirPath, Serie.label, Serie.serie)

		for _, xmlUrl := range Serie.xmlUrlList {

			xmlUrlBatcher.stack(xmlUrl)

			if xmlUrlBatcher.full() {
				xmlUrlBatcher.pump(Serie.pipeline)
			}
		}

		xmlUrlBatcher.pump(Serie.pipeline)

		Serie.pipeline.doneBatch <- true
	}()
}

type seriePipes struct {
	files       chan sitemapFile
	names       chan string
	doneSave    chan bool
	doneArchive chan bool
	doneBatch   chan bool
}

func (P *seriePipes) init() {
	P.files = make(chan sitemapFile)
	P.names = make(chan string)

	P.doneSave = make(chan bool, 1)
	P.doneArchive = make(chan bool, 1)
	P.doneBatch = make(chan bool, 1)
}

func (P *seriePipes) wait() {
	<-P.doneSave
	<-P.doneArchive
	<-P.doneBatch
}
