package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/patrickhener/go-bhtool/db"
)

const ver string = "v0.0.2"

var (
	uri      string
	user     string
	pass     string
	domain   string
	userlist string
	tls      bool
)

var generalHelp = `
  Usage: go-bhtool [--neo4j-uri] [--neo4j-user] [--neo4j-pass] [--tls] [command] [--help]

  Version: ` + ver + ` (` + runtime.Version() + `)

  Defaults:
    neo4j-uri: 	bolt://localhost:7687
    neo4j-user:	neo4j
    neo4j-pass:	admin
    tls:	false

  Commands:
    own		- mark multiple users as owned
    owned	- get a list of owned users

  Read more:
    https://github.com/patrickhener/go-bhtool

`

func main() {
	// Flags
	flag.StringVar(&uri, "neo4j-uri", "bolt://localhost:7687", "")
	flag.StringVar(&user, "neo4j-user", "neo4j", "")
	flag.StringVar(&pass, "neo4j-pass", "neo4j", "")
	flag.BoolVar(&tls, "tls", false, "")

	version := flag.Bool("version", false, "")
	v := flag.Bool("v", false, "")
	help := flag.Bool("help", false, "")
	h := flag.Bool("h", false, "")
	flag.Usage = func() {}
	flag.Parse()

	if len(os.Args) <= 1 {
		fmt.Print(generalHelp)
		os.Exit(0)
	}

	if *version || *v {
		fmt.Println(ver)
		os.Exit(0)
	}

	if *help || *h {
		fmt.Print(generalHelp)
		os.Exit(0)
	}

	args := flag.Args()

	// Now test for db connectivity
	// Connect to neo4j
	neo4jCon := &db.Neo4jDB{}

	if err := neo4jCon.Connect(uri, user, pass, tls); err != nil {
		log.Printf("Error connecting to neo4j instance: %+v", err)
		os.Exit(1)
	}

	subcmd := ""
	if len(args) > 0 {
		subcmd = args[0]
		args = args[1:]
	}

	switch subcmd {
	case "own":
		own(args, neo4jCon)
	case "owned":
		owned(args, neo4jCon)
	default:
		fmt.Print(help)
		os.Exit(0)
	}
}

var commonHelp = `
  Version:
    ` + ver + ` (` + runtime.Version() + `)

  Read more:
    https://github.com/patrickhener/go-bhtool

`

var ownHelp = `
  Usage: go-bhtool own [options] [users...]

  Mark multiple users as owned

  Options:

    --userlist, Path to list of file with users - one per line
    --domain, Domain to add to users where there is no domain

  Examples:

  * Import a list of users and add domain when missing

  go-bhtool own --userlist /path/to/myuserlist.txt --domain contoso.com

  * Import two users without a list

  go-bhtool own user1@contoso.com user2@contoso.com
` + commonHelp

func own(args []string, db *db.Neo4jDB) {
	var usersToAdd []string = make([]string, 0)

	flags := flag.NewFlagSet("own", flag.ContinueOnError)
	flags.StringVar(&userlist, "userlist", "", "")
	flags.StringVar(&domain, "domain", "", "")
	flags.Usage = func() {
		fmt.Print(ownHelp)
		os.Exit(0)
	}
	flags.Parse(args)

	// Read userlist if defined
	if userlist != "" {
		file, err := os.Open(userlist)
		if err != nil {
			log.Printf("Error reading file @ %s: %+v", userlist, err)
			os.Exit(1)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			// Read line from file
			line := scanner.Text()
			// All upper case
			addUser := strings.ToUpper(line)
			// Add domain if flag is defined
			if domain != "" {
				// If there is no @ in line
				if !strings.Contains(line, "@") {
					addUser += "@" + strings.ToUpper(domain)
				}
			}

			usersToAdd = append(usersToAdd, addUser)
		}

		if scanner.Err() != nil {
			log.Printf("Error when scanning file: %+v", scanner.Err())
		}
	}

	// Read positional flags and add to usersToAdd
	users := flags.Args()
	for _, u := range users {
		// All upper case
		addUser := strings.ToUpper(u)
		// Add domain if flag is defined
		if domain != "" {
			// If there is no @ in line
			if !strings.Contains(u, "@") {
				addUser += "@" + strings.ToUpper(domain)
			}
		}

		usersToAdd = append(usersToAdd, addUser)
	}

	if err := db.Own(usersToAdd); err != nil {
		log.Printf("Error when trying to add users to neo4j database: %+v", err)
		os.Exit(1)
	}
}

var ownedHelp = `
  Usage: go-bhtool owned

  Get a list of owned users
` + commonHelp

func owned(args []string, db *db.Neo4jDB) {
	flags := flag.NewFlagSet("owned", flag.ContinueOnError)
	flags.Usage = func() {
		fmt.Print(ownedHelp)
		os.Exit(0)
	}
	flags.Parse(args)

	if err := db.Owned(); err != nil {
		log.Printf("There was an error fetching owned users from neo4j database: %+v", err)
		os.Exit(1)
	}
}

/*

}*/
