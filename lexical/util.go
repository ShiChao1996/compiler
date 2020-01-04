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
 *     Initial: 2018/01/17        shichao
 */

package lexical

import (
	"bytes"
	"errors"
	"unicode"
	"unicode/utf8"
)

var (
	ErrInvalidNum = errors.New(" is not a valid number")
	ErrInvalidID  = errors.New(" is not a valid identify name")
)

func IsLetter(r rune) bool {
	return r >= 0x41 && r <= 0x7A
}

func IsDigit(r rune) bool {
	return r >= 0x30 && r <= 0x39
}

func isDot(rune rune) bool {
	return rune == 0x2e
}

func IsValidNumber(n []byte) (bool, error) { //error == nil means n is valid number,bool == true means it's integer
	var (
		total    = 0
		dotCount = 0
		err      error
	)
	len := len(n)
	for {
		r, size := utf8.DecodeRune(n)
		if !IsDigit(r) {
			if isDot(r) {
				dotCount += 1
			} else {
				err = ErrInvalidNum
			}
		}
		n = n[size:]
		total += size
		if total == len {
			break
		}
	}
	if err != nil || dotCount > 1 {
		return false, err
	}
	if dotCount == 0 {
		return true, nil
	}
	return false, nil
}

func IsvalidIdentify(n []byte) (bool, error) {
	l := len(n)
	total := 0
	for {
		r, size := utf8.DecodeRune(n)
		if !IsLetter(r) {
			return false, ErrInvalidID
		}
		n = n[size:]
		total += size
		if total == l {
			return true, nil
		}
	}
}

var tokenStr = "( ) + - * / . , := : ; <= <> < >= > ="
var tokens = bytes.Fields([]byte(tokenStr))
var lineCount = 0
var lineDelim = []byte("\n")

func ContainsToken(s []byte) (int, int, bool) {
	for _, c := range tokens {
		if bytes.Contains(s, c) {
			index := bytes.Index(s, c)
			return index, len(c), true
		}
	}
	return -1, 0, false
}

func MyIsSpace(r rune) bool {
	if uint32(r) <= unicode.MaxLatin1 {
		switch r {
		case '\n':
			return false
		}
	}
	return unicode.IsSpace(r)
}

func LineCount() int {
	return lineCount
}
