## User Data Exporter
This is a simple tool to export user data from Hiring Motion.
It is intended to be used as a starting point for a more complex export process, if required.

### Usage
1. Clone this repository
2. Run `go mod download` to download dependencies
3. Add your Hiring Motion access token to `.import.env`. Use the network inspector to get the token from the request headers.
4. Run `go run import.go` to import the data to `users.csv`
