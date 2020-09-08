package hw03_frequency_analysis //nolint:golint,stylecheck

import (
	"regexp"
	"sort"
	"strings"
)

var removePunctuationRegexp = regexp.MustCompile(`(?i)[^a-zа-яё0-9\- ]`)

func splitByWords(text string) []string {
	text = strings.ToLower(text)
	text = removePunctuationRegexp.ReplaceAllString(text, " ")
	words := strings.Fields(text)
	return words
}

// WordStat хранит, сколько раз встречается слов.
type WordStat struct {
	word  string
	count int
}

// WordStatList возвращает слайс для хранения статистики слов.
type WordStatList []WordStat

func (w WordStatList) sortDesc() {
	sort.SliceStable(w, func(i, j int) bool {
		return w[i].count > w[j].count
	})
}

func (w WordStatList) topN(n int) []string {
	w.sortDesc()

	result := make([]string, 0, n)
	for _, ws := range w {
		result = append(result, ws.word)
		if len(result) == n {
			return result
		}
	}
	return result
}

func newWordStatList(stat map[string]int) WordStatList {
	wordStatList := make(WordStatList, 0, len(stat))
	for word, count := range stat {
		ws := WordStat{
			word:  word,
			count: count,
		}
		wordStatList = append(wordStatList, ws)
	}
	return wordStatList
}

func Top10(text string) []string {
	words := splitByWords(text)
	stat := make(map[string]int)
	for _, word := range words {
		if word == "-" {
			continue
		}
		stat[word]++
	}

	wordStatList := newWordStatList(stat)
	top := wordStatList.topN(10)
	// fmt.Printf("-> %#v\n", top)
	return top
}
