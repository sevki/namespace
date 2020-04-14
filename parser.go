// Copyright 2020 Sevki <s@sevki.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package namespace parses name space description files
// https://plan9.io/magic/man2html/6/namespace
package namespace

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// Parse takes a namespace file and returns a collection
// of operations.
func Parse(r io.Reader) (File, error) {
	scanner := bufio.NewScanner(r)

	cmds := []cmd{}

	for scanner.Scan() {
		buf := scanner.Bytes()
		if len(buf) <= 0 {
			continue
		}
		r := buf[0]
		// Blank lines and lines with # as the first non–space character are ignored.
		if r == '#' || r == ' ' {
			continue
		}
		cmd, err := parseLine(scanner.Text())
		if err != nil {
			return nil, err
		}
		cmds = append(cmds, cmd)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	return cmds, nil
}

func parseLine(line string) (cmd, error) {

	var arg string
	args := strings.Fields(line)
	poparg := func() {
		arg = args[0]
		args = args[1:]
	}
	poparg()
	trap := syzcall(0)

	c := cmd{
		syscall: trap,
		flag:    REPL,
	}

	switch arg {
	case "bind":
		c.syscall = BIND
	case "mount":
		c.syscall = MOUNT
	case "unmount":
		c.syscall = UNMOUNT
	case "clear":
		c.syscall = RFORK
	case "cd":
		c.syscall = CHDIR
	case ".":
		c.syscall = INCLUDE
	case "import":
		c.syscall = IMPORT
	default:
		panic(arg)
	}

	// we don't have to chck the size of the second index here as string.FIELDS doesn't return "" strings.
	if len(args) > 0 && args[0][0] == '-' {
		poparg()
		var r byte
	PARSE_FLAG:
		if len(arg) > 0 {
			r, arg = arg[0], arg[1:]
			switch r {
			case 'a':
				c.flag |= AFTER
			case 'b':
				c.flag |= BEFORE
			case 'c':
				c.flag |= CREATE
			case 'C':
				c.flag |= CACHE
			default:
			}
			goto PARSE_FLAG
		}
	}

	c.args = args

	return c, nil
}
