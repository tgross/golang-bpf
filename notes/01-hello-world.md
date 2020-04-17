# Experiment 01: Hello World

There isn't an example `funccount.bt`, but it's simple to make one for
our application. For this exercise we'll just use the one-liner syntax
because these aren't very interested and aren't too complicated.

First attempt: too many functions to trace!

```
$  sudo bpftrace -e 'uprobe:./bin/helloworld:* {printf("ok\n");}'
Can't attach to 1773 probes because it exceeds the current limit of 512 probes.
You can increase the limit through the BPFTRACE_MAX_PROBES environment variable, but BE CAREFUL since a high number of probes attached can cause your system to crash.

$ objdump -tT ./bin/helloworld | wc -l
objdump: ./bin/helloworld: not a dynamic object
2639
```

Wildcard syntax seems to have trouble resolving the symbol, let's come back to that later.

```
$ sudo bpftrace -e 'uprobe:./bin/helloworld:"fmt.*" {printf("ok\n");}'
Attaching 1 probe...
Could not resolve symbol: ./bin/helloworld:fmt.*
```

Let's find the symbol we want:

```
$ objdump -t ./bin/helloworld | grep -i println
000000000048aee0 g     F .text  00000000000000f2 fmt.Fprintln
0000000000491120 g     F .text  00000000000001e4 fmt.(*pp).doPrintln
```

```
sudo bpftrace -e 'uprobe:./bin/helloworld:"fmt.Fprintln" { printf("ok\n"); }'
Attaching 1 probe...
ok
ok
ok
ok
^C
```

Let's count how many times we hit the function by taking the `func`
built-in for the function name and adding it to a map of counts, which
we print at the end. But this will fail to validate (we can dig into
the BPF instructions with the `-v` flag):

```
$ sudo bpftrace -v -e 'uprobe:./bin/helloworld:"fmt.Fprintln" { @counts[func] = count() } END { print(@counts, 10); }'
Attaching 2 probes...

Program ID: 165
...

Error log:
0: (79) r1 = *(u64 *)(r1 +128)
1: (71) r3 = *(u8 *)(r1 +9)
R1 invalid mem access 'inv'
processed 2 insns (limit 1000000) max_states_per_insn 0 total_states 0 peak_states 0 mark_read 0

Error loading program: uprobe:./bin/helloworld:fmt.Fprintln
```

It looks like `func` built-in variable is nil because it isn't exposed
to `uprobes`? Let's come back to that one later too. For now let's
just use the probe name. Once that's running, we'll run
`./bin/helloworld` in another terminal window a handful of times:

```
$ sudo bpftrace -e 'uprobe:./bin/helloworld:"fmt.Fprintln" { @counts[probe] = count() } END {     printf("\n"); print(@counts, 10); }'
Attaching 2 probes...
^C
@counts[uprobe:./bin/helloworld:fmt.Fprintln]: 5
```

Easy enough, let's move on to a still-minimal but more interesting
example.

## References

- [Brendan Gregg's blog on golang function
  tracing](http://www.brendangregg.com/blog/2017-01-31/golang-bcc-bpf-function-tracing.html)
- [`bpftrace` reference
  guide](https://github.com/iovisor/bpftrace/blob/master/docs/reference_guide.md)
- [`bpftrace` sample
  tools](https://github.com/iovisor/bpftrace/tree/master/tools)
- [`bpftrace` issue
  #1098](https://github.com/iovisor/bpftrace/issues/1098): "can't find
  go symbols"
- [`funccount.py`](https://github.com/iovisor/bcc/blob/master/tools/funccount.py)
