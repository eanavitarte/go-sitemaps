package sitemaps

import (
	"compress/gzip"
	"log"
	"os"
)

type sitemapFile struct {
	name       string
	filePath   string
	xmlContent []byte
}

func (F *sitemapFile) save(compression CompressOption) {
	// log.Println(F.name)

	switch compression {
	case NoCompress:
		F.writeToFile(F.filePath)
	case Gzip:
		F.writeToGzip(F.filePath + ".gz")
	case Both:
		F.writeToFile(F.filePath)
		F.writeToGzip(F.filePath + ".gz")
	default:
		log.Println("[Error]: Invalid compression option")
	}
}

func (F *sitemapFile) writeToFile(filePath string) {
	f, err := os.Create(filePath)
	if err != nil {
		log.Println("[Error]:", err)
		return
	}
	defer f.Close()

	_, err = f.Write(F.xmlContent)
	if err != nil {
		log.Println("[Error]:", err)
	}
}

func (F *sitemapFile) writeToGzip(filePath string) {
	f, err := os.Create(filePath)
	if err != nil {
		log.Println("[Error]:", err)
		return
	}
	defer f.Close()

	zip, err := gzip.NewWriterLevel(f, gzip.BestCompression)
	if err != nil {
		log.Println("[Error]:", err)
		return
	}
	defer zip.Close()

	_, err = zip.Write(F.xmlContent)
	if err != nil {
		log.Println("[Error]:", err)
	}
}
