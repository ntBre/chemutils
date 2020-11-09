#lang racket

(struct vib (freq contribs)
  #:transparent
  #:mutable)

(define (process-freqs freqs)
  (map string->number freqs))

(define (empty-string? str)
  (not (non-empty-string? str)))

(define (dot? str)
  (string-contains? str "."))

(define (read-output in)
  (define in-lxm? #f)
  (define head? #f)
  (define freqs null)
  (define ht (make-hash))
  (for ((l (in-lines in)))
    (cond
      ((string-contains? l "LXM MATRIX") (set! in-lxm? #t) (set! head? #t))
      ((string-contains? l "LX MATRIX") (set! in-lxm? #f))
      ((and in-lxm? (empty-string? l)) (set! head? #t))
      ((and in-lxm? (not (string-contains? l "---")) (dot? l))
       (cond
         (head? (set! head? #f)
                (set! freqs (append freqs (process-freqs (string-split l)))))
         (else
          (let* ((slc (string-split l))
                 (fst (car slc)))
            (hash-set! ht fst (append (hash-ref ht fst null) (cdr slc)))))))))
  (values freqs ht))

(define-values (freqs contribs) (call-with-input-file "spectro.out" read-output))

(set! contribs
      (let ((keys (map number->string (sort (map string->number (hash-keys contribs)) <))))
        (for/list ((k keys))
          (hash-ref contribs k))))

(define vibs null)
(define (make-vibs freqs contribs)
  (set! vibs (append vibs (list (vib (car freqs) (map car contribs)))))
  (cond
    ((null? (cdr freqs)) null)
    (else (make-vibs (cdr freqs) (map cdr contribs)))))

(make-vibs freqs contribs)
(for ((line vibs))
  (displayln line))

;; extract information from spectro output
;; - need original geometry
;; - LXM matrix for vibrations

;; make an image of the geometry
;; - svg or png or something - gif for freqs
;; - display that image with the gui framework

;; animate the image using LXM matrix vibrations

;; when you hit "LXM MATRIX" you are in the right part
;; want to gather the line with the frequencies, put those in a list called freqs
;; there is another set of these after a blank line
;; after that, want to gather contributions into a 2d array
;; can gather the first part easily by processing each line at a time
;; but after the blank line, you need a way to append to the earlier lines
