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
 *     Initial: 2018/01/23        ShiChao
 */

package grammer

import (
	"fmt"
	"strconv"
)

type Quotation struct {
	id   int
	op   string
	a    string
	b    string
	addr string
}

type QuotationList struct {
	List  []*Quotation
	CurID int
}

func (q *QuotationList) stringId(off int) string {
	return strconv.Itoa(q.CurID + off)
}

func (q *QuotationList) Add(op, a, b, addr string) {
	qt := &Quotation{
		q.CurID,
		op,
		a,
		b,
		addr,
	}
	q.CurID += 1
	q.List = append(q.List, qt)
}

func (q *QuotationList) Print() {
	for _, l := range q.List {
		fmt.Printf("%-6d (%-6s %-6s %-6s %s)\n", l.id, l.op, l.a, l.b, l.addr)
	}
}

func (q *QuotationList) SetAddr(tar, val int) {
	q.List[tar-100].addr = strconv.Itoa(val)
}
