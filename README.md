# go-lisp


# Description

A LISP interpreter hosted on go-lang. 
Initially I am focusing on the Scheme dialect since that's where I have the most
experience with lisp

## Scheme

### REPL

Interactive loop to evaluate expressions. Currently at a hello world stage
```lisp
(format t "hello world")
```
Literals resolve to their own value
```lisp
("hello")
```
prints `hello`

```lisp
(1234)
```
prints `1234`

