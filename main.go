package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

type Statistics struct {
	CacheDirectory           string    `json:"cache_directory"`
	PrimaryConfig            string    `json:"primary_config"`
	SecondaryConfigReadonly  string    `json:"secondary_config_readonly"`
	StatsTime                time.Time `json:"stats_time"`
	StatsZeroTime            time.Time `json:"stats_zero_time"`
	CacheHitDirect           int       `json:"cache_hit_direct"`
	CacheHitPreprocessed     int       `json:"cache_hit_preprocessed"`
	CacheMiss                int       `json:"cache_miss"`
	CacheHitRate             float64   `json:"cache_hit_rate"`
	CalledForLink            int       `json:"called_for_link"`
	CalledForPreprocessing   int       `json:"called_for_preprocessing"`
	UnsupportedCodeDirective int       `json:"unsupported_code_directive"`
	NoInputFile              int       `json:"no_input_file"`
	CleanupsPerformed        int       `json:"cleanups_performed"`
	FilesInCache             int       `json:"files_in_cache"`
	CacheSize                string    `json:"cache_size"`
	CacheSizeBytes           int       `json:"cache_size_bytes"`
	MaxCacheSize             string    `json:"max_cache_size"`
	MaxCacheSizeBytes        int       `json:"max_cache_size_bytes"`
}

var rules = map[string]*regexp.Regexp{
	"cacheDirectory":           regexp.MustCompile(`cache directory\s+(.+)`),
	"primaryConfig":            regexp.MustCompile(`primary config\s+(.+)`),
	"secondaryConfigReadonly":  regexp.MustCompile(`secondary config\s+(\(readonly\)\s+)?(.+)`),
	"statsZeroTime":            regexp.MustCompile(`stats zero time\s+(.*)`),
	"cacheHitDirect":           regexp.MustCompile(`cache hit \(direct\)\s+(\d+)`),
	"cacheHitPreprocessed":     regexp.MustCompile(`cache hit \(preprocessed\)\s+(\d+)`),
	"cacheMiss":                regexp.MustCompile(`cache miss\s+(\d+)`),
	"cacheHitRate":             regexp.MustCompile(`cache hit rate\s+(\d+(\.\d+)?) %`),
	"calledForLink":            regexp.MustCompile(`called for link\s+(\d+)`),
	"calledForPreprocessing":   regexp.MustCompile(`called for preprocessing\s+(\d+)`),
	"unsupportedCodeDirective": regexp.MustCompile(`unsupported code directive\s+(\d+)`),
	"noInputFile":              regexp.MustCompile(`no input file\s+(\d+)`),
	"cleanupsPerformed":        regexp.MustCompile(`cleanups performed\s+(\d+)`),
	"filesInCache":             regexp.MustCompile(`files in cache\s+(\d+)`),
	"cacheSize":                regexp.MustCompile(`cache size\s+(.+)`),
	"maxCacheSize":             regexp.MustCompile(`max cache size\s+(.+)`),
}

func (s *Statistics) Parse(text string) {
	s.CacheDirectory = rules["cacheDirectory"].FindStringSubmatch(text)[1]
	s.PrimaryConfig = rules["primaryConfig"].FindStringSubmatch(text)[1]
	s.SecondaryConfigReadonly = rules["secondaryConfigReadonly"].FindStringSubmatch(text)[2]

	// now's the time
	s.StatsTime = time.Now()

	// assume stats originate from the local host
	statsZeroTime := rules["statsZeroTime"].FindStringSubmatch(text)[1]
	s.StatsZeroTime, _ = time.ParseInLocation("Mon Jan 2 15:04:05 2006", statsZeroTime, s.StatsTime.Location())

	s.CacheHitDirect, _ = strconv.Atoi(rules["cacheHitDirect"].FindStringSubmatch(text)[1])

	s.CacheHitPreprocessed, _ = strconv.Atoi(rules["cacheHitPreprocessed"].FindStringSubmatch(text)[1])
	s.CacheMiss, _ = strconv.Atoi(rules["cacheMiss"].FindStringSubmatch(text)[1])
	s.CacheHitRate, _ = strconv.ParseFloat(rules["cacheHitRate"].FindStringSubmatch(text)[1], 64)
	s.CalledForLink, _ = strconv.Atoi(rules["calledForLink"].FindStringSubmatch(text)[1])
	s.CalledForPreprocessing, _ = strconv.Atoi(rules["calledForPreprocessing"].FindStringSubmatch(text)[1])
	s.UnsupportedCodeDirective, _ = strconv.Atoi(rules["unsupportedCodeDirective"].FindStringSubmatch(text)[1])
	s.NoInputFile, _ = strconv.Atoi(rules["noInputFile"].FindStringSubmatch(text)[1])
	s.CleanupsPerformed, _ = strconv.Atoi(rules["cleanupsPerformed"].FindStringSubmatch(text)[1])
	s.FilesInCache, _ = strconv.Atoi(rules["filesInCache"].FindStringSubmatch(text)[1])

	s.CacheSize = rules["cacheSize"].FindStringSubmatch(text)[1]
	s.MaxCacheSize = rules["maxCacheSize"].FindStringSubmatch(text)[1]
}

func main() {
	stat, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if stat.Mode()&os.ModeNamedPipe == 0 {
		// TODO add flags, read from stdin / file(s)
		// TODO add help
		panic("No data piped to stdin")
	}

	var text string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text += scanner.Text() + "\n"
	}
	stats := Statistics{}
	stats.Parse(text)
	statsJson, _ := json.Marshal(stats)
	fmt.Println(string(statsJson))
}
