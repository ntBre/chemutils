*** hno CCSD(T)/TZ-F12 QFF
memory, 995, m;
gthresh,energy=1.d-12,zero=1.d-22,oneint=1.d-22,twoint=1.d-22;
gthresh,optgrad=1.d-8,optstep=1.d-8;
nocompress;
geomtyp=xyz
bohr
geometry={
3
       121       121
H        0.0000000000        1.6767635949       -1.9385346468
N        0.0000000000       -0.1677750491       -1.1476238104
O        0.0000000000        0.0327792442        1.1355688662
}
 
basis=av5z
  {hf,maxit=500;accu,20;}
{ccsd(t),nocheck,maxit=250;orbital,IGNORE_ERROR;}
