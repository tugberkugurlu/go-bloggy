package main

import "html/template"

type SpeakingActivity struct {
	Title       string
	Activity    string
	ImageURL    string
	City        string
	Country     string
	DisplayDate string
	Extras      []SpeakingActivityExtra

	EmbededHTMLData template.HTML
}

type SpeakingActivityExtra struct {
	Name           string
	Link           string
	IconCSSClasses string
}

var speakingActivities = []SpeakingActivity{
	{
		Title:           "Essentials for Building and Leading Highly Effective Development Teams",
		Activity:        "DevConf",
		ImageURL:        "https://tugberkugurlu.blob.core.windows.net/speaking/2018-09-devconf-krakow/30236914487_dba11cff33_o.jpg",
		City:            "Kraków",
		Country:         "Poland",
		DisplayDate:     "27-28 Sep 2018",
		EmbededHTMLData: template.HTML(`<iframe class="embed-responsive-item" src="https://www.youtube.com/embed/qpfFus69pN8" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>`),
		Extras: []SpeakingActivityExtra{
			{
				Name:           "Slides",
				Link:           "https://speakerdeck.com/tourismgeek/essentials-for-building-and-leading-highly-effective-development-teams",
				IconCSSClasses: "area chart",
			},
			{
				Name:           "Conference Profile",
				Link:           "https://web.archive.org/web/20181108073211/http://devconf.pl/speakers/tugberk-ugurlu/",
				IconCSSClasses: "external link",
			},
			{
				Name:           "Chosen Tweet",
				Link:           "https://twitter.com/tourismgeek/status/1045282398181761026",
				IconCSSClasses: "twitter",
			},
		},
	},

	{
		Title:           "Let the Uncertainty be Your Friend: Finding Your Path in a Wiggly Road",
		Activity:        "NewCrafts",
		ImageURL:        "https://tugberkugurlu.blob.core.windows.net/speaking/2018-05-newcrafts-paris/DdZ3x-fWkAcRn_c.jpg",
		City:            "Paris",
		Country:         "France",
		DisplayDate:     "17-18 May 2018",
		EmbededHTMLData: template.HTML(`<iframe class="embed-responsive-item" src="https://player.vimeo.com/video/275529797?title=0&byline=0&portrait=0" webkitallowfullscreen mozallowfullscreen allowfullscreen></iframe>`),
		Extras: []SpeakingActivityExtra{
			{
				Name:           "Slides",
				Link:           "https://speakerdeck.com/tourismgeek/let-the-uncertainty-be-your-friend-finding-your-path-in-a-wiggly-road",
				IconCSSClasses: "area chart",
			},
			{
				Name:           "Conference Profile",
				Link:           "http://www.ncrafts.io/speaker/tourismgeek",
				IconCSSClasses: "external link",
			},
			{
				Name:           "Chosen Tweet",
				Link:           "https://twitter.com/tourismgeek/status/997107891063721985",
				IconCSSClasses: "twitter",
			},
		},
	},
}
