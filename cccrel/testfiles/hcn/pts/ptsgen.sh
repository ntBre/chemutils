#!/bin/bash

mkdir mt mtc dk dkr av5z avqz avtz
for i in */; do cp /ddn/home4/r2666/chem/rotation/${i%/}.py $i/.; cp file07 $i/.; (cd $i; ./${i%/}.py hcn N C H); done
