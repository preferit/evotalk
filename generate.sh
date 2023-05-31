#!/bin/bash

gen() {
    expr="s|ex/$1|rebel|g"
    tree ex/$1 | sed "$expr" > docs/ex$1.tree
}

gen 01
gen 02
gen 03
gen 04
gen 05
gen 06
