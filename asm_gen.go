// Copyright 2016 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License license that can be found in
// the LICENSE file.

// +build ignore

package main

import "github.com/tmthrgd/asm"

const header = `// Copyright 2016 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License license that can be found in
// the LICENSE file.
//
// This file is auto-generated - do not modify

// +build amd64,!gccgo,!appengine
`

func memsetAsm(a *asm.Asm) {
	a.NewFunction("memsetAsm")
	a.NoSplit()

	dst := a.Argument("dst", 8)
	length := a.Argument("len", 8)
	value := a.Argument("value", 4)

	a.Start()

	hugeloop := a.NewLabel("hugeloop")
	hugeloop_nt := hugeloop.Suffix("nt")
	hugeloop_nt_preheader := hugeloop_nt.Suffix("preheader")
	bigloop := a.NewLabel("bigloop")
	loop := a.NewLabel("loop")
	ret := a.NewLabel("ret")
	ret_after_y0 := ret.Suffix("after_y0")
	ret_after_nt := ret.Suffix("after_nt")

	di, si, cx := asm.DI, asm.SI, asm.BX

	a.Movq(di, dst)
	a.Movq(cx, length)
	a.Movb(si, value)

	a.Cmpq(asm.Constant(16), cx)
	a.Jb(loop)

	a.Pinsrb(asm.X0, si, asm.Constant(0))
	a.Pxor(asm.X1, asm.X1)
	a.Pshufb(asm.X0, asm.X1)

	a.Cmpb(asm.Constant(1), asm.Data("runtime·support_avx"))
	a.Jne(bigloop)

	a.Cmpq(asm.Constant(64), cx)
	a.Jb(bigloop)

	a.Vinsertf128(asm.Y0, asm.Y0, asm.X0, asm.Constant(1))

	a.Cmpq(asm.Constant(0x1000000), cx)
	a.Jae(hugeloop_nt_preheader)

	a.Label(hugeloop)

	for i := -32; i >= -64; i -= 32 {
		a.Vmovdqu(asm.Address(di, cx, asm.SX1, i), asm.Y0)
	}

	a.Subq(cx, asm.Constant(64))
	a.Jz(ret_after_y0)

	a.Cmpq(asm.Constant(64), cx)
	a.Jae(hugeloop)

	a.Vzeroupper()

	a.Cmpq(asm.Constant(16), cx)
	a.Jb(loop)

	a.Label(bigloop)

	a.Movou(asm.Address(di, cx, asm.SX1, -16), asm.X0)

	a.Subq(cx, asm.Constant(16))
	a.Jz(ret)

	a.Cmpq(asm.Constant(16), cx)
	a.Jae(bigloop)

	a.Label(loop)

	a.Movb(asm.Address(di, cx, asm.SX1, -1), si)

	a.Decq(cx)
	a.Jnz(loop)

	a.Label(ret)
	a.Ret()

	a.Label(ret_after_y0)
	a.Vzeroupper()
	a.Ret()

	a.Label(hugeloop_nt_preheader)

	// Align to 32 byte boundary
	a.Vmovdqu(asm.Address(di, cx, asm.SX1, -32), asm.Y0)
	a.Addq(cx, di)
	a.Andq(cx, asm.Constant(^uint64(31)))
	a.Subq(cx, di)

	a.Label(hugeloop_nt)

	for i := -32; i >= -128; i -= 32 {
		a.Vmovntdq(asm.Address(di, cx, asm.SX1, i), asm.Y0)
	}

	a.Subq(cx, asm.Constant(128))
	a.Jz(ret_after_nt)

	a.Cmpq(asm.Constant(128), cx)
	a.Jae(hugeloop_nt)

	a.Sfence()
	a.Vzeroupper()

	a.Cmpq(asm.Constant(16), cx)
	a.Jae(bigloop)
	a.Jmp(loop)

	a.Label(ret_after_nt)
	a.Sfence()
	a.Vzeroupper()
	a.Ret()
}

func main() {
	if err := asm.Do("memset_amd64.s", header, memsetAsm); err != nil {
		panic(err)
	}
}
