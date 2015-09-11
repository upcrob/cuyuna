package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

func getHtml(markdownFile string) string {
	data, err := ioutil.ReadFile(markdownFile)
	if err != nil {
		fmt.Println("error reading file")
		os.Exit(1)
	}
	unsafe := blackfriday.MarkdownCommon(data)
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)

	return string(html)
}

func customMarkupHtml(html string) string {
	replaced := replaceCustomBlock(html, "IMPORTANT")
	replaced = replaceCustomBlock(replaced, "NOTE")
	replaced = replaceCustomBlock(replaced, "TIP")
	replaced = replaceCustomBlock(replaced, "WARNING")
	replaced = replaceCodeBlock(replaced)
	return replaced
}

func getCss(theme string) string {
	if theme == "simple" {
		return `p {
				font-family: "Palatino Linotype", Palatino, Book Antiqua;
				font-size: 14px;
			}

			h1, h2, h3, h4, h5, h6 {
				font-family: "Palatino Linotype", Palatino, Book Antiqua;
			}

			a:link {
				color: #0055cc;
				text-decoration: none;
			}

			a:visited {
				color: #0055cc;
				text-decoration: none;
			}

			a:hover {
				color: #0055cc;
				text-decoration: underline;
			}

			a:active {
				color: #0055cc;
				text-decoration: underline;
			}

			code {
				font-size: 14px;
				font-family: Courier New;
				background-color:rgba(0, 0, 0, 0.05);
			}

			table {
				font-family: "Palatino Linotype", Palatino, Book Antiqua;
				border-collapse: collapse;
			}

			td, th {
				font-size: 14px;
				border: 1px solid #172c4c;
				padding: 3px 10px 3px 10px;
			}

			th {
				font-size: 15px;
				text-align: left;
				padding-top: 5px;
				padding-bottom: 4px;
				background-color: #506899;
				color: #fff;
			}

			td {
				color: #000;
				background-color: #eef0f5;
			}

			.codeblock {
				display: block;
				white-space: pre-wrap;
				background-color: transparent;

				padding-left: 15px;
				padding-right: 5px;
				padding-top: 2px;
				padding-bottom: 2px;

				border-left: 5px;
				border-left-style: solid;
				border-left-color: #00FF44;

				margin-top: 20px;
				margin-bottom: 20px;
				margin-left: 30px;
			}

			.note, .warning, .important, .tip {
				font-family: "Palatino Linotype", Palatino, Book Antiqua;
				font-size: 18px;
				font-weight: bold;
				vertical-align: middle;
				border-right: 5px solid #506899;
				padding-left: 20px;
				padding-right: 3px;
				margin-right: 5px;
				padding-top: 3px;
				padding-bottom: 3px;
				width: 1px;
			}

			.notecontent, .warningcontent, .importantcontent, .tipcontent {
				font-family: "Palatino Linotype", Palatino, Book Antiqua;
				font-size: 14px;
				padding-left: 15px;
				padding-right: 2px;
				padding-top: 3px;
				padding-bottom: 3px;
			}

			.note {
				border-right: 5px solid #506899;
			}

			.warning {
				border-right: 5px solid #e65c00;
			}

			.important {
				border-right: 5px solid #cc0000;
			}

			.tip {
				border-right: 5px solid #e6e600;
			}

			.customBlockTable {
				display: table;
				width: 100%%;
				margin-top: 10px;
				margin-bottom: 10px;
			}
			`
	} else {
		return ""
	}
}

func replaceCustomBlock(text string, label string) string {
	reg, _ := regexp.Compile("(<p>" + label + ":)(.*?)(</p>)")
	replaceText := "<div class=\"customBlockTable\"><div style=\"display: table-row\"><div class=\"" + strings.ToLower(label) + "\" style=\"display: table-cell;\">" +
		label + "</div><div class=\"" + strings.ToLower(label) + "content\" style=\"display: table-cell\">$2</div></div></div>"
	return reg.ReplaceAllString(text, replaceText)
}

func replaceCodeBlock(text string) string {
	reg, _ := regexp.Compile("<pre><code>")
	replaceText := "<pre><code class=\"codeblock\">"
	return reg.ReplaceAllString(text, replaceText)
}