// Copyright 2014 The Gman Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in LICENSE file.

// man2md
// Package for converting man (*roff) pages to Mandown format.

package man2md

import (
	"bufio"
	"io"
)

func man2md(reader io.Reader, writer io.Writer) (err error) {
	bufferedReader := bufio.NewReader(reader)
	bufferedWriter := bufio.NewWriter(writer)

	eof := false
	for !eof {
		line, err := bufferedReader.ReadString('\n')
		if err == io.EOF {
			err = nil
			eof = true
		} else if err != nil {
			return err
		}

		// Skip commands
		if len(line) > 0 && line[0] != '.' {
			_, err := bufferedWriter.WriteString(line)
			if err != nil {
				return err
			}
		}
	}

	return err
}
