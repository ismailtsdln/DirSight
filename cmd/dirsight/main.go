package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"text/tabwriter"
	"time"

	"github.com/ismailtsdln/DirSight/internal/engine"
	"github.com/ismailtsdln/DirSight/internal/wordlist"
)

const (
	banner = `
%s  ____  _      ____  _       _     _   
 |  _ \(_)_ __/ ___|(_) __ _| |__ | |_ 
 | | | | | '__\___ \| |/ _' | '_ \| __|
 | |_| | | |   ___) | | (_| | | | | |_ 
 |____/|_|_|  |____/|_|\__, |_| |_|\__|
                       |___/           
    %sAdvanced Directory Discovery & Bypass%s
	`
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[37m"
	colorBold   = "\033[1m"
)

func main() {
	fmt.Printf(banner, colorCyan, colorBold+colorPurple, colorReset)
	target := flag.String("u", "", "Target URL (e.g., https://example.com)")
	wordlistPath := flag.String("w", "", "Path to wordlist file")
	threads := flag.Int("t", 10, "Number of threads")
	timeout := flag.Duration("timeout", 10*time.Second, "Request timeout")
	proxy := flag.String("proxy", "", "Proxy URL (e.g., http://127.0.0.1:8080)")
	insecure := flag.Bool("k", false, "Allow insecure server connections when using SSL")
	expand := flag.Bool("expand", false, "Expand wordlist with 403 bypass variations")
	jsonPath := flag.String("json", "", "Export results to JSON file")

	flag.Parse()

	if *target == "" || *wordlistPath == "" {
		fmt.Println("Usage: dirsight -u <url> -w <wordlist>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("%s[*] DirSight starting on target:%s %s\n", colorBlue, colorReset, *target)

	// Initialize Client
	client, err := engine.NewClient(*timeout, *proxy, *insecure)
	if err != nil {
		fmt.Printf("%s[!] Error initializing client:%s %v\n", colorRed, colorReset, err)
		os.Exit(1)
	}

	// Initialize Loader and load wordlist
	loader := &wordlist.Loader{}
	list, err := loader.LoadFromFile(*wordlistPath)
	if err != nil {
		fmt.Printf("Error loading wordlist: %v\n", err)
		os.Exit(1)
	}

	if *expand {
		fmt.Printf("%s[*] Expanding wordlist with bypass variations...%s\n", colorYellow, colorReset)
		list = loader.ExpandWithBypasses(list)
	}

	fmt.Printf("%s[*] Loaded %d paths to scan%s\n", colorBlue, len(list), colorReset)

	// Perform initial request to check for WAF
	fmt.Printf("%s[*] Checking for WAF...%s\n", colorYellow, colorReset)
	initialReq, _ := http.NewRequest("GET", *target, nil)
	initialResp, err := client.Do(initialReq)
	if err == nil {
		waf := engine.DetectWAF(initialResp)
		if waf != engine.WAFNone {
			fmt.Printf("%s[!] WAF Detected:%s %s\n", colorRed, colorReset, waf)
		} else {
			fmt.Printf("%s[*] No common WAF detected.%s\n", colorGreen, colorReset)
		}
		initialResp.Body.Close()
	}

	// Initialize Scanner
	scanner := engine.NewScanner(client, *threads)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start scanning
	go scanner.Scan(ctx, *target, list)

	fmt.Printf("%s[*] Scanning...%s\n", colorBlue, colorReset)

	var allResults []engine.Result
	filter := &engine.Filter{ExcludeStatusCodes: []int{404}}
	var processed uint64
	total := uint64(len(list))

	// Result Writer
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprintln(w, colorBold+"STATUS\tLENGTH\tURL\tMETHOD"+colorReset)

	// Process results
	resultsDone := make(chan bool)
	go func() {
		for result := range scanner.Results {
			atomic.AddUint64(&processed, 1)
			if filter.ShouldShow(result) {
				statusColor := colorGreen
				if result.StatusCode >= 300 && result.StatusCode < 400 {
					statusColor = colorCyan
				} else if result.StatusCode >= 400 {
					statusColor = colorYellow
				}
				// Clear the current line before printing a result to avoid progress bar artifacts
				fmt.Print("\r\033[K")
				fmt.Fprintf(w, "%s%d%s\t%d\t%s\t%s\n", statusColor, result.StatusCode, colorReset, result.Length, result.URL, result.Method)
				w.Flush()
				allResults = append(allResults, result)
			}
		}
		resultsDone <- true
	}()

	// Progress indicator
	progressDone := make(chan bool)
	ticker := time.NewTicker(200 * time.Millisecond)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				p := atomic.LoadUint64(&processed)
				percentage := float64(p) / float64(total) * 100
				fmt.Printf("\r%s[*] Progress: [%-20s] %.2f%% (%d/%d)%s", colorGray, strings.Repeat("=", int(percentage/5)), percentage, p, total, colorReset)
			case <-resultsDone:
				fmt.Printf("\r%s[*] Progress: [%-20s] 100.00%% (%d/%d)%s\n", colorGreen, strings.Repeat("=", 20), total, total, colorReset)
				progressDone <- true
				return
			}
		}
	}()

	<-progressDone

	if *jsonPath != "" {
		// ... (keep the same JSON logic but maybe with colors)
		file, err := os.Create(*jsonPath)
		if err != nil {
			fmt.Printf("\n%s[!] Error creating JSON file:%s %v\n", colorRed, colorReset, err)
		} else {
			defer file.Close()
			encoder := json.NewEncoder(file)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(allResults); err != nil {
				fmt.Printf("\n%s[!] Error encoding JSON:%s %v\n", colorRed, colorReset, err)
			} else {
				fmt.Printf("\n%s[*] Results exported to %s%s\n", colorBlue, *jsonPath, colorReset)
			}
		}
	}

	fmt.Printf("%s[+] Scan completed. %d interesting results found.%s\n", colorGreen, len(allResults), colorReset)
}
