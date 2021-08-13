*** hno CCSD(T)/TZ-F12 QFF
memory, 995, m;
gthresh,energy=1.d-12,zero=1.d-22,oneint=1.d-22,twoint=1.d-22;
gthresh,optgrad=1.d-8,optstep=1.d-8;
nocompress;
geomtyp=xyz
bohr
geometry={
3
        83        83
N        0.0000000000       -0.1575577763       -1.1549269869
H        0.0000000000        1.6691664896       -1.9394309548
O        0.0000000000        0.0301590766        1.1437683507
}
 
basis=avqz
  {hf,maxit=500;accu,20;}
{ccsd(t),nocheck,maxit=250;orbital,IGNORE_ERROR;}