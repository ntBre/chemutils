*** hno CCSD(T)/TZ-F12 QFF
memory, 995, m;
gthresh,energy=1.d-12,zero=1.d-22,oneint=1.d-22,twoint=1.d-22;
gthresh,optgrad=1.d-8,optstep=1.d-8;
nocompress;
geomtyp=xyz
bohr
geometry={
3
        44        44
N        0.0000000000       -0.1563506460       -1.1600219139
H        0.0000000000        1.6647929857       -1.9335362004
O        0.0000000000        0.0333254503        1.1429685233
}
 
basis=avqz
  {hf,maxit=500;accu,20;}
{ccsd(t),nocheck,maxit=250;orbital,IGNORE_ERROR;}
