# Experiment 02: Minimal Example

Build and run our minimal example application:

```
$ make minimal && ./bin/minimal
mkdir -p ./bin
cd targets/minimal && go build -gcflags '-l' -o /opt/gopath/src/github.com/hashicorp/hackweek/bin/minimal
55
4
18
world hello
5 bytes: [31 32 37 2e 30]
```

Let's look for our functions:

```
$ objdump -t ./bin/minimal | grep 'main\.read'
00000000004df360 g     O .rodata        0000000000000028 main.read.stkobj
0000000000494980 g     F .text  000000000000035a main.read
$ objdump -t ./bin/minimal | grep 'main\.add'
$ objdump -t ./bin/minimal | grep 'main\.swap'
```

Oops they're inlined. Let's tweak that in our makefile by adding `-gcflags '-l'`:

```
$ objdump -t ./bin/minimal | grep 'main\.read'
00000000004df360 g     O .rodata        0000000000000018 main.read.stkobj
0000000000494c30 g     F .text  0000000000000223 main.read
$ objdump -t ./bin/minimal | grep 'main\.add'
0000000000494e60 g     F .text  0000000000000013 main.add
$ objdump -t ./bin/minimal | grep 'main\.swap'
0000000000494e80 g     F .text  0000000000000029 main.swap
```

Let's try out our last exercise's oneliner again, running `./bin/minimal` a handful of times in another terminal window:

```
$ sudo bpftrace -e 'uprobe:./bin/minimal:"main.read" { @counts[probe] = count() } END { printf("\n")
; print(@counts, 10); }'
Attaching 2 probes...
^C
@counts[uprobe:./bin/minimal:main.read]: 5
```

Ok, now let's dig into using uprobes for the arguments. Our trace scripts are going to quickly get long enough where one-liners aren't practical, so we'll write these in their own `.bt` script files.

Here's the probe for our first attempt:

```
uprobe:./bin/minimal:"main.add"
{
    @counts[arg0] = count();
    @counts[arg1] = count();
}
```

Let's run that once:

```
$ sudo ./scripts/02-minimal-01-args.bt
Attaching 3 probes...
Counting arguments to add... Hit Ctrl-C to end.
^C
Top 10 arguments by count:
@counts[1]: 1
@counts[824633787088]: 1
@counts[42]: 1
@counts[7]: 1
@counts[824633778200]: 2
```

What the heck is this?! We can see the first argument fine but the 2nd argument seems to be gone and replaced with... is that a pointer? Let's dump the registers to see what's going on:

```
$ gdb ./bin/minimal
...
(gdb) break main.add
Breakpoint 1 at 0x494e60: file /opt/gopath/src/golang-bpf/targets/minimal/minimal.go, line 29.
(gdb) run
Starting program: /opt/gopath/src/golang-bpf/bin/minimal
[New LWP 2170]
[New LWP 2171]
[New LWP 2172]

Thread 1 "minimal" hit Breakpoint 1, main.add (x=42, y=13, ~r2=<optimized out>)
    at /opt/gopath/src/golang-bpf/targets/minimal/minimal.go:29
29              return x + y
(gdb) info registers
rax            0xd                 13
rbx            0xc000090ed8        824634314456
rcx            0x0                 0
rdx            0x2                 2
rsi            0x2a                42
rdi            0xc0000842a0        824634262176
rbp            0xc000090f78        0xc000090f78
rsp            0xc000090e60        0xc000090e60
r8             0x0                 0
r9             0x0                 0
r10            0x4dc958            5097816
r11            0x1                 1
r12            0xffffffffffffffff  -1
r13            0x2b                43
r14            0x2a                42
r15            0x200               512
rip            0x494e60            0x494e60 <main.add>
eflags         0x202               [ IF ]
cs             0x33                51
ss             0x2b                43
ds             0x0                 0
es             0x0                 0
fs             0x0                 0
gs             0x0                 0
```

