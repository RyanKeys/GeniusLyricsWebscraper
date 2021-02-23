package lyricfinder

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

var songlst = []Song {}

var foldername string = "Scraped-Data" 


func Run() {
	// Instantiate default collector only allowed to visit 'genius.com'.
    c := colly.NewCollector(colly.AllowedDomains("genius.com"))
    // detailCollector is just a second collector to grab different data from nested objects.
    detailCollector := c.Clone()

    // Work in Progress: Get more songs per request. Selects the "Load More" div on the site.
    // c.OnHTML("#top-songs > div > div.PageGridCenter-q0ues6-0.Charts__LoadMore-sc-1re0f44-1.eDwRUT > div", func(e *colly.HTMLElement) {
    //     c.Visit(e.Request.Method)
    // })

    //Checks if links leads to song lyrics:
    c.OnHTML("a[href]", func(e *colly.HTMLElement) {
        // Handles empty elements
        if len(e.Text) >= 1 {
            // for each song in the top 100. TODO: See Work in Progress above.
            for i := 1; i < 100; i++ {
                // Checks that the first text element is the same as the current index. AKA Makes sure you're only getting the 'top songs' by grabbing them in order. 
                if string(e.Text[0]) == fmt.Sprint(i) {
                    //Run a sub-scraper 'detailCollector' on the found link.
                    link := e.Attr("href")
                    detailCollector.Visit(link)
                }
            }
        }
    })

    //Handles the collection of lyrics from the links found by 'c.OnHTML("a[href]", func(e *colly.HTMLElement){}'
    detailCollector.OnHTML("p", func(e *colly.HTMLElement) {
        //If Verse, Chorus, or Intro is found in any of the <p> elements; create song data and append it to list of found song lyrics. 
        if strings.Contains(e.Text, "Verse") || strings.Contains(e.Text, "Chorus") || strings.Contains(e.Text, "Intro") {
            // Creates a 'Song{}' object, and saves the url, title, and lyrical content of the current <p> tag.
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
        //Writes Song data stored in songlst to the current user's Desktop.
        for i := 0; i < len(songlst); i++ {
            //Converts 'Song{}' struct into a string.
            songdata := string(songlst[i].url + "\n" + songlst[i].title + "\n\n" + string(songlst[i].lyrics))
            // Writes 'songdata' string to a folder on your desktop. File name will be the artist and song name followed by lyrics. Example: 'Daft-Punk-get-luck-lyrics.txt'
            writeFileFromString(songlst[i].title + ".txt", songdata, getPathToDesktop(foldername))
        }
        fmt.Println("Finished", r.Request.URL)
        
    })

    
    // Start scraping
    c.Visit("https://genius.com/#top-songs")
    
    
}

// Private func to get current user's 'Desktop' folder location. Requires a name 'foldername' for the desired folder. Example: 'Scraped-Data' 
func getPathToDesktop(foldername string) string {
    session, _ := user.Current()
    return session.HomeDir + "/" + "Desktop/" + foldername
}

//Private func that needs: 'filename' a name for the file you are about to create, 'data' the string of data that will be the contents of your new file, and 'dirPath' the location where you want to write the files. 
func writeFileFromString(filename string, data string, dirPath string) {
    bytesToWrite := []byte(data)
    _, err := os.Stat(dirPath)
    if os.IsNotExist(err) {
        fmt.Println(err)
        os.Mkdir(dirPath, 0755)
    }
    
    err = ioutil.WriteFile(dirPath + "/" + filename, bytesToWrite, 0644)
	if err != nil {
		panic(err)
	}
}