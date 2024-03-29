*** hno CCSD(T)/TZ-F12 QFF
memory, 995, m;
gthresh,energy=1.d-12,zero=1.d-22,oneint=1.d-22,twoint=1.d-22;
gthresh,optgrad=1.d-8,optstep=1.d-8;
nocompress;
geomtyp=xyz
bohr
geometry={
3
       110       110
N        0.0000000000       -0.1628622964       -1.1530925039
H        0.0000000000        1.6733011401       -1.9395592532
O        0.0000000000        0.0313289463        1.1420621661
}
 
basis=vtz-dk
dkroll=1
  {hf,maxit=500;accu,20;}
{ccsd(t),nocheck,maxit=250;orbital,IGNORE_ERROR;}
