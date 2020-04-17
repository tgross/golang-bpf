# BPF for Go

*A work-in-progress repository of guidelines for using `bpftrace` with
go applications*

## Project Goals

- document a reusable workaround for `bpftrace` uretprobes
- produce a collection of "snippets" for using `bpftrace` w/ golang
  applications
- document the process by which I worked out all the workarounds and
  snippet described above
- demonstrate usage of `bpftrace` on
  [Nomad](https://www.nomadproject.io/), which is the project I'm
  working on full-time
- produce some kind of e-book to share this is a nicer format than
  this messy repo

**Non-Goals**
- dealing with arbitrary BPF programs beyond `bpftrace` (i.e. no
  `bcc-tool`, `gobpf`, etc.)

## Progress

**Examples:**

- [x] Hello world
- [x] Minimal stack arguments tracing (minimal)
- [x] Single integer return tracing (minimal)
- [x] Multiple return tracing (minimal)
- [ ] Function latency example (worker, web)
- [ ] Working with goroutines (worker, web)
- [ ] Unpacking structs and slices from results (worker, web)
- [ ] Mixing application-level and OS-level tracing (Nomad)

**Documentation:**
- [ ] Introduce BPF and `bpftrace` (probably borrow from my [DevOps
      Data Philly
      workshop](https://github.com/tgross/ebpf-workshop-dodphilly2019))
- [ ] Explain why golang is weird and special for BPF vs C
- [ ] Very light introduction to reading assembly of golang programs
- [ ] Curated pointers to better reference material

## Setup

### Tool Install

Use [Vagrant](https://www.vagrantup.com/) with the `Vagrantfile` at the
root of this repo.

- run `vagrant up` to start the machine and provision it
- starts with `bento/ubuntu-19.10` machine
- adds scripts from Nomad for golang development
- add scripts for installing `bpftrace` and tools
- run `vagrant ssh` to shell in

### Repo Contents

- `Vagrantfile` is the definition of the virtual machine used to run
  the experiments.
- `GNUmakefile` is the [makefile](https://www.gnu.org/software/make/)
  for building the workloads used in the experiments.

Directories in this repo:

- `./notes` contains notes of my exploration of using `bpftrace` to
  examine golang applications. They include plenty of false starts and
  are generally rough-draft quality.
- `./provisioning` contains shell scripts used by the Vagrantfile at
  the root of the repo to stand up an Ubuntu box with a recent Linux
  kernel and the tools needed to run through the exercises in the
  notes.
- `./scripts` contains some example `bpftrace` scripts used by the
  examples in the notes.
- `./targets` contains some example workloads to trace.

Directories created by this repo:

- `./.vagrant` is the state of the Vagrant-managed virtual machine.
- `./bin` is the binaries of the example workloads that the makefile
  will build.
- `./bpftrace` is the [source
  repo](https://github.com/iovisor/bpftrace) for bpftrace.

### Target Workloads

1. a hello world to repro [Brendan Gregg's
   blog](http://www.brendangregg.com/blog/2017-01-31/golang-bcc-bpf-function-tracing.html)
   but w/ bpftrace
2. a minimal application that makes some function calls
3. a worker application that has concurrency (goroutines), defers,
   multiple returns, and slightly more complex parameter and results
   parsing.
4. [TBD] a web application that does JSON parsing (reflection), calls
   into sqlite (CGO)
5. Nomad

## References

- [Brendan Gregg's blog on golang function
  tracing](http://www.brendangregg.com/blog/2017-01-31/golang-bcc-bpf-function-tracing.html)
- [`bpftrace` reference
  guide](https://github.com/iovisor/bpftrace/blob/master/docs/reference_guide.md)
- [`bpftrace` sample
  tools](https://github.com/iovisor/bpftrace/tree/master/tools)
- [`bpftrace` source repo](https://github.com/iovisor/bpftrace)
