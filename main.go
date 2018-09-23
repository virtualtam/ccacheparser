package main

import (
	"encoding/json"
	"fmt"
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

const text = `
cache directory                     /home/virtualtam/.ccache
primary config                      /home/virtualtam/.ccache/ccache.conf
secondary config      (readonly)    /etc/ccache.conf
stats zero time                     Sun Sep 23 01:18:52 2018
cache hit (direct)                    73
cache hit (preprocessed)               4
cache miss                           207
cache hit rate                     27.11 %
called for link                       28
called for preprocessing             170
unsupported code directive             4
no input file                         58
cleanups performed                     0
files in cache                       639
cache size                          12.1 MB
max cache size                      15.0 GB
`

func main() {
	stats := Statistics{}
	stats.Parse(text)
	statsJson, _ := json.Marshal(stats)
	fmt.Println(string(statsJson))
}
