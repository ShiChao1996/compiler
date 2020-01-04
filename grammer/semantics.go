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
	"github.com/shichao1996compiler/language"
	lex "github.com/shichao1996compiler/lexical"
	log "github.com/shichao1996compiler/logger"
	"strconv"
	"strings"
)

var actions *ActionMap

func init() {
	actions = &ActionMap{
		make(map[string]func(code *lex.CodeLine, codeSlice []*lex.CodeLine) error),
		make([]string, 0),
		&QuotationList{
			make([]*Quotation, 0),
			100,
		},
		make(map[string]interface{}),
		make([]error, 0),
	}
	actions.Add(":=", actions.Assign)
	actions.Add("+", actions.Plus)
	actions.Add("-", actions.Dec)
	actions.Add("*", actions.Multiply)
	actions.Add("/", actions.Div)
	actions.Add("if then", actions.IfThen)
	actions.Add("if then else", actions.IfThenElse)
	actions.Add("while do", actions.WhileDo)
	actions.Add(">", actions.GreatOp)
	actions.Add(">=", actions.GreatEqualOp)
	actions.Add("<", actions.LowOp)
	actions.Add("<=", actions.LowEqualOp)
	actions.Add("=", actions.EqualOp)
	actions.Add("id", actions.IdFunc)
	actions.Add("i", actions.IFunc)
	actions.Add(": integer", actions.Declare)
	actions.Add(": bool", actions.Declare)
	actions.Add(": real", actions.Declare)
	actions.Add(",", actions.Declare)
	actions.Add("if then end", actions.IfThen)

}

type ActionMap struct {
	Action     map[string]func(code *lex.CodeLine, codeSlice []*lex.CodeLine) error
	keys       []string
	quoList    *QuotationList
	variations map[string]interface{}
	errs       []error
}

func (am *ActionMap) SelectFunc(code []byte) func(code *lex.CodeLine, codeSlice []*lex.CodeLine) error {
	if key, ok := am.containKey(string(code)); ok {
		return am.Action[key]
	}
	if isId, _ := lex.IsvalidIdentify(code); isId {
		return am.Action["id"]
	}
	if isNum, _ := lex.IsValidNumber(code); isNum {
		return am.Action["i"]
	}
	return nil
}

func (am *ActionMap) containKey(code string) (key string, ok bool) {
	for _, s := range am.keys {
		if s == filterToken(code) {
			log.Info("success matched: " + code)
			return s, true
		}
	}
	return "", false
}

func filterToken(s string) string {
	slice := strings.Fields(s)
	var filtered []string
	for _, s := range slice {
		if language.IsBuiltIn([]byte(s)) {
			filtered = append(filtered, s)
		}
	}
	str := strings.Join(filtered, " ")
	return str
}

func (am *ActionMap) Add(key string, f func(code *lex.CodeLine, codeSlice []*lex.CodeLine) error) bool {
	if am.Contain(key) {
		return false
	}
	am.Action[key] = f
	am.keys = append(am.keys, key)
	return true
}

func (am *ActionMap) Get(key string) (func(code *lex.CodeLine, codeSlice []*lex.CodeLine) error, error) {
	for k := range am.Action {
		if k == key {
			return am.Action[key], nil
		}
	}
	return nil, errors.New("err Attr not found")
}

func (am *ActionMap) Set(k string, f func(code *lex.CodeLine, codeSlice []*lex.CodeLine) error) bool {
	if !am.Contain(k) {
		return false
	}
	am.Action[k] = f
	return true
}

func (am *ActionMap) Contain(key string) bool {
	for k := range am.Action {
		if key == k {
			return true
		}
	}
	return false
}

func (am *ActionMap) ContainVariation(key string) bool {
	for k := range am.variations {
		if key == k {
			return true
		}
	}
	return false
}

func (am *ActionMap) AddVariation(key string) {
	if !am.ContainVariation(key) {
		am.variations[key] = ""
	}
}

func (am *ActionMap) GetVariation(key string) (interface{}, error) {
	for k := range am.variations {
		if k == key {
			return am.variations[key], nil
		}
	}
	return nil, errors.New("err variation not found")
}

func (am *ActionMap) RecordErr(err error) {
	am.errs = append(am.errs, err)
}

func (am *ActionMap) PrintErr() {
	for _, e := range am.errs {
		log.Error("there are err(s) found: ", e.Error())
	}
}

func (am *ActionMap) BackPatch(t, v string) {
	strSlice := strings.Fields(t)
	for _, str := range strSlice {
		tar, _ := strconv.Atoi(str)
		val, _ := strconv.Atoi(v)

		if tar > am.quoList.CurID {
			return
		}
		log.Info("backpatch num " + strconv.Itoa(tar) + " : " + strconv.Itoa(val))
		am.quoList.SetAddr(tar, val)
	}
}
