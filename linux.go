package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
)

func prepare_function_linux(reader *bufio.Reader, return_reg string, stack_pointer string, base_pointer string, int_size int, func_name string, bracket_exit bool) error {
	var size_word string
	var cell_reg string
	var push_pop string
	var err error

	if int_size == 8 {
		size_word = "byte"
		cell_reg = "dl"
	} else if int_size == 16 {
		size_word = "word"
		cell_reg = "dx"
	} else if int_size == 32 {
		size_word = "dword"
		cell_reg = "edx"
	} else if int_size == 64 {
		size_word = "qword"
		cell_reg = "rdx"
	} else {
		println("Handle invalid size")
	}
	push_pop = cell_reg
	if int_size == 8 {
		push_pop = "dx"
	} else if target == "elf64" && int_size == 32 {
		push_pop = "rdx"
	}

	function := bytes.NewBuffer(make([]byte, 0))
	if func_name == "_start" {
		_, err = function.WriteString("_start:\nmov " + address_reg + ", buffer\n")
	} else if target == "elf32" {
		_, err = function.WriteString(func_name + ":\npush " + base_pointer + "\nmov " + base_pointer + ", " + stack_pointer + "\npush " + address_reg + "\npush " + push_pop + "\nmov " + address_reg + ", [" + base_pointer + "+8]\n")
	} else if target == "elf64" {
		_, err = function.WriteString(func_name + ":\npush " + push_pop + "\n")
	}

	err = parse_linux(reader, int_size, size_word, cell_reg, function, return_reg, bracket_exit)
	if err != nil {
		return err
	}

	if func_name == "_start" {
		_, err = function.WriteString("jmp quit\n")
	} else if target == "elf32" {
		_, err = function.WriteString("pop " + push_pop + "\npop " + address_reg + "\npop " + base_pointer + "\nret\n")
	} else if target == "elf64" {
		_, err = function.WriteString("pop " + push_pop + "\nret\n")
	}
	func_contents = append(func_contents, function)
	global_funcs = append(global_funcs, func_name)
	return nil
}

func functions_linux(reader *bufio.Reader, size int, return_reg string, stack_pointer string, base_pointer string) error {
	for {
		b, err := reader.ReadByte()

		if err != nil {
			if err.Error() == "EOF" {
				break
			} else {
				panic(err)
			}
		}
		if b == ')' {
			println("Handle unbalanced brackets")
			return errors.New("Unbalanced brackets")
		}

		if b == '#' {
			var func_name string
			var func_size string
			for {
				b, err = reader.ReadByte()
				if err != nil {
					if err.Error() == "EOF" {
						println("Handle unexpected EOF")
						return errors.New("Unexpected EOF while parsing function declaration")
					} else {
						panic(err)
					}
				}
				if b == ':' {
					break
				}
				func_size += string(b)
			}
			int_size, err := strconv.Atoi(func_size)
			check(err)
			for {
				b, err = reader.ReadByte()
				if err != nil {
					if err.Error() == "EOF" {
						println("Handle unexpected EOF")
						return errors.New("Unexpected EOF while parsing function")
					} else {
						panic(err)
					}
				}
				if b == '(' {
					break
				}
				func_name += string(b)
			}

			err = prepare_function_linux(reader, return_reg, stack_pointer, base_pointer, int_size, func_name, true)
			check(err)
		}

		if err != nil {
			panic(err)
		}
	}
	return nil
}

