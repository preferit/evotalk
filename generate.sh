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
gen 07_1
gen 07_2
gen 07_3
gen 07_4
gen 07_5
gen 08_1
gen 08_2
gen 08_3

