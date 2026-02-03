package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

// ANSI color codes
const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"
	Bold    = "\033[1m"
	Dim     = "\033[2m"
)

var verbose bool

var DEFAULT_SOURCES = []string{
	"https://raw.githubusercontent.com/proxifly/free-proxy-list/main/proxies/protocols/http/data.txt",
	"https://raw.githubusercontent.com/proxifly/free-proxy-list/main/proxies/protocols/https/data.txt",
	"https://raw.githubusercontent.com/TheSpeedX/PROXY-List/master/http.txt",
	"https://raw.githubusercontent.com/iplocate/free-proxy-list/main/protocols/http.txt",
	"https://raw.githubusercontent.com/ClearProxy/checked-proxy-list/main/http/raw/all.txt",
	"https://raw.githubusercontent.com/ALIILAPRO/Proxy/main/http.txt",
	"https://raw.githubusercontent.com/roosterkid/openproxylist/main/HTTPS_RAW.txt",
	"https://raw.githubusercontent.com/vakhov/fresh-proxy-list/master/http.txt",
	"https://raw.githubusercontent.com/vakhov/fresh-proxy-list/master/https.txt",
	"https://raw.githubusercontent.com/monosans/proxy-list/refs/heads/main/proxies/http.txt",
	"https://raw.githubusercontent.com/mmpx12/proxy-list/refs/heads/master/http.txt",
	"https://raw.githubusercontent.com/mmpx12/proxy-list/refs/heads/master/https.txt",
	"https://raw.githubusercontent.com/Zaeem20/FREE_PROXIES_LIST/refs/heads/master/http.txt",
	"https://raw.githubusercontent.com/Zaeem20/FREE_PROXIES_LIST/refs/heads/master/https.txt",
	"https://raw.githubusercontent.com/sunny9577/proxy-scraper/refs/heads/master/proxies.txt",
	"https://raw.githubusercontent.com/mzyui/proxy-list/refs/heads/main/http.txt",
	"https://raw.githubusercontent.com/elliottophellia/proxylist/refs/heads/master/results/http/global/http_checked.txt",
	"https://raw.githubusercontent.com/officialputuid/KangProxy/refs/heads/KangProxy/http/http.txt",
	"https://raw.githubusercontent.com/officialputuid/KangProxy/refs/heads/KangProxy/https/https.txt",
	"https://raw.githubusercontent.com/databay-labs/free-proxy-list/refs/heads/master/http.txt",
	"https://raw.githubusercontent.com/claude89757/free_https_proxies/refs/heads/main/https_proxies.txt",
	"https://raw.githubusercontent.com/claude89757/free_https_proxies/refs/heads/main/isz_https_proxies.txt",
	"https://raw.githubusercontent.com/r00tee/Proxy-List/refs/heads/main/Https.txt",
	"https://raw.githubusercontent.com/fyvri/fresh-proxy-list/archive/storage/classic/http.txt",
	"https://raw.githubusercontent.com/vmheaven/VMHeaven-Free-Proxy-Updated/refs/heads/main/http.txt",
	"https://raw.githubusercontent.com/theriturajps/proxy-list/refs/heads/main/proxies.txt",
	"https://raw.githubusercontent.com/ProxyScraper/ProxyScraper/refs/heads/main/http.txt",
	"https://raw.githubusercontent.com/trio666/proxy-checker/refs/heads/main/http.txt",
	"https://raw.githubusercontent.com/trio666/proxy-checker/refs/heads/main/https.txt",
	"https://raw.githubusercontent.com/Anonym0usWork1221/Free-Proxies/refs/heads/main/proxy_files/http_proxies.txt",
	"https://raw.githubusercontent.com/Anonym0usWork1221/Free-Proxies/refs/heads/main/proxy_files/https_proxies.txt",
}