func parse_linux(reader *bufio.Reader, size int, size_word string, cell_reg string, writer *bytes.Buffer, return_reg string, bracket_exit bool) error {
	brackets := make([]int, 0)
	multi_cell_op := 0
	multi_move_op := 0
	eof := false
	cell_loaded := false
	size_bytes = size / 8

	for {
		b, err := reader.ReadByte()

		if err != nil {
			if err.Error() == "EOF" {
				if bracket_exit {
					println("Handle unexpected EOF")
				} else {
					eof = true
				}
			} else {
				panic(err)
			}
		}
		if b == ')' && bracket_exit {
			eof = true
		}

		cell_loaded = false

		if multi_move_op != 0 && b != '>' && b != '<' {
			if multi_move_op < 0 {
				_, err = writer.WriteString("sub " + address_reg + ", " + strconv.Itoa(0-multi_move_op) + "\n")
			} else if multi_move_op > 0 {
				_, err = writer.WriteString("add " + address_reg + ", " + strconv.Itoa(multi_move_op) + "\n")
			}
			multi_move_op = 0
		}
		if multi_cell_op != 0 && b != '+' && b != '-' {
			if multi_cell_op == 1 {
				_, err = writer.WriteString("inc " + size_word + " [" + address_reg + "]\n")
			} else if multi_cell_op == -1 {
				_, err = writer.WriteString("dec " + size_word + " [" + address_reg + "]\n")
			} else if multi_cell_op < 0 {
				_, err = writer.WriteString("mov " + cell_reg + ", " + size_word + " [" + address_reg + "]\nsub " + cell_reg + ", " + strconv.Itoa(0-multi_cell_op) + "\nmov [" + address_reg + "], " + cell_reg + "\n")
				cell_loaded = true
			} else if multi_cell_op > 0 {
				_, err = writer.WriteString("mov " + cell_reg + ", " + size_word + " [" + address_reg + "]\nadd " + cell_reg + ", " + strconv.Itoa(multi_cell_op) + "\nmov [" + address_reg + "], " + cell_reg + "\n")
				cell_loaded = true
			}
			multi_cell_op = 0
		}

		if eof {
			if len(brackets) != 0 {
				println("Handle unbalanced brackets")
				break
			} else {
				break
			}
		}

		if b == '+' {
			multi_cell_op++
		} else if b == '-' {
			multi_cell_op--
		} else if b == '>' {
			multi_move_op += size_bytes
		} else if b == '<' {
			multi_move_op -= size_bytes
		} else if b == '.' {
			_, err = writer.WriteString("call dot\n")
		} else if b == ',' {
			_, err = writer.WriteString("call comma\n")
		} else if b == '[' {
			str_num := strconv.Itoa(loop_number)
			if !cell_loaded {
				_, err = writer.WriteString("mov " + cell_reg + ", [" + address_reg + "]\n")
			}
			_, err = writer.WriteString("test " + cell_reg + ", " + cell_reg + "\njz __loop" + str_num + "e\n__loop" + str_num + "s:\n")
			brackets = append(brackets, loop_number)
			loop_number++
		} else if b == ']' {
			if len(brackets) > 0 {
				str_num := strconv.Itoa(brackets[len(brackets)-1])
				brackets = brackets[:len(brackets)-1]
				if !cell_loaded {
					_, err = writer.WriteString("mov " + cell_reg + ", [" + address_reg + "]\n")
				}
				_, err = writer.WriteString("test " + cell_reg + ", " + cell_reg + "\njnz __loop" + str_num + "s\n__loop" + str_num + "e:\n")
			} else {
				err = errors.New("Unbalanced brackets")
			}
		} else if b == '$' {
			if cell_loaded {
				_, err = writer.WriteString("mov " + return_reg + ", " + cell_reg + "\n")
			} else {
				_, err = writer.WriteString("mov " + return_reg + ", [" + address_reg + "]\n")
			}
		}

		if err != nil {
			panic(err)
		}
	}
	return nil
}

func build_linux(length int, size int, files []string, out_file string, buffer *[]byte, parse_mode int) {

	o, err := os.Create(out_file)
	check(err)
	writer := bufio.NewWriter(o)

	var define_type string
	var ra, rb, rc, rd, sp, bp string

	if size == 32 {
		define_type = "dd"
		ra = "eax"
		rb = "ebx"
		rc = "ecx"
		rd = "edx"
		sp = "esp"
		bp = "ebp"
	} else if size == 64 {
		define_type = "dq"
		ra = "rax"
		rb = "rbx"
		rc = "rcx"
		rd = "rdx"
		sp = "rsp"
		bp = "rbp"
	} else {
		println("Handle invalid size")
	}

	start := "section .data\nbuffer times " + strconv.Itoa(length) + " " + define_type + " 0\nsection .text\n"
	dot_func := "dot:\npush " + ra + "\npush " + rd + "\nmov " + rd + ", 1\nmov " + rc + ", " + address_reg + "\nmov " + rb + ", 1\nmov " + ra + ", 4\nint 80h\npop " + rd + "\npop " + ra + "\nret\n"
	comma_func := "comma:\npush " + ra + "\npush " + rd + "\nmov " + rd + ", 1\nmov " + rc + ", " + address_reg + "\nmov " + rb + ", 1\nmov " + ra + ", 3\nint 80h\npop " + rd + "\npop " + ra + "\nret\n"
	quit_func := "quit:\nmov " + rb + ", " + ra + "\nmov " + ra + ", 1\nint 80h"

	_, err = writer.WriteString(start)
	check(err)

	if parse_mode == 0 {
		var io_readers []io.Reader
		for _, file := range files {
			f, err := os.Open(file)
			check(err)
			io_readers = append(io_readers, f)
		}
		readers := io.MultiReader(io_readers...)
		final_reader := bufio.NewReader(readers)
		err = prepare_function_linux(final_reader, ra, sp, bp, size, "_start", false)
		check(err)
	} else if parse_mode == 1 {
		for _, file := range files {
			path := strings.Split(file, "/")
			end_name := path[len(path)-1]
			split := strings.Split(end_name, ".")

			if len(split) == 3 {
				name := split[0]
				str_size := split[1]

				int_size, err := strconv.Atoi(str_size)
				check(err)

				f, err := os.Open(file)
				check(err)
				reader := bufio.NewReader(f)
				err = prepare_function_linux(reader, ra, sp, bp, int_size, name, false)
				check(err)
			} else {
				println("Handle bad name")
			}
		}
	} else if parse_mode == 2 {
		for _, file := range files {
			f, err := os.Open(file)
			check(err)
			reader := bufio.NewReader(f)
			err = functions_linux(reader, size, ra, sp, bp)
			check(err)
		}
	}

	for i := 0; i < len(global_funcs); i++ {
		_, err = writer.WriteString("global " + global_funcs[i] + "\n")
		check(err)
	}

	for i := 0; i < len(func_contents); i++ {
		_, err = writer.ReadFrom(func_contents[i])
		check(err)
	}

	_, err = writer.WriteString(dot_func + comma_func + quit_func)
	check(err)

	err = writer.Flush()
	check(err)
}
