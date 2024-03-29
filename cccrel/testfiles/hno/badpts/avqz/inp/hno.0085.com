*** hno CCSD(T)/TZ-F12 QFF
memory, 995, m;
gthresh,energy=1.d-12,zero=1.d-22,oneint=1.d-22,twoint=1.d-22;
gthresh,optgrad=1.d-8,optstep=1.d-8;
nocompress;
geomtyp=xyz
bohr
geometry={
3
        85        85
N        0.0000000000       -0.1563567151       -1.1573417013
H        0.0000000000        1.6699384650       -1.9428440464
O        0.0000000000        0.0281860401        1.1495961567
}
 
basis=avqz
  {hf,maxit=500;accu,20;}
{ccsd(t),nocheck,maxit=250;orbital,IGNORE_ERROR;}
