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

 no=1.20786572 ANG
 nh=1.05203411 ANG
 hno=108.1889935 DEG

basis=sto-3g
  {hf,maxit=500;wf,charge=0,spin=0;accu,20;}
