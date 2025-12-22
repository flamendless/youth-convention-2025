# Requirements

1. [Go](https://go.dev)
2. Setup `GOPATH`, see [doc](https://go.dev/wiki/SettingGOPATH)
3. [Git](https://git-scm.com)

# Setup

1. `git clone` the project of course
2. `cd youth-convention-2025/`
3. `go install`
4. `chmod +x run.sh`
5. Edit `run.sh` accordingly:
    - Update `GOOS` based on your operating system (or just delete the line)
    - Replace `vivaldi` with `chrome` or whatever browser you use
6. `./run.sh deps` to install other dependencies
7. `./run.sh client` to run the server and launch the site on your browser

The project already support hot-reloading via `air` so changes in the codebase should automatically update the client.