// ProxyInfo holds detailed information about a working proxy
type ProxyInfo struct {
	Proxy       string     `json:"proxy"`
	TestedAt    string     `json:"tested_at"`
	IP          string     `json:"ip"`
	Success     bool       `json:"success"`
	Type        string     `json:"type"`
	Continent   string     `json:"continent"`
	Country     string     `json:"country"`
	CountryCode string     `json:"country_code"`
	Region      string     `json:"region"`
	RegionCode  string     `json:"region_code"`
	City        string     `json:"city"`
	Latitude    float64    `json:"latitude"`
	Longitude   float64    `json:"longitude"`
	IsEU        bool       `json:"is_eu"`
	Postal      string     `json:"postal"`
	Connection  Connection `json:"connection"`
	Timezone    Timezone   `json:"timezone"`
}

type Connection struct {
	ASN    int    `json:"asn"`
	Org    string `json:"org"`
	ISP    string `json:"isp"`
	Domain string `json:"domain"`
}

type Timezone struct {
	ID          string `json:"id"`
	Abbr        string `json:"abbr"`
	IsDST       bool   `json:"is_dst"`
	Offset      int    `json:"offset"`
	UTC         string `json:"utc"`
	CurrentTime string `json:"current_time"`
}

// Logging helpers
func logInfo(msg string, args ...any) {
	if verbose {
		timestamp := time.Now().Format("15:04:05")
		fmt.Printf("%s[%s]%s %s%s%s ", Dim, timestamp, Reset, Cyan, msg, Reset)
		for i := 0; i < len(args); i += 2 {
			if i+1 < len(args) {
				fmt.Printf("%s%v%s=%v ", Yellow, args[i], Reset, args[i+1])
			}
		}
		fmt.Println()
	}
}

func logSuccess(msg string, args ...any) {
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("%s[%s]%s %s✓ %s%s ", Dim, timestamp, Reset, Green+Bold, msg, Reset)
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			fmt.Printf("%s%v%s=%v ", Yellow, args[i], Reset, args[i+1])
		}
	}
	fmt.Println()
}

func logError(msg string, args ...any) {
	timestamp := time.Now().Format("15:04:05")
	fmt.Printf("%s[%s]%s %s✗ %s%s ", Dim, timestamp, Reset, Red+Bold, msg, Reset)
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			fmt.Printf("%s%v%s=%v ", Yellow, args[i], Reset, args[i+1])
		}
	}
	fmt.Println()
}

func logWarn(msg string, args ...any) {
	if verbose {
		timestamp := time.Now().Format("15:04:05")
		fmt.Printf("%s[%s]%s %s⚠ %s%s ", Dim, timestamp, Reset, Yellow+Bold, msg, Reset)
		for i := 0; i < len(args); i += 2 {
			if i+1 < len(args) {
				fmt.Printf("%s%v%s=%v ", Yellow, args[i], Reset, args[i+1])
			}
		}
		fmt.Println()
	}
}

func logDebug(msg string, args ...any) {
	if verbose {
		timestamp := time.Now().Format("15:04:05")
		fmt.Printf("%s[%s] %s ", Dim, timestamp, msg)
		for i := 0; i < len(args); i += 2 {
			if i+1 < len(args) {
				fmt.Printf("%v=%v ", args[i], args[i+1])
			}
		}
		fmt.Println(Reset)
	}
}

func printProgress(current, total, working int, proxy string) {
	if !verbose {
		percent := float64(current) / float64(total) * 100
		barWidth := 30
		filledWidth := int(float64(barWidth) * float64(current) / float64(total))
		bar := strings.Repeat("█", filledWidth) + strings.Repeat("░", barWidth-filledWidth)

		// Truncate proxy for display - fixed width
		displayProxy := proxy
		if len(displayProxy) > 30 {
			displayProxy = displayProxy[:27] + "..."
		}
		// Pad to fixed width to overwrite old content
		displayProxy = fmt.Sprintf("%-30s", displayProxy)

		// Clear line with ANSI escape code and print
		fmt.Printf("\r\033[K%s[%s]%s %s%s%s %s%5.1f%%%s %s[%d/%d]%s %sWorking: %s%d%s %s%s%s",
			Dim, time.Now().Format("15:04:05"), Reset,
			Cyan, bar, Reset,
			Bold, percent, Reset,
			Dim, current, total, Reset,
			Green, Bold, working, Reset,
			Dim, displayProxy, Reset)
	}
}

