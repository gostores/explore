package feeds

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"time"
)

const ans = "http://www.w3.org/2005/Atom"
const ins = "http://www.itunes.com/dtds/podcast-1.0.dtd"

// private wrapper around the PodcastFeed which gives us the <rss>..</rss> xml
type podcastFeedXml struct {
	XMLName     xml.Name `xml:"rss"`
	Version     string   `xml:"version,attr"`
	XmlnsAtom   string   `xml:"xmlns:atom,attr"`
	XmlnsItunes string   `xml:"xmlns:itunes,attr"`
	Channel     *PodcastFeed
}

// Itunes
type IAuthor struct {
	XMLName xml.Name `xml:"itunes:owner"`
	Name    string   `xml:"itunes:name"`
	Email   string   `xml:"itunes:email"`
}

type ICategory struct {
	XMLName xml.Name `xml:"itunes:category"`
	Text    string   `xml:"text,attr"`
}

type ISummary struct {
	XMLName xml.Name `xml:"itunes:summary"`
	Text    string   `xml:",cdata"`
}

type IImage struct {
	XMLName xml.Name `xml:"itunes:image"`
	Href    string   `xml:"href,attr"`
}

// Podcast
type PodcastImage struct {
	XMLName     xml.Name `xml:"image"`
	Url         string   `xml:"url"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description,omitempty"`
	Width       int      `xml:"width,omitempty"`
	Height      int      `xml:"height,omitempty"`
}

type PodcastAtomLink struct {
	XMLName xml.Name `xml:"atom:link"`
	Href    string   `xml:"href,attr"`
	Rel     string   `xml:"rel,attr"`
	Type    string   `xml:"type,attr"`
}

type PodcastTextInput struct {
	XMLName     xml.Name `xml:"textInput"`
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
	Name        string   `xml:"name"`
	Link        string   `xml:"link"`
}

type PodcastEnclosure struct {
	XMLName xml.Name `xml:"enclosure"`
	Url     string   `xml:"url,attr"`
	Length  string   `xml:"length,attr"`
	Type    string   `xml:"type,attr"`
}

type PodcastItem struct {
	XMLName            xml.Name `xml:"item"`
	Title              string   `xml:"title"`       // required
	Link               string   `xml:"link"`        // required
	Description        string   `xml:"description"` // required
	Category           string   `xml:"category,omitempty"`
	Comments           string   `xml:"comments,omitempty"`
	Guid               string   `xml:"guid,omitempty"`    // Id used
	PubDate            string   `xml:"pubDate,omitempty"` // created or updated
	Source             string   `xml:"source,omitempty"`
	Author             string   `xml:"author,omitempty"`
	Enclosure          *PodcastEnclosure
	IAuthor            string `xml:"itunes:author,omitempty"`
	ISubtitle          string `xml:"itunes:subtitle,omitempty"`
	IDuration          string `xml:"itunes:duration,omitempty"`
	IExplicit          string `xml:"itunes:explicit,omitempty"`
	IIsClosedCaptioned string `xml:"itunes:isClosedCaptioned,omitempty"`
	IOrder             string `xml:"itunes:order,omitempty"`
	ISummary           *ISummary
	IImage             *IImage
}

type PodcastFeed struct {
	XMLName        xml.Name `xml:"channel"`
	Title          string   `xml:"title"`       // required
	Link           string   `xml:"link"`        // required
	Description    string   `xml:"description"` // required
	Category       string   `xml:"category,omitempty"`
	Cloud          string   `xml:"cloud,omitempty"`
	Copyright      string   `xml:"copyright,omitempty"`
	Docs           string   `xml:"docs,omitempty"`
	Generator      string   `xml:"generator,omitempty"`
	Language       string   `xml:"language,omitempty"`
	LastBuildDate  string   `xml:"lastBuildDate,omitempty"`  // updated used
	ManagingEditor string   `xml:"managingEditor,omitempty"` // Author used
	PubDate        string   `xml:"pubDate,omitempty"`        // created or updated
	Rating         string   `xml:"rating,omitempty"`
	WebMaster      string   `xml:"webMaster,omitempty"`
	Ttl            int      `xml:"ttl,omitempty"`
	SkipHours      string   `xml:"skipHours,omitempty"`
	SkipDays       string   `xml:"skipDays,omitempty"`
	Image          *PodcastImage
	TextInput      *PodcastTextInput
	AtomLink       *PodcastAtomLink
	Items          []*PodcastItem
	IAuthor        string `xml:"itunes:author,omitempty"`
	ISubtitle      string `xml:"itunes:subtitle,omitempty"`
	IBlock         string `xml:"itunes:block,omitempty"`
	IDuration      string `xml:"itunes:duration,omitempty"`
	IExplicit      string `xml:"itunes:explicit,omitempty"`
	IComplete      string `xml:"itunes:complete,omitempty"`
	INewFeedURL    string `xml:"itunes:new-feed-url,omitempty"`
	ISummary       *ISummary
	IImage         *IImage
	IOwner         *IAuthor
	ICategory      *ICategory
}

