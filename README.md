# Mindbreaker
An optimising compiler for Brainfuck written in Go. It is capable of compiling into Linux NASM assembly code, which can then be compiled to an executable or an object file to be linked to other code.

## Usage
To compile the program, use `go build main.go linux.go`. This compiles the main argument-parsing code and the code to compile for Linux. After this I would recommend renaming the executable to something like `mindbreaker` (or using the `-o` option when compiling).

### Commandline Arguments
`-o`: Output file name (default `out.asm`)
`--target`: The target format (`elf32` or `elf64`). This specifies the kind of assembly to generate
`--buffer_size`: The size of the array to store values in
**File Format Options**
- `--functional`: Brainfuck programs are written in the form of functions like `#integer_size:function_name(code)` where `integer_size` is the size in bits of each integer in the array (acceptable values: `8`, `16`, `32`, `64`). These can be called from other programs if a linker is used
- `--flat`: Files are named in the format `function_name.integer_size.bf`. Each file is a function with the same purpose as `functional`
- `--raw`: Raw Brainfuck; creates a program by concatenating all the input file content and compiling it

**TODO**: Write language guide
