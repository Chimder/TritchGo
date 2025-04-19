package main

import (
	"log"
	"net/url"
	"strings"
)

func main() {
	// s1 := []int{1, 3, 2, 10, 9, 8}
	// s2 := []int{4, 5, 6}
	// s3 := append(s1, s2...)
	// slices.Sort(s3)
	manga := "https://mangadex.org/manga/new?id=1092323&loxx=1"
	urlparse, err := url.Parse(manga)
	if err != nil {
		log.Printf("errr %v", err)
	}
	log.Print(urlparse)
	log.Print(urlparse.Host)
	log.Print(urlparse.Scheme)
	log.Print(urlparse.Query().Get("loxx"))
	log.Print(urlparse.Path)
	new := strings.Split(urlparse.Path, "/")
	log.Print(new)
	// fmt.Print(os.Args)
	// fmt.Println(s3)
}
