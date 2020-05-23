# gomon
A live code reloader for Go code, no configuration required

## Instalation
Installation of gomon is easy: ```go get github.com/anacanm/gomon```
This installs gomon into your $GOPATH/bin directory. 

If go binaries are not globally executable for you, you'll want to create a permanent alias that executes $GOPATH/bin/gomon.
[Here](https://jonsuh.com/blog/bash-command-line-shortcuts/) is a helpful article on how to create aliases.

You'll end up adding something like the following to your .zshrc or other terminal profile: 

  ```alias gomon="/Users/anacan/go/bin/gomon"```
  
## Usage
gomon is easy to use: cd to the root of your go project, and run ```gomon dir/to/start``` 
 
  Example: if your executable code is in cmd/, then cd into the root of your project and run ```gomon cmd/```

If the go code that you want to run is in the root of your go project, just run ```gomon``` by itself 

### Flags
By default, gomon only watches files that end with .go

If you want gomon to reload when other types of files are changed, add their file extension as a flag (don't include the ".")

Example: executing gomon to run the code in cmd/, reloading whenever a .go, .hbs, or .css file is changed:
```bash 
gomon cmd -hbs -css
```
