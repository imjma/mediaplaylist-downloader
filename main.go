package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/grafov/m3u8"
)

func main() {
	f, err := os.Open("media.m3u8")
	if err != nil {
		panic(err)
	}
	p, listType, err := m3u8.DecodeFrom(bufio.NewReader(f), true)
	if err != nil {
		panic(err)
	}

	switch listType {
	case m3u8.MEDIA:
		mediapl := p.(*m3u8.MediaPlaylist)
		for _, seg := range mediapl.Segments {
			if seg != nil {
				u, err := url.Parse(seg.URI)
				if err != nil {
					panic(err)
				}
				fileName := path.Base(u.Path)
				fmt.Println(fileName)
				err = DownloadFile(fileName, seg.URI)
				if err != nil {
					panic(err)
				}
			}
		}
	case m3u8.MASTER:
		// masterpl := p.(*m3u8.MasterPlaylist)
	}

}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
