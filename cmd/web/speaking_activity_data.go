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
		EmbededHTMLData: template.HTML(`<div class="embed-responsive embed-responsive-16by9">
	<iframe class="embed-responsive-item" src="https://www.youtube.com/embed/qpfFus69pN8" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>
</div>`),
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
		EmbededHTMLData: template.HTML(`<div class="embed-responsive embed-responsive-16by9">
	<iframe class="embed-responsive-item" src="https://player.vimeo.com/video/275529797?title=0&byline=0&portrait=0" webkitallowfullscreen mozallowfullscreen allowfullscreen></iframe>
</div>`),
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

	{
		Title:           "Levelling up to Become a Technical Lead",
		Activity:        "DDD Scotland",
		ImageURL:        "https://tugberkugurlu.blob.core.windows.net/speaking/2018-02-ddd-scotland-glasgow/DVrgv-JX4AMjraj.jpg",
		City:            "Glasgow",
		Country:         "Scotland",
		DisplayDate:     "10 Feb 2018",
		EmbededHTMLData: template.HTML(`<script async class="speakerdeck-embed" data-id="78340f5c21f34357943c8259f323614e" data-ratio="1.77777777777778" src="//speakerdeck.com/assets/embed.js"></script>`),
		Extras: []SpeakingActivityExtra{
			{
				Name:           "Slides",
				Link:           "https://speakerdeck.com/tourismgeek/levelling-up-to-become-a-technical-lead",
				IconCSSClasses: "area chart",
			},
			{
				Name:           "Chosen Tweet",
				Link:           "https://twitter.com/tourismgeek/status/961963028651610117",
				IconCSSClasses: "twitter",
			},
		},
	},

	{
		Title:           "Getting Into the Zero Downtime Deployment World",
		Activity:        "NDC Oslo",
		ImageURL:        "https://tugberkugurlu.blob.core.windows.net/speaking/2016-06-ndc-oslo/02B6DF5D-7DEF-428E-9385-5E1539A6B665.png",
		City:            "Oslo",
		Country:         "Norway",
		DisplayDate:     "6-10 June 2016",
		EmbededHTMLData: template.HTML(`<div class="embed-responsive embed-responsive-16by9">
	<iframe class="embed-responsive-item" src="https://player.vimeo.com/video/171317249?title=0&byline=0&portrait=0" webkitallowfullscreen mozallowfullscreen allowfullscreen></iframe>
</div>`),
		Extras: []SpeakingActivityExtra{
			{
				Name:           "Slides",
				Link:           "https://speakerdeck.com/tourismgeek/let-the-uncertainty-be-your-friend-finding-your-path-in-a-wiggly-road",
				IconCSSClasses: "area chart",
			},
			{
				Name:           "Code",
				Link:           "https://github.com/tugberkugurlu/AspNetCoreSamples/tree/ndcoslo2016/haproxy-zero-downtime-sample",
				IconCSSClasses: "code",
			},
			{
				Name:           "Conference Profile",
				Link:           "https://web.archive.org/web/20160530194908/https://ndcoslo.com/speaker/tugberk-ugurlu/",
				IconCSSClasses: "external link",
			},
			{
				Name:           "Announcment Article",
				Link:           "http://www.tugberkugurlu.com/archive/off-to-oslo-for-ndc-developer-conference",
				IconCSSClasses: "rss",
			},
			{
				Name:           "Recap Article",
				Link:           "http://www.tugberkugurlu.com/archive/ndc-oslo-2016-in-a-nutshell",
				IconCSSClasses: "rss",
			},
		},
	},
}
