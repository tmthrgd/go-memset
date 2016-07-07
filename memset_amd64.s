// Copyright 2016 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License license that can be found in
// the LICENSE file.
//
// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build amd64,!gccgo,!appengine

#include "textflag.h"

GLOBL memsetShuf<>(SB),RODATA,$16

// func memsetAsm(*byte, uint64, byte)
TEXT ·memsetAsm(SB),NOSPLIT,$0
	MOVQ data+0(FP), DI
	MOVQ len+8(FP), CX
	MOVB value+16(FP), SI

	CMPQ CX, $16
	JB loop

	PINSRB $0, SI, X0
	PSHUFB memsetShuf<>(SB), X0

	CMPB runtime·support_avx(SB), $1
	JNE bigloop

	CMPQ CX, $64
	JB bigloop

	// VINSERTF128 $1, X0, Y0, Y0
	BYTE $0xc4; BYTE $0xe3; BYTE $0x7d; BYTE $0x18; BYTE $0xc0; BYTE $0x01

	CMPQ CX, $0x1000000
	JGE hugeloop_nt_preheader

hugeloop:
	VMOVDQU Y0, -32(DI)(CX*1)
	VMOVDQU Y0, -64(DI)(CX*1)

	SUBQ $64, CX
	JZ ret_after_y0

	CMPQ CX, $64
	JGE hugeloop

	VZEROUPPER

	CMPQ CX, $16
	JB loop

bigloop:
	MOVOU X0, -16(DI)(CX*1)

	SUBQ $16, CX
	JZ ret

	CMPQ CX, $16
	JGE bigloop

loop:
	MOVB SI, -1(DI)(CX*1)

	DECQ CX
	JNZ loop

ret:
	RET

ret_after_y0:
	VZEROUPPER
	RET

hugeloop_nt_preheader:
	// Align to 32 byte boundary
	VMOVDQU Y0, -32(DI)(CX*1)
	ADDQ DI, CX
	ANDQ $~31, CX
	SUBQ DI, CX

hugeloop_nt:
	VMOVNTDQ Y0, -32(DI)(CX*1)
	VMOVNTDQ Y0, -64(DI)(CX*1)
	VMOVNTDQ Y0, -96(DI)(CX*1)
	VMOVNTDQ Y0, -128(DI)(CX*1)

	SUBQ $128, CX
	JZ ret_after_nt

	CMPQ CX, $128
	JGE hugeloop_nt

	SFENCE
	VZEROUPPER

	CMPQ CX, $16
	JGE bigloop
	JMP loop

ret_after_nt:
	SFENCE
	VZEROUPPER
	RET
