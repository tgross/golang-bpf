# Experiment 04: A Slightly More Complex Example

The source for `./target/worker` runs work across multiple goroutines,
has some state, and passes around structs and slices. Let's run it:

```
$ make clean worker && ./bin/worker
rm -rf ./bin
mkdir -p ./bin
cd targets/worker && go build -o /opt/gopath/src/golang-bpf/bin/worker
1035
$
```

Let's practice getting the arguments via uprobes on the worker's
`work` method. First, we need the symbol:

```
$ objdump -t ./bin/worker | grep work
./bin/worker:     file format elf64-x86-64
0000000000564da0 g     O .bss   0000000000000190 runtime.work
000000000041b6c0 g     F .text  00000000000002c8 runtime.scanframeworker
0000000000420dd0 g     F .text  0000000000000054 runtime.(*workbuf).checknonempty
0000000000420e30 g     F .text  0000000000000054 runtime.(*workbuf).checkempty
00000000004919e0 g     F .text  0000000000000199 main.(*Worker).work
```

Then we need to remember that the actual 1st arg for a method will be
a pointer to the struct, so to get the 2nd param we want the _3rd_
argument, or `sarg2`:

```
$ sudo bpftrace -e 'uprobe:./bin/worker:"main.(*Worker).work"
{ printf("times: %d\n", sarg2) }'
Attaching 1 probe...
times: 9
times: 0
times: 1
times: 2
times: 3
times: 4
times: 5
times: 6
times: 7
times: 8
```


Next let's try to get the results of the `Inc` method. First, we need
the symbol:

```
$ objdump -t ./bin/worker | grep Inc
0000000000429fc0 g     F .text  00000000000000b6 runtime.mSysStatInc
00000000004917b0 g     F .text  00000000000000c7 main.(*State).Inc
```

Let's just make sure we the syntax right first:

```
$ sudo bpftrace -e 'uprobe:./bin/worker:"main.(*State).Inc" { @count++ } END{ printf("%d\n", @count); }'
Attaching 2 probes...
^C43


@count: 43
```

And now let's try to get the results. First we'll disassemble the
function; the `-Mintel` flag here displays the assembly in my
preferred syntax, and the `-S` flag interleaves the source lines which
is nice for navigating the listing.

```
$ objdump --disassemble="main.(*State).Inc" -Mintel -S ./bin/worker

./bin/worker:     file format elf64-x86-64


Disassembly of section .text:

00000000004918a0 <main.(*State).Inc>:
type State struct {
        Counts int
        lock   sync.Mutex
}

func (s *State) Inc() int {
  4918a0:       64 48 8b 0c 25 f8 ff    mov    rcx,QWORD PTR fs:0xfffffffffffffff8
  4918a7:       ff ff
  4918a9:       48 3b 61 10             cmp    rsp,QWORD PTR [rcx+0x10]
  4918ad:       0f 86 97 00 00 00       jbe    49194a <main.(*State).Inc+0xaa>
  4918b3:       48 83 ec 30             sub    rsp,0x30
  4918b7:       48 89 6c 24 28          mov    QWORD PTR [rsp+0x28],rbp
  4918bc:       48 8d 6c 24 28          lea    rbp,[rsp+0x28]
  4918c1:       0f 57 c0                xorps  xmm0,xmm0
  4918c4:       0f 11 44 24 18          movups XMMWORD PTR [rsp+0x18],xmm0
  4918c9:       c6 44 24 0f 00          mov    BYTE PTR [rsp+0xf],0x0
  4918ce:       48 c7 44 24 40 00 00    mov    QWORD PTR [rsp+0x40],0x0
  4918d5:       00 00
        s.lock.Lock()
  4918d7:       48 8b 44 24 38          mov    rax,QWORD PTR [rsp+0x38]
  4918dc:       84 00                   test   BYTE PTR [rax],al
  4918de:       48 8d 48 08             lea    rcx,[rax+0x8]
  4918e2:       48 89 4c 24 10          mov    QWORD PTR [rsp+0x10],rcx
  4918e7:       48 89 0c 24             mov    QWORD PTR [rsp],rcx
  4918eb:       e8 d0 96 fd ff          call   46afc0 <sync.(*Mutex).Lock>
        defer s.lock.Unlock()
  4918f0:       48 8d 05 a9 64 03 00    lea    rax,[rip+0x364a9]        # 4c7da0 <go.func.*+0x42e>
  4918f7:       48 89 44 24 20          mov    QWORD PTR [rsp+0x20],rax
  4918fc:       48 8b 44 24 10          mov    rax,QWORD PTR [rsp+0x10]
  491901:       48 89 44 24 18          mov    QWORD PTR [rsp+0x18],rax
  491906:       c6 44 24 0f 01          mov    BYTE PTR [rsp+0xf],0x1
        s.Counts++
  49190b:       48 8b 44 24 38          mov    rax,QWORD PTR [rsp+0x38]
  491910:       48 8b 08                mov    rcx,QWORD PTR [rax]
  491913:       48 ff c1                inc    rcx
  491916:       48 89 08                mov    QWORD PTR [rax],rcx
        return s.Counts
  491919:       48 89 4c 24 40          mov    QWORD PTR [rsp+0x40],rcx
  49191e:       c6 44 24 0f 00          mov    BYTE PTR [rsp+0xf],0x0
  491923:       48 8b 44 24 18          mov    rax,QWORD PTR [rsp+0x18]
  491928:       48 89 04 24             mov    QWORD PTR [rsp],rax
  49192c:       e8 af 99 fd ff          call   46b2e0 <sync.(*Mutex).Unlock>
  491931:       48 8b 6c 24 28          mov    rbp,QWORD PTR [rsp+0x28]
  491936:       48 83 c4 30             add    rsp,0x30
  49193a:       c3                      ret
  49193b:       e8 e0 c2 f9 ff          call   42dc20 <runtime.deferreturn>
  491940:       48 8b 6c 24 28          mov    rbp,QWORD PTR [rsp+0x28]
  491945:       48 83 c4 30             add    rsp,0x30
  491949:       c3                      ret
  49194a:       e8 d1 82 fc ff          call   459c20 <runtime.morestack_noctxt>
  49194f:       e9 4c ff ff ff          jmp    4918a0 <main.(*State).Inc>
```

