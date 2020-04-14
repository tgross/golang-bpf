# BPF for Nomad

## Goals

- make reusable workaround for `bpftrace` uretprobes
- document reusable workflows for using `bpftrace` w/ golang applications
- demonstrate usage of `bpftrace` on Nomad
- build `nomad-debug` to deploy `bpftrace` scripts vs arbitrary Nomad workloads

## Non-Goals

- dealing with arbitrary BPF programs beyond `bpftrace` (i.e. no `bcc-tool`, `gobpf`, etc.)

## Tool Install

Use `Vagrantfile` in this directory
- starts with `bento/ubuntu-19.10` machine
- adds scripts from Nomad for golang development
- add scripts for installing bpftrace and tools

## Experiments

### Target Workloads

1. a hello world to repro [Brendan Gregg's blog](http://www.brendangregg.com/blog/2017-01-31/golang-bcc-bpf-function-tracing.html) but w/ bpftrace
2. a minimal application that makes some function calls, defers, etc.
3. a web application that does JSON parsing (reflection), calls into sqlite (CGO)
4. Nomad

## References

- [Brendan Gregg's blog on golang function tracing](http://www.brendangregg.com/blog/2017-01-31/golang-bcc-bpf-function-tracing.html)
- [`bpftrace` reference guide](https://github.com/iovisor/bpftrace/blob/master/docs/reference_guide.md)
- [`bpftrace` sample tools](https://github.com/iovisor/bpftrace/tree/master/tools)
