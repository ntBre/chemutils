;; TODO reset geometry when slider changes, otherwise can trap on wrong side
;; TODO offsets for labels aren't going to work when rotated
;; TODO coule probably DRY out draw-x,y,z but not sure how exactly. macro?
;; could use parameter to return pen to default
#lang racket/gui

(require racket/gui/base)

(define *bg-color* (make-color 255 255 255))
(define *xangle* (cos (/ (* 5 pi) 4.)))
(define *yangle* (sin (/ (* 5 pi) 4.)))
(define *axis-scale* 0.75)
(define *atom-scale* 0.25)
(define *refresh-rate* 0.5)

(struct vib (freq contribs)
  #:transparent
  #:mutable)

(define (process-freqs freqs)
  (map string->number freqs))

(define (empty-string? str)
  (not (non-empty-string? str)))

(define (dot? str)
  (string-contains? str "."))

(define ptable
  (hash
   "1.0078250" "H"
   "12.0000000" "C"
   "14.0030740" "N"
   "15.9949146" "O"
   ))

(define pcolor
  (hash
   "H" (make-color 221 228 240)
   "C" (make-color 39 45 54)
   "N" (make-color 10 100 245)
   "O" (make-color 255 0 0)))

(define (atom-color atom)
  (hash-ref pcolor atom))

(define (read-output in)
  (let ((in-lxm? #f)
        (head? #f)
        (geom? #f)
        (freqs null)
        (geom null)
        (skip 0)
        (ht (make-hash)))
  (for ((l (in-lines in)))
    (cond
      ((string-contains? l "LXM MATRIX") (set! in-lxm? #t) (set! head? #t))
      ((string-contains? l "LX MATRIX") (set! in-lxm? #f))
      ((string-contains? l "MOLECULAR CARTESIAN GEOMETRY") (set! geom? #t) (set! skip 2))
      ((> skip 0) (set! skip (- skip 1)))
      ((and geom? (empty-string? l)) (set! geom? #f))
      (geom? (set! geom (cons (cdr (string-split l)) geom)))
      ((and in-lxm? (empty-string? l)) (set! head? #t))
      ((and in-lxm? (not (string-contains? l "---")) (dot? l))
       (cond
         (head? (set! head? #f)
                (set! freqs (append freqs (process-freqs (string-split l)))))
         (else
          (let* ((slc (string-split l))
                 (fst (car slc)))
            (hash-set! ht fst (append (hash-ref ht fst null) (cdr slc)))))))))
  (values freqs ht (reverse geom))))

(define-values (freqs contribs geom) (call-with-input-file "spectro.out" read-output))

(set! contribs
      (let ((keys (map number->string (sort (map string->number (hash-keys contribs)) <))))
        (for/list ((k keys))
          (map string->number (hash-ref contribs k)))))

(define (make-vibs freqs contribs)
  (cond
    ((null? freqs) null)
    (else (cons (vib (car freqs) (map car contribs))
                (make-vibs (cdr freqs) (map cdr contribs))))))

(define vibs (make-vibs freqs contribs))

(define (print-vibs)
  (for ((line vibs))
    (displayln line)))

(define (vib-choices)
  (for/list ((v vibs))
    (vib-freq v)))

(define atoms
  (map (lambda (a)
         (hash-ref ptable a))
       (flatten (map (lambda (l)
                       (take-right l 1)) geom))))

(define coords
  (map (lambda (l) (map string->number (take l 3))) geom))

(define my-frame%
  (class frame%
    (define/override (on-subwindow-char rec event)
      (let* ((sel (send list-box get-selections2))
             (sel-max (- (send list-box number-of-visible-items) 1)) 
             (sli (send slider get-value)))
        (cond
          ((equal? (send event get-key-code) #\q)
           (exit))
          ((equal? (send event get-key-code) 'escape)
           (send list-box clear-select))
          ((equal? (send event get-key-code) #\j)
           (cond
             ((null? sel) (send list-box select 0))
             ((< sel sel-max)
              (send list-box select (+ 1 sel)))))
          ((equal? (send event get-key-code) #\k)
           (cond
             ((null? sel) (send list-box select sel-max))
             ((> sel 0) (send list-box select (- sel 1)))))
          ((equal? (send event get-key-code) #\l)
           (when (< sli slide-max)
             (send slider set-value (+ (send slider get-value) 5))))
          ((equal? (send event get-key-code) #\h)
           (when (> sli slide-min)
             (send slider set-value (- (send slider get-value) 5)))))))
    (super-new)))

(define frame (new my-frame%
                   (label "spectroscope")
                   (width 500)
                   (height 500)))

(define panel (new horizontal-panel%
                   (style '(border))
                   (parent frame)
                   (alignment '(center center))))

(define left-panel (new vertical-panel%
                        (parent panel)))

(define (center canvas)
  (let ((width (send canvas get-width))
        (height (send canvas get-height)))
    (values (/ width 2.) (/ height 2.))))

(define (extent canvas)
  (values
   (send canvas get-width)
   (send canvas get-height)))

(define dash-pen (new pen% (style 'long-dash)))
(define def-pen (new pen%))
(define def-brush (new brush%))

(define (x-help-lines dc w h)
  (let ((wend (- w (* *axis-scale* w)))
        (hend h))
    (send dc set-pen dash-pen)
    (send dc draw-line
          w h
          wend hend))
  (let ((wend w)
        (hend (+ h (* *axis-scale* h))))
    (send dc draw-line
          w h
          wend hend))
  (send dc set-pen def-pen))

(define (draw-x dc w h)
  ;; (x-help-lines dc w h)
  (let-values (((wbeg hbeg) (cart->2d canvas 0 0 0 *axis-scale*))
               ((wend hend) (cart->2d canvas 1 0 0 *axis-scale*)))
    (send dc draw-line
          wbeg hbeg
          wend hend)
    (define-values (woff hoff d a) (send dc get-text-extent "x"))
    (send dc draw-text "x" (- wend woff) hend)))

(define (draw-y dc w h)
  (let-values (((wbeg hbeg) (cart->2d canvas 0 0 0 *axis-scale*))
               ((wend hend) (cart->2d canvas 0 1 0 *axis-scale*)))
    (send dc draw-line
          wbeg hbeg
          wend hend)
    (define-values (woff hoff d a) (send dc get-text-extent "y"))
    (send dc draw-text "y" (+ wend woff) (- hend (/ hoff 2)))))

(define (draw-z dc w h)
  (let-values (((wbeg hbeg) (cart->2d canvas 0 0 0 *axis-scale*))
               ((wend hend) (cart->2d canvas 0 0 1 *axis-scale*)))
    (send dc draw-line
          wbeg hbeg
          wend hend)
    (define-values (woff hoff d a) (send dc get-text-extent "z"))
    (send dc draw-text "z" (- wend (/ woff 2)) (- hend hoff))))

(define (cart->2d canvas x y z scale)
  (let-values (((cw ch) (center canvas))
               ((mw mh) (extent canvas)))
  (values
   (+ cw (* (+ y (* x *xangle*)) (- mw cw) scale))
   (- ch (* (+ z (* x *yangle*)) (- mh ch) scale)))))

(define (draw-atom canvas dc atom x y z)
    (let-values (((w h) (cart->2d canvas x y z *atom-scale*)))
      (send dc set-brush (atom-color atom) 'solid)
      (send dc draw-ellipse w h 20 20))
  (send dc set-brush def-brush))

(define (draw-axes canvas dc)
  ;; TODO update these functions when I introduce rotation
  (define-values (w h) (center canvas))
  (draw-x dc w h)
  (values (draw-y dc w h) (draw-z dc w h)))

(define (draw-geom canvas dc atoms coords)
  (for ((atom atoms) (coord coords))
    (apply draw-atom canvas dc atom coord)))

(define my-list-box%
  (class list-box%
    (define/public (clear-select)
      (let ((test (send this get-selections2)))
      (unless (null? test) (send this select test #f))))
    (define/public (get-selections2)
      (let* ((sel (send this get-selections)))
        (if (not (null? sel)) (car sel) null)))
    (super-new)))

(define canvas (new canvas%
                    (parent left-panel)
                    (min-width 500)
                    (min-height 500)
                    (paint-callback
                     (lambda (canvas dc)
                       (draw-canvas canvas dc)))))

(send canvas focus)
(define dc (send canvas get-dc))
(send canvas set-canvas-background *bg-color*)

(define slide-min 0)
(define slide-max 100)
(define slide-init 50)

(define (magnitude)
  (/ (send slider get-value) slide-init))

(define slider (new slider%
                    (label "Magnitude")
                    (parent left-panel)
                    (min-value slide-min)
                    (max-value slide-max)
                    (init-value 50)))

(define right-panel (new vertical-panel%
                         (parent panel)))

(define list-box (new my-list-box%
                      (parent right-panel)
                      (min-width 100)
                      (label #f)
                      (choices (map number->string (vib-choices)))
                      (columns '("Frequencies"))
                      (style '(single column-headers))))

(send frame show #t)

(define (draw-canvas canvas dc)
  (draw-axes canvas dc)
  (draw-geom canvas dc atoms coords))

(define (resplit lst)
  (cond
    ((null? lst) null)
    (else (cons (take lst 3) (resplit (drop lst 3))))))

(define vibrate
  (let ((n 0)
        (steps (list + - - +)))
    (lambda (contribs (reset? #f))
      (when reset? (set! n 0))
      (let ((op (list-ref steps (remainder n 4)))
            (mag (magnitude)))
        (set! coords (resplit (map op
                                   (flatten coords)
                                   (map (lambda (c)
                                          (* mag c)) contribs)))))
      (set! n (add1 n)))))

(define ref-coords coords)

(define prev
  (let ((hold null))
    (lambda (next)
      (begin0
          (equal? hold next)
        (set! hold next)))))

(define (loop)
  (let* ((sel (send list-box get-selections2))
         (diff? (prev sel)))
    (unless diff? ;; if new selection made, reset coords
      (set! coords ref-coords))
    (unless (null? sel) ;; only vibrate when there is a selection
      (if diff? 
          (vibrate (vib-contribs (list-ref vibs sel)) #f)
          (vibrate (vib-contribs (list-ref vibs sel)) #t))))
  (send canvas refresh)
  (send canvas on-paint)
  (sleep/yield *refresh-rate*)
  (loop))

(loop)

;; maintain an internal counter in draw-geom, cycle through (+contrib,
;; 0, -contrib) on repeated calls; take intensity from slider

;; extract information from spectro output - CHECK
;; - need original geometry - split into coords and atoms
;; - LXM matrix for vibrations - in vibs

;; make a GUI window

;; map from weights in geom to atom colors

;; make an image of the geometry
;; - svg or png or something - gif for freqs
;; - display that image with the gui framework

;; animate the image using LXM matrix vibrations
