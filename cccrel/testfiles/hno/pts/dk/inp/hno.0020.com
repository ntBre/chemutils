*** hno CCSD(T)/TZ-F12 QFF
memory, 995, m;
gthresh,energy=1.d-12,zero=1.d-22,oneint=1.d-22,twoint=1.d-22;
gthresh,optgrad=1.d-8,optstep=1.d-8;
nocompress;
geomtyp=xyz
bohr
geometry={
3
        20        20
H        0.0000000000        1.6593518710       -1.9263638113
N        0.0000000000       -0.1540844149       -1.1588767900
O        0.0000000000        0.0365003339        1.1346510104
}
 
basis=vtz-dk
dkroll=0
  {hf,maxit=500;accu,20;}
{ccsd(t),nocheck,maxit=250;orbital,IGNORE_ERROR;}
