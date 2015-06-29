# go-outdated

Go-outdated helps to find outdated packages hosted on github.com in your golang project.

![Dashboard](https://raw.githubusercontent.com/firstrow/go-outdated/master/sample.png)
(sample-image)

## Installation
Install the library with as usual:
``` bash
go get -u github.com/firstrow/go-outdated
```

## Usage
Note: To use this library without limitations, you should create GitHub access token.
``` bash
cd $GOPATH/path/to/your/project
go-outdated
```

## GitHub access token
GitHub API has requests limit. You should create access token and pass it to `go-outdated`
``` bash
go-outdated -token=YOUR_PRIVATE_GITHUB_TOKEN
```
How to create tokens: https://help.github.com/articles/creating-an-access-token-for-command-line-use/

## Aliases
Yes, typing each time access token in command line is not easy. You can create command line alias in your `.zshrc` or `.bashrc` files. Example:
``` bash
alias go-outdated='go-outdated -token=YOUR_PRIVATE_GITHUB_TOKEN'
```

## License:
The MIT License (MIT) 
http://opensource.org/licenses/MIT