// Function to fetch proxy list and save lines to Slice
func getList(proxyList string, custom bool) []string {
	var proxies []string

	if custom {
		// Use custom file
		content, err := os.Open(proxyList)
		if err != nil {
			logError("Failed to open proxy list", "file", proxyList, "error", err)
			os.Exit(1)
		}
		defer content.Close()
		scanner := bufio.NewScanner(content)

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				proxies = append(proxies, line)
			}
		}
		logSuccess("Discovered proxies", "count", len(proxies), "file", proxyList)
	} else {
		// Scrape from DEFAULT_SOURCES
		logInfo("Scraping proxies from default sources", "sources", len(DEFAULT_SOURCES))

		client := &http.Client{Timeout: 30 * time.Second}
		seen := make(map[string]bool) // Deduplicate proxies

		for i, source := range DEFAULT_SOURCES {
			if verbose {
				logInfo("Fetching source", "current", i+1, "total", len(DEFAULT_SOURCES), "url", source)
			} else {
				printProgress(i+1, len(DEFAULT_SOURCES), len(proxies), source)
			}

			resp, err := client.Get(source)
			if err != nil {
				logDebug("Failed to fetch source", "url", source, "error", err)
				continue
			}

			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				logDebug("Failed to read response", "url", source, "error", err)
				continue
			}

			lines := strings.Split(string(body), "\n")
			added := 0
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" && !seen[line] {
					seen[line] = true
					proxies = append(proxies, line)
					added++
				}
			}
			logDebug("Added proxies from source", "count", added, "url", source)
		}

		if !verbose {
			fmt.Println() // New line after progress
		}
		logSuccess("Scraped proxies from default sources", "total", len(proxies), "sources", len(DEFAULT_SOURCES))
	}

	return proxies
}

// loadProxies loads proxies from file or scrapes from default sources, dedupes, then shuffles them
func loadProxies(proxyFile string) []string {
	var proxies []string
	if proxyFile == "" {
		// No file provided, scrape from default sources
		proxies = getList("", false)
	} else {
		// Use custom file
		proxies = getList(proxyFile, true)
	}

	// Deduplicate the list
	seen := make(map[string]bool)
	dedupedProxies := make([]string, 0, len(proxies))
	for _, proxy := range proxies {
		if !seen[proxy] {
			seen[proxy] = true
			dedupedProxies = append(dedupedProxies, proxy)
		}
	}

	if len(proxies) != len(dedupedProxies) {
		logSuccess("Deduplicated proxy list", "before", len(proxies), "after", len(dedupedProxies), "removed", len(proxies)-len(dedupedProxies))
	}

	rand.Shuffle(len(dedupedProxies), func(i, j int) {
		dedupedProxies[i], dedupedProxies[j] = dedupedProxies[j], dedupedProxies[i]
	})
	return dedupedProxies
}

