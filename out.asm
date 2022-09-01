section .data
buffer times 8192 dd 0
section .text
global test_func
global adder
global __start
test_func:
push ebp
mov ebp, esp
push edi
push dx
mov edi, [ebp+8]
mov dl, byte [edi]
add dl, 32
mov [edi], dl
test dl, dl
jz __loop0e
__loop0s:
call dot
add edi, 1
mov dl, [edi]
test dl, dl
jnz __loop0s
__loop0e:
pop dx
pop edi
pop ebp
ret
adder:
push ebp
mov ebp, esp
push edi
push edx
mov edi, [ebp+8]
mov edx, [edi]
test edx, edx
jz __loop1e
__loop1s:
dec dword [edi]
add edi, 4
inc dword [edi]
sub edi, 4
mov edx, [edi]
test edx, edx
jnz __loop1s
__loop1e:
add edi, 4
mov eax, [edi]
pop edx
pop edi
pop ebp
ret
__start:
push ebp
mov ebp, esp
push edi
push edx
mov edi, [ebp+8]
add edi, 4
mov edx, dword [edi]
add edx, 9
mov [edi], edx
test edx, edx
jz __loop2e
__loop2s:
sub edi, 4
mov edx, dword [edi]
add edx, 8
mov [edi], edx
add edi, 4
dec dword [edi]
mov edx, [edi]
test edx, edx
jnz __loop2s
__loop2e:
sub edi, 4
call dot
add edi, 4
mov edx, dword [edi]
add edx, 7
mov [edi], edx
test edx, edx
jz __loop3e
__loop3s:
sub edi, 4
mov edx, dword [edi]
add edx, 4
mov [edi], edx
add edi, 4
dec dword [edi]
mov edx, [edi]
test edx, edx
jnz __loop3s
__loop3e:
sub edi, 4
inc dword [edi]
call dot
mov edx, dword [edi]
add edx, 7
mov [edi], edx
call dot
call dot
mov edx, dword [edi]
add edx, 3
mov [edi], edx
call dot
add edi, 12
mov edx, dword [edi]
add edx, 8
mov [edi], edx
test edx, edx
jz __loop4e
__loop4s:
sub edi, 4
mov edx, dword [edi]
add edx, 4
mov [edi], edx
add edi, 4
dec dword [edi]
mov edx, [edi]
test edx, edx
jnz __loop4s
__loop4e:
sub edi, 4
call dot
add edi, 12
mov edx, dword [edi]
add edx, 10
mov [edi], edx
test edx, edx
jz __loop5e
__loop5s:
sub edi, 4
mov edx, dword [edi]
add edx, 9
mov [edi], edx
add edi, 4
dec dword [edi]
mov edx, [edi]
test edx, edx
jnz __loop5s
__loop5e:
sub edi, 4
mov edx, dword [edi]
sub edx, 3
mov [edi], edx
call dot
sub edi, 16
call dot
mov edx, dword [edi]
add edx, 3
mov [edi], edx
call dot
mov edx, dword [edi]
sub edx, 6
mov [edi], edx
call dot
mov edx, dword [edi]
sub edx, 8
mov [edi], edx
call dot
add edi, 8
inc dword [edi]
call dot
add edi, 4
mov edx, dword [edi]
add edx, 10
mov [edi], edx
call dot
add edi, 32
mov edx, dword [edi]
add edx, 5
mov [edi], edx
add edi, 4
call comma
mov eax, [edi]
sub edi, 4
mov edx, [edi]
test edx, edx
jz __loop6e
__loop6s:
add edi, 4
call dot
sub edi, 4
dec dword [edi]
mov edx, [edi]
test edx, edx
jnz __loop6s
__loop6e:
pop edx
pop edi
pop ebp
ret
dot:
push eax
push edx
mov edx, 1
mov ecx, edi
mov ebx, 1
mov eax, 4
int 80h
pop edx
pop eax
ret
comma:
push eax
push edx
mov edx, 1
mov ecx, edi
mov ebx, 1
mov eax, 3
int 80h
pop edx
pop eax
ret
quit:
mov ebx, eax
mov eax, 1
int 80h