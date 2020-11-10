#lang racket/gui

(require racket/gui/base)

(define *bg-color* (make-color 255 255 255))
(define *xangle* (cos (/ (* 5 pi) 4.)))
(define *yangle* (sin (/ (* 5 pi) 4.)))
(define *axis-scale* 0.75)

(struct vib (freq contribs)
  #:transparent
  #:mutable)

(define (process-freqs freqs)
  (map string->number freqs))

(define (empty-string? str)
  (not (non-empty-string? str)))

(define (dot? str)
  (string-contains? str "."))

;; TODO should probably map this to color instead of name
(define ptable
  (hash
   "1.0078250" "H"
   "12.0000000" "C"
   "14.0030740" "N"
   "15.9949146" "O"
   ))


(define (read-output in)
  (define in-lxm? #f)
  (define head? #f)
  (define geom? #f)
  (define freqs null)
  (define geom null)
  (define skip 0)
  (define ht (make-hash))
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
  (values freqs ht (reverse geom)))

(define-values (freqs contribs geom) (call-with-input-file "spectro.out" read-output))

(set! contribs
      (let ((keys (map number->string (sort (map string->number (hash-keys contribs)) <))))
        (for/list ((k keys))
          (hash-ref contribs k))))

;; need to look up how to do a recursive accumulator
(define vibs null)
(define (make-vibs freqs contribs)
  (set! vibs (append vibs (list (vib (car freqs) (map car contribs)))))
  (cond
    ((null? (cdr freqs)) (void))
    (else (make-vibs (cdr freqs) (map cdr contribs)))))

(make-vibs freqs contribs)

(define (print-vibs)
  (for ((line vibs))
    (displayln line)))

(define (vib-choices)
  (for/list ((v vibs))
    (vib-freq v)))

(define atoms
  (map (lambda (a)
         (hash-ref ptable a)) (flatten (map (lambda (l)
                                              (take-right l 1)) geom))))

(define coords
  (map (lambda (l) (map string->number (take l 3))) geom))

(define frame (new frame%
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

(define dash-pen (new pen% (style 'long-dash)))
(define def-pen (new pen%))

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
  (let ((wend (+ w (* w *axis-scale* *xangle*)))
        (hend (- h (* h *axis-scale* *yangle*))))
    (send dc draw-line
          w h
          wend hend)
    (define-values (woff hoff d a) (send dc get-text-extent "x"))
    (send dc draw-text "x" (- wend woff) hend)))

(define (draw-y dc w h)
  (let ((wend (+ w (* w *axis-scale*)))
        (hend h))
    (send dc draw-line
          w h
          wend hend)
    (define-values (woff hoff d a) (send dc get-text-extent "y"))
    ;; add full offset in width and subtract half in height
    (send dc draw-text "y" (+ wend woff) (- hend (/ hoff 2)))
    wend))

(define (draw-z dc w h)
  (let ((wend w)
        (hend (- h (* h *axis-scale*)))
        (maxh (+ h (* h *axis-scale*))))
    (send dc draw-line
          w h
          wend hend )
    (define-values (woff hoff d a) (send dc get-text-extent "z"))
    (send dc draw-text "z" (- wend (/ woff 2)) (- hend hoff))
    maxh))

(define (square x)
  (* x x))

(define (vec-len x y z)
  (sqrt (foldl (lambda (a b)
                 (+ (square a) b)) 0.0 (list x y z))))

(define (cart->2d x y z)
  (displayln (list x y z)))

(define (draw-atom dc canvas maxw maxh atom x y z)
  (let* ((len (vec-len x y z))
         (vec (map (lambda (a)
                     (/ a len)) (list x y z))))
    ;; (define-values 
    ;; TODO convert the vectors to 2d coordinates and draw the atoms
    (display atom)
    (apply cart->2d vec)))

(define (draw-axes canvas dc)
  ;; TODO update these functions when I introduce rotation
  (define-values (w h) (center canvas))
  (draw-x dc w h)
  (values (draw-y dc w h) (draw-z dc w h)))

;; (displayln (car coords))
(define (draw-geom canvas dc maxw maxh)
  ;; need to scale the whole geometry
  (for ((atom atoms) (coord coords))
    (apply draw-atom canvas dc maxw maxh atom coord)))

(define my-canvas%
  (class canvas%
    (define/override (on-char event)
      (cond
        ((equal? (send event get-key-code) #\q) (exit))))
    (super-new)))

(define (draw-canvas canvas dc)
  (define-values (maxw maxh) (draw-axes canvas dc))
  (draw-geom canvas dc maxw maxh))

(define canvas (new my-canvas%
                    (parent left-panel)
                    (min-width 500)
                    (min-height 500)
                    (paint-callback
                     (lambda (canvas dc)
                       (draw-canvas canvas dc)))))

(send canvas focus)

(send canvas set-canvas-background *bg-color*)

(define slider (new slider%
                    (label "Magnitude")
                    (parent left-panel)
                    (min-value 0)
                    (max-value 100)
                    (init-value 50)))

(define right-panel (new vertical-panel%
                         (parent panel)))

(define list-box (new list-box%
                      (parent right-panel)
                      (label #f)
                      (choices (map number->string (vib-choices)))
                      (columns '("Frequencies"))
                      (style '(single column-headers))))

(send frame show #t)

;; extract information from spectro output - CHECK
;; - need original geometry - split into coords and atoms
;; - LXM matrix for vibrations - in vibs

;; make a GUI window

;; map from weights in geom to atom colors

;; make an image of the geometry
;; - svg or png or something - gif for freqs
;; - display that image with the gui framework

;; animate the image using LXM matrix vibrations
