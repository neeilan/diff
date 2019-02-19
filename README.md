# diff: Automatic (Symbolic?) Differentiation in Go

[![Build Status](https://travis-ci.org/neeilan/diff.svg?branch=master)](https://travis-ci.org/neeilan/diff)


We use "reverse mode automatic differentiation. What does this mean? If we have some function of the form f(g(x)) we work form the outside in when we calculate the derivative.
Hence, we'll calculate f' first, then g'.

To run tests locally:
```
go test -v diff/difflib
```
