*** hno CCSD(T)/TZ-F12 QFF
memory, 995, m;
gthresh,energy=1.d-12,zero=1.d-22,oneint=1.d-22,twoint=1.d-22;
gthresh,optgrad=1.d-8,optstep=1.d-8;
nocompress;
geomtyp=xyz
bohr
geometry={
3
        51        51
H        0.0000000000        1.6650243097       -1.9213150022
N        0.0000000000       -0.1638476599       -1.1418309818
O        0.0000000000        0.0405911403        1.1125563930
}
 
basis=avqz
  {hf,maxit=500;accu,20;}
{ccsd(t),nocheck,maxit=250;orbital,IGNORE_ERROR;}
