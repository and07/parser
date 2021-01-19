package main

var ruleData string = `
{
	"url":"https://brainberries.co/",
	"all":false,
	"rule":[
	   {
		  "type":"link",
		  "name":"link",
		  "attr":null,
		  "parent":null,
		  "page":"https://brainberries.co/",
		  "path":"//*/header[@class='entry-header']/h2[@class='entry-title']/a",
		  "children":[
			 {
				"type":"text",
				"name":"title1",
				"attr":null,
				"parent":"link",
				"page":"https://brainberries.co/interesting/8-amazing-celebs-with-their-own-businesses/",
				"path":"//*/h1[@class='entry-title']"
			 },
			 {
				"type":"html",
				"name":"test",
				"attr":null,
				"parent":"link",
				"page":"https://brainberries.co/interesting/8-amazing-celebs-with-their-own-businesses/",
				"path":"//*/div[@class='entry-content']"
			 }
		  ]
	   }
	],
	"host":"brainberries.co",
	"name":"",
	"domen":"https://brainberries.co",
	"limit":"1",
	"path_type":0
 }
`

var ruleData1 string = `
{
	"url":"https://brainberries.co/interesting/8-amazing-celebs-with-their-own-businesses/",
	"all":false,
	"rule":[
			 {
				"type":"text",
				"name":"title1",
				"attr":null,
				"parent":"link",
				"page":"https://brainberries.co/interesting/8-amazing-celebs-with-their-own-businesses/",
				"path":"//*/h1[@class='entry-title']"
			 },
			 {
				"type":"html",
				"name":"test",
				"attr":null,
				"parent":"link",
				"page":"https://brainberries.co/interesting/8-amazing-celebs-with-their-own-businesses/",
				"path":"//*/div[@class='entry-content']"
			 },
			  {
				"type": "attr",
				"name": "img1",
				"attr": "src",
				"parent": null,
				"page": "https://brainberries.co/interesting/8-amazing-celebs-with-their-own-businesses/",
				"path": "//*/img[@class='wp-image-83713']"
			  },
			  {
				"type": "attr",
				"name": "img2",
				"attr": "src",
				"parent": null,
				"page": "https://brainberries.co/interesting/8-amazing-celebs-with-their-own-businesses/",
				"path": "//*/img[@class='wp-image-83714']"
			  },
			  {
				"type": "attr",
				"name": "img3",
				"attr": "src",
				"parent": null,
				"page": "https://brainberries.co/interesting/8-amazing-celebs-with-their-own-businesses/",
				"path": "//*/img[@class='wp-image-83715']"
			  },
			  {
				"type": "attr",
				"name": "img4",
				"attr": "src",
				"parent": null,
				"page": "https://brainberries.co/interesting/8-amazing-celebs-with-their-own-businesses/",
				"path": "//*/img[@class='wp-image-83716']"
			  },
			  {
				"type": "attr",
				"name": "img5",
				"attr": "src",
				"parent": null,
				"page": "https://brainberries.co/interesting/8-amazing-celebs-with-their-own-businesses/",
				"path": "//*/img[@class='wp-image-83717']"
			  },
			  {
				"type": "attr",
				"name": "img6",
				"attr": "src",
				"parent": null,
				"page": "https://brainberries.co/interesting/8-amazing-celebs-with-their-own-businesses/",
				"path": "//*/img[@class='wp-image-83718']"
			  },
			  {
				"type": "attr",
				"name": "img7",
				"attr": "src",
				"parent": null,
				"page": "https://brainberries.co/interesting/8-amazing-celebs-with-their-own-businesses/",
				"path": "//*/img[@class='wp-image-83719']"
			  },
			  {
				"type": "attr",
				"name": "img8",
				"attr": "src",
				"parent": null,
				"page": "https://brainberries.co/interesting/8-amazing-celebs-with-their-own-businesses/",
				"path": "//*/img[@class='wp-image-83720']"
			  }
	],
	"host":"brainberries.co",
	"name":"",
	"domen":"https://brainberries.co/interesting/8-amazing-celebs-with-their-own-businesses/",
	"limit":"1",
	"path_type":0
 }
`
