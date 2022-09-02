package db

import (
	"fmt"
	"os"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

const (
	queryUserOwned            string = "MATCH (u:User) WHERE u.name = $name RETURN u.owned AS owned"
	queryAlreadyOwned         string = "MATCH (u:User) WHERE u.owned = True RETURN u.name AS name"
	queryOwn                  string = "MATCH (u:User) WHERE u.name = $name SET u.owned=True RETURN u.name AS name"
	queryComputerOwned        string = "MATCH (u:Computer) WHERE u.name = $name RETURN u.owned AS owned"
	queryAlreadyOwnedComputer string = "MATCH (u:Computer) WHERE u.owned = True RETURN u.name AS name"
	queryOwnComputer          string = "MATCH (u:Computer) WHERE u.name = $name SET u.owned=True RETURN u.name AS name"
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

// Own will modify given objects and add the "owned" flag
func (n *Neo4jDB) Own(objects []string, what string) error {
	alreadyOwned := 0
	owned := 0
	switch what {
	case "user":
		for _, u := range objects {
			// if object is not already owned
			if !n.checkowned(u, what) {
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
	case "computer":
		for _, u := range objects {
			// if object is not already owned
			if !n.checkowned(u, what) {
				result, err := n.Session.Run(queryOwnComputer, map[string]interface{}{
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
			fmt.Printf("[*] %d computer was marked as owned\n", owned)
		} else {
			fmt.Printf("[*] %d computer were marked as owned\n", owned)
		}
		if alreadyOwned == 1 {
			fmt.Printf("[*] There was %d computer already marked as owned\n", alreadyOwned)
		} else {
			fmt.Printf("[*] There were %d computer already marked as owned\n", alreadyOwned)
		}
		fmt.Println("[*] Happy Graphing!")
	default:
		fmt.Println("Nothing happend")
		os.Exit(0)
	}

	return nil
}

// This function is used to check if a user is already owned or not
func (n *Neo4jDB) checkowned(object string, what string) bool {
	var result neo4j.Result
	var err error
	switch what {
	case "user":
		result, err = n.Session.Run(queryUserOwned, map[string]interface{}{
			"name": object,
		})
		if err != nil {
			return false
		}

	case "computer":
		result, err = n.Session.Run(queryComputerOwned, map[string]interface{}{
			"name": object,
		})
		if err != nil {
			return false
		}
	default:
		return false
	}

	if result.Next() {
		owned, ok := result.Record().Get("owned")
		if !ok {
			return false
		}
		ownedBool, ok := owned.(bool)
		if !ok {
			return false
		}
		return ownedBool
	}
	return false
}

// Owned will print out all objects marked as owned
func (n *Neo4jDB) Owned(what string) error {
	var result neo4j.Result
	var err error
	switch what {
	case "user":
		result, err = n.Session.Run(queryAlreadyOwned, nil)
		if err != nil {
			return err
		}

	case "computer":
		result, err = n.Session.Run(queryAlreadyOwnedComputer, nil)
		if err != nil {
			return err
		}
	default:
		fmt.Println("Nothing happend")
		os.Exit(0)
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
