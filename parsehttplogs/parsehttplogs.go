package main

/*
Parses logs for corral.com and all the other domains since they use the same logs.

This is meant to be run once a day. I run it shortly after midnight. If run more than once you'll get duplicate data

Since log rotation happens by file size not date I look at current log and anything that was rotated today or yesterday.
I look through those files and grab everything that has yesterdays date. This works because traffic to my site is so small.

After I parse everything I get the country code from the IP address so I can map that info for the traffic map.
*/

import (
	// Standard packages
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	// My packages
	"common"
)

// LocationCodes - structure used for parsing GEO location json
type LocationCodes struct {
	GeoCountryCode string `json:"countryCode"`
	GeoRegionCode  string `json:"region"`
}

// LocationCodesFromIP - used for list of IPs i've already got GEO info for
type LocationCodesFromIP struct {
	CountryCode string
	RegionCode  string
}

func sanitizeRequest(originalRequest string) string {
	// sanitizeRequest - I think I relplace the single quotes to solve some MySQL issue. I don't remember
	// what breaks when if I leave this out. Truncating is just because of column size

	request := strings.Replace(originalRequest, "'", "\\'", -1)
	if len(request) > 250 {
		tempRunes := []rune(request)
		request = string(tempRunes[0:250])
	}
	return request
}

func LogFileList() []string {
	// LogFileList - returns list of log files to parse through

	yyyymmddFmt := "20060102"
	todaysDate := time.Now().Format(yyyymmddFmt)
	yesterdaysDate := time.Now().AddDate(0, 0, -1).Format(yyyymmddFmt)

	currentLogFile := "/var/log/httpd/access_log"
	logFileRotatedToday := currentLogFile + "-" + todaysDate
	logFileRotatedYesterday := currentLogFile + "-" + yesterdaysDate

	// always return current log file
	fileList := []string{currentLogFile}

	// if log file rotated today add that to list
	if _, err := os.Stat(logFileRotatedToday); err == nil {
		fileList = append(fileList, logFileRotatedToday)
	}
	// if log file rotated yesterday add that to list
	if _, err := os.Stat(logFileRotatedYesterday); err == nil {
		fileList = append(fileList, logFileRotatedYesterday)
	}
	return fileList
}

func RegionCodes(ipAddress string) (string, string) {
	// Calls an API that returns GEO info about an IP address. I only use the country and region codes

	geoURL := "http://ip-api.com/json/" + ipAddress

	geoClient := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}

	req, err := http.NewRequest(http.MethodGet, geoURL, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr := geoClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	geoCodes := LocationCodes{}
	jsonErr := json.Unmarshal(body, &geoCodes)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	if geoCodes.GeoRegionCode == "" {
		geoCodes.GeoRegionCode = "Unknown"
	}

	if geoCodes.GeoCountryCode == "" {
		geoCodes.GeoCountryCode = "Unknown"
	}
	// Total hack - put here because public API will block me if I spam them.
	// Since this is a batch job that runs once a day I'm OK with it.
	time.Sleep(1 * time.Second)

	return geoCodes.GeoCountryCode, geoCodes.GeoRegionCode
}

func main() {

	geoFromIP := make(map[string]LocationCodesFromIP)

	// Build up regex to get yesterdays data from and get submatches
	// sub matches: 1 - ip, 2 - time, 3 - request, 4 - result
	ddmonyyyyFmt := "02/Jan/2006"
	yesterdaysDate := time.Now().AddDate(0, 0, -1).Format(ddmonyyyyFmt)
	logLineRegExp := "([(\\d\\.)]+) - - \\[(" + yesterdaysDate + ".*?) \\+0000\\] \"GET (.*?) HTTP/1.\\d\" (\\d+)"
	r, _ := regexp.Compile(logLineRegExp)

	// excludeList - Excludes things I don't care about or traffic where I'm not interested in
	excludeList := ".jpg |.gif |.png |.css |.js |.ico |robots.txt|searchme.html|54.172.87.69|76.126.34.55"

	logFiles := LogFileList()

	db := common.OpenDB()
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO logdatatest (ip,requestDate,request,responseCode,CountryCode,regionCode) VALUES (?, STR_TO_DATE(?,'%d/%b/%Y:%T'), ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	for _, logFile := range logFiles {

		fileP, err := os.Open(logFile)
		if err != nil {
			log.Fatal(err)
		}
		defer fileP.Close()

		singleLogLine := bufio.NewScanner(fileP)
		if err := singleLogLine.Err(); err != nil {
			log.Fatal(err)
		}
		for singleLogLine.Scan() {
			if logLineValues := r.FindStringSubmatch(singleLogLine.Text()); len(logLineValues) != 0 {

				r2, _ := regexp.Compile(excludeList)
				if excludedValues := r2.FindStringSubmatch(singleLogLine.Text()); len(excludedValues) != 0 {
					continue
				}
				// If I already have GEO info for an IP I don't need to get it again
				if _, keyExists := geoFromIP[logLineValues[1]]; !keyExists {
					countryCode, regionCode := RegionCodes(logLineValues[1])
					geoFromIP[logLineValues[1]] = LocationCodesFromIP{countryCode, regionCode}
				}
				countryCode := geoFromIP[logLineValues[1]].CountryCode
				regionCode := geoFromIP[logLineValues[1]].RegionCode
				request := sanitizeRequest(logLineValues[3])
				_, err := stmt.Exec(logLineValues[1], logLineValues[2], request, logLineValues[4], countryCode, regionCode)
				if err != nil {
					log.Fatal(err)
				}
			}
		}

	}
}
