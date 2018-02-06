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
	"fmt"
	"bytes"
	"strconv"
)

var (
	pt *priorityTable
)

type priorityTable struct {
	table     [][]int
	undefined int
	greatter  int
	lower     int
	equal     int
}

func newPriorityTable() *priorityTable {
	return &priorityTable{
		undefined: 100,
		greatter:  1,
		lower:     -1,
		equal:     0,
	}
}

func (pt *priorityTable) init() {
	pt.table = make([][]int, len(Vt.eles)+2)
	for i := range pt.table {
		pt.table[i] = make([]int, len(Vt.eles)+2)
		for j := range pt.table[i] {
			pt.table[i][j] = pt.undefined
		}
	}
}

func (pt *priorityTable) addRelation(line, column []byte, relation int) {
	indexLine := Vt.index(line)
	indexColumn := Vt.index(column)
	if indexColumn == -1 || indexLine == -1 {
		return
	}
	if pt.table[indexLine][indexColumn] != pt.undefined && pt.table[indexLine][indexColumn] != relation {
		fmt.Println("==================================dup!!!")
		fmt.Println(string(line),string(column))
		return
	}
	pt.table[indexLine][indexColumn] = relation
}

func (pt *priorityTable) compare(line, column []byte) int {
	indexLine := Vt.index(line)
	indexColumn := Vt.index(column)
	if indexColumn == -1 || indexLine == -1 {
		return pt.undefined
	}
	return pt.table[indexLine][indexColumn]
}

func (pt *priorityTable) great(line, column []byte) bool {
	return pt.compare(line, column) == pt.greatter
}

func (pt *priorityTable) low(line, column []byte) bool {
	return pt.compare(line, column) == pt.lower
}

func (pt *priorityTable) eq(line, column []byte) bool {
	return pt.compare(line, column) == pt.equal
}

func (pt *priorityTable) print(n int) {
	fmtStr := "%-" + strconv.Itoa(n) + "s"
	fmt.Printf(fmtStr, " ")
	for _, vt := range Vt.eles {
		fmt.Printf(fmtStr, string(vt))
	}
	fmt.Print("\n")
	for i := range pt.table {
		if i < len(Vt.eles) {
			//fmt.Print(string(Vt.eles[i]) + "  ")
			fmt.Printf(fmtStr, string(Vt.eles[i]))
		}
		for j := range pt.table[i] {
			switch pt.table[i][j] {
			case pt.greatter:
				fmt.Printf(fmtStr, ">")
			case pt.lower:
				fmt.Printf(fmtStr, "<")
			case pt.equal:
				fmt.Printf(fmtStr, "=")
			default:
				fmt.Printf(fmtStr, " ")
			}
		}
		fmt.Print("\n")
	}
}

func genPriorityTable() {
	lines := regulateGenerations()
	for _, line := range lines {
		line = line[1:]
		words := bytes.Fields(line)
		l := len(words) - 1
		for i := 0; i < l; i++ {
			if Vt.contains(words[i]) && Vn.contains(words[i+1]) {
				addLowerRelation(words[i], FirstVT[words[i+1][0]])
			}
			if Vn.contains(words[i]) && Vt.contains(words[i+1]) {
				addGreatterRelation(LastVT[words[i][0]], words[i+1])
			}
			if Vt.contains(words[i]) && Vt.contains(words[i+1]) {
				pt.addRelation(words[i], words[i+1], pt.equal)
			}
		}
		for i := 0; i < l-1; i++ {
			if Vt.contains(words[i]) && Vn.contains(words[i+1]) && Vt.contains(words[i+2]) {
				pt.addRelation(words[i], words[i+2], pt.equal)
			}
		}
	}
}

func addLowerRelation(a []byte, firstVt *Set) {
	for _, vt := range firstVt.eles {
		pt.addRelation(a, vt, pt.lower)
	}
}
func addGreatterRelation(lastVt *Set, b []byte) {
	for _, vt := range lastVt.eles {
		pt.addRelation(vt, b, pt.greatter)
	}
}
