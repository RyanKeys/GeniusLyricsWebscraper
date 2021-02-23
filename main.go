package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"

	"github.com/gocolly/colly"
)

type Song struct {
    url string
    title string
    lyrics string
}

var foldername string = "Scraped-Data" 

var songlst = []Song {}

// main() contains code adapted from example found in Colly's docs:
// http://go-colly.org/docs/examples/basic/
func main() {
    // Instantiate default collector only allowed to visit 'genius.com'.
    c := colly.NewCollector(colly.AllowedDomains("genius.com"))
    // detailCollector is just a second collector to grab different data from nested objects.
    detailCollector := c.Clone()

    c.OnHTML("a[href]", func(e *colly.HTMLElement) {
        if len(e.Text) >= 1 {
            for i := 1; i < 100; i++ {
                if string(e.Text[0]) == fmt.Sprint(i) {
                    link := e.Attr("href")
                    detailCollector.Visit(link)                    
                }
            }
        }
    })

    detailCollector.OnHTML("p", func(e *colly.HTMLElement) {
        if strings.Contains(e.Text, "Verse") || strings.Contains(e.Text, "Chorus") || strings.Contains(e.Text, "Intro") {
            s := Song{url:"genius.com" + e.Request.URL.Path, title: e.Request.URL.Path[1:len(e.Request.URL.Path)], lyrics: e.Text}
            fmt.Printf("Found Song: %s", s.url + "\n")
            songlst = append(songlst, s)
        }
    })

    c.OnRequest(func(r *colly.Request) {
        fmt.Println("Visiting", r.URL)
    })

    c.OnError(func(_ *colly.Response, err error) {
        fmt.Println("Something went wrong:", err)
    })

    c.OnResponse(func(r *colly.Response) {
        fmt.Println("Visited", r.Request.URL)
    })

    c.OnScraped(func(r *colly.Response) {
        fmt.Println("Finished", r.Request.URL)
        
    })

    

    // Start scraping
    c.Visit("https://genius.com/#top-songs")
    c.Wait()
    for i := 0; i < len(songlst); i++ {
        songdata := string(songlst[i].url + "\n" + songlst[i].title + "\n" + string(songlst[i].lyrics))
        WriteFileFromString(songlst[i].title + ".txt", songdata)
    }
    
}


func WriteFileFromString(filename string, data string) {
    session, _ := user.Current()
	bytesToWrite := []byte(data)
    dirPath := session.HomeDir + "/" + "Desktop/" + foldername
    os.Mkdir(dirPath, 0755)
    err := ioutil.WriteFile(dirPath + "/" + filename, bytesToWrite, 0644)
	if err != nil {
		panic(err)
	}
}