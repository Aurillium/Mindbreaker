package main

// TODO
// Don't immediately create the output file
// Handle errors
// Count lines and columns for errors
// Validate function names

import (
	"bytes"
	"flag"
)

var loop_number int
var global_funcs []string
var func_contents []*bytes.Buffer
var target string
var address_reg string

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func nTrue(b ...bool) int {
	n := 0
	for _, v := range b {
		if v {
			n++
		}
	}
	return n
}

func main() {
	var size int
	var output string
	var raw, flat, functional bool

	flag.IntVar(&size, "buffer_size", 8192, "size of the buffer in a pure Brainfuck program")
	flag.StringVar(&target, "target", "elf32", "build target type")
	flag.StringVar(&output, "o", "out.asm", "output file")
	flag.BoolVar(&raw, "raw", false, "these files are a whole raw program")
	flag.BoolVar(&flat, "flat", false, "these files are functions named by the file name")
	flag.BoolVar(&functional, "functional", false, "these files have functions")
	flag.Parse()
	files := flag.Args()
	
	if len(files) == 0 {
		println("You need to specify one or more files to compile.")
	}

	if nTrue(raw, flat, functional) > 1 {
		println("Arguments 'raw', 'flat', and 'functional' are mutually exclusive.")
		return
	}
	var parse_mode int = 2
	if raw {
		parse_mode = 0
	} else if flat {
		parse_mode = 1
	} else if functional {
		parse_mode = 2
	}

	if target == "elf32" {
		address_reg = "edi"
		build_linux(size, 32, files, output, nil, parse_mode)
	} else if target == "elf64" {
		address_reg = "rdi"
		build_linux(size, 64, files, output, nil, parse_mode)
	} else {
		println("Handle invalid target")
	}
}
