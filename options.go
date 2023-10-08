package sitemaps

type CompressOption int

const (
	NoCompress CompressOption = iota
	Gzip
	Both
)
