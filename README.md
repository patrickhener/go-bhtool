A helper tool for BloodHound lets you mark a batch of users as owned. Afterwards using "Shortest Path from Owned Principals" is recommended.

# Usage

Build it and then run `go-bhtool --help` to see the available options.

```bash
❯ ./go-bhtool --help

  Usage: go-bhtool [--neo4j-uri] [--neo4j-user] [--neo4j-pass] [command] [--help]

  Version: v0.0.1 (go1.16)

  Defaults:
    neo4j-uri:  bolt://localhost:7687
    neo4j-user: neo4j
    neo4j-pass: admin

  Commands:
    own         - mark multiple users as owned
    owned       - get a list of owned users

  Read more:
    https://github.com/patrickhener/go-bhtool
```