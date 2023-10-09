# go-sitemaps [![GoDoc](https://godoc.org/github.com/eanavitarte/go-sitemaps?status.png)](https://godoc.org/github.com/eanavitarte/go-sitemaps)

Generates and parse sitemaps and index files using the sitemaps.org protocol.

### Example:

Creates a new group of sitemaps that used a common name.

~~~ go
serie := new(SitemapsSerie)

serie.Add(url1)

serie.Add(url2) //...

serie.Process()

~~~

You can fine tune their behaviour by adding options:

~~~ go
serie := new(SitemapsSerie)

serie.Configure(SitemapsSerieOpts{
	Compress: Both,
	Label:    "blog",
	DirPath:  "tmp",
})

~~~

Options support compressed and non-compressed version:
- blog-1.xml
- blog-1.xml.gz

### Indexes:

For creating an index you can do:

~~~ go
sitemaps := append(serie.List(), serie2.List()...)

index := sitemap.CreateIndexBySlice(sitemaps, "https://example.com/")

err := sitemap.CreateSitemapIndex("index.xml.gz", index)
if err != nil {
	log.Fatal(err)
}

~~~

If the sitemap exceed the limit of 50k urls, new sitemaps will have a numeric suffix to the name. Example:
- blog-1.xml.gz
- blog-2.xml.gz

### URL Format
Urls follow sitemaps.org format:

~~~ go
url1 := xml.XMLURL{
	Loc:     "https://example.com",
	ChangeFreq: sitemap.Hourly,
    LastMod: &now, // now := time.Now()
    Priority: 0.9,
}
~~~

### Ping Search Engines
It allows to ping search engines:

~~~ go
sitemap.PingSearchEngines("https://example.com/index.xml.gz")
~~~

### TODO
- Allow adding stylesheets to Sitemaps Series
- Add a logger, so you can see what it was processed
- Wrap xml internal package for using just once import statement

### Thanks

Parts of this library were inspired (if not outright copied) from StudioSol's excellent [`sitemap`](https://github.com/StudioSol/sitemap) library.
