# How to use Kaia tests

[Kaia tests](https://github.com/kaiachain/kaia-tests) is not currently included
here due to its relatively large size.  It will be added as a git submodule
later.

In order to use Kaia tests, you would need to clone it as `testdata`; or
clone it somewhere and make a symbolic link to it.  This document assumes
Kaia tests is cloned outside the Kaia source tree, and the following
instructions describe how to use Kaia tests with `go test`.


## 1. Clone Kaia tests

Clone Kaia tests in the location where you'd like to have it.  Let's say we
clone it in `$HOME/workspace`.

```
$ cd $HOME/workspace
$ git clone git@github.com:kaiachain/kaia-tests.git
```


## 2. Create a symbolic link

We assume Kaia source tree is located in
`$HOME/workspace/go/src/github.com/kaiachain/kaia`.

```
$ cd $HOME/workspace/go/src/github.com/kaiachain/kaia/tests
$ ln -s $HOME/workspace/kaia-tests testdata
```


## 3-1. Run all test cases

Inside `tests` directory, you can run all test cases simply using `go test`.

```
$ go test
```


## 3-2. Run a specific test

There are five test sets available in Kaia tests, which can be put after
`-run` flag.
- Blockchain
   - NOTE: all test cases in BlockchainTests are skipped at this moment because
     of the change in block header.
- State
- Transition
- VM
- RLP

```
$ go test -run VM
```


## 3-3. Run in verbose mode

`-v` flag can be used to run tests in verbose mode.

```
$ go test -run VM -v
```
