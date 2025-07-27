// bitasm_amd64.s
#include "textflag.h"

TEXT ·Bit8(SB), NOSPLIT, $0-9
    MOVB    a+0(FP), AL
    MOVQ    i+8(FP), CX
    MOVL    $1, DX
    SHLL    CX, DX
    TESTB   AL, DL
    SETNE   AX
    MOVB    AL, ret+16(FP)
    RET
