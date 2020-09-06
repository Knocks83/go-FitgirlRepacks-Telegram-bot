/*
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package fitgirl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// BaseURL = Fitgirl repacks base site
const BaseURL = "https://fitgirl-repacks.site/wp-json/wp/v2/"

/* Contains the allowed website from which taking the urls
var because arrays cannot be constant
*/
var AllowedWebsites = []string{"1337x.to", "newkatcr.co", "hermietkreeft.site", "rutor.info", "tapochek.net", "pirated.me", "pasteit.top"}

// SearchResult = The struct that will contain the search result
type SearchResult struct {
	ID   uint16
	Name string
}

// DownloadLink =  The struct that will contain the link and the name of the host
type DownloadLink struct {
	Name string
	Link string
}

// Check if a string is contained in a string array
func Contains(s []string, e string) bool {
	for _, a := range s {
		if strings.Contains(e, a) {
			return true
		}
	}
	return false
}

func Search(query string) []SearchResult {
	// Create the array that'll contain the results
	searchRes := make([]SearchResult, 0)

	// Make GET request and handle errors
	resp, err := http.Get(BaseURL + "posts?search=" + url.QueryEscape(query))
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// Read body content and decode it (it's a JSON)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var res []map[string]interface{}

	json.Unmarshal(body, &res)

	// For each result, create a SearchResult object containing its ID and name
	for i := range res {
		id := uint16(res[i]["id"].(float64))
		name := html.UnescapeString(res[i]["title"].(map[string]interface{})["rendered"].(string))
		result := SearchResult{id, name}

		searchRes = append(searchRes, result)
	}

	return searchRes
}

func FindDownloadLinks(id uint16) []DownloadLink {
	// Create the array that'll contain all the download links
	downloadLinks := make([]DownloadLink, 0)

	// Make GET request and handle errors
	resp, err := http.Get(BaseURL + "posts/" + strconv.Itoa(int(id)))
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// Read body content and decode it (it's a JSON)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	var res map[string]interface{}

	json.Unmarshal(body, &res)

	// Get the page content (contains all the download links)
	pageContent := res["content"].(map[string]interface{})["rendered"].(string)

	// I don't know how to do it with regex
	// Remove the part before the links

	pageContent = strings.ToLower(pageContent)
	if strings.Contains(pageContent, "<h3>download mirrors:</h3>") {
		pageContent = strings.Split(pageContent, "<h3>download mirrors:</h3>")[1]
	} else if strings.Contains(pageContent, "<h3>download mirrors</h3>") {
		pageContent = strings.Split(pageContent, "<h3>download mirrors</h3>")[1]
	}

	// Remove the part after the links
	pageContent = strings.Split(pageContent, "<h3>")[0]

	// Parse the remaining string as html node
	doc, err := html.Parse(strings.NewReader(pageContent))
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var findText func(n *html.Node, buf *bytes.Buffer)
	findText = func(n *html.Node, buf *bytes.Buffer) {
		if n.Type == html.TextNode {
			buf.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findText(c, buf)
		}
	}

	var extractLinks func(*html.Node)
	extractLinks = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					// Parse the URL and extract the host, to append it to the list
					u, err := url.Parse(a.Val)
					if err != nil {
						panic(err)
					}
					// Check if it's an allowed host
					if !Contains(AllowedWebsites, u.Host) /*u.Host !in allowed hosts*/ {
						return
					}
					// Find the a tag text
					text := &bytes.Buffer{}
					findText(n, text)

					// I don't know how to name this one so I'll just skip it
					if text.String() == ".torrent file only" {
						return
					}

					// If it's a magnet, in the host name just write magnet
					if u.Scheme == "magnet" {
						downloadLinks = append(downloadLinks, DownloadLink{u.Scheme, a.Val})
					} else {
						downloadLinks = append(downloadLinks, DownloadLink{text.String(), a.Val})
					}
					return
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractLinks(c)
		}
	}
	// Find the first UL tag, so it'll reduce the amount of tags to parse
	var findFirstUl func(*html.Node)
	findFirstUl = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "ul" {
			extractLinks(n)
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findFirstUl(c)
		}
	}

	findFirstUl(doc)

	return downloadLinks
}

func FindGameTitle(id uint16) string {
	// Create the array that'll contain the results
	var gameTitle string

	// Make GET request and handle errors
	resp, err := http.Get(BaseURL + "posts/" + strconv.Itoa(int(id)))
	if err != nil {
		fmt.Println(err)
		return ""
	}

	// Read body content and decode it (it's a JSON)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	var res map[string]interface{}

	json.Unmarshal(body, &res)

	gameTitle = html.UnescapeString(res["title"].(map[string]interface{})["rendered"].(string))

	return gameTitle
}
