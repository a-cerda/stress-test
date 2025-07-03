# Stress test for GET API's

This program provides a simple CLI interface for stress testing GET API's (POST support planned).

## Requierements

1. GO version 1.22 or later
    - On Fedora and Amazon Linux: `sudo dnf install golang`
    - On Ubuntu based distros: `sudo apt install golang`
    - For MacOS and Windows refer to [https://go.dev/dl/]()  
That's it.

## Running the test

1. Put your API url onto API_PATH, e.g for linux and MacOS:
    - `export API_PATH=<Your API URL>`
2. Make sure you have a text.txt file with possible queries:
    - `touch text.txt`
    - Queries must have the form `endpoint?param1=value1&param2=value2` etc.
3. Configure how many users you want to simulate, inside stress-test.go there's an array with tests (default 1000)
4. Build or run directly: `go build stress-test.go` or `go run stress-test.go`

## Questions or suggestions

- Please open an Issue if you have a question, feel free to open a pull request if you have a solution.
- POST support is planned for a future iteration.
