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

package language

import (
	"bytes"
)

var (
	keys   = "and begin bool do else end false if integer not or program real then true var while 标志符 整数 实数"
	Tokens = "( ) + - * / . , : ; := = <= < <> > >="

	/*Grammer = `
	X-> # E #
	E-> E + T | T
	T-> T * F | F
	F-> P . F | P
	P-> ( E ) | i
	`*/

	Grammer = `
X -> # P #
P -> program L
L -> S | id , L | id : K | var L ; G
K -> integer | bool | real
G -> begin S end
S -> id := E | if B then S else S | while B do S
B -> id < I | id > I
E -> id + I | id - I
I -> i | id | ( E ) | E
`

	keySlice   [][]byte
	tokenSlice [][]byte
)

func init() {
	keySlice = bytes.Fields([]byte(keys))
	tokenSlice = bytes.Fields([]byte(Tokens))
}

func IsKey(word []byte) (bool, int) {
	for i, c := range keySlice {
		if bytes.Equal(word, c) {
			return true, i + 1
		}
	}
	return false, 0
}

func IsToken(word []byte) (bool, int) {
	for i, c := range tokenSlice {
		if bytes.Equal(word, c) {
			return true, i + 21
		}
	}
	return false, 0
}

func IsBuiltIn(word []byte) bool {
	isToken, _ := IsToken(word)
	isKey, _ := IsKey(word)
	return isToken || isKey
}
