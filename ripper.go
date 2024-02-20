package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"sync"

	"github.com/schollz/progressbar/v3"
)

// variables declaration section
var f_requests int
var f_threads int
var f_txt_export bool
var f_verbose bool
var f_very_verbose bool

// init func :
// cool braille anime figlet
// Aika - Go - =CLOUDRIPPER= + pp anime
// parse flags n switches
func init() {
	fmt.Println("/!\\ *the display will be broken if the console isn't at least as large as the logo*")
	fmt.Println(" _   __       _   __       _   __       _   __       _   __       _   __       _   __      ")
	fmt.Println("_________________________________________________________________________________________")
	figtext := `_______         _____  _     _ ______   ______ _____  _____   _____  _______  ______
 |       |      |     | |     | |     \ |_____/   |   |_____] |_____] |______ |_____/
 |_____  |_____ |_____| |_____| |_____/ |    \_ __|__ |       |       |______ |    \_`
	fmt.Println("\x1b[38;5;208m", figtext, "\x1b[0m")
	fmt.Println("_____________________________________________________________________made by \x1b[31mhanabi \x1b[0m______")
	fmt.Println(" _   __       _   __       _   __       _   __       _   __       _   __       _   __      ")
	fmt.Println(" |   /  .---' |   /  .---' |   /  .---' |   /  .---' |   /  .---' |   /  .---' |   /  .---'")
	fmt.Println(" `  /         `  /         `  /         `  /         `  /         `  /         `  /        ")
	fmt.Println("\x1b[38;5;208m  \\/           \\/           \\/           \\/           \\/           \\/           \\/\x1b[0m")
	fmt.Println("")
	flag.IntVar(&f_requests, "r", 0, "number of total requests")
	flag.IntVar(&f_threads, "t", 0, "number of threads")
	flag.BoolVar(&f_txt_export, "x", false, "export/append results to export.txt file")
	flag.BoolVar(&f_verbose, "v", false, "enable verbose mode")
	flag.BoolVar(&f_very_verbose, "vv", false, "enable very verbose mode")
	flag.Parse()
}

type Sync struct {
	mu           sync.Mutex
	url_match    []string
	req_total    int
	req_match    int
	req_notfound int
	rgx          regexp.Regexp
	bar          progressbar.ProgressBar
}

// main func
// check flag switches, and exec the core func with goroutines (using n threads), c'est quoi deja, go, channels, mutexes ?
func main() {
	if f_requests <= 0 {
		fmt.Println("[-] requests number not set. using default value (100 requests)")
		f_requests = 100
	}
	if f_threads <= 0 {
		fmt.Println("[-] threads number not set. using default value (1 thread)")
		f_threads = 1
	}

	c := Sync{
		req_total:    0,
		req_match:    0,
		req_notfound: 0,
		rgx:          *regexp.MustCompile(`/s-[a-zA-Z0-9]{11}`),
		bar: *progressbar.NewOptions(f_requests,
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionShowBytes(true),
			progressbar.OptionSetWidth(25),
			progressbar.OptionSetDescription("private tracks found :       |    progress :"),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "[red]=[reset]",
				SaucerHead:    "[red]>[reset]",
				SaucerPadding: " ",
				BarStart:      "[red][[reset]",
				BarEnd:        "[red]][reset]",
			})),
	}

	//init waitgroup channel
	var wg sync.WaitGroup
	//let's Go(routine)
	for i := 0; i < f_threads; i++ {
		//+1 delta
		wg.Add(1)

		if f_very_verbose {
			fmt.Println("THREAD", i, " <=> ", f_requests/f_threads, "REQ")
		}

		go c.ripper(f_requests/f_threads, &wg)
	}
	//wait for all the goroutines to finish (une fois que Add(x) == Done(x))
	wg.Wait()
	fmt.Println("\n[\x1b[38;5;208m+\x1b[0m] done !")
	if f_txt_export {
		fmt.Println("\n[\x1b[38;5;208m+\x1b[0m] exporting...")
		txt_export(&c)
		fmt.Println("\n[\x1b[38;5;208m+\x1b[0m]\x1b[38;5;208m results written in 'export.txt'\x1b[0m")
	}
}

// fonction de la livebar
// concretement c'est un taf de gitan :
// la progressbar initiale est constamment loaded, et je rajoute un print par dessus pour afficher les matchs
// car j'arrive pas a avoir plusieurs lignes actualisés, autre que celle de la progressbar
func (c *Sync) livestream() {
	c.bar.Add(1)
	fmt.Print("\r", "\x1b[38;5;208mprivate tracks found : \x1b[0m", c.req_match)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandGen() string {
	b := make([]byte, 5)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return "https://on.soundcloud.com/" + string(b)
}

// GO RIPPER FUNC (generate random, catch redirection link, regex and save in temp array)
func (c *Sync) ripper(req_number int, wg *sync.WaitGroup) {
	//+1 done routine, defer
	defer wg.Done()

	//ratata requests
	for i := 0; i < req_number; i++ {
		//mutex lock :fire:
		c.mu.Lock()
		c.req_total++
		c.livestream()
		c.mu.Unlock()

		//webreq -> catch redirect -> regex redirect link -> into array if match
		req, err := http.NewRequest("GET", RandGen(), nil)
		if err != nil {
			if f_very_verbose {
				fmt.Println("error making request")
			}
		}
		client := new(http.Client)
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return errors.New("Redirect")
		}

		response, _ := client.Do(req)
		if response != nil && response.StatusCode == http.StatusFound { //status code 302
			location := response.Header.Get("Location")
			matched := c.rgx.FindString(location)
			if matched != "" {
				if f_very_verbose {
					println(location)
				}
				c.mu.Lock()
				c.url_match = append(c.url_match, location)
				c.req_match++
				c.mu.Unlock()
			}

		} else {
			c.mu.Lock()
			c.req_notfound++
			c.mu.Unlock()
		}
	}
}

// EXPORT FUNC
func txt_export(c *Sync) error {
	//opens export file (or create it if not present) | O_APPEND & CREATE
	file, err := os.OpenFile("export.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// make da writer
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	// foreach 2.0 Golang Edition
	for _, line := range c.url_match {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

//FP FILTER FUNC
