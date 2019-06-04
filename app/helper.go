package app

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/likexian/whois-go"
	"github.com/orestrepov/metadatahost/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"regexp"
)

// HTTP client to consume API services
func HTTPClient(APIURL string) ([]byte, error) {

	req, err := http.NewRequest(http.MethodGet, APIURL, nil)
	if err != nil {
		return nil, err
	}
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Query to recovery the whois ip data using whois-go library
func WhoisQuery(ip string) (*model.Whois, error) {

	var whoisModel model.Whois

	// Get the whois content for an ip
	whoisRaw, err := whois.Whois(ip)
	if err != nil {
		return nil, err
	}

	// Get the ip country value using reg exp
	rCountryLine, _ := regexp.Compile("[Cc]ountry.*:\\s*(.+)")
	sCountryLine := rCountryLine.FindString(whoisRaw)
	rCountry, _ := regexp.Compile("[Cc]ountry.*:\\s*")
	sCountry := rCountry.ReplaceAllString(sCountryLine, "")
	whoisModel.Country = sCountry

	// Get the ip owner value using reg exp
	rOrgNameLine, _ := regexp.Compile("[Oo]rgName.*:\\s*(.+)")
	sOrgNameLine := rOrgNameLine.FindString(whoisRaw)
	rOrgName, _ := regexp.Compile("[Oo]rgName.*:\\s*")
	sOrgName := rOrgName.ReplaceAllString(sOrgNameLine, "")
	whoisModel.OrgName = sOrgName

	if sCountry == "" && sOrgName == "" {
		return nil, errors.New("can't get the who information from the ip")
	}

	return &whoisModel, nil
}

// Scraping a Web page to get the title and logo image URL
func ScrapingWebPage(URL string) (*model.WebPage, error) {

	var webPage model.WebPage
	webPage.Title = ""
	webPage.Logo = ""

	// Request the HTML page
	res, err := http.Get("http://" + URL)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		logrus.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		logrus.Error(err)
	}

	doc.Find("head").Each(func(index int, item *goquery.Selection) {
		// Recovery page Title
		webPage.Title = item.Find("title").Contents().Text()
		// Recovery the page Logo image url
		item.Find("link").Each(func(index int, item *goquery.Selection) {
			linkTag := item
			link, _ := linkTag.Attr("href")
			rLogo, _ := regexp.Compile("(http.*\\.png)")
			if rLogo.MatchString(link) && webPage.Logo == "" {
				webPage.Logo = link
			}
		})
	})

	return &webPage, nil
}
