package fangraphs

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/austinlukaschewski/go-work-in-sports/internal/model"
	"github.com/austinlukaschewski/go-work-in-sports/internal/xlog"
	"github.com/gocolly/colly/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var titleCaser cases.Caser = cases.Title(language.AmericanEnglish)

func GetJobs(keywords []string, limit int) (jobs []model.Job, err error) {
	c := colly.NewCollector()

	c.OnHTML("div.post", func(post *colly.HTMLElement) {
		if len(jobs) >= limit {
			return
		}

		var postKeywordSlice []string

		titleAnchor := post.DOM.Find("h2.instagraphstitle > a")
		postKeywordSlice = append(postKeywordSlice, strings.ToLower(titleAnchor.Text()))

		post.DOM.Find("div.fullpostentry > p:nth-child(2) > a").Each(func(i int, a *goquery.Selection) {
			postKeywordSlice = append(postKeywordSlice, strings.ToLower(a.Text()))
		})

		var keyword string
		passes := slices.ContainsFunc(postKeywordSlice, func(e string) bool {
			return slices.ContainsFunc(keywords, func(f string) bool {
				keyword = e
				return strings.Contains(e, f)
			})
		})

		if !passes {
			return
		}

		j := model.Job{
			Title: titleCaser.String(keyword),
			Org:   strings.TrimPrefix(strings.Split(titleAnchor.Text(), "–")[0], "Job Posting: "),
			Link:  titleAnchor.AttrOr("href", ""),
		}

		datePostedString := post.DOM.Find("div.postmeta > div:nth-child(2)").Text()
		date, err := time.Parse("January 2, 2006", datePostedString)
		if err == nil {
			j.Date = datePostedString
			duration := time.Since(date)
			j.SincePostedDate = fmt.Sprintf("%.0f days ago", duration.Hours()/24)
		}

		duration := time.Since(date)
		if duration.Hours() < 24 {
			j.SincePostedDate = fmt.Sprintf("%.0f hours ago", duration.Hours())
		} else {
			j.SincePostedDate = fmt.Sprintf("%.0f days ago", duration.Hours()/24)
		}

		jobs = append(jobs, j)
	})

	if err := c.Visit("https://blogs.fangraphs.com/category/job-postings/"); err != nil {
		xlog.Error("Error visiting Fangraphs job postings", "error", err)
		return jobs, err
	}

	return jobs, err
}
