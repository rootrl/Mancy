## Introduction

#### Mancy is a file watcher and supported to auto upload the changes to remote server via ssh/sftp.

I start this project for learning golang and solve some problem in our develop environment.

## Usage
1.  go get github.com/rootrl/Mancy
2.  go build
3.  run "Mancy" or "Mancy.exe"(windows)

## Todo
* The variable "sftpClient" is used as a global variable, that's not cool, Maybe better to define a sftp struct and provides set of methods