// testProxy tests a single proxy and returns ProxyInfo if it works, nil otherwise
func testProxy(proxyStr string, timeout int) *ProxyInfo {
	//originalProxy := proxyStr
	// Only add http:// if the proxy doesn't already have a scheme
	if !strings.HasPrefix(proxyStr, "http://") && !strings.HasPrefix(proxyStr, "https://") {
		proxyStr = "http://" + proxyStr
	}

	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		logWarn("Invalid proxy format, skipping", "proxy", proxyStr, "error", err)
		return nil
	}

	tr := &http.Transport{
		MaxIdleConns:       100,
		IdleConnTimeout:    time.Duration(timeout) * time.Second,
		DisableCompression: true,
		Proxy:              http.ProxyURL(proxyURL),
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(timeout) * time.Second,
	}

	resp, err := client.Get("https://ipwho.is/")
	if err != nil {
		logDebug("Failed", "proxy", proxyStr, "error", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		logDebug("Bad status", "proxy", proxyStr, "status", resp.Status)
		return nil
	}

	// Read and parse the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logWarn("Failed to read response body", "proxy", proxyStr, "error", err)
		return nil
	}

	// Parse JSON response from ipwho.is
	var ipData struct {
		IP          string  `json:"ip"`
		Success     bool    `json:"success"`
		Type        string  `json:"type"`
		Continent   string  `json:"continent"`
		Country     string  `json:"country"`
		CountryCode string  `json:"country_code"`
		Region      string  `json:"region"`
		RegionCode  string  `json:"region_code"`
		City        string  `json:"city"`
		Latitude    float64 `json:"latitude"`
		Longitude   float64 `json:"longitude"`
		IsEU        bool    `json:"is_eu"`
		Postal      string  `json:"postal"`
		Connection  struct {
			ASN    int    `json:"asn"`
			Org    string `json:"org"`
			ISP    string `json:"isp"`
			Domain string `json:"domain"`
		} `json:"connection"`
		Timezone struct {
			ID          string `json:"id"`
			Abbr        string `json:"abbr"`
			IsDST       bool   `json:"is_dst"`
			Offset      int    `json:"offset"`
			UTC         string `json:"utc"`
			CurrentTime string `json:"current_time"`
		} `json:"timezone"`
	}

	if err := json.Unmarshal(body, &ipData); err != nil {
		logWarn("Failed to parse JSON response", "proxy", proxyStr, "error", err)
		return nil
	}

	if !ipData.Success {
		logDebug("IP lookup failed", "proxy", proxyStr)
		return nil
	}

	logSuccess("Working proxy found", "proxy", proxyStr, "ip", ipData.IP, "country", ipData.Country, "city", ipData.City)

	return &ProxyInfo{
		Proxy:       proxyStr,
		TestedAt:    time.Now().UTC().Format(time.RFC3339),
		IP:          ipData.IP,
		Success:     ipData.Success,
		Type:        ipData.Type,
		Continent:   ipData.Continent,
		Country:     ipData.Country,
		CountryCode: ipData.CountryCode,
		Region:      ipData.Region,
		RegionCode:  ipData.RegionCode,
		City:        ipData.City,
		Latitude:    ipData.Latitude,
		Longitude:   ipData.Longitude,
		IsEU:        ipData.IsEU,
		Postal:      ipData.Postal,
		Connection: Connection{
			ASN:    ipData.Connection.ASN,
			Org:    ipData.Connection.Org,
			ISP:    ipData.Connection.ISP,
			Domain: ipData.Connection.Domain,
		},
		Timezone: Timezone{
			ID:          ipData.Timezone.ID,
			Abbr:        ipData.Timezone.Abbr,
			IsDST:       ipData.Timezone.IsDST,
			Offset:      ipData.Timezone.Offset,
			UTC:         ipData.Timezone.UTC,
			CurrentTime: ipData.Timezone.CurrentTime,
		},
	}
}

// saveJSON saves the working proxies to a JSON file
func saveJSON(workingProxies []ProxyInfo) {
	if len(workingProxies) == 0 {
		logInfo("No working proxies to save")
		return
	}
	jsonData, err := json.MarshalIndent(workingProxies, "", "  ")
	if err != nil {
		logError("Failed to marshal JSON", "error", err)
		return
	}
	if err := os.WriteFile("working_proxies.json", jsonData, 0644); err != nil {
		logError("Failed to write JSON file", "error", err)
	} else {
		logInfo("Saved detailed proxy info", "file", "working_proxies.json", "count", len(workingProxies))
	}
}