type Podcast struct {
	*Feed
}

func checkEnclosureType(t string) string {
	switch t {
	case "m4a":
		return "audio/x-m4a"
	case "m4v":
		return "video/x-m4v"
	case "mp4":
		return "video/mp4"
	case "mp3":
		return "audio/mpeg"
	case "mov":
		return "video/quicktime"
	case "pdf":
		return "application/pdf"
	case "epub":
		return "document/x-epub"
	}
	return "application/octet-stream"
}

// create a new RssItem with a generic Item struct's data
func newPodcastItem(i *Item) *PodcastItem {
	item := &PodcastItem{
		Title:       i.Title,
		Description: i.Description,
		Link:        i.Link.Href,
		Guid:        i.Id,
		PubDate:     anyTimeFormat(time.RFC1123Z, i.Created, i.Updated),
		IAuthor:     i.Author.Name,
		ISubtitle:   i.Itunes.Subtitle,
		ISummary:    &ISummary{Text: i.Description},
		IImage:      &IImage{Href: i.Itunes.Image},
	}
	if i.Source != nil {
		item.Source = i.Source.Href
	}

	if i.Itunes.AudioSize > 0 || i.Itunes.AudioType != "" {
		item.Enclosure = &PodcastEnclosure{Url: i.Itunes.AudioHref, Type: checkEnclosureType(i.Itunes.AudioType), Length: strconv.FormatInt(i.Itunes.AudioSize, 10)}
	}
	if i.Author != nil {
		item.Author = i.Author.Name
	}

	return item
}

// create a new PodcastFeed with a generic Feed struct's data
func (r *Podcast) PodcastFeed() *PodcastFeed {
	pub := anyTimeFormat(time.RFC1123Z, r.Created, r.Updated)
	build := anyTimeFormat(time.RFC1123Z, r.Updated)
	author := ""
	if r.Author != nil {
		author = r.Author.Email
		if len(r.Author.Name) > 0 {
			author = fmt.Sprintf("%s (%s)", r.Author.Email, r.Author.Name)
		}
	}

	channel := &PodcastFeed{
		Title:          r.Title,
		Link:           r.Link.Href,
		Description:    r.Description,
		ManagingEditor: author,
		PubDate:        pub,
		LastBuildDate:  build,
		Copyright:      r.Copyright,
		Language:       r.Itunes.Language,
		Image:          &PodcastImage{Title: r.Title, Link: r.Link.Href, Url: r.Itunes.Logo},
		AtomLink:       &PodcastAtomLink{Href: r.Link.Href, Rel: r.Link.Rel, Type: r.Link.Type},
		IAuthor:        r.Itunes.Author,
		ISubtitle:      r.Subtitle,
		IBlock:         r.Itunes.Block,
		IDuration:      r.Itunes.Duration,
		IExplicit:      r.Itunes.Explicit,
		IComplete:      r.Itunes.Complete,
		INewFeedURL:    r.Itunes.NewFeedURL,
		ISummary:       &ISummary{Text: r.Description},
		IImage:         &IImage{Href: r.Itunes.Logo},
		IOwner:         &IAuthor{Name: r.Itunes.Author, Email: r.Itunes.Email},
		ICategory:      &ICategory{Text: r.Itunes.Category},
	}
	for _, i := range r.Items {
		channel.Items = append(channel.Items, newPodcastItem(i))
	}
	return channel
}

// return an XML-Ready object for an Podcast object
func (r *Podcast) FeedXml() interface{} {
	return r.PodcastFeed().FeedXml()
}

// return an XML-ready object for an PodcastFeed object
func (r *PodcastFeed) FeedXml() interface{} {
	return &podcastFeedXml{
		Version:     "2.0",
		XmlnsAtom:   ans,
		XmlnsItunes: ins,
		Channel:     r,
	}
}
