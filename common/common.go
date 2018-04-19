package common

// Contains functions that are used in multiple programs or thought I was going to use in multiple programs

// Standard packages
import (
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
)

// Downloaded packages
import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func DisplayableName(name string, modifier int) string {
	// DisplayableName - Handles the hacky way many years go I dealt with how artist name is stored in the DB vs
	// how it is displayed and haven't come up with something better so I just rewrite this in
	// every language I use.
	//
	// 0 - leave alone unless there is a comma then swap first 2 words around comma, leave rest alone
	//        "Cale, John / Lou Reed / Nico" -> "John Cale / Lou Reed / Nico"
	//        "Ridgway, Stan" -> "Stan Ridgway"
	//        "Wire" -> "Wire"
	// 1 - leave alone even though there is a comma in the name
	//        "Not Drowning, Waving" -> "Not Drowning, Waving"
	// 2 - Deals with someone having a middle initial
	//        "Bach, Johann Sebastian" -> "Johann Sebastian Bach"
	// 4 - add "The " at beginning
	//        "Velvet Underground" -> "The Velvet Underground"

	if modifier == 4 {
		return fmt.Sprintf("The %s", name)
	} else if modifier == 1 {
		return name
	} else if strings.Contains(name, ",") && (strings.Count(name, " ") == 1 || modifier == 2) {
		n := strings.Split(name, ", ")
		return fmt.Sprintf("%s %s", n[1], n[0])
	} else if strings.Contains(name, ",") && strings.Count(name, " ") > 1 {
		n := strings.Split(name, ", ")
		r := strings.SplitN(n[1], " ", 2)
		return fmt.Sprintf("%s %s %s", r[0], n[0], r[1])
	}
	return name
}

func IsAlphaNumeric(s string) bool {
	// IsAlphaNumeric - Just sanitizing user entered string

	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			return false
		}
	}
	return true
}

func TruncateString(str string, maxLen int) string {
	// TruncateString - truncates artist name so hopefully it's not wider than the cover image
	// maxLen is usually 11

	shortString := str
	if len(str) > maxLen {
		shortString = str[0:maxLen-3] + "..."
	}
	return shortString
}

func DBCredentials() (string, string, string) {
	// DBCredentials - returns db credentials from environment

	dbname := os.Getenv("DBNAME")
	if len(dbname) == 0 {
		log.Fatal("DBNAME not defined")
	}
	dbpass := os.Getenv("DBPASS")
	if len(dbpass) == 0 {
		log.Fatal("DBPASS not defined")
	}
	dbuser := os.Getenv("DBUSER")
	if len(dbuser) == 0 {
		log.Fatal("DBUSER not defined")
	}
	return dbname, dbpass, dbuser
}

func OpenDB() *sql.DB {
	// OpenDB - Opens a connection to the DB and returns pointer to it

	dbname, dbpass, dbuser := DBCredentials()

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		dbuser, dbpass, "127.0.0.1", "3306", dbname))

	if err != nil {
		log.Fatal(err)
	}
	return db
}
