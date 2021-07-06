*** AlONC, CCSD(T)/aug-cc-pV5Z QFF
memory, 995, m;
gthresh,energy=1.d-12,zero=1.d-22,oneint=1.d-22,twoint=1.d-22;
gthresh,optgrad=1.d-8,optstep=1.d-8;
nocompress;
geometry={

N
O 1 no
H 1 nh 2 hno
}

 no=1.2 ANG
 nh=0.8 ANG
 hno=105.0 DEG

basis=av5z,al=aug-cc-pV(5+d)Z
  {hf,maxit=500;wf,charge=0,spin=0;accu,20;}
  {ccsd(t),maxit=250;wf,charge=0,spin=0;orbital,IGNORE_ERROR;}
 {optg,grms=1.d-8,srms=1.d-8}
