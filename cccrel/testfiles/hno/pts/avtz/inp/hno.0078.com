*** hno CCSD(T)/TZ-F12 QFF
memory, 995, m;
gthresh,energy=1.d-12,zero=1.d-22,oneint=1.d-22,twoint=1.d-22;
gthresh,optgrad=1.d-8,optstep=1.d-8;
nocompress;
geomtyp=xyz
bohr
geometry={
3
        78        78
H        0.0000000000        1.6682414571       -1.9299469849
N        0.0000000000       -0.1632684515       -1.1566815234
O        0.0000000000        0.0367947845        1.1360389173
}
 
basis=avtz
  {hf,maxit=500;accu,20;}
{ccsd(t),nocheck,maxit=250;orbital,IGNORE_ERROR;}
