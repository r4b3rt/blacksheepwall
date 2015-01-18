package bsw

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// YandexAPI uses Yandex XML API and the 'rhost' search operator to find subdomains of a
// given domain.
func YandexAPI(domain, apiURL, serverAddr string) (string, Results, error) {
	task := "Yandex API"
	results := Results{}
	xmlTemplate := "<?xml version='1.0' encoding='UTF-8'?><request><query>%s</query><sortby>rlv</sortby><maxpassages>1</maxpassages><page>0</page><groupings><groupby attr=\" \" mode=\"flat\" groups-on-page=\"100\" docs-in-group=\"1\" /></groupings></request>"
	parts := strings.Split(domain, ".")
	var query = "rhost:" + parts[1] + "." + parts[0] + ".*"
	postBody := fmt.Sprintf(xmlTemplate, query)
	resp, err := http.Post(apiURL, "text/xml", strings.NewReader(postBody))
	if err != nil {
		return task, results, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return task, results, nil
	}
	domainSet := make(map[string]bool)
	doc.Find("domain").Each(func(_ int, s *goquery.Selection) {
		domain := s.Text()
		if domainSet[domain] {
			return
		}
		ip, err := LookupName(domain, serverAddr)
		if err == nil {
			results = append(results, Result{Source: task, IP: ip, Hostname: domain})
		}
	})
	return task, results, nil
}
