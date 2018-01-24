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
 *     Initial: 2018/01/17        ShiChao
 */

package lexical

import (
	"bytes"
	"unicode/utf8"
	"compiler/language"
	"strconv"
	"fmt"
	"errors"
	log "compiler/logger"
)

type lexErr struct {
	where int
	which []byte
	err   error
}

func (le *lexErr) Print() {
	log.Error(string(le.which) + le.err.Error())
}

type Lex struct {
	CodeList    []*CodeLine
	Errs        []lexErr
	currentLine int
	codeIndex   int
}

func NewLex() *Lex {
	l := &Lex{}
	l.addLine(NewCodeLine(-1, []byte("#"))) // add first #
	return l
}

func (l *Lex) Analyse(s string) []error {
	str := bytes.Fields([]byte(s))

	for _, word := range str {
		var sep []byte
		for i, length, ok := ContainsToken(word); ok; i, length, ok = ContainsToken(word) {
			sep = word[:i]
			l.handleWord(sep)
			l.handleWord(word[i:i+length])
			word = word[i+length:]
		}
		l.handleWord(word)
	}

	l.addLine(NewCodeLine(-1, []byte("#"))) // add last #

	fmt.Printf("=========there are %d errors =========\n", len(l.Errs))
	for _, c := range l.Errs {
		c.Print()
	}
	return nil
}

func (l *Lex) handleWord(word []byte) {
	var line *CodeLine
	for {
		c, _ := utf8.DecodeRune(word)
		if IsLetter(c) {
			if isKey, index := language.IsKey(word); isKey { //关键字
				//line = l.getLine(index, string(word))
				line = NewCodeLine(index, word)
				break
			}

			if isID, err := IsvalidIdentify(word); isID { //标志符
				line = NewCodeLine(18, word)
			} else {
				l.recordErr(word, err)
			}
			break
		}
		if IsDigit(c) { //数字
			isInteger, err := IsValidNumber(word)
			if err != nil {
				l.recordErr(word, err)
				break
			}
			if isInteger {
				line = NewCodeLine(19, word)
			} else {
				line = NewCodeLine(20, word)
			}
			break
		}

		if isToken, i := language.IsToken(word); isToken { //符号
			line = NewCodeLine(i, word)
			break
		}

		break
	}
	if line != nil {
		l.addLine(line)
	}
}

func (l *Lex) recordErr(word []byte, err error) {
	e := lexErr{
		where: l.currentLine,
		which: word,
		err:   err,
	}
	l.Errs = append(l.Errs, e)
}

func (l *Lex) addLine(line *CodeLine) {
	l.CodeList = append(l.CodeList, line)
}

func (l *Lex) getLine(tp int, value string) []byte {
	return []byte(strconv.Itoa(tp) + " " + value)
}

func (l *Lex) PrintCode() {
	for _, c := range l.CodeList {
		fmt.Println(strconv.Itoa(c.Tp), string(c.Value))
	}
}

func (l *Lex) NextCode() *CodeLine {
	if !l.HasNext() {
		return &CodeLine{}
	}
	code := l.CodeList[l.codeIndex]
	l.codeIndex += 1
	return code
}

func (l *Lex) Peek() *CodeLine {
	if !l.HasNext() {
		return &CodeLine{}
	}
	code := l.CodeList[l.codeIndex]
	return code
}

func (l *Lex) HasNext() bool {
	return !(l.codeIndex == len(l.CodeList))
}

func (l *Lex) Sprintf() string {
	code := ""
	for i := l.codeIndex; i < len(l.CodeList)-1; i++ {
		code += string(l.CodeList[i].Value) + " "
	}
	return code
}

type CodeLine struct {
	Tp     int //type
	Value  []byte
	Symbol []byte
	Attr   map[string]string
}

var (
	tpId     = 18
	tpInt    = 19
	tpFloat  = 20
	symbolId = []byte("id")
	symbolI  = []byte("i")
)

func (c *CodeLine) AddAttr(key, v string) bool {
	if c.ContainAttr(key) {
		return false
	}
	c.Attr[key] = v
	return true
}

func (c *CodeLine) GetAttr(key string) (string, error) {
	for k := range c.Attr {
		if k == key {
			return c.Attr[key], nil
		}
	}
	//log.Error(string(c.Value)+" " + key + " not found")
	return "", errors.New("err Attr not found")
}

func (c *CodeLine) SetAttr(k, v string) bool {
	if !c.ContainAttr(k) {
		return false
	}
	c.Attr[k] = v
	return true
}

func (c *CodeLine) ContainAttr(key string) bool {
	for k := range c.Attr {
		if key == k {
			return true
		}
	}
	return false
}

func NewCodeLine(tp int, value []byte) *CodeLine {
	var symbol []byte
	switch tp {
	case tpId:
		symbol = symbolId
	case tpInt, tpFloat:
		symbol = symbolI
	default:
		symbol = value
	}

	return &CodeLine{
		tp,
		value,
		symbol,
		make(map[string]string),
	}
}
