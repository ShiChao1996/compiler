/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Co., Ltd.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

/*
 * Revision History:
 *     Initial: 2018/01/19        ShiChao
 */

package grammer

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/shichao1996compiler/language"
	"github.com/shichao1996compiler/lexical"
	"strings"
)

var (
	grammer            string
	regularGenerations [][]byte
	Vn                 Set
	Vt                 Set
	FirstVT            map[byte]*Set
	LastVT             map[byte]*Set
	LastVnFlag         map[byte]bool
	FirstVnFlag        map[byte]bool
)

func init() {
	grammer = language.Grammer
	regularGenerations = regulateGenerations()
	initVtVn()
	Vt.Print()
	InitFirstVTAndLastVT()

	/*for k, v := range FirstVT {
		fmt.Println("=======first vt of ", string(k))
		for _,c := range v.eles  {
			fmt.Printf("%-6s", string(c))
		}
		fmt.Println("")
	}

	for k, v := range LastVT {
		fmt.Println("=======-----------Last vt of ", string(k))
		for _,c := range v.eles  {
			fmt.Printf("%-6s", string(c))
		}
		fmt.Println("")
	}*/

	pt = newPriorityTable()
	pt.init()
	genPriorityTable()
	pt.print(7)
}

type Set struct {
	eles [][]byte
}

func (s *Set) add(b []byte) {
	if !s.contains(b) {
		s.eles = append(s.eles, b)
	}
}

func (s *Set) Print() {
	for _, c := range s.eles {
		fmt.Println(string(c))
	}
}

func (s *Set) contains(b []byte) bool {
	for _, c := range s.eles {
		if bytes.Equal(b, c) {
			return true
		}
	}
	return false
}

func (s *Set) index(b []byte) int {
	for i, c := range s.eles {
		if bytes.Equal(b, c) {
			return i
		}
	}
	return -1
}

func initVtVn() {
	grammerStr := strings.Replace(grammer, "->", " ", -1)
	grammerStr = strings.Replace(grammerStr, "|", " ", -1)

	s := []byte(grammerStr)
	words := bytes.Fields(s)
	for _, c := range words {
		if len(c) == 1 && isCaptainLetter(c[0]) {
			Vn.add(c)
		} else {
			Vt.add(c)
		}
	}
}

func InitFirstVTAndLastVT() {
	FirstVT = make(map[byte]*Set)
	LastVT = make(map[byte]*Set)
	LastVnFlag = make(map[byte]bool)
	FirstVnFlag = make(map[byte]bool)

	for _, vn := range Vn.eles {
		FirstVT[vn[0]] = &Set{}
		LastVT[vn[0]] = &Set{}
		FirstVnFlag[vn[0]] = false
		LastVnFlag[vn[0]] = false
	}

	lines := regularGenerations

	for _, vn := range Vn.eles {
		for _, vn := range Vn.eles {
			FirstVnFlag[vn[0]] = false
		}
		findFirstVt(lines, vn, vn[0])

		for _, vn := range Vn.eles {
			LastVnFlag[vn[0]] = false
		}
		findLastVt(lines, vn, vn[0])
	}
}

func findFirstVt(lines [][]byte, start []byte, vn byte) {
	if FirstVnFlag[start[0]] {
		return
	}
	FirstVnFlag[start[0]] = true // 产生式有右递归，防止死循环
	newLines := linesStartWith(lines, start)

	for _, line := range newLines {

		line = line[1:]

		words := bytes.Fields(line)
		if Vn.contains(words[0]) {
			if len(words) > 1 && Vt.contains(words[1]) {
				FirstVT[vn].add(words[1])
			}
			findFirstVt(lines, words[0], vn)
		}
		if Vt.contains(words[0]) {
			FirstVT[vn].add(words[0])
		}
	}
}

func findLastVt(lines [][]byte, start []byte, vn byte) {
	if LastVnFlag[start[0]] {
		return
	}
	LastVnFlag[start[0]] = true // 产生式有右递归，防止死循环
	newLines := linesStartWith(lines, start)

	for _, line := range newLines {
		words := bytes.Fields(line)
		lastOne := words[len(words)-1]

		if Vt.contains(lastOne) {
			LastVT[vn].add(lastOne)
			continue
		}
		if Vn.contains(lastOne) {
			if len(words) > 1 {
				lastSecond := words[len(words)-2]
				if Vt.contains(lastSecond) {
					LastVT[vn].add(lastSecond)
				}
			}
			findLastVt(lines, lastOne, vn)
		}
	}
}

func linesStartWith(lines [][]byte, s []byte) [][]byte {
	var newlines [][]byte
	for _, line := range lines {
		if bytes.HasPrefix(line, s) {
			newlines = append(newlines, line)
		}
	}
	return newlines
}

func isCaptainLetter(b byte) bool {
	return b >= 0x41 && b <= 0x5A
}

func regulateGenerations() [][]byte {
	str := strings.Replace(grammer, "->", " ", -1)
	b := []byte(str)
	lines := bytes.Split(b, []byte("\n"))
	var res [][]byte
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		start := line[0]
		newLines := bytes.Split(line[1:], []byte("|"))
		for _, l := range newLines {
			if len(l) == 0 {
				continue
			}
			nl := []byte(string(start) + " " + string(l))
			res = append(res, nl)
		}
	}

	for i := range res {
		res[i] = bytes.Join(bytes.Fields(res[i]), []byte(" "))
	}

	return res
}

/*
func FindVtParent(code []byte) ([]byte, error) {
	fmt.Println("==----------")
	var (
		err   error
		res   []byte
		count = 0
	)

	for res = code; err == nil; count += 1 {
		res, err = findParent(res)
	}
	fmt.Println("==----------")
	fmt.Println(string(res))
	if count == 1 {
		err = errors.New("vn not found")
	} else {
		err = nil
	}

	return append([]byte("100 "), res...), err
}*/

func FindVtParent(code []byte) (*lexical.CodeLine, error) {
	for _, line := range regularGenerations {
		l := line[2:]
		if bytes.Equal(l, code) {
			return lexical.NewCodeLine(100, line[:1]), nil
		}

		if len(code) == len(l) {
			codewords := bytes.Fields(code)
			lineWords := bytes.Fields(l)
			flag := true
			for i := 0; i < len(code); i++ {
				if bytes.Equal(codewords[i], lineWords[i]) || Vn.contains(codewords[i]) && Vt.contains(lineWords[i]) {
					continue
				} else {
					flag = false
				}
			}
			if flag {
				return lexical.NewCodeLine(100, line[:1]), nil
			}
		}
	}

	return nil, errors.New("vn not found")
}
