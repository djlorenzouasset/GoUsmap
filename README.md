# GoUsmap

A fast and reliable `.usmap` file parser for Go, supporting Oodle, Brotli, and ZStandard compression methods.

[![Release](https://img.shields.io/github/release/djlorenzouasset/GoUsmap)]()
[![GoMod](https://img.shields.io/github/go-mod/go-version/djlorenzouasset/GoUsmap?style=flat)](https://pkg.go.dev/github.com/djlorenzouasset/GoUsmap)

## Features

- Multiple compression formats: Oodle, Brotli, ZStandard
- Support mostly all usmap versions
- Type-safe API with proper error handling

## Installation
```bash
go get github.com/djlorenzouasset/GoUsmap
```

**Requirements:**
- Go 1.24 or higher
- Oodle DLL (only required for Oodle-compressed mappings)

## Usage

### Parsing Brotli/ZStandard compressed mappings

```go
package main

import (
    "fmt"
    "github.com/djlorenzouasset/GoUsmap"
)

func main() {
    usmap, err := gousmap.ParseFromFile("path/to/mappings.usmap", nil)
    if err != nil {
        fmt.Printf("error: %v", err)
        return
    }

    fmt.Printf("ZStandard/Brotli compressed mapping: %s\n", usmap.ToString())
}
```

### Parsing Oodle-compressed mappings

For Oodle compression, you need to provide the Oodle DLL:

```go
package main

import (
    "fmt"
    "github.com/djlorenzouasset/GoUsmap"
)

func main() {
    oodleInst, err := gousmap.CreateOodleInstance("path/to/oo2core_9_win64.dll")
    if err != nil {
        fmt.Printf("error: %v", err)
        return
    }

    // pass Oodle instance used to uncompress data
    usmap, err := gousmap.ParseFromFile("path/to/mappings.usmap", oodleInst)
    if err != nil {
        fmt.Printf("error: %v", err)
        return
    }

    fmt.Printf("Oodle-compressed mapping: %s\n", usmap.ToString())
}
```

## API Reference

### Main Types

- `Usmap` - Main structure containing names, enums, and schemas
- `UsmapEnum` - Enum definition with name and values
- `UsmapSchema` - Class/struct schema with properties
- `UsmapProperty` - Property definition with type information

### Functions

- `ParseFromFile(filePath string, oodleInstance *Oodle) (*Usmap, error)` - Parse from file
- `CreateOodleInstance(dllPath string) (*Oodle, error)` - Load Oodle DLL (Windows only)

## Credits

This project is a Go port of [Usmap.NET](https://github.com/NotOfficer/Usmap.NET) by [NotOfficer](https://github.com/NotOfficer).

## Contributing

Bug reports and pull requests are welcome! This project was created to learn Go and understand the usmap format.
