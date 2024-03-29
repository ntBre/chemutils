*** hno CCSD(T)/TZ-F12 QFF
memory, 995, m;
gthresh,energy=1.d-12,zero=1.d-22,oneint=1.d-22,twoint=1.d-22;
gthresh,optgrad=1.d-8,optstep=1.d-8;
nocompress;
geomtyp=xyz
bohr
geometry={
3
       115       115
H        0.0000000000        1.6755207386       -1.9292183510
N        0.0000000000       -0.1727349890       -1.1470332540
O        0.0000000000        0.0389820404        1.1256620140
}
 
basis=avqz
  {hf,maxit=500;accu,20;}
{ccsd(t),nocheck,maxit=250;orbital,IGNORE_ERROR;}
