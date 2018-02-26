# Duplicate Image Search in GO

[![Build Status](https://travis-ci.org/abates/disgo.svg?branch=develop)](https://travis-ci.org/abates/disgo) [![GoDoc](https://godoc.org/github.com/abates/disgo?status.png)](https://godoc.org/github.com/abates/disgo) [![Coverage Status](https://coveralls.io/repos/github/abates/disgo/badge.svg?branch=develop)](https://coveralls.io/github/abates/disgo?branch=develop)

This package is still a work in progress, but it works well in my own testing.

### Example

```Go
package main

import "fmt"

import "github.com/abates/disgo"

func main() {
  // Create a disgo database with the default radix index
  db, _ := disgo.New()

  // load an image into the database and get the hash back
  file, _ := os.Open("test.png")
  hash, _ := db.AddFile(file)
  fmt.Printf("Image Hash: %08x\n", hash)

  // search for all hashes with a Hamming distance of 3
  // or less
  matches, _ := db.SearchByHash(hash, 3)
  fmt.Printf("Matches: %v\n", matches)
}

```

### TODO
- [ ] make radix index save/load functions thread safe
- [ ] add record storage (e.g. file path) to database

