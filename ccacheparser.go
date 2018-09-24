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
	matches := rules["cacheDirectory"].FindStringSubmatch(text)
	if len(matches) == 2 {
		s.CacheDirectory = matches[1]
	}

	matches = rules["primaryConfig"].FindStringSubmatch(text)
	if len(matches) == 2 {
		s.PrimaryConfig = matches[1]
	}

	matches = rules["secondaryConfigReadonly"].FindStringSubmatch(text)
	if len(matches) == 2 {
		s.SecondaryConfigReadonly = matches[1]
	} else if len(matches) == 3 {
		s.SecondaryConfigReadonly = matches[2]
	}

	// now's the time
	s.StatsTime = time.Now()

	// assume stats originate from the local host
	statsZeroTime := rules["statsZeroTime"].FindStringSubmatch(text)[1]

	s.StatsZeroTime, _ = time.ParseInLocation("Mon Jan 2 15:04:05 2006", statsZeroTime, s.StatsTime.Location())

	matches = rules["cacheHitDirect"].FindStringSubmatch(text)
	if len(matches) == 2 {
		s.CacheHitDirect, _ = strconv.Atoi(matches[1])
	}

	matches = rules["cacheHitPreprocessed"].FindStringSubmatch(text)
	if len(matches) == 2 {
		s.CacheHitPreprocessed, _ = strconv.Atoi(matches[1])
	}

	matches = rules["cacheMiss"].FindStringSubmatch(text)
	if len(matches) == 2 {
		s.CacheMiss, _ = strconv.Atoi(matches[1])
	}

	matches = rules["cacheHitRate"].FindStringSubmatch(text)
	if len(matches) == 3 {
		s.CacheHitRate, _ = strconv.ParseFloat(matches[1], 64)
	}

	matches = rules["calledForLink"].FindStringSubmatch(text)
	if len(matches) == 2 {
		s.CalledForLink, _ = strconv.Atoi(matches[1])
	}

	matches = rules["calledForPreprocessing"].FindStringSubmatch(text)
	if len(matches) == 2 {
		s.CalledForPreprocessing, _ = strconv.Atoi(matches[1])
	}

	matches = rules["unsupportedCodeDirective"].FindStringSubmatch(text)
	if len(matches) == 2 {
		s.UnsupportedCodeDirective, _ = strconv.Atoi(matches[1])
	}

	matches = rules["noInputFile"].FindStringSubmatch(text)
	if len(matches) == 2 {
		s.NoInputFile, _ = strconv.Atoi(matches[1])
	}

	matches = rules["cleanupsPerformed"].FindStringSubmatch(text)
	if len(matches) == 2 {
		s.CleanupsPerformed, _ = strconv.Atoi(matches[1])
	}

	matches = rules["filesInCache"].FindStringSubmatch(text)
	if len(matches) == 2 {
		s.FilesInCache, _ = strconv.Atoi(matches[1])
	}

	matches = rules["cacheSize"].FindStringSubmatch(text)
	if len(matches) == 2 {
		s.CacheSize = matches[1]
	}

	matches = rules["maxCacheSize"].FindStringSubmatch(text)
	if len(matches) == 2 {
		s.MaxCacheSize = matches[1]
	}
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
