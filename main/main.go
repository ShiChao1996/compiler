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

package main

import (
	"compiler/lexical"
	_ "fmt"
	_ "compiler/language"
	"compiler/grammer"
	"fmt"
)

func main() {
	/*code :=
		`program
		var a,b,b: integer;
		begin
		while a < b do
		if b > 0 then b := b - a else b := a + 1
		end`*/
	//code := `(i + i) + i * i . i #`
	code := `
program
while a < b do
if c < d then x := y + z
end`
	lex := lexical.NewLex()
	lex.Analyse(code)
	lex.PrintCode()
	fmt.Println("=====================================")
	if len(lex.Errs) != 0 {
		return
	}
	grammer := &grammer.Parser{}
	grammer.Analyse(lex)
}
