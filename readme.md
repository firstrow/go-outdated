# go-outdated

Go-outdated helps to find outdated packages hosted on github.com in your golang project.

![Dashboard](https://raw.githubusercontent.com/firstrow/go-outdated/master/sample.png)

## Installation
Install the library with:
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

## Re-entering API token
Yes, typing each time access token in command-line is not easy. Use git config to re-use you token:
``` bash
git config --global github.token YOUR_PRIVATE_GITHUB_TOKEN
go-outdate # token will be taken from git config
```
Github doc: https://github.com/blog/180-local-github-config

## Aliases
If you for some reason do not want to configure you local git config, You can create command-line alias in your `.zshrc` or `.bashrc` files. Example:
``` bash
alias go-outdated='go-outdated -token=YOUR_PRIVATE_GITHUB_TOKEN'
```

## Todo
- Cache
- Refactor

## Links
How-to create API tokens: https://help.github.com/articles/creating-an-access-token-for-command-line-use/  
See also GUI alternative: https://github.com/shurcooL/Go-Package-Store  

## License:
The MIT License (MIT) 
http://opensource.org/licenses/MIT
