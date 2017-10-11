package stint

import (
	"bytes"
	"encoding/xml"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"

	log "github.com/Sirupsen/logrus"
)

type Time struct {
	time.Time
}

func (t *Time) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var value string

	d.DecodeElement(&value, &start)
	parse, err := time.Parse(time.RFC1123Z, value)
	if err != nil {
		return err
	}

	*t = Time{parse}
	return nil
}

type Item struct {
	Title         string `xml:"title"`
	Category      string `xml:"category"`
	Link          string `xml:"link"`
	PubDate       Time   `xml:"pubDate"`
	ContentLength string `xml:"contentLength"`
	InfoHash      string `xml:"infoHash"`
	MagnetURI     string `xml:"magnetURI"`
	Seeds         uint   `xml:"seeds"`
	Peers         uint   `xml:"peers"`
	Verified      uint   `xml:"verified"`
	FileName      string `xml:"fileName"`
}

type Feed struct {
	Title         string `xml:"title"`
	Link          string `xml:"link"`
	Description   string `xml:"description"`
	LastBuildDate Time   `xml:"lastBuildDate"`

	Items []Item `xml:"item"`
}

func Do(path string, level string) error {
	switch level {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	default:
		log.SetLevel(log.WarnLevel)

	}
	log.Debugf("Pulling configuration")
	c, err := conf(path)
	if err != nil {
		log.WithError(err).Errorf("Error while reading the configuration file")
		return err
	}

	log.Debugf("Fetching RSS feed from url: %s", c.URL)
	r, err := http.Get(c.URL)
	if err != nil {
		log.WithError(err).Errorf("Error while fetching RSS feed")
		return err
	}
	defer r.Body.Close()

	rss := struct {
		XMLName xml.Name `xml:"rss"`
		Channel Feed     `xml:"channel"`
	}{}

	log.Debugf("Decoding RSS feed")
	decoder := xml.NewDecoder(r.Body)
	err = decoder.Decode(&rss)
	if err != nil {
		log.WithError(err).Errorf("Error while decoding RSS feed")
		return err
	}

	filtered := []Item{}

	// fuck this loop as much as my sleepy self can.
	// optimisation only comes when it gets working.
	log.Infof("Found total items: %d", len(rss.Channel.Items))
	for _, item := range rss.Channel.Items {
		log.WithFields(log.Fields{
			"link": item.Link,
			"time": item.PubDate,
			"size": item.ContentLength,
		}).Debugf("Parsing item: %s", item.Title)
		for t, q := range c.Shows {
			if matchTitle(item.Title, t) && matchQuality(item.Title, q) {
				log.Infof("Matched item: %s", item.Title)
				filtered = append(filtered, item)
			}
		}
	}

	if len(filtered) == 0 {
		log.Warnf("No items found based on provided filters")
		return nil
	}

	log.Debugf("Preparing CLI command arguments: %s", c.Args)
	tmpl, err := template.New("cmd").Parse(c.Args)
	if err != nil {
		log.WithError(err).Errorf("Error while preparing CLI command arguments: %s", c.Args)
		return err
	}

	for _, f := range filtered {
		log.WithFields(log.Fields{
			"title": f.Title,
		}).Debugf("Preparing CLI command: %s", c.Cmd)
		buf := new(bytes.Buffer)
		if err := tmpl.Execute(buf, f); err != nil {
			log.WithError(err).Errorf("Error while preparing CLI command: %s %s", c.Cmd, c.Args)
			return err
		}

		log.Debugf("Executing: %s %s", c.Cmd, buf.String())
		cmd := exec.Command(c.Cmd, strings.Split(buf.String(), " ")...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.WithError(err).Errorf("Error while running CLI command")
			return err
		}
	}

	return nil
}
