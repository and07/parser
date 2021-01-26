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
			 },
			 {
				"type": "attr",
				"name": "img1",
				"attr": "src",
				"parent": "link",
				"page": "https://brainberries.co/interesting/8-amazing-celebs-with-their-own-businesses/",
				"path": "//*/figure[@class='wp-block-image size-large'][1]/img"
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
