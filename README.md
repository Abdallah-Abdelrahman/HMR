# HMR (Hot module replacement)

Is a feature to update the application in real time without refreshing the browser. the project is written specifically for developers
working with markup files (html, etc..) that are not using any framework like React, Angular, Vue, etc.

## Installation:
- clone the repo
- navigate to the project's directory
- run `go mod tidy`

## How to use HMR
- build  the binary by running `go build -o hmr`
- make `symlink` to the binary by running `ln -s <project-absolute-path>/hmr /usr/local/bin/hmr`
- run the binary from anywhere in your file system `hmr <path-to-dir-where-your-(html|css|js)-files-are>`.
- put this [script](ws.html) in your html file.
