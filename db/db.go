package db

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

const (
	queryUserOwned    string = "MATCH (u:User) WHERE u.name = $name RETURN u.owned AS owned"
	queryAlreadyOwned string = "MATCH (u:User) WHERE u.owned = True RETURN u.name AS name"
	queryOwn          string = "MATCH (u:User) WHERE u.name = $name SET u.owned=True RETURN u.name AS name"
)

// Neo4jDB will hold the neo4j connection details and the session object
type Neo4jDB struct {
	Driver  neo4j.Driver
	Session neo4j.Session
}

// Connect will connect to the database and test the connection
func (n *Neo4jDB) Connect(uri, user, pass string, tls bool) error {
	var err error
	n.Driver, err = neo4j.NewDriver(uri, neo4j.BasicAuth(user, pass, ""), func(c *neo4j.Config) {
		if tls {
			c.Encrypted = true
		} else {
			c.Encrypted = false
		}
	})
	if err != nil {
		return err
	}

	if err := n.Driver.VerifyConnectivity(); err != nil {
		return err
	}

	n.Session, err = n.Driver.NewSession(neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	if err != nil {
		return err
	}

	return nil
}

// Own will modify given user and add the "owned" flag
func (n *Neo4jDB) Own(users []string) error {
	alreadyOwned := 0
	owned := 0
	for _, u := range users {
		// if user is not already owned
		if !n.checkowned(u) {
			result, err := n.Session.Run(queryOwn, map[string]interface{}{
				"name": u,
			})
			if err != nil {
				return err
			}

			if result.Next() {
				fmt.Printf("[+] %s marked as owned\n", u)
				owned++
			} else {
				fmt.Printf("[-] %s does not exist\n", u)
			}

		} else {
			fmt.Printf("[-] %s already marked as owned\n", u)
			alreadyOwned++
		}
	}

	// Statistics
	fmt.Println("[*] Operation finished")
	if owned == 1 {
		fmt.Printf("[*] %d user was marked as owned\n", owned)
	} else {
		fmt.Printf("[*] %d user were marked as owned\n", owned)
	}
	if alreadyOwned == 1 {
		fmt.Printf("[*] There was %d user already marked as owned\n", alreadyOwned)
	} else {
		fmt.Printf("[*] There were %d user already marked as owned\n", alreadyOwned)
	}
	fmt.Println("[*] Happy Graphing!")

	return nil
}

// This function is used to check if a user is already owned or not
func (n *Neo4jDB) checkowned(user string) bool {
	result, err := n.Session.Run(queryUserOwned, map[string]interface{}{
		"name": user,
	})
	if err != nil {
		return false
	}

	if result.Next() {
		owned, ok := result.Record().Get("owned")
		if !ok {
			return false
		}
		if owned != nil {
			return true
		}
	}
	return false
}

// Owned will print out all user marked as owned
func (n *Neo4jDB) Owned() error {
	result, err := n.Session.Run(queryAlreadyOwned, nil)
	if err != nil {
		return err
	}

	for result.Next() {
		record := result.Record()
		name, ok := record.Get("name")
		if ok {
			fmt.Println(name)
		}

	}

	if result.Err() != nil {
		return result.Err()
	}

	return nil
}
