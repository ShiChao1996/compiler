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
 *     Initial: 2018/01/19        Shichao
 */

package grammer

import (
	"bytes"
	"fmt"
	lex "github.com/shichao1996compiler/lexical"
	log "github.com/shichao1996compiler/logger"
)

type Parser struct {
	codeStack  codeStack
	startStack numStack
	pt         priorityTable
	progress   []string
}

func (g *Parser) codeToReduction() ([]*lex.CodeLine, int, int) {
	starter := g.codeStack.findStartIndex()
	ending := g.codeStack.findEndIndex()

	if starter == -1 || ending == -1 || starter > ending {
		log.Error("err no starter!!!")
		return nil, 0, 0
	}
	if starter == ending {
		return g.codeStack.stk[starter : ending+1], starter, ending
	}
	return g.codeStack.stk[starter:ending], starter, ending
}

func (g *Parser) shouldReduction() bool {
	prev, _ := g.codeStack.headingVt(1)
	next, _ := g.codeStack.headingVt(0)
	//fmt.Println("should reduction: ",strconv.Itoa(prev.Tp),string(prev.Value),strconv.Itoa(next.Tp),string(next.Value))
	switch pt.compare(prev.Symbol, next.Symbol) {
	case 1:
		return true
	default:
		return false
	}
}

func (g *Parser) Analyse(l *lex.Lex) {
	var next *lex.CodeLine
	var action string

	g.codeStack.push(l.NextCode())
	//fmt.Println("----------start reduction----------")
	for l.HasNext() {
		next = l.NextCode()
		prev := g.codeStack.top()
		g.codeStack.push(next)

		prevSymbol := prev.Symbol
		nextSymbol := next.Symbol

		//fmt.Println("compare: ", string(prevSymbol), string(nextSymbol), strconv.Itoa(pt.compare(prevSymbol, nextSymbol)))
		switch pt.compare(prevSymbol, nextSymbol) {
		case pt.greatter: // 归约
			action = "reduction"
			g.Reduction()

		case pt.equal: // 推入下一个
			action = "input next"
			continue

		case pt.lower: // 记录起始下标
			action = "lower"
		default:
			action = "default"
			g.HandleErr()
		}

		g.progress = append(g.progress, fmt.Sprintf("%-80s %90s %40s", g.codeStack.sprintf(), l.Sprintf(), action))
		if action == "default" {
			break
		}
	}

	//fmt.Println("----------reduction ending----------")
	//g.codeStack.print()
	g.printProgress()
	fmt.Println("----------reduction ending----------")
	actions.quoList.Print()
	actions.PrintErr()
}

func (g *Parser) HandleErr() {
	cur := g.codeStack.top()
	log.Error("Maybe the word \"" + string(cur.Value) + "\" shouldn't show up here")
}

func (g *Parser) Reduction() {
	//fmt.Println("====================================new Reduction")
	for g.shouldReduction() {
		//fmt.Println("---reduction start---")
		g.reduction()
	}
}

func (g *Parser) reduction() {

	codeSlice, start, end := g.codeToReduction()
	//fmt.Printf("start: %d   end: %d \n", start, end)
	g.progress = append(g.progress, fmt.Sprintf("%-80s %90s %40s", g.codeStack.sprintf(), "", "reduction"))

	var lines [][]byte

	for _, c := range codeSlice {
		lines = append(lines, c.Value)
	}

	code := bytes.Join(lines, []byte(" "))
	code = bytes.TrimPrefix([]byte(code), []byte(" "))
	//g.codeStack.print()

	//fmt.Println("wait to reduction: ", string(code))
	starter := lex.NewCodeLine(100, []byte("P"))

	emitFunc := actions.SelectFunc(code)
	if emitFunc != nil {
		emitFunc(starter, codeSlice)
	}
	/*value, err := starter.GetAttr("value")
	if err == nil {
		log.Debug("==========this is a.value: " + value)
	}*/
	g.codeStack.replace(start, end, starter)

}

func (g *Parser) printProgress() {
	for _, s := range g.progress {
		fmt.Printf(s + "\n")
	}
}
