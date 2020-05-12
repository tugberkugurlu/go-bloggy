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
		Title:       "Essentials for Building and Leading Highly Effective Development Teams",
		Activity:    "DevConf",
		ImageURL:    "https://tugberkugurlu.blob.core.windows.net/speaking/2018-09-devconf-krakow/30236914487_dba11cff33_o.jpg",
		City:        "Kraków",
		Country:     "Poland",
		DisplayDate: "27-28 Sep 2018",
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
		Title:       "Let the Uncertainty be Your Friend: Finding Your Path in a Wiggly Road",
		Activity:    "NewCrafts",
		ImageURL:    "https://tugberkugurlu.blob.core.windows.net/speaking/2018-05-newcrafts-paris/DdZ3x-fWkAcRn_c.jpg",
		City:        "Paris",
		Country:     "France",
		DisplayDate: "17-18 May 2018",
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
		Title:       "Getting Into the Zero Downtime Deployment World",
		Activity:    "NDC Oslo",
		ImageURL:    "https://tugberkugurlu.blob.core.windows.net/speaking/2016-06-ndc-oslo/02B6DF5D-7DEF-428E-9385-5E1539A6B665.png",
		City:        "Oslo",
		Country:     "Norway",
		DisplayDate: "6-10 June 2016",
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

	{
		Title:       "Latest SQL Compare features and support for SQL Server 2017",
		Activity:    "SQL in the City Streamed",
		ImageURL:    "https://tugberkugurlu.blob.core.windows.net/speaking/2017-12-sql-in-the-city-cambridge/E9F8A9A0-111E-4D1F-AD1B-18B9A26B7DB3.png",
		City:        "Cambridge",
		Country:     "England",
		DisplayDate: "13 Dec 2017",
		EmbededHTMLData: template.HTML(`<div class="embed-responsive embed-responsive-16by9">
	<iframe class="embed-responsive-item" src="https://www.youtube.com/embed/9a6gVaX192g" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>
</div>`),
		Extras: []SpeakingActivityExtra{
			{
				Name:           "Conference Profile",
				Link:           "https://www.red-gate.com/hub/events/sqlinthecity/2017",
				IconCSSClasses: "external link",
			},
			{
				Name:           "Announcment Article",
				Link:           "http://www.tugberkugurlu.com/archive/speaking-at-sql-in-the-city-2017-register-now",
				IconCSSClasses: "rss",
			},
		},
	},

	{
		Title:       "Docker Changes the Way You Develop and Release Your Scalable Solutions",
		Activity:    "I T.A.K.E Unconference",
		ImageURL:    "https://tugberkugurlu.blob.core.windows.net/speaking/2016-05-itakeunconf-bucharest/Ci40PKzUkAAm6S5-small.jpg",
		City:        "Bucharest",
		Country:     "Romania",
		DisplayDate: "19-20 May 2016",
		EmbededHTMLData: template.HTML(`<div class="embed-responsive embed-responsive-16by9">
	<iframe class="embed-responsive-item" src="https://www.youtube.com/embed/nswAIlgD4fU" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>
</div>`),
		Extras: []SpeakingActivityExtra{
			{
				Name:           "Slides",
				Link:           "https://speakerdeck.com/tourismgeek/docker-changes-the-way-you-develop-and-release-your-scalable-solutions",
				IconCSSClasses: "area chart",
			},
			{
				Name:           "Conference Profile",
				Link:           "http://itakeunconf.com/speakers/tugberk-ugurlu/",
				IconCSSClasses: "external link",
			},
			{
				Name:           "Announcment Article",
				Link:           "http://itakeunconf.com/docker-zero-downtime-deployment-rules/",
				IconCSSClasses: "rss",
			},
		},
	},

	{
		Title:       "Getting Into the Zero Downtime Deployment World",
		Activity:    "Dev Day",
		ImageURL:    "https://tugberkugurlu.blob.core.windows.net/speaking/2016-09-dev-day-krakow/7E202F4A-548B-44C5-B3C7-A2F2BAFA09F3.png",
		City:        "Kraków",
		Country:     "Poland",
		DisplayDate: "14-16 Sep 2016",
		EmbededHTMLData: template.HTML(`<div class="embed-responsive embed-responsive-16by9">
	<iframe class="embed-responsive-item" src="https://www.youtube.com/embed/mTZSvK6I3Xs" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>
</div>`),
		Extras: []SpeakingActivityExtra{
			{
				Name:           "Slides",
				Link:           "https://speakerdeck.com/tourismgeek/getting-into-the-zero-downtime-deployment-world-1",
				IconCSSClasses: "area chart",
			},
			{
				Name:           "Conference Profile",
				Link:           "https://web.archive.org/web/20161003064950/http://devday.pl/",
				IconCSSClasses: "external link",
			},
		},
	},

	{
		Title:           "Architecting Polyglot-Persistent Solutions",
		Activity:        "DevConf",
		ImageURL:        "https://tugberkugurlu.blob.core.windows.net/bloggyimages/c9bbd57d-1275-4517-b5b1-d85b022f3279.jpg",
		City:            "Johannesburg",
		Country:         "South Africa",
		DisplayDate:     "8 March 2016",
		EmbededHTMLData: template.HTML(`<script async class="speakerdeck-embed" data-id="3474381924404553afdfbd7c8010fcfa" data-ratio="1.77777777777778" src="//speakerdeck.com/assets/embed.js"></script>`),
		Extras: []SpeakingActivityExtra{
			{
				Name:           "Slides",
				Link:           "https://speakerdeck.com/tourismgeek/architecting-polyglot-persistent-solutions",
				IconCSSClasses: "area chart",
			},
			{
				Name:           "Conference Profile",
				Link:           "https://web.archive.org/web/20160318163414/https://www.devconf.co.za/",
				IconCSSClasses: "external link",
			},
			{
				Name:           "Chosen Tweet",
				Link:           "https://twitter.com/tourismgeek/status/707123586197295104",
				IconCSSClasses: "twitter",
			},
			{
				Name:           "Announcment Article",
				Link:           "http://www.tugberkugurlu.com/archive/upcoming-conferences-and-talks",
				IconCSSClasses: "rss",
			},
			{
				Name:           "Recap Article",
				Link:           "http://www.tugberkugurlu.com/archive/my-summary-of-devconf-2016",
				IconCSSClasses: "rss",
			},
		},
	},

	{
		Title:       "Profiling .NET server applications",
		Activity:    "Umbraco UK Festival",
		ImageURL:    "https://tugberkugurlu.blob.core.windows.net/speaking/2015-10-umbraco-uk-festival/A7F75261-4584-4D79-9EC8-D3FA216731A5.png",
		City:        "London",
		Country:     "England",
		DisplayDate: "30 Oct 2015",
		EmbededHTMLData: template.HTML(`<div class="embed-responsive embed-responsive-16by9">
	<iframe class="embed-responsive-item" src="https://www.youtube.com/embed/OEZKXRWDv60" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>
</div>`),
		Extras: []SpeakingActivityExtra{
			{
				Name:           "Slides",
				Link:           "https://speakerdeck.com/tourismgeek/profiling-net-server-applications",
				IconCSSClasses: "area chart",
			},
			{
				Name:           "Conference Profile",
				Link:           "http://www.tugberkugurlu.com/archive/my-talk-on-profiling--net-server-applications-from-umbraco-uk-festival-2015",
				IconCSSClasses: "external link",
			},
		},
	},

	{
		Title:       "ASP.NET 5: How to Get Your Cheese Back",
		Activity:    "Progressive .NET Tutorials",
		ImageURL:    "https://tugberkugurlu.blob.core.windows.net/speaking/2015-07-progressive-dot-net-london/19495238401_e4ed74fa5c_z.jpg",
		City:        "London",
		Country:     "England",
		DisplayDate: "1 July 2015",
		EmbededHTMLData: template.HTML(`<div>
	<a href="https://skillsmatter.com/skillscasts/6401-aspdot-net-5-how-to-get-your-cheese-back#video" target="_blank">
    	<img src="https://tugberkugurlu.blob.core.windows.net/speaking/2015-07-progressive-dot-net-london/cheese-back-CBEEBCB1-34C3-4D22-822B-55D7498A0141.png" alt="Tugberk Ugurlu @@ Progressive .NET Tutorials, London - 2015">
	</a>
</div>`),
		Extras: []SpeakingActivityExtra{
			{
				Name:           "Slides",
				Link:           "https://speakerdeck.com/tourismgeek/asp-dot-net-5-how-to-get-your-cheese-back",
				IconCSSClasses: "area chart",
			},
			{
				Name:           "Code",
				Link:           "https://github.com/british-proverbs/british-proverbs-mvc-6/tree/0ca8143e19a76a43cceeafdcecd28f69007a9108",
				IconCSSClasses: "code",
			},
			{
				Name:           "Conference Profile",
				Link:           "https://skillsmatter.com/legacy_profile/tugberk-ugurlu",
				IconCSSClasses: "external link",
			},
			{
				Name:           "Announcment Article",
				Link:           "http://www.tugberkugurlu.com/archive/upcoming-conferences-that-i-am-speaking-at",
				IconCSSClasses: "rss",
			},
			{
				Name:           "Recap Article",
				Link:           "http://www.tugberkugurlu.com/archive/progressive--net-tutorials-2015-and-recording-videos-of-my-asp-net-5-talks",
				IconCSSClasses: "rss",
			},
		},
	},
}
