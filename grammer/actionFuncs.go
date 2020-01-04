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
	"errors"
	"fmt"
	lex "github.com/shichao1996compiler/lexical"
	log "github.com/shichao1996compiler/logger"
	"strconv"
)

var (
	empty      = "_"
	tempICount = 0
)

func genTemp() string {
	res := "T" + strconv.Itoa(tempICount)
	tempICount += 1
	return res
}

func genAddr(a []byte) []byte {
	if isNum, _ := lex.IsValidNumber([]byte(a)); isNum {
		return []byte(genTemp())
	}
	return a
}

func printCode(codeSlice []*lex.CodeLine) {
	for _, c := range codeSlice {
		log.Error(string(c.Value) + " ")
	}
	fmt.Print("\n")
}

func (am *ActionMap) Declare(code *lex.CodeLine, codeSlice []*lex.CodeLine) error {
	tar := codeSlice[0]
	src := codeSlice[2]

	srcValue, _ := src.GetAttr("value")
	code.AddAttr("value", srcValue)
	code.Tp = src.Tp
	if am.ContainVariation(string(tar.Value)) {
		am.RecordErr(errors.New("Duplicate declared variable: " + string(tar.Value)))
	} else {
		am.AddVariation(string(tar.Value))
	}

	return nil
}

func (am *ActionMap) Assign(code *lex.CodeLine, codeSlice []*lex.CodeLine) error {
	tar := codeSlice[0]
	src := codeSlice[2]
	srcValue, _ := src.GetAttr("value")
	srcQuad, _ := src.GetAttr("quad")
	code.AddAttr("quad", srcQuad)

	am.quoList.Add(":=", string(srcValue), empty, string(tar.Value))
	code.AddAttr("value", srcValue)
	code.Tp = src.Tp
	return nil
}

func (am *ActionMap) Plus(code *lex.CodeLine, codeSlice []*lex.CodeLine) error {
	tar := codeSlice[0]
	src := codeSlice[2]

	addr := string(genAddr(src.Value))
	code.AddAttr("quad", am.quoList.stringId(0))

	am.quoList.Add("+", string(src.Value), string(tar.Value), addr)
	code.AddAttr("value", addr)

	return nil
}

func (am *ActionMap) Dec(code *lex.CodeLine, codeSlice []*lex.CodeLine) error {
	tar := codeSlice[0]
	src := codeSlice[2]

	addr := string(genAddr(src.Value))
	code.AddAttr("quad", am.quoList.stringId(0))

	am.quoList.Add("-", string(src.Value), string(tar.Value), addr)
	code.AddAttr("value", addr)

	return nil
}

func (am *ActionMap) Multiply(code *lex.CodeLine, codeSlice []*lex.CodeLine) error {
	tar := codeSlice[0]
	src := codeSlice[2]

	addr := string(genAddr(src.Value))
	am.quoList.Add("*", string(src.Value), string(tar.Value), addr)
	code.AddAttr("value", addr)

	return nil
}

func (am *ActionMap) Div(code *lex.CodeLine, codeSlice []*lex.CodeLine) error {
	tar := codeSlice[0]
	src := codeSlice[2]

	addr := string(genAddr(src.Value))
	am.quoList.Add("/", string(src.Value), string(tar.Value), addr)
	code.AddAttr("value", addr)

	return nil
}

func (am *ActionMap) IfThen(code *lex.CodeLine, codeSlice []*lex.CodeLine) error {
	cond := codeSlice[1]
	do := codeSlice[3]

	condTrue, _ := cond.GetAttr("true")
	condFalse, _ := cond.GetAttr("false")
	quad, _ := do.GetAttr("quad")

	next, _ := do.GetAttr("next")
	am.BackPatch(condTrue, quad)
	//code.AddAttr("quad", am.quoList.stringId(0))
	code.AddAttr("quad", am.quoList.stringId(0))
	code.AddAttr("next", condFalse+" "+next)
	return nil
}

