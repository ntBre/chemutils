*** hno CCSD(T)/TZ-F12 QFF
memory, 995, m;
gthresh,energy=1.d-12,zero=1.d-22,oneint=1.d-22,twoint=1.d-22;
gthresh,optgrad=1.d-8,optstep=1.d-8;
nocompress;
geomtyp=xyz
bohr
geometry={
3
        26        26
N        0.0000000000       -0.1541638920       -1.1414636276
H        0.0000000000        1.6614004550       -1.9279841806
O        0.0000000000        0.0345312269        1.1188582172
}
 
basis=vtz-dk
dkroll=0
  {hf,maxit=500;accu,20;}
{ccsd(t),nocheck,maxit=250;orbital,IGNORE_ERROR;}