// checkerMain orchestrates the proxy checking process with concurrent workers
func checkerMain(proxyFile string, timeout int, numWorkers int) {
	proxies := loadProxies(proxyFile)

	// Open/create checked.txt for writing working proxies
	checkedFile, err := os.OpenFile("checked.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logError("Failed to open checked.txt", "error", err)
		os.Exit(1)
	}
	defer checkedFile.Close()

	// Thread-safe storage for results
	var workingProxies []ProxyInfo
	var mu sync.Mutex // protects workingProxies, checkedFile, and counters
	var tested int

	// Channels
	jobs := make(chan string, numWorkers*2) // buffered job queue
	var wg sync.WaitGroup

	// Handle Ctrl+C to save JSON before exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println()
		logWarn("Interrupt received, saving progress...")
		mu.Lock()
		saveJSON(workingProxies)
		mu.Unlock()
		os.Exit(0)
	}()

	// Spawn workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for proxyStr := range jobs {
				// Test the proxy
				proxyInfo := testProxy(proxyStr, timeout)

				// Update shared state safely
				mu.Lock()
				tested++
				current := tested
				working := len(workingProxies)

				if proxyInfo != nil {
					checkedFile.WriteString(proxyInfo.Proxy + "\n")
					workingProxies = append(workingProxies, *proxyInfo)
					working = len(workingProxies)
					saveJSON(workingProxies)
				}

				if !verbose {
					printProgress(current, len(proxies), working, proxyStr)
				}
				mu.Unlock()
			}
		}(i)
	}

	// Send all proxies to the job queue
	for _, proxy := range proxies {
		jobs <- proxy
	}
	close(jobs) // signal no more jobs

	// Wait for all workers to finish
	wg.Wait()

	if !verbose {
		fmt.Println()
	}
	saveJSON(workingProxies)
	logSuccess("Check complete", "working", len(workingProxies), "total", len(proxies))
}

func worker(id int, jobs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		fmt.Printf("Worker %d processing: %s\n", id, job)
		time.Sleep(500 * time.Millisecond) // simulate work
	}
}
func printBanner() {
	banner := []string{
		" ███████████  ███████████      ███████    █████ █████ █████ █████    █████   ███   █████ █████ ███████████",
		"░░███░░░░░███░░███░░░░░███   ███░░░░░███ ░░███ ░░███ ░░███ ░░███    ░░███   ░███  ░░███ ░░███ ░█░░░░░░███ ",
		" ░███    ░███ ░███    ░███  ███     ░░███ ░░███ ███   ░░███ ███      ░███   ░███   ░███  ░███ ░     ███░  ",
		" ░██████████  ░██████████  ░███      ░███  ░░█████     ░░█████       ░███   ░███   ░███  ░███      ███    ",
		" ░███░░░░░░   ░███░░░░░███ ░███      ░███   ███░███     ░░███        ░░███  █████  ███   ░███     ███     ",
		" ░███         ░███    ░███ ░░███     ███   ███ ░░███     ░███         ░░░█████░█████░    ░███   ████     █",
		" █████        █████   █████ ░░░███████░   █████ █████    █████          ░░███ ░░███      █████ ███████████",
		"░░░░░        ░░░░░   ░░░░░    ░░░░░░░    ░░░░░ ░░░░░    ░░░░░            ░░░   ░░░      ░░░░░ ░░░░░░░░░░░",
	}

	// Color gradient top to bottom
	colors := []string{Magenta, Magenta, Cyan, Cyan, Cyan, Blue, Blue, Dim}

	fmt.Println()
	for i, line := range banner {
		fmt.Printf("%s%s%s\n", colors[i], line, Reset)
	}

	fmt.Println()
	fmt.Printf("%s%s   🧙 The Fastest HTTP/S Proxy Checker/Scraper 🧙%s\n", Green+Bold, strings.Repeat(" ", 15), Reset)
	fmt.Printf("%s%s@Syn2Much%s\n\n", Dim, strings.Repeat(" ", 50), Reset)
}

func main() {
	var numWorkers int
	var timeout int
	var verboseInput string
	printBanner()

	var fileName string
	fmt.Printf("%sEnter proxy file name (leave empty to scrape from default sources):%s ", Cyan, Reset)
	fmt.Scanln(&fileName)

	fmt.Printf("%sAmount of Workers:%s ", Cyan, Reset)
	fmt.Scanln(&numWorkers)
	if numWorkers < 1 {
		numWorkers = 10 // sensible default
	}

	fmt.Printf("%sTimeout (seconds):%s ", Cyan, Reset)
	fmt.Scanln(&timeout)
	if timeout < 1 {
		timeout = 8 // sensible default
	}

	fmt.Printf("%sVerbose mode? (y/N):%s ", Cyan, Reset)
	fmt.Scanln(&verboseInput)
	verbose = strings.ToLower(verboseInput) == "y" || strings.ToLower(verboseInput) == "yes"

	// Pass numWorkers to checkerMain
	checkerMain(fileName, timeout, numWorkers)

	fmt.Println("All done")
}