There's actually _two_ paths through this code because of the lock,
but the return value location is the same either way. At address
`491916`, we see the result of the increment instruction being moved
out of `rcx` and into `rax`. That's at relative position
`main.(*State).Inc+0x76`. (You can use Python as a quickie hex
calculator `python3 -c "print(hex(0x491916 - 0x4918a0))"`)

```
$ sudo bpftrace -e 'uprobe:./bin/worker:"main.(*State).Inc"+0x76 { printf("%d\n", *reg("ax")) }'
Attaching 1 probe...
0
1
2
3
4
5
6
7
8
9
...
^C
```

Let's look at a more complex return structure. The `results` method on
`Worker` returns a slice of pointers to structs. Note that if we
didn't have the `-l` flag, this very simple method would be inlined
again, which would make it harder to get the results. We'll try this
again with inlining later.

```
$ objdump --disassemble="main.(*Worker).results" -Mintel -S ./bin/worker

./bin/worker:     file format elf64-x86-64


Disassembly of section .text:

0000000000491b80 <main.(*Worker).results>:
                w.counts = append(w.counts, &count{i: result})
        }
}

func (w *Worker) results() []*count {
        return w.counts
  491b80:       48 8b 44 24 08          mov    rax,QWORD PTR [rsp+0x8]
  491b85:       48 8b 48 10             mov    rcx,QWORD PTR [rax+0x10]
  491b89:       48 8b 50 08             mov    rdx,QWORD PTR [rax+0x8]
  491b8d:       48 8b 40 18             mov    rax,QWORD PTR [rax+0x18]
  491b91:       48 89 54 24 10          mov    QWORD PTR [rsp+0x10],rdx
  491b96:       48 89 4c 24 18          mov    QWORD PTR [rsp+0x18],rcx
  491b9b:       48 89 44 24 20          mov    QWORD PTR [rsp+0x20],rax
  491ba0:       c3                      ret
```

We expect the results to come back on the stack, so we can hook our
probe at the `ret` instruction at `491ba0` (or
`main.(*Worker).results+0x20`) and then walk the stack from there,
just as we did in the minimal example application.

We know that golang slices are a pointer to the 1st element of the
backing array, an integer length, and an integer capacity. If we run
this through `gdb`, put a breakpoint at the start of `results`, hit
`continue` a few times, we can inspect what this looks like:

```
(gdb) x $rsp+0x20
0xc000090ef0:   4  # <-- capacity

(gdb) x $rsp+0x18
0xc000090ee8:   3  # <-- length

(gdb) x/a $rsp+0x10
0xc000090ee0:   0xc00009c1a0  # <-- pointer
```

That pointer `0xc00009c1a0` should point to an array of pointers, each
of which points to one of our `count` structs. Let's see if we can
unpack that with bpftrace.

Now before we go pointer chasing, it's a good idea to make sure we're
at least hooking the right thing:

```
$ sudo bpftrace -e 'uprobe:./bin/worker:"main.(*Worker).results"+0x20
> {
>     printf("len: %d, cap: %d\n", *(reg("sp")+0x18), *(reg("sp")+0x20));
> }'
Attaching 1 probe...
len: 0, cap: 0
len: 1, cap: 1
len: 2, cap: 2
len: 3, cap: 4
len: 4, cap: 4
len: 5, cap: 8
len: 6, cap: 8
len: 7, cap: 8
len: 8, cap: 8
len: 9, cap: 16
```

This is kinda cool, because you can see how the runtime is expanding
the capacity of the slice by powers of 2 as the length
increases. Let's see if we can run through the structs.


**TODO**


## References

- [Go Slices: usage and
  internals](https://blog.golang.org/slices-intro)
