package xml

import (
	"errors"
)

const (
	XMLNS         = "http://www.sitemaps.org/schemas/sitemap/0.9"
	XMLNSMOBILE   = "http://www.google.com/schemas/sitemap-mobile/1.0"
	PREAMBLE      = "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"
	MAXURLSETSIZE = 50000
	MAXFILESIZE   = 50 * 1024 * 1024 // 50mb
)

var (
	ErrMaxUrlSetSize = errors.New("exceeded maximum number of URLs allowed in sitemap")
	ErrMaxFileSize   = errors.New("exceeded maximum file size of a sitemap (10mb)")
	ISMOBILE         = new(struct{})
)

type ChangeFreq string

const (
	Always  ChangeFreq = "always"
	Hourly  ChangeFreq = "hourly"
	Daily   ChangeFreq = "daily"
	Weekly  ChangeFreq = "weekly"
	Monthly ChangeFreq = "monthly"
	Yearly  ChangeFreq = "yearly"
	Never   ChangeFreq = "never"
)