Oh right... in AMD64, we'd expect to see arg0 in the `rdi` register and arg1 in the `rsi` register, but golang doesn't respect the AMD64 calling convention because it's golang and they have to be unique. In golang arguments are passed on the stack, so let's inspect the stack. (We're e`x`amining `3` `g`iant words (64-bit words) as `d`ecimal numbers.)

```
(gdb) x/3dg $rsp
0xc000090e60:   4804642 42
0xc000090e70:   13
```

Fortunately for us, the bpftrace folks already have a [`sarg0, sarg1,... sargN` builtin](https://github.com/iovisor/bpftrace/blob/master/docs/reference_guide.md#1-builtins) specifically to deal with languages that pass arguments on the stack. Which I suspect is literally only golang. Of course.

```
uprobe:./bin/minimal:"main.add"
{
    @counts[sarg0] = count();
    @counts[sarg1] = count();
}
```

Unfortunately for us, the version of bpftrace that ships for our test box doesn't have that option yet:

```
$ sudo ./scripts/02-minimal-02-sargs.bt
Unknown identifier: 'sarg0'
Unknown identifier: 'sarg1'
```

This was added in [0.9.3](https://github.com/iovisor/bpftrace/releases/tag/v0.9.3). So we're going to do something incredibly gross, and [copy the bpftrace binary](https://github.com/iovisor/bpftrace/blob/master/INSTALL.md#copying-bpftrace-binary-from-docker) out of the most recent Docker image published by upstream.

```
$ bpftrace -V
bpftrace v0.10.0

$ sudo ./scripts/02-minimal-02-sargs.bt
Attaching 3 probes...
Counting arguments to add... Hit Ctrl-C to end.
^CCould not resolve symbol: /proc/self/exe:END_trigger
```

Hrm, looks like there's a bug (maybe related to [#954](https://github.com/iovisor/bpftrace/issues/954)?) and 0.10.0 was released only 2 days ago (as of this writing), so let's try to grab the last of the 0.9.4? We can build this ourselves, but let's be lazy and look at https://quay.io/repository/iovisor/bpftrace?tab=tags first.

```
$ docker run -v $(pwd)/output:/output quay.io/iovisor/bpftrace:v0.9.4 /bin/bash -c "cp /usr/bin/bpftr
ace /output"
Unable to find image 'quay.io/iovisor/bpftrace:v0.9.4' locally
v0.9.4: Pulling from iovisor/bpftrace
5c939e3a4d10: Pull complete
c63719cdbe7a: Pull complete
19a861ea6baf: Pull complete
651c9d2d6c4f: Pull complete
129a8cf84190: Pull complete
f50efa378520: Pull complete
Digest: sha256:ef46a1d7ccc13b51bd9b643eadfcd922b92ab1fc928d6963c4325d94f29bb854
Status: Downloaded newer image for quay.io/iovisor/bpftrace:v0.9.4

$ ./output/bpftrace -V
bpftrace v0.9.4-dirty

$ sudo ./output/bpftrace ./scripts/02-minimal-02-sargs.bt
Attaching 3 probes...
Counting arguments to add... Hit Ctrl-C to end.
^C
Top 10 arguments by count:
@counts[3]: 1
@counts[42]: 1
@counts[7]: 1
@counts[1]: 1
@counts[11]: 1
@counts[13]: 1
```

Success! From this point on we'll use our "dirty" 0.9.4 version as the installed version and see how that goes for us. The wonders of living on the bleeding edge.

Ok, so now let's exercise filters. Here we're filtering for whenever one of the arguments is over 10.

```
uprobe:./bin/minimal:"main.add" /sarg0 > 10 || sarg1 > 10/
{
    // make sure you end printf with a newline or the last
    // call won't flush!
    printf("got an argument > 10 (sarg0: %d, sarg1: %d)\n", sarg0, sarg1);
}
```

```
$ sudo ./scripts/02-minimal-03-filtered-sargs.bt
Attaching 2 probes...
Waiting... Hit Ctrl-C to end.
got an argument > 10 (sarg0: 42, sarg1: 13)
got an argument > 10 (sarg0: 7, sarg1: 11)
^C
```

Integers are fun, but let's read strings too.

```
uprobe:./bin/minimal:"main.swap"
{
    printf("swapping \"%s\" and \"%s\"\n", str(sarg0), str(sarg1));
}
```

```
$ sudo ./scripts/02-minimal-04-strings.bt
Attaching 2 probes...
Waiting... Hit Ctrl-C to end.
swapping "helloint16int32int64panicscav sleepslicesse41sse42ssse3uint8wor" and ""
```

Oh dear, that's not right... golang doesn't pass null-terminated char arrays around like C, but probably annotates them with the length. Let's look at this in gdb again:

```
# note some relevant config options:
$ cat ~/.gdbinit
set disassembly-flavor intel
add-auto-load-safe-path /usr/local/go/src/runtime/runtime-gdb.py

$ gdb ./bin/minimal
(gdb) b main.swap
Breakpoint 1 at 0x494e80: file /opt/gopath/src/github.com/hashicorp/hackweek/targets/minimal/minimal.go, line 34.
(gdb) r
Starting program: /opt/gopath/src/github.com/hashicorp/hackweek/bin/minimal
[New LWP 3431]
[New LWP 3432]
[New LWP 3433]
[New LWP 3434]
55
4
18

Thread 1 "minimal" hit Breakpoint 1, main.swap (x=..., y=..., ~r2=..., ~r3=...)
    at /opt/gopath/src/github.com/hashicorp/hackweek/targets/minimal/minimal.go:34
34              return y, x


(gdb) info args
x = 0x4c4319 "hello"
y = 0x4c4355 "world"
~r2 = <optimized out>
~r3 = <optimized out>
```

Let's look at those arguments:

```
(gdb) x/s 0x4c4319
0x4c4319:       "helloint16int32int64panicscav sleepslicesse41sse42ssse3uint8worldwrite Value addr= base  code= ctxt: curg= goid  jobs= list= m->p= next= p->m= prev= span= varp=% util(...)\n, i = , not 390625<-chanArab"...
(gdb) x/s 0x4c4355
0x4c4355:       "worldwrite Value addr= base  code= ctxt: curg= goid  jobs= list= m->p= next= p->m= prev= span= varp=% util(...)\n, i = , not 390625<-chanArabicBrahmiCarianChakmaCommonCopticFormatGOROOTGothicHangulHatr"...
```

Ok, so they're not null-terminated, just as we thought. Let's look at the disassembly:

```
(gdb) disas
Dump of assembler code for function main.swap:
=> 0x0000000000494e80 <+0>:     mov    rax,QWORD PTR [rsp+0x18]
   0x0000000000494e85 <+5>:     mov    QWORD PTR [rsp+0x28],rax
   0x0000000000494e8a <+10>:    mov    rax,QWORD PTR [rsp+0x20]
   0x0000000000494e8f <+15>:    mov    QWORD PTR [rsp+0x30],rax
   0x0000000000494e94 <+20>:    mov    rax,QWORD PTR [rsp+0x8]
   0x0000000000494e99 <+25>:    mov    QWORD PTR [rsp+0x38],rax
   0x0000000000494e9e <+30>:    mov    rax,QWORD PTR [rsp+0x10]
   0x0000000000494ea3 <+35>:    mov    QWORD PTR [rsp+0x40],rax
   0x0000000000494ea8 <+40>:    ret
End of assembler dump.
```

Note I'm not an assembly guru! But it looks to me like we're moving 4 arguments around. The 1st and 3rd are pointers to the start of the string arguments:

```
(gdb) x/a $rsp+0x18
0xc000098e78:   0x4c4355
(gdb) x/a $rsp+0x8
0xc000098e68:   0x4c4319
```

Whereas the 2nd and 4th are integers:

```
(gdb) x/d $rsp+0x20
0xc000098e80:   5
(gdb) x/d $rsp+0x10
0xc000098e70:   5
```

Hey, do those look like the length of the string argument to you? Let's update our script:

```
uprobe:./bin/minimal:"main.swap"
{
    printf("swapping \"%s\" and \"%s\"\n", str(sarg0, sarg1), str(sarg2, sarg3));
}
```

And now we can read string arguments!

```
$ sudo ./scripts/02-minimal-04-strings.bt
Attaching 2 probes...
Waiting... Hit Ctrl-C to end.
swapping "hello" and "world"
^C
```

At this point there's a lot more to do in terms of reading structs, slices, etc. But let's move wrap up our minimal example by trying uretprobes in the next section.
