*** hcn CCSD(T)/TZ-F12 QFF
memory, 995, m;
gthresh,energy=1.d-12,zero=1.d-22,oneint=1.d-22,twoint=1.d-22;
gthresh,optgrad=1.d-8,optstep=1.d-8;
nocompress;
geomtyp=xyz
bohr
geometry={
3
        41        41
N       -0.0046256470        0.0000000675        1.1309517820
C        0.0096312361       -0.0000000018       -1.0573000722
H       -0.0050055891       -0.0000000657       -3.0707216918
}
 
basis=av5z
  {hf,maxit=500;accu,20;}
{ccsd(t),nocheck,maxit=250;orbital,IGNORE_ERROR;}
