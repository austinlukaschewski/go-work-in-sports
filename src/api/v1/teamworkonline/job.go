package teamworkonline

import (
	"fmt"
	"strings"

	"github.com/austinlukaschewski/go-work-in-sports/internal/model"
	"github.com/austinlukaschewski/go-work-in-sports/internal/xlog"
	"github.com/gocolly/colly/v2"
)

func GetJobs(keywords []string, limit int) (jobs []model.Job, err error) {
	c := colly.NewCollector()

	c.OnHTML("div.browse-jobs-card__content", func(content *colly.HTMLElement) {
		var job model.Job
		a := content.DOM.Find("a.browse-jobs-card__content--title")
		for _, word := range keywords {
			if strings.Contains(strings.ToLower(a.Text()), strings.ToLower(word)) {
				job.Title = a.Text()
				job.Link = fmt.Sprintf("https://teamworkonline.com%s", a.AttrOr("href", ""))
				job.Org = content.DOM.Find("div.browse-jobs-card__content--organization").Text()

				timeDiv := content.DOM.Find("div.browse-jobs-card__content--bottom_row > div.browse-jobs-card__content--time")
				timeInfo := timeDiv.Find("div.browse-jobs-card__scoreboard").Text()
				timeDuration := timeDiv.Find("div.trending__content--small").Text()

				job.SincePostedDate = fmt.Sprintf("%s %s", strings.TrimLeft(timeInfo, "0"), timeDuration)

				jobs = append(jobs, job)
				break
			}
		}
	})

	for i := 1; i <= limit; i++ {
		if err := c.Visit(fmt.Sprintf("https://www.teamworkonline.com/jobs-in-sports?page=%d", i)); err != nil {
			xlog.Error("Error getting teamworkonline jobs", "error", err)
			return jobs, err
		}
	}

	return jobs, nil
}
