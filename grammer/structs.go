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
 *     Initial: 2018/01/22        ShiChao
 */

package grammer

import (
	"fmt"
	lex "compiler/lexical"
	"strconv"
)

type codeStack struct {
	stk   []*lex.CodeLine
	index int
}

func (s *codeStack) push(l *lex.CodeLine) {
	s.stk = append(s.stk, l)
	s.index += 1
}

func (s *codeStack) pop() *lex.CodeLine {
	res := s.stk[len(s.stk)-1]
	s.stk = s.stk[:len(s.stk)-1]
	s.index -= 1
	return res
}

func (s *codeStack) popN(n int) {
	s.stk = s.stk[:len(s.stk)-n]
	s.index -= n
}

func (s *codeStack) top() *lex.CodeLine {
	return s.stk[len(s.stk)-1]
}

func (s *codeStack) replace(start, end int, newCode *lex.CodeLine) {
	if start == end {
		s.stk[start] = newCode
		return
	}
	temp := append(s.stk[:start], newCode)
	s.stk = append(temp, s.stk[end:]...)
}

func (s *codeStack) headingVt(n int) (*lex.CodeLine, int) {
	code := s.top()
	var (
		i     int
		count int
		index int
	)
	for i = len(s.stk) - 1; count <= n; i-- {
		if i < 0 {
			fmt.Println("heading out of range", strconv.Itoa(n))
			break
		}
		if s.stk[i].Tp != 100 {
			code = s.stk[i]
			index = i
			count += 1
		}
	}
	return code, index
}

func (s *codeStack) bottomVt(n int) (*lex.CodeLine, int) {
	code := s.top()
	var (
		i     int
		count int
		index int
	)
	for i = 0; count <= n; i++ {
		if s.stk[i].Tp != 100 {
			code = s.stk[i]
			index = i
			count += 1
		}
	}
	return code, index
}

func (s *codeStack) findStartIndex() int {
	var (
		next      *lex.CodeLine
		prev      *lex.CodeLine
		count     int
		index     = -1
		prevIndex int
	)

	for i := len(s.stk) - 1; i >= 0; i-- {
		next, _ = s.headingVt(count)
		prev, prevIndex = s.headingVt(count + 1)
		if pt.compare(prev.Symbol, next.Symbol) == pt.lower {
			//fmt.Printf("prev: %s ;   next: %s \n", string(prev.Value), string(next.Value))
			index = prevIndex + 1
			//fmt.Printf("starter: %d  %s \n", index, string(s.stk[index].Value))
			break
		}
		count += 1
	}
	return index
}

func (s *codeStack) findEndIndex() int {
	var (
		next      *lex.CodeLine
		prev      *lex.CodeLine
		count     int
		index     = -1
		nextIndex int
	)

	for i := len(s.stk) - 2; i >= 0; i-- { // len-2 跳过最后一个
		prev, _ = s.headingVt(count + 1)
		next, nextIndex = s.headingVt(count)

		if pt.compare(prev.Symbol, next.Symbol) == pt.greatter {
			//fmt.Printf("prev: %s ;   next: %s \n", string(prev.Value), string(next.Value))
			index = nextIndex
			//fmt.Printf("end: %d  %s \n", index, string(s.stk[index].Value))
			break
		}
		count += 1
	}
	return index
}

func (s *codeStack) print() {
	for i, a := range s.stk {
		v, _ := a.GetAttr("value")
		if v==""{
			v = "_"
		}
		fmt.Println(strconv.Itoa(i), strconv.Itoa(a.Tp), string(a.Value), v)
	}
}

func (s *codeStack) sprintf() string {
	code := ""
	for _, a := range s.stk {
		code += string(a.Symbol) + " "
	}
	return code
}

type numStack struct {
	stk []int
}

func (s *numStack) push(l int) {
	s.stk = append(s.stk, l)
}

func (s *numStack) pop() {
	i := len(s.stk) - 1
	if i < 0 {
		return
	}
	s.stk = s.stk[:i]
}

func (s *numStack) print() {
	fmt.Println("current numStack:===")
	for _, v := range s.stk {
		fmt.Println(strconv.Itoa(v))
	}
}

func (s *numStack) top() int {
	//todo
	i := len(s.stk) - 1
	if i < 0 {
		return 0
	}
	return s.stk[i]
}
