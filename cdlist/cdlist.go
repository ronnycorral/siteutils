package main

// Creates the files that contain the html that displays 3 CD covers based on the type of list
// It's either 3 random covers or the 3 CDs that have the latest date_added

import (
	// standard packages
	"bytes"
	"flag"
	"fmt"
	"html"
	"log"
	"os"

	// packages written by me
	"common"
)

func main() {

	var (
		artistid     int
		cdartist     string
		namemodifier int
		coverfile    string
		cdtitle      string
	)
	var htmlCode bytes.Buffer
	var orderBy string
	var listFile string

	domainPtr := flag.String("domain", "", "domain CD list is for")
	listTypePtr := flag.String("type", "", "type of CD list (random or latest)")
	flag.Parse()

	domain := *domainPtr
	listType := *listTypePtr

	if len(domain) == 0 {
		flag.Usage()
		log.Fatal("ERROR: A domain name is required")
	}

	// Just sanitizing the data since it is user input
	if common.IsAlphaNumeric(domain) == false {
		flag.Usage()
		log.Fatal("ERROR: I only buy domains with numbers and letters")
	}

	if _, err := os.Stat(fmt.Sprintf("/var/www/cdfiles/%s", domain)); os.IsNotExist(err) {
		flag.Usage()
		log.Fatal("ERROR: I don't use that domain:", domain)
	}

	switch listType {
	case "random":
		orderBy = "rand()"
		listFile = "randcd.txt"
	case "latest":
		orderBy = "date_added desc"
		listFile = "whatsnewcd.txt"
	default:
		flag.Usage()
		log.Fatal("ERROR: type either missing or not valid type: ", listType)

	}

	sql := fmt.Sprintf("SELECT distinct cdtitle.artist_id as artistid, cdartist.artist as cdartist, cdartist.name_action as namemodifier, cdtitle.cover_file as coverfile, cdtitle.title as cdtitle FROM cdtitle right join cdartist on cdartist.artist_id = cdtitle.artist_id where cdartist.artist != 'Compilations' and cdartist.artist != 'Soundtracks' order by %s limit 3", orderBy)

	db := common.OpenDB()
	defer db.Close()

	rows, err := db.Query(sql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&artistid, &cdartist, &namemodifier, &coverfile, &cdtitle)
		if err != nil {
			log.Fatal(err)
		}
		print_this_name := html.EscapeString(common.TruncateString(common.DisplayableName(cdartist, namemodifier), 11))
		print_this_title := html.EscapeString(cdtitle)
		path_to_cover := html.EscapeString(fmt.Sprintf("/images/cdcovers/%s/%s", coverfile[0:1], coverfile))
		switch listType {
		case "random":
			htmlCode.WriteString(fmt.Sprintf("<li><a href=\"/cdcoll.php?artist_id=%d\"><img class=\"roundcorner\" src=\"%s\" height=\"64\" width=\"64\" alt=\"%s\" title=\"%s\"></a><p>%s</p></li>\n", artistid, path_to_cover, print_this_title, print_this_title, print_this_name))
		case "latest":
			htmlCode.WriteString(fmt.Sprintf("<li><img class=\"roundcorner imagepointer\" src=\"%s\" height=\"64\" width=\"64\" alt=\"%s\" title=\"%s\" onclick=\"getArtistData(%d,'%s')\"><p>%s</p></li>\n", path_to_cover, print_this_title, print_this_title, artistid, domain, print_this_name))
		}
	}
	// This section needs to be updated to include error checking and first writing to temp file so I don't blow
	// away current good file with possible currupted data.
	f, err := os.Create(fmt.Sprintf("/var/www/cdfiles/%s/%s", domain, listFile))
	_, err = f.WriteString(htmlCode.String())
	f.Sync()
}