func (am *ActionMap) IfThenElse(code *lex.CodeLine, codeSlice []*lex.CodeLine) error {
	cond := codeSlice[1]
	do1 := codeSlice[3]
	do2 := codeSlice[5]

	condTrue, _ := cond.GetAttr("true")
	condFalse, _ := cond.GetAttr("false")
	quad1, _ := do1.GetAttr("quad")
	quad2, _ := do2.GetAttr("quad")
	next1, _ := do1.GetAttr("next")
	next2, _ := do2.GetAttr("next")
	nextN, _ := do1.GetAttr("Nnext")

	am.BackPatch(condTrue, quad1)
	am.BackPatch(condFalse, quad2)
	code.AddAttr("quad", am.quoList.stringId(0))

	code.AddAttr("next", next1+" "+next2+" "+nextN)
	return nil
}

func checkErr(des string, err error) {
	if err != nil {
		log.Error(des + " not found")
	}

}

func (am *ActionMap) WhileDo(code *lex.CodeLine, codeSlice []*lex.CodeLine) error {
	cond := codeSlice[1]
	do := codeSlice[3]
	doNext, _ := do.GetAttr("next")

	quad1, _ := cond.GetAttr("quad")

	//quad2, _ := do.GetAttr("quad")

	condFalse, _ := cond.GetAttr("false")
	condTrue, _ := cond.GetAttr("true")

	//log.Error("next: "+doNext+"  quad1: "+ quad1+"  quad2: "+quad2+"  true: "+condTrue+"  false: "+condFalse)
	code.AddAttr("next", condFalse)
	am.BackPatch(doNext, quad1)
	am.BackPatch(condTrue, "102")
	code.AddAttr("quad", am.quoList.stringId(0))
	am.quoList.Add("j", empty, empty, quad1)
	return nil
}

func (am *ActionMap) GreatOp(code *lex.CodeLine, codeSlice []*lex.CodeLine) error {
	am.compareEmiter(">", code, codeSlice)
	return nil
}

func (am *ActionMap) GreatEqualOp(code *lex.CodeLine, codeSlice []*lex.CodeLine) error {
	am.compareEmiter(">=", code, codeSlice)

	return nil
}

func (am *ActionMap) LowOp(code *lex.CodeLine, codeSlice []*lex.CodeLine) error {
	am.compareEmiter("<", code, codeSlice)

	return nil
}

func (am *ActionMap) LowEqualOp(code *lex.CodeLine, codeSlice []*lex.CodeLine) error {
	am.compareEmiter("<=", code, codeSlice)

	return nil
}

func (am *ActionMap) EqualOp(code *lex.CodeLine, codeSlice []*lex.CodeLine) error {
	am.compareEmiter("=", code, codeSlice)
	return nil
}

func (am *ActionMap) IdFunc(code *lex.CodeLine, codeSlice []*lex.CodeLine) error {
	src := codeSlice[0]
	code.Value = src.Value
	code.AddAttr("value", string(code.Value))
	return nil
}

func (am *ActionMap) IFunc(code *lex.CodeLine, codeSlice []*lex.CodeLine) error {
	src := codeSlice[0]
	code.Value = src.Value
	code.AddAttr("value", string(code.Value))
	return nil
}

func (am *ActionMap) compareEmiter(op string, code *lex.CodeLine, codeSlice []*lex.CodeLine) {
	var opstr = "j" + op

	id1 := codeSlice[0]
	id2 := codeSlice[2]

	code.AddAttr("quad", am.quoList.stringId(0))
	id1.AddAttr("value", string(id1.Value))
	code.AddAttr("true", am.quoList.stringId(0))
	code.AddAttr("false", am.quoList.stringId(1))
	value1, _ := id1.GetAttr("value")
	value2, _ := id2.GetAttr("value")

	am.quoList.Add(opstr, value1, value2, "0")
	am.quoList.Add("j", empty, empty, "0")
}
