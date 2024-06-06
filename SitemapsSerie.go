package sitemaps

import (
	"fmt"
	"log"
	"time"

	"github.com/eanavitarte/go-sitemaps/xml"
)

type SitemapsSerieOpts struct {
	Compress  CompressOption
	Label     string
	DirPath   string
	PublicUrl string
	IsMobile  bool
	DoEmpty   bool
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
	Serie.doEmpty = opts.DoEmpty
	Serie.publicUrl = opts.PublicUrl
}

type SitemapsSerie struct {
	label   string
	dirPath string
	// styleSheet   string
	serie        int
	xmlUrlList   []xml.XMLURL
	isMobile     bool
	sitemapsList []string
	pipeline     seriePipes
	publicUrl    string

	compress CompressOption
	doEmpty  bool
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

func (Serie *SitemapsSerie) Index() {
	index := new(xml.XMLIndex)

	if len(Serie.sitemapsList) <= 0 {
		return
	}

	for _, smFileName := range Serie.sitemapsList {
		lastModified := time.Now()
		index.XMLSitemaps = append(index.XMLSitemaps, xml.XMLSitemap{Loc: Serie.publicUrl + "/" + smFileName, LastMod: &lastModified})
	}

	indexXmlContent, err := index.RenderXML()
	if err != nil {
		log.Println("[Error] rendering Index File:", err)
		return
	}

	indexFile := &sitemapFile{
		filePath:   fmt.Sprintf("%s/%s-index.xml", Serie.dirPath, Serie.label),
		xmlContent: indexXmlContent,
	}

	indexFile.save(Serie.compress)
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

		if len(Serie.xmlUrlList) > 0 || Serie.doEmpty {
			for _, xmlUrl := range Serie.xmlUrlList {

				xmlUrlBatcher.stack(xmlUrl)

				if xmlUrlBatcher.full() {
					xmlUrlBatcher.pump(Serie.pipeline)
				}
			}
			xmlUrlBatcher.pump(Serie.pipeline)
		}

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
