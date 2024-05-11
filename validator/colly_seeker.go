package validator

import (
	"errors"
	"net/http"
	"time"

	"github.com/gocolly/colly/v2"
)

type Seeker struct {
	Config   CheckConfig
	InputURL string
	Client   *colly.Collector
	Logs     *SeekerLogs
}

type SeekerLogs struct {
	StartAt           time.Time
	FirstByteDuration time.Duration
	Duration          time.Duration
	ContentLength     int64
	Content           *[]byte
	Redirects         *[]http.Request
	ContentType       string
	StatusCode        int
}

func NewSeeker(config CheckConfig, collector *colly.Collector) (seeker *Seeker) {

	seeker = &Seeker{
		Config: config,
	}

	if collector != nil {
		seeker.Client = collector
	} else {
		seeker.Client = colly.NewCollector(
			colly.MaxDepth(seeker.Config.MaxRedirects),
			colly.MaxBodySize(int(seeker.Config.MaxSize)),
			colly.IgnoreRobotsTxt(),
			colly.AllowURLRevisit(),
			colly.Async(true),
			colly.UserAgent("RouteHub-Link-Validator"),
			colly.TraceHTTP(),
		)
	}

	seeker.Logs = &SeekerLogs{
		Redirects: &[]http.Request{},
	}

	seeker.Client.SetRequestTimeout(seeker.Config.HTTPClientTimeout)

	seeker.Client.SetRedirectHandler(func(req *http.Request, via []*http.Request) error {
		if len(via) >= seeker.Config.MaxRedirects {
			return errors.New(CheckErrorMessages[ErrMaxRedirectsExceeded])
		}

		redirects := append(*seeker.Logs.Redirects, *req)

		seeker.Logs.Redirects = &redirects
		return nil
	})

	seeker.Client.OnRequest(func(r *colly.Request) {
		seeker.Logs.StartAt = time.Now()
	})

	seeker.Client.OnResponse(func(r *colly.Response) {
		seeker.Logs.FirstByteDuration = r.Trace.FirstByteDuration
		seeker.Logs.Duration = r.Trace.ConnectDuration
		seeker.Logs.ContentLength = int64(len(r.Body))
		seeker.Logs.Content = &r.Body
		seeker.Logs.ContentType = r.Headers.Get("Content-Type")
		seeker.Logs.StatusCode = r.StatusCode
	})

	return seeker
}

func (s Seeker) GetCollector() *colly.Collector {
	return s.Client
}

func (s Seeker) Seek(inputURL string) (err error) {
	s.InputURL = inputURL
	err = s.Client.Visit(s.InputURL)
	s.Client.Wait()

	if s.Logs.StatusCode != 200 {
		err = errors.New(CheckErrorMessages[ErrUnreachable])
	}
	return err
}
