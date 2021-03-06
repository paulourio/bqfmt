// Code generated by gocc; DO NOT EDIT.

package lexer

import (
	"io/ioutil"
	"unicode/utf8"

	"github.com/paulourio/bqfmt/zetasql/token"
)

const (
	NoState    = -1
	NumStates  = 1151
	NumSymbols = 2870
)

type Lexer struct {
	src     []byte
	pos     int
	line    int
	column  int
	Context token.Context
}

func NewLexer(src []byte) *Lexer {
	lexer := &Lexer{
		src:     src,
		pos:     0,
		line:    1,
		column:  1,
		Context: nil,
	}
	return lexer
}

// SourceContext is a simple instance of a token.Context which
// contains the name of the source file.
type SourceContext struct {
	Filepath string
}

func (s *SourceContext) Source() string {
	return s.Filepath
}

func NewLexerFile(fpath string) (*Lexer, error) {
	src, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	lexer := NewLexer(src)
	lexer.Context = &SourceContext{Filepath: fpath}
	return lexer, nil
}

func (l *Lexer) Scan() (tok *token.Token) {
	tok = &token.Token{}
	if l.pos >= len(l.src) {
		tok.Type = token.EOF
		tok.Pos.Offset, tok.Pos.Line, tok.Pos.Column = l.pos, l.line, l.column
		tok.Pos.Context = l.Context
		return
	}
	start, startLine, startColumn, end := l.pos, l.line, l.column, 0
	tok.Type = token.INVALID
	state, rune1, size := 0, rune(-1), 0
	for state != -1 {
		if l.pos >= len(l.src) {
			rune1 = -1
		} else {
			rune1, size = utf8.DecodeRune(l.src[l.pos:])
			l.pos += size
		}

		nextState := -1
		if rune1 != -1 {
			nextState = TransTab[state](rune1)
		}
		state = nextState

		if state != -1 {

			switch rune1 {
			case '\n':
				l.line++
				l.column = 1
			case '\r':
				l.column = 1
			case '\t':
				l.column += 4
			default:
				l.column++
			}

			switch {
			case ActTab[state].Accept != -1:
				tok.Type = ActTab[state].Accept
				end = l.pos
			case ActTab[state].Ignore != "":
				start, startLine, startColumn = l.pos, l.line, l.column
				state = 0
				if start >= len(l.src) {
					tok.Type = token.EOF
				}

			}
		} else {
			if tok.Type == token.INVALID {
				end = l.pos
			}
		}
	}
	if end > start {
		l.pos = end
		tok.Lit = l.src[start:end]
	} else {
		tok.Lit = []byte{}
	}
	tok.Pos.Offset, tok.Pos.Line, tok.Pos.Column = start, startLine, startColumn
	tok.Pos.Context = l.Context

	return
}

func (l *Lexer) Reset() {
	l.pos = 0
}

/*
Lexer symbols:
0: '_'
1: '_'
2: '_'
3: 'f'
4: 'F'
5: 'a'
6: 'A'
7: 'l'
8: 'L'
9: 's'
10: 'S'
11: 'e'
12: 'E'
13: 't'
14: 'T'
15: 'r'
16: 'R'
17: 'u'
18: 'U'
19: 'e'
20: 'E'
21: 'a'
22: 'A'
23: 'b'
24: 'B'
25: 'o'
26: 'O'
27: 'r'
28: 'R'
29: 't'
30: 'T'
31: 'a'
32: 'A'
33: 'c'
34: 'C'
35: 'c'
36: 'C'
37: 'e'
38: 'E'
39: 's'
40: 'S'
41: 's'
42: 'S'
43: 'a'
44: 'A'
45: 'c'
46: 'C'
47: 't'
48: 'T'
49: 'i'
50: 'I'
51: 'o'
52: 'O'
53: 'n'
54: 'N'
55: 'a'
56: 'A'
57: 'd'
58: 'D'
59: 'd'
60: 'D'
61: 'a'
62: 'A'
63: 'g'
64: 'G'
65: 'g'
66: 'G'
67: 'r'
68: 'R'
69: 'e'
70: 'E'
71: 'g'
72: 'G'
73: 'a'
74: 'A'
75: 't'
76: 'T'
77: 'e'
78: 'E'
79: 'a'
80: 'A'
81: 'l'
82: 'L'
83: 'l'
84: 'L'
85: 'a'
86: 'A'
87: 'l'
88: 'L'
89: 't'
90: 'T'
91: 'e'
92: 'E'
93: 'r'
94: 'R'
95: 'a'
96: 'A'
97: 'n'
98: 'N'
99: 'a'
100: 'A'
101: 'l'
102: 'L'
103: 'y'
104: 'Y'
105: 'z'
106: 'Z'
107: 'e'
108: 'E'
109: 'a'
110: 'A'
111: 'n'
112: 'N'
113: 'd'
114: 'D'
115: 'a'
116: 'A'
117: 'n'
118: 'N'
119: 'o'
120: 'O'
121: 'n'
122: 'N'
123: 'y'
124: 'Y'
125: 'm'
126: 'M'
127: 'i'
128: 'I'
129: 'z'
130: 'Z'
131: 'a'
132: 'A'
133: 't'
134: 'T'
135: 'i'
136: 'I'
137: 'o'
138: 'O'
139: 'n'
140: 'N'
141: 'a'
142: 'A'
143: 'r'
144: 'R'
145: 'r'
146: 'R'
147: 'a'
148: 'A'
149: 'y'
150: 'Y'
151: 'a'
152: 'A'
153: 's'
154: 'S'
155: 'a'
156: 'A'
157: 's'
158: 'S'
159: 'c'
160: 'C'
161: 'a'
162: 'A'
163: 's'
164: 'S'
165: 's'
166: 'S'
167: 'e'
168: 'E'
169: 'r'
170: 'R'
171: 't'
172: 'T'
173: 'a'
174: 'A'
175: 't'
176: 'T'
177: 'b'
178: 'B'
179: 'a'
180: 'A'
181: 't'
182: 'T'
183: 'c'
184: 'C'
185: 'h'
186: 'H'
187: 'b'
188: 'B'
189: 'e'
190: 'E'
191: 'g'
192: 'G'
193: 'i'
194: 'I'
195: 'n'
196: 'N'
197: 'b'
198: 'B'
199: 'e'
200: 'E'
201: 't'
202: 'T'
203: 'w'
204: 'W'
205: 'e'
206: 'E'
207: 'e'
208: 'E'
209: 'n'
210: 'N'
211: 'b'
212: 'B'
213: 'i'
214: 'I'
215: 'g'
216: 'G'
217: 'd'
218: 'D'
219: 'e'
220: 'E'
221: 'c'
222: 'C'
223: 'i'
224: 'I'
225: 'm'
226: 'M'
227: 'a'
228: 'A'
229: 'l'
230: 'L'
231: 'b'
232: 'B'
233: 'i'
234: 'I'
235: 'g'
236: 'G'
237: 'n'
238: 'N'
239: 'u'
240: 'U'
241: 'm'
242: 'M'
243: 'e'
244: 'E'
245: 'r'
246: 'R'
247: 'i'
248: 'I'
249: 'c'
250: 'C'
251: 'b'
252: 'B'
253: 'r'
254: 'R'
255: 'e'
256: 'E'
257: 'a'
258: 'A'
259: 'k'
260: 'K'
261: 'b'
262: 'B'
263: 'y'
264: 'Y'
265: 'c'
266: 'C'
267: 'a'
268: 'A'
269: 'l'
270: 'L'
271: 'l'
272: 'L'
273: 'c'
274: 'C'
275: 'a'
276: 'A'
277: 's'
278: 'S'
279: 'c'
280: 'C'
281: 'a'
282: 'A'
283: 'd'
284: 'D'
285: 'e'
286: 'E'
287: 'c'
288: 'C'
289: 'a'
290: 'A'
291: 's'
292: 'S'
293: 'e'
294: 'E'
295: 'c'
296: 'C'
297: 'a'
298: 'A'
299: 's'
300: 'S'
301: 't'
302: 'T'
303: 'c'
304: 'C'
305: 'h'
306: 'H'
307: 'e'
308: 'E'
309: 'c'
310: 'C'
311: 'k'
312: 'K'
313: 'c'
314: 'C'
315: 'l'
316: 'L'
317: 'a'
318: 'A'
319: 'm'
320: 'M'
321: 'p'
322: 'P'
323: 'e'
324: 'E'
325: 'd'
326: 'D'
327: 'c'
328: 'C'
329: 'l'
330: 'L'
331: 'o'
332: 'O'
333: 'n'
334: 'N'
335: 'e'
336: 'E'
337: 'c'
338: 'C'
339: 'l'
340: 'L'
341: 'u'
342: 'U'
343: 's'
344: 'S'
345: 't'
346: 'T'
347: 'e'
348: 'E'
349: 'r'
350: 'R'
351: 'c'
352: 'C'
353: 'o'
354: 'O'
355: 'l'
356: 'L'
357: 'u'
358: 'U'
359: 'm'
360: 'M'
361: 'n'
362: 'N'
363: 'c'
364: 'C'
365: 'o'
366: 'O'
367: 'l'
368: 'L'
369: 'u'
370: 'U'
371: 'm'
372: 'M'
373: 'n'
374: 'N'
375: 's'
376: 'S'
377: 'c'
378: 'C'
379: 'o'
380: 'O'
381: 'm'
382: 'M'
383: 'm'
384: 'M'
385: 'i'
386: 'I'
387: 't'
388: 'T'
389: 'c'
390: 'C'
391: 'o'
392: 'O'
393: 'n'
394: 'N'
395: 'n'
396: 'N'
397: 'e'
398: 'E'
399: 'c'
400: 'C'
401: 't'
402: 'T'
403: 'i'
404: 'I'
405: 'o'
406: 'O'
407: 'n'
408: 'N'
409: 'c'
410: 'C'
411: 'o'
412: 'O'
413: 'n'
414: 'N'
415: 's'
416: 'S'
417: 't'
418: 'T'
419: 'a'
420: 'A'
421: 'n'
422: 'N'
423: 't'
424: 'T'
425: 'c'
426: 'C'
427: 'o'
428: 'O'
429: 'n'
430: 'N'
431: 's'
432: 'S'
433: 't'
434: 'T'
435: 'r'
436: 'R'
437: 'a'
438: 'A'
439: 'i'
440: 'I'
441: 'n'
442: 'N'
443: 't'
444: 'T'
445: 'c'
446: 'C'
447: 'o'
448: 'O'
449: 'n'
450: 'N'
451: 't'
452: 'T'
453: 'i'
454: 'I'
455: 'n'
456: 'N'
457: 'u'
458: 'U'
459: 'e'
460: 'E'
461: 'c'
462: 'C'
463: 'o'
464: 'O'
465: 'p'
466: 'P'
467: 'y'
468: 'Y'
469: 'c'
470: 'C'
471: 'r'
472: 'R'
473: 'o'
474: 'O'
475: 's'
476: 'S'
477: 's'
478: 'S'
479: 'c'
480: 'C'
481: 'u'
482: 'U'
483: 'r'
484: 'R'
485: 'r'
486: 'R'
487: 'e'
488: 'E'
489: 'n'
490: 'N'
491: 't'
492: 'T'
493: 'd'
494: 'D'
495: 'a'
496: 'A'
497: 't'
498: 'T'
499: 'a'
500: 'A'
501: 'd'
502: 'D'
503: 'a'
504: 'A'
505: 't'
506: 'T'
507: 'a'
508: 'A'
509: 'b'
510: 'B'
511: 'a'
512: 'A'
513: 's'
514: 'S'
515: 'e'
516: 'E'
517: 'd'
518: 'D'
519: 'a'
520: 'A'
521: 't'
522: 'T'
523: 'e'
524: 'E'
525: 'd'
526: 'D'
527: 'a'
528: 'A'
529: 't'
530: 'T'
531: 'e'
532: 'E'
533: 't'
534: 'T'
535: 'i'
536: 'I'
537: 'm'
538: 'M'
539: 'e'
540: 'E'
541: 'd'
542: 'D'
543: 'e'
544: 'E'
545: 'c'
546: 'C'
547: 'i'
548: 'I'
549: 'm'
550: 'M'
551: 'a'
552: 'A'
553: 'l'
554: 'L'
555: 'd'
556: 'D'
557: 'e'
558: 'E'
559: 'c'
560: 'C'
561: 'l'
562: 'L'
563: 'a'
564: 'A'
565: 'r'
566: 'R'
567: 'e'
568: 'E'
569: 'd'
570: 'D'
571: 'e'
572: 'E'
573: 'f'
574: 'F'
575: 'i'
576: 'I'
577: 'n'
578: 'N'
579: 'e'
580: 'E'
581: 'r'
582: 'R'
583: 'd'
584: 'D'
585: 'e'
586: 'E'
587: 'l'
588: 'L'
589: 'e'
590: 'E'
591: 't'
592: 'T'
593: 'e'
594: 'E'
595: 'd'
596: 'D'
597: 'e'
598: 'E'
599: 's'
600: 'S'
601: 'c'
602: 'C'
603: 'd'
604: 'D'
605: 'e'
606: 'E'
607: 's'
608: 'S'
609: 'c'
610: 'C'
611: 'r'
612: 'R'
613: 'i'
614: 'I'
615: 'b'
616: 'B'
617: 'e'
618: 'E'
619: 'd'
620: 'D'
621: 'e'
622: 'E'
623: 's'
624: 'S'
625: 'c'
626: 'C'
627: 'r'
628: 'R'
629: 'i'
630: 'I'
631: 'p'
632: 'P'
633: 't'
634: 'T'
635: 'o'
636: 'O'
637: 'r'
638: 'R'
639: 'd'
640: 'D'
641: 'e'
642: 'E'
643: 't'
644: 'T'
645: 'e'
646: 'E'
647: 'r'
648: 'R'
649: 'm'
650: 'M'
651: 'i'
652: 'I'
653: 'n'
654: 'N'
655: 'i'
656: 'I'
657: 's'
658: 'S'
659: 't'
660: 'T'
661: 'i'
662: 'I'
663: 'c'
664: 'C'
665: 'd'
666: 'D'
667: 'i'
668: 'I'
669: 's'
670: 'S'
671: 't'
672: 'T'
673: 'i'
674: 'I'
675: 'n'
676: 'N'
677: 'c'
678: 'C'
679: 't'
680: 'T'
681: 'd'
682: 'D'
683: 'o'
684: 'O'
685: 'd'
686: 'D'
687: 'r'
688: 'R'
689: 'o'
690: 'O'
691: 'p'
692: 'P'
693: 'e'
694: 'E'
695: 'l'
696: 'L'
697: 's'
698: 'S'
699: 'e'
700: 'E'
701: 'e'
702: 'E'
703: 'l'
704: 'L'
705: 's'
706: 'S'
707: 'e'
708: 'E'
709: 'i'
710: 'I'
711: 'f'
712: 'F'
713: 'e'
714: 'E'
715: 'n'
716: 'N'
717: 'd'
718: 'D'
719: 'e'
720: 'E'
721: 'n'
722: 'N'
723: 'f'
724: 'F'
725: 'o'
726: 'O'
727: 'r'
728: 'R'
729: 'c'
730: 'C'
731: 'e'
732: 'E'
733: 'd'
734: 'D'
735: 'e'
736: 'E'
737: 'r'
738: 'R'
739: 'r'
740: 'R'
741: 'o'
742: 'O'
743: 'r'
744: 'R'
745: 'e'
746: 'E'
747: 'x'
748: 'X'
749: 'c'
750: 'C'
751: 'e'
752: 'E'
753: 'p'
754: 'P'
755: 't'
756: 'T'
757: 'e'
758: 'E'
759: 'x'
760: 'X'
761: 'c'
762: 'C'
763: 'e'
764: 'E'
765: 'p'
766: 'P'
767: 't'
768: 'T'
769: 'i'
770: 'I'
771: 'o'
772: 'O'
773: 'n'
774: 'N'
775: 'e'
776: 'E'
777: 'x'
778: 'X'
779: 'e'
780: 'E'
781: 'c'
782: 'C'
783: 'u'
784: 'U'
785: 't'
786: 'T'
787: 'e'
788: 'E'
789: 'e'
790: 'E'
791: 'x'
792: 'X'
793: 'i'
794: 'I'
795: 's'
796: 'S'
797: 't'
798: 'T'
799: 's'
800: 'S'
801: 'e'
802: 'E'
803: 'x'
804: 'X'
805: 'p'
806: 'P'
807: 'l'
808: 'L'
809: 'a'
810: 'A'
811: 'i'
812: 'I'
813: 'n'
814: 'N'
815: 'e'
816: 'E'
817: 'x'
818: 'X'
819: 'p'
820: 'P'
821: 'o'
822: 'O'
823: 'r'
824: 'R'
825: 't'
826: 'T'
827: 'e'
828: 'E'
829: 'x'
830: 'X'
831: 't'
832: 'T'
833: 'e'
834: 'E'
835: 'r'
836: 'R'
837: 'n'
838: 'N'
839: 'a'
840: 'A'
841: 'l'
842: 'L'
843: 'e'
844: 'E'
845: 'x'
846: 'X'
847: 't'
848: 'T'
849: 'r'
850: 'R'
851: 'a'
852: 'A'
853: 'c'
854: 'C'
855: 't'
856: 'T'
857: 'f'
858: 'F'
859: 'i'
860: 'I'
861: 'l'
862: 'L'
863: 'e'
864: 'E'
865: 's'
866: 'S'
867: 'f'
868: 'F'
869: 'i'
870: 'I'
871: 'l'
872: 'L'
873: 'l'
874: 'L'
875: 'f'
876: 'F'
877: 'i'
878: 'I'
879: 'l'
880: 'L'
881: 't'
882: 'T'
883: 'e'
884: 'E'
885: 'r'
886: 'R'
887: '_'
888: '_'
889: 'f'
890: 'F'
891: 'i'
892: 'I'
893: 'e'
894: 'E'
895: 'l'
896: 'L'
897: 'd'
898: 'D'
899: 's'
900: 'S'
901: 'f'
902: 'F'
903: 'i'
904: 'I'
905: 'l'
906: 'L'
907: 't'
908: 'T'
909: 'e'
910: 'E'
911: 'r'
912: 'R'
913: 'f'
914: 'F'
915: 'i'
916: 'I'
917: 'r'
918: 'R'
919: 's'
920: 'S'
921: 't'
922: 'T'
923: 'f'
924: 'F'
925: 'o'
926: 'O'
927: 'l'
928: 'L'
929: 'l'
930: 'L'
931: 'o'
932: 'O'
933: 'w'
934: 'W'
935: 'i'
936: 'I'
937: 'n'
938: 'N'
939: 'g'
940: 'G'
941: 'f'
942: 'F'
943: 'o'
944: 'O'
945: 'r'
946: 'R'
947: 'e'
948: 'E'
949: 'i'
950: 'I'
951: 'g'
952: 'G'
953: 'n'
954: 'N'
955: 'f'
956: 'F'
957: 'o'
958: 'O'
959: 'r'
960: 'R'
961: 'm'
962: 'M'
963: 'a'
964: 'A'
965: 't'
966: 'T'
967: 'f'
968: 'F'
969: 'r'
970: 'R'
971: 'o'
972: 'O'
973: 'm'
974: 'M'
975: 'f'
976: 'F'
977: 'u'
978: 'U'
979: 'l'
980: 'L'
981: 'l'
982: 'L'
983: 'f'
984: 'F'
985: 'u'
986: 'U'
987: 'n'
988: 'N'
989: 'c'
990: 'C'
991: 't'
992: 'T'
993: 'i'
994: 'I'
995: 'o'
996: 'O'
997: 'n'
998: 'N'
999: 'g'
1000: 'G'
1001: 'e'
1002: 'E'
1003: 'n'
1004: 'N'
1005: 'e'
1006: 'E'
1007: 'r'
1008: 'R'
1009: 'a'
1010: 'A'
1011: 't'
1012: 'T'
1013: 'e'
1014: 'E'
1015: 'd'
1016: 'D'
1017: 'g'
1018: 'G'
1019: 'r'
1020: 'R'
1021: 'a'
1022: 'A'
1023: 'n'
1024: 'N'
1025: 't'
1026: 'T'
1027: 'g'
1028: 'G'
1029: 'r'
1030: 'R'
1031: 'o'
1032: 'O'
1033: 'u'
1034: 'U'
1035: 'p'
1036: 'P'
1037: '_'
1038: '_'
1039: 'r'
1040: 'R'
1041: 'o'
1042: 'O'
1043: 'w'
1044: 'W'
1045: 's'
1046: 'S'
1047: 'g'
1048: 'G'
1049: 'r'
1050: 'R'
1051: 'o'
1052: 'O'
1053: 'u'
1054: 'U'
1055: 'p'
1056: 'P'
1057: 'h'
1058: 'H'
1059: 'a'
1060: 'A'
1061: 'v'
1062: 'V'
1063: 'i'
1064: 'I'
1065: 'n'
1066: 'N'
1067: 'g'
1068: 'G'
1069: 'h'
1070: 'H'
1071: 'i'
1072: 'I'
1073: 'd'
1074: 'D'
1075: 'd'
1076: 'D'
1077: 'e'
1078: 'E'
1079: 'n'
1080: 'N'
1081: 'i'
1082: 'I'
1083: 'g'
1084: 'G'
1085: 'n'
1086: 'N'
1087: 'o'
1088: 'O'
1089: 'r'
1090: 'R'
1091: 'e'
1092: 'E'
1093: 'i'
1094: 'I'
1095: 'm'
1096: 'M'
1097: 'm'
1098: 'M'
1099: 'e'
1100: 'E'
1101: 'd'
1102: 'D'
1103: 'i'
1104: 'I'
1105: 'a'
1106: 'A'
1107: 't'
1108: 'T'
1109: 'e'
1110: 'E'
1111: 'i'
1112: 'I'
1113: 'm'
1114: 'M'
1115: 'm'
1116: 'M'
1117: 'u'
1118: 'U'
1119: 't'
1120: 'T'
1121: 'a'
1122: 'A'
1123: 'b'
1124: 'B'
1125: 'l'
1126: 'L'
1127: 'e'
1128: 'E'
1129: 'i'
1130: 'I'
1131: 'm'
1132: 'M'
1133: 'p'
1134: 'P'
1135: 'o'
1136: 'O'
1137: 'r'
1138: 'R'
1139: 't'
1140: 'T'
1141: 'i'
1142: 'I'
1143: 'n'
1144: 'N'
1145: 'i'
1146: 'I'
1147: 'n'
1148: 'N'
1149: 'c'
1150: 'C'
1151: 'l'
1152: 'L'
1153: 'u'
1154: 'U'
1155: 'd'
1156: 'D'
1157: 'e'
1158: 'E'
1159: 'i'
1160: 'I'
1161: 'n'
1162: 'N'
1163: 'd'
1164: 'D'
1165: 'e'
1166: 'E'
1167: 'x'
1168: 'X'
1169: 'i'
1170: 'I'
1171: 'n'
1172: 'N'
1173: 'n'
1174: 'N'
1175: 'e'
1176: 'E'
1177: 'r'
1178: 'R'
1179: 'i'
1180: 'I'
1181: 'n'
1182: 'N'
1183: 'o'
1184: 'O'
1185: 'u'
1186: 'U'
1187: 't'
1188: 'T'
1189: 'i'
1190: 'I'
1191: 'n'
1192: 'N'
1193: 's'
1194: 'S'
1195: 'e'
1196: 'E'
1197: 'r'
1198: 'R'
1199: 't'
1200: 'T'
1201: 'i'
1202: 'I'
1203: 'n'
1204: 'N'
1205: 't'
1206: 'T'
1207: 'e'
1208: 'E'
1209: 'r'
1210: 'R'
1211: 's'
1212: 'S'
1213: 'e'
1214: 'E'
1215: 'c'
1216: 'C'
1217: 't'
1218: 'T'
1219: 'i'
1220: 'I'
1221: 'n'
1222: 'N'
1223: 't'
1224: 'T'
1225: 'e'
1226: 'E'
1227: 'r'
1228: 'R'
1229: 'v'
1230: 'V'
1231: 'a'
1232: 'A'
1233: 'l'
1234: 'L'
1235: 'i'
1236: 'I'
1237: 'n'
1238: 'N'
1239: 'v'
1240: 'V'
1241: 'o'
1242: 'O'
1243: 'k'
1244: 'K'
1245: 'e'
1246: 'E'
1247: 'r'
1248: 'R'
1249: 'i'
1250: 'I'
1251: 's'
1252: 'S'
1253: 'i'
1254: 'I'
1255: 's'
1256: 'S'
1257: 'o'
1258: 'O'
1259: 'l'
1260: 'L'
1261: 'a'
1262: 'A'
1263: 't'
1264: 'T'
1265: 'i'
1266: 'I'
1267: 'o'
1268: 'O'
1269: 'n'
1270: 'N'
1271: 'i'
1272: 'I'
1273: 't'
1274: 'T'
1275: 'e'
1276: 'E'
1277: 'r'
1278: 'R'
1279: 'a'
1280: 'A'
1281: 't'
1282: 'T'
1283: 'e'
1284: 'E'
1285: 'j'
1286: 'J'
1287: 'o'
1288: 'O'
1289: 'i'
1290: 'I'
1291: 'n'
1292: 'N'
1293: 'j'
1294: 'J'
1295: 's'
1296: 'S'
1297: 'o'
1298: 'O'
1299: 'n'
1300: 'N'
1301: 'k'
1302: 'K'
1303: 'e'
1304: 'E'
1305: 'y'
1306: 'Y'
1307: 'l'
1308: 'L'
1309: 'a'
1310: 'A'
1311: 'n'
1312: 'N'
1313: 'g'
1314: 'G'
1315: 'u'
1316: 'U'
1317: 'a'
1318: 'A'
1319: 'g'
1320: 'G'
1321: 'e'
1322: 'E'
1323: 'l'
1324: 'L'
1325: 'a'
1326: 'A'
1327: 's'
1328: 'S'
1329: 't'
1330: 'T'
1331: 'l'
1332: 'L'
1333: 'e'
1334: 'E'
1335: 'a'
1336: 'A'
1337: 'v'
1338: 'V'
1339: 'e'
1340: 'E'
1341: 'l'
1342: 'L'
1343: 'e'
1344: 'E'
1345: 'f'
1346: 'F'
1347: 't'
1348: 'T'
1349: 'l'
1350: 'L'
1351: 'e'
1352: 'E'
1353: 'v'
1354: 'V'
1355: 'e'
1356: 'E'
1357: 'l'
1358: 'L'
1359: 'l'
1360: 'L'
1361: 'i'
1362: 'I'
1363: 'k'
1364: 'K'
1365: 'e'
1366: 'E'
1367: 'l'
1368: 'L'
1369: 'i'
1370: 'I'
1371: 'm'
1372: 'M'
1373: 'i'
1374: 'I'
1375: 't'
1376: 'T'
1377: 'l'
1378: 'L'
1379: 'o'
1380: 'O'
1381: 'a'
1382: 'A'
1383: 'd'
1384: 'D'
1385: 'l'
1386: 'L'
1387: 'o'
1388: 'O'
1389: 'o'
1390: 'O'
1391: 'p'
1392: 'P'
1393: 'm'
1394: 'M'
1395: 'a'
1396: 'A'
1397: 't'
1398: 'T'
1399: 'c'
1400: 'C'
1401: 'h'
1402: 'H'
1403: 'm'
1404: 'M'
1405: 'a'
1406: 'A'
1407: 't'
1408: 'T'
1409: 'c'
1410: 'C'
1411: 'h'
1412: 'H'
1413: 'e'
1414: 'E'
1415: 'd'
1416: 'D'
1417: 'm'
1418: 'M'
1419: 'a'
1420: 'A'
1421: 't'
1422: 'T'
1423: 'e'
1424: 'E'
1425: 'r'
1426: 'R'
1427: 'i'
1428: 'I'
1429: 'a'
1430: 'A'
1431: 'l'
1432: 'L'
1433: 'i'
1434: 'I'
1435: 'z'
1436: 'Z'
1437: 'e'
1438: 'E'
1439: 'd'
1440: 'D'
1441: 'm'
1442: 'M'
1443: 'a'
1444: 'A'
1445: 'x'
1446: 'X'
1447: 'm'
1448: 'M'
1449: 'e'
1450: 'E'
1451: 's'
1452: 'S'
1453: 's'
1454: 'S'
1455: 'a'
1456: 'A'
1457: 'g'
1458: 'G'
1459: 'e'
1460: 'E'
1461: 'm'
1462: 'M'
1463: 'i'
1464: 'I'
1465: 'n'
1466: 'N'
1467: 'm'
1468: 'M'
1469: 'o'
1470: 'O'
1471: 'd'
1472: 'D'
1473: 'e'
1474: 'E'
1475: 'l'
1476: 'L'
1477: 'm'
1478: 'M'
1479: 'o'
1480: 'O'
1481: 'd'
1482: 'D'
1483: 'u'
1484: 'U'
1485: 'l'
1486: 'L'
1487: 'e'
1488: 'E'
1489: 'n'
1490: 'N'
1491: 'o'
1492: 'O'
1493: 't'
1494: 'T'
1495: 'n'
1496: 'N'
1497: 'u'
1498: 'U'
1499: 'l'
1500: 'L'
1501: 'l'
1502: 'L'
1503: 'n'
1504: 'N'
1505: 'u'
1506: 'U'
1507: 'l'
1508: 'L'
1509: 'l'
1510: 'L'
1511: 's'
1512: 'S'
1513: 'n'
1514: 'N'
1515: 'u'
1516: 'U'
1517: 'm'
1518: 'M'
1519: 'e'
1520: 'E'
1521: 'r'
1522: 'R'
1523: 'i'
1524: 'I'
1525: 'c'
1526: 'C'
1527: 'o'
1528: 'O'
1529: 'f'
1530: 'F'
1531: 'f'
1532: 'F'
1533: 's'
1534: 'S'
1535: 'e'
1536: 'E'
1537: 't'
1538: 'T'
1539: 'o'
1540: 'O'
1541: 'n'
1542: 'N'
1543: 'o'
1544: 'O'
1545: 'n'
1546: 'N'
1547: 'l'
1548: 'L'
1549: 'y'
1550: 'Y'
1551: 'o'
1552: 'O'
1553: 'p'
1554: 'P'
1555: 't'
1556: 'T'
1557: 'i'
1558: 'I'
1559: 'o'
1560: 'O'
1561: 'n'
1562: 'N'
1563: 's'
1564: 'S'
1565: 'o'
1566: 'O'
1567: 'r'
1568: 'R'
1569: 'o'
1570: 'O'
1571: 'r'
1572: 'R'
1573: 'd'
1574: 'D'
1575: 'e'
1576: 'E'
1577: 'r'
1578: 'R'
1579: 'o'
1580: 'O'
1581: 'u'
1582: 'U'
1583: 't'
1584: 'T'
1585: 'o'
1586: 'O'
1587: 'u'
1588: 'U'
1589: 't'
1590: 'T'
1591: 'e'
1592: 'E'
1593: 'r'
1594: 'R'
1595: 'o'
1596: 'O'
1597: 'v'
1598: 'V'
1599: 'e'
1600: 'E'
1601: 'r'
1602: 'R'
1603: 'o'
1604: 'O'
1605: 'v'
1606: 'V'
1607: 'e'
1608: 'E'
1609: 'r'
1610: 'R'
1611: 'w'
1612: 'W'
1613: 'r'
1614: 'R'
1615: 'i'
1616: 'I'
1617: 't'
1618: 'T'
1619: 'e'
1620: 'E'
1621: 'p'
1622: 'P'
1623: 'a'
1624: 'A'
1625: 'r'
1626: 'R'
1627: 't'
1628: 'T'
1629: 'i'
1630: 'I'
1631: 't'
1632: 'T'
1633: 'i'
1634: 'I'
1635: 'o'
1636: 'O'
1637: 'n'
1638: 'N'
1639: 'p'
1640: 'P'
1641: 'e'
1642: 'E'
1643: 'r'
1644: 'R'
1645: 'c'
1646: 'C'
1647: 'e'
1648: 'E'
1649: 'n'
1650: 'N'
1651: 't'
1652: 'T'
1653: 'p'
1654: 'P'
1655: 'i'
1656: 'I'
1657: 'v'
1658: 'V'
1659: 'o'
1660: 'O'
1661: 't'
1662: 'T'
1663: 'p'
1664: 'P'
1665: 'o'
1666: 'O'
1667: 'l'
1668: 'L'
1669: 'i'
1670: 'I'
1671: 'c'
1672: 'C'
1673: 'i'
1674: 'I'
1675: 'e'
1676: 'E'
1677: 's'
1678: 'S'
1679: 'p'
1680: 'P'
1681: 'o'
1682: 'O'
1683: 'l'
1684: 'L'
1685: 'i'
1686: 'I'
1687: 'c'
1688: 'C'
1689: 'y'
1690: 'Y'
1691: 'p'
1692: 'P'
1693: 'r'
1694: 'R'
1695: 'e'
1696: 'E'
1697: 'c'
1698: 'C'
1699: 'e'
1700: 'E'
1701: 'd'
1702: 'D'
1703: 'i'
1704: 'I'
1705: 'n'
1706: 'N'
1707: 'g'
1708: 'G'
1709: 'p'
1710: 'P'
1711: 'r'
1712: 'R'
1713: 'i'
1714: 'I'
1715: 'm'
1716: 'M'
1717: 'a'
1718: 'A'
1719: 'r'
1720: 'R'
1721: 'y'
1722: 'Y'
1723: 'p'
1724: 'P'
1725: 'r'
1726: 'R'
1727: 'i'
1728: 'I'
1729: 'v'
1730: 'V'
1731: 'a'
1732: 'A'
1733: 't'
1734: 'T'
1735: 'e'
1736: 'E'
1737: 'p'
1738: 'P'
1739: 'r'
1740: 'R'
1741: 'i'
1742: 'I'
1743: 'v'
1744: 'V'
1745: 'i'
1746: 'I'
1747: 'l'
1748: 'L'
1749: 'e'
1750: 'E'
1751: 'g'
1752: 'G'
1753: 'e'
1754: 'E'
1755: 'p'
1756: 'P'
1757: 'r'
1758: 'R'
1759: 'i'
1760: 'I'
1761: 'v'
1762: 'V'
1763: 'i'
1764: 'I'
1765: 'l'
1766: 'L'
1767: 'e'
1768: 'E'
1769: 'g'
1770: 'G'
1771: 'e'
1772: 'E'
1773: 's'
1774: 'S'
1775: 'p'
1776: 'P'
1777: 'r'
1778: 'R'
1779: 'o'
1780: 'O'
1781: 'c'
1782: 'C'
1783: 'e'
1784: 'E'
1785: 'd'
1786: 'D'
1787: 'u'
1788: 'U'
1789: 'r'
1790: 'R'
1791: 'e'
1792: 'E'
1793: 'p'
1794: 'P'
1795: 'r'
1796: 'R'
1797: 'o'
1798: 'O'
1799: 't'
1800: 'T'
1801: 'o'
1802: 'O'
1803: 'p'
1804: 'P'
1805: 'u'
1806: 'U'
1807: 'b'
1808: 'B'
1809: 'l'
1810: 'L'
1811: 'i'
1812: 'I'
1813: 'c'
1814: 'C'
1815: 'q'
1816: 'Q'
1817: 'u'
1818: 'U'
1819: 'a'
1820: 'A'
1821: 'l'
1822: 'L'
1823: 'i'
1824: 'I'
1825: 'f'
1826: 'F'
1827: 'y'
1828: 'Y'
1829: 'r'
1830: 'R'
1831: 'a'
1832: 'A'
1833: 'i'
1834: 'I'
1835: 's'
1836: 'S'
1837: 'e'
1838: 'E'
1839: 'r'
1840: 'R'
1841: 'a'
1842: 'A'
1843: 'n'
1844: 'N'
1845: 'g'
1846: 'G'
1847: 'e'
1848: 'E'
1849: 'r'
1850: 'R'
1851: 'e'
1852: 'E'
1853: 'a'
1854: 'A'
1855: 'd'
1856: 'D'
1857: 'r'
1858: 'R'
1859: 'e'
1860: 'E'
1861: 'f'
1862: 'F'
1863: 'e'
1864: 'E'
1865: 'r'
1866: 'R'
1867: 'e'
1868: 'E'
1869: 'n'
1870: 'N'
1871: 'c'
1872: 'C'
1873: 'e'
1874: 'E'
1875: 's'
1876: 'S'
1877: 'r'
1878: 'R'
1879: 'e'
1880: 'E'
1881: 'm'
1882: 'M'
1883: 'o'
1884: 'O'
1885: 't'
1886: 'T'
1887: 'e'
1888: 'E'
1889: 'r'
1890: 'R'
1891: 'e'
1892: 'E'
1893: 'm'
1894: 'M'
1895: 'o'
1896: 'O'
1897: 'v'
1898: 'V'
1899: 'e'
1900: 'E'
1901: 'r'
1902: 'R'
1903: 'e'
1904: 'E'
1905: 'n'
1906: 'N'
1907: 'a'
1908: 'A'
1909: 'm'
1910: 'M'
1911: 'e'
1912: 'E'
1913: 'r'
1914: 'R'
1915: 'e'
1916: 'E'
1917: 'p'
1918: 'P'
1919: 'e'
1920: 'E'
1921: 'a'
1922: 'A'
1923: 't'
1924: 'T'
1925: 'r'
1926: 'R'
1927: 'e'
1928: 'E'
1929: 'p'
1930: 'P'
1931: 'e'
1932: 'E'
1933: 'a'
1934: 'A'
1935: 't'
1936: 'T'
1937: 'a'
1938: 'A'
1939: 'b'
1940: 'B'
1941: 'l'
1942: 'L'
1943: 'e'
1944: 'E'
1945: 'r'
1946: 'R'
1947: 'e'
1948: 'E'
1949: 'p'
1950: 'P'
1951: 'l'
1952: 'L'
1953: 'a'
1954: 'A'
1955: 'c'
1956: 'C'
1957: 'e'
1958: 'E'
1959: '_'
1960: '_'
1961: 'f'
1962: 'F'
1963: 'i'
1964: 'I'
1965: 'e'
1966: 'E'
1967: 'l'
1968: 'L'
1969: 'd'
1970: 'D'
1971: 's'
1972: 'S'
1973: 'r'
1974: 'R'
1975: 'e'
1976: 'E'
1977: 'p'
1978: 'P'
1979: 'l'
1980: 'L'
1981: 'a'
1982: 'A'
1983: 'c'
1984: 'C'
1985: 'e'
1986: 'E'
1987: 'r'
1988: 'R'
1989: 'e'
1990: 'E'
1991: 's'
1992: 'S'
1993: 'p'
1994: 'P'
1995: 'e'
1996: 'E'
1997: 'c'
1998: 'C'
1999: 't'
2000: 'T'
2001: 'r'
2002: 'R'
2003: 'e'
2004: 'E'
2005: 's'
2006: 'S'
2007: 't'
2008: 'T'
2009: 'r'
2010: 'R'
2011: 'i'
2012: 'I'
2013: 'c'
2014: 'C'
2015: 't'
2016: 'T'
2017: 'r'
2018: 'R'
2019: 'e'
2020: 'E'
2021: 's'
2022: 'S'
2023: 't'
2024: 'T'
2025: 'r'
2026: 'R'
2027: 'i'
2028: 'I'
2029: 'c'
2030: 'C'
2031: 't'
2032: 'T'
2033: 'i'
2034: 'I'
2035: 'o'
2036: 'O'
2037: 'n'
2038: 'N'
2039: 'r'
2040: 'R'
2041: 'e'
2042: 'E'
2043: 't'
2044: 'T'
2045: 'u'
2046: 'U'
2047: 'r'
2048: 'R'
2049: 'n'
2050: 'N'
2051: 'r'
2052: 'R'
2053: 'e'
2054: 'E'
2055: 't'
2056: 'T'
2057: 'u'
2058: 'U'
2059: 'r'
2060: 'R'
2061: 'n'
2062: 'N'
2063: 's'
2064: 'S'
2065: 'r'
2066: 'R'
2067: 'e'
2068: 'E'
2069: 'v'
2070: 'V'
2071: 'o'
2072: 'O'
2073: 'k'
2074: 'K'
2075: 'e'
2076: 'E'
2077: 'r'
2078: 'R'
2079: 'i'
2080: 'I'
2081: 'g'
2082: 'G'
2083: 'h'
2084: 'H'
2085: 't'
2086: 'T'
2087: 'r'
2088: 'R'
2089: 'o'
2090: 'O'
2091: 'l'
2092: 'L'
2093: 'l'
2094: 'L'
2095: 'b'
2096: 'B'
2097: 'a'
2098: 'A'
2099: 'c'
2100: 'C'
2101: 'k'
2102: 'K'
2103: 'r'
2104: 'R'
2105: 'o'
2106: 'O'
2107: 'l'
2108: 'L'
2109: 'l'
2110: 'L'
2111: 'u'
2112: 'U'
2113: 'p'
2114: 'P'
2115: 'r'
2116: 'R'
2117: 'o'
2118: 'O'
2119: 'w'
2120: 'W'
2121: 'r'
2122: 'R'
2123: 'o'
2124: 'O'
2125: 'w'
2126: 'W'
2127: 's'
2128: 'S'
2129: 'r'
2130: 'R'
2131: 'u'
2132: 'U'
2133: 'n'
2134: 'N'
2135: 's'
2136: 'S'
2137: 'a'
2138: 'A'
2139: 'f'
2140: 'F'
2141: 'e'
2142: 'E'
2143: '_'
2144: 'c'
2145: 'C'
2146: 'a'
2147: 'A'
2148: 's'
2149: 'S'
2150: 't'
2151: 'T'
2152: 's'
2153: 'S'
2154: 'c'
2155: 'C'
2156: 'h'
2157: 'H'
2158: 'e'
2159: 'E'
2160: 'm'
2161: 'M'
2162: 'a'
2163: 'A'
2164: 's'
2165: 'S'
2166: 'e'
2167: 'E'
2168: 'a'
2169: 'A'
2170: 'r'
2171: 'R'
2172: 'c'
2173: 'C'
2174: 'h'
2175: 'H'
2176: 's'
2177: 'S'
2178: 'e'
2179: 'E'
2180: 'c'
2181: 'C'
2182: 'u'
2183: 'U'
2184: 'r'
2185: 'R'
2186: 'i'
2187: 'I'
2188: 't'
2189: 'T'
2190: 'y'
2191: 'Y'
2192: 's'
2193: 'S'
2194: 'e'
2195: 'E'
2196: 'l'
2197: 'L'
2198: 'e'
2199: 'E'
2200: 'c'
2201: 'C'
2202: 't'
2203: 'T'
2204: 's'
2205: 'S'
2206: 'h'
2207: 'H'
2208: 'o'
2209: 'O'
2210: 'w'
2211: 'W'
2212: 's'
2213: 'S'
2214: 'i'
2215: 'I'
2216: 'm'
2217: 'M'
2218: 'p'
2219: 'P'
2220: 'l'
2221: 'L'
2222: 'e'
2223: 'E'
2224: 's'
2225: 'S'
2226: 'n'
2227: 'N'
2228: 'a'
2229: 'A'
2230: 'p'
2231: 'P'
2232: 's'
2233: 'S'
2234: 'h'
2235: 'H'
2236: 'o'
2237: 'O'
2238: 't'
2239: 'T'
2240: 's'
2241: 'S'
2242: 'o'
2243: 'O'
2244: 'u'
2245: 'U'
2246: 'r'
2247: 'R'
2248: 'c'
2249: 'C'
2250: 'e'
2251: 'E'
2252: 's'
2253: 'S'
2254: 'q'
2255: 'Q'
2256: 'l'
2257: 'L'
2258: 's'
2259: 'S'
2260: 't'
2261: 'T'
2262: 'a'
2263: 'A'
2264: 'b'
2265: 'B'
2266: 'l'
2267: 'L'
2268: 'e'
2269: 'E'
2270: 's'
2271: 'S'
2272: 't'
2273: 'T'
2274: 'a'
2275: 'A'
2276: 'r'
2277: 'R'
2278: 't'
2279: 'T'
2280: 's'
2281: 'S'
2282: 't'
2283: 'T'
2284: 'o'
2285: 'O'
2286: 'r'
2287: 'R'
2288: 'e'
2289: 'E'
2290: 'd'
2291: 'D'
2292: 's'
2293: 'S'
2294: 't'
2295: 'T'
2296: 'o'
2297: 'O'
2298: 'r'
2299: 'R'
2300: 'i'
2301: 'I'
2302: 'n'
2303: 'N'
2304: 'g'
2305: 'G'
2306: 's'
2307: 'S'
2308: 't'
2309: 'T'
2310: 'r'
2311: 'R'
2312: 'u'
2313: 'U'
2314: 'c'
2315: 'C'
2316: 't'
2317: 'T'
2318: 's'
2319: 'S'
2320: 'y'
2321: 'Y'
2322: 's'
2323: 'S'
2324: 't'
2325: 'T'
2326: 'e'
2327: 'E'
2328: 'm'
2329: 'M'
2330: '_'
2331: '_'
2332: 't'
2333: 'T'
2334: 'i'
2335: 'I'
2336: 'm'
2337: 'M'
2338: 'e'
2339: 'E'
2340: 's'
2341: 'S'
2342: 'y'
2343: 'Y'
2344: 's'
2345: 'S'
2346: 't'
2347: 'T'
2348: 'e'
2349: 'E'
2350: 'm'
2351: 'M'
2352: 't'
2353: 'T'
2354: 'a'
2355: 'A'
2356: 'b'
2357: 'B'
2358: 'l'
2359: 'L'
2360: 'e'
2361: 'E'
2362: 't'
2363: 'T'
2364: 'a'
2365: 'A'
2366: 'b'
2367: 'B'
2368: 'l'
2369: 'L'
2370: 'e'
2371: 'E'
2372: 's'
2373: 'S'
2374: 'a'
2375: 'A'
2376: 'm'
2377: 'M'
2378: 'p'
2379: 'P'
2380: 'l'
2381: 'L'
2382: 'e'
2383: 'E'
2384: 't'
2385: 'T'
2386: 'a'
2387: 'A'
2388: 'r'
2389: 'R'
2390: 'g'
2391: 'G'
2392: 'e'
2393: 'E'
2394: 't'
2395: 'T'
2396: 't'
2397: 'T'
2398: 'e'
2399: 'E'
2400: 'm'
2401: 'M'
2402: 'p'
2403: 'P'
2404: 't'
2405: 'T'
2406: 'e'
2407: 'E'
2408: 'm'
2409: 'M'
2410: 'p'
2411: 'P'
2412: 'o'
2413: 'O'
2414: 'r'
2415: 'R'
2416: 'a'
2417: 'A'
2418: 'r'
2419: 'R'
2420: 'y'
2421: 'Y'
2422: 't'
2423: 'T'
2424: 'h'
2425: 'H'
2426: 'e'
2427: 'E'
2428: 'n'
2429: 'N'
2430: 't'
2431: 'T'
2432: 'i'
2433: 'I'
2434: 'm'
2435: 'M'
2436: 'e'
2437: 'E'
2438: 't'
2439: 'T'
2440: 'i'
2441: 'I'
2442: 'm'
2443: 'M'
2444: 'e'
2445: 'E'
2446: 's'
2447: 'S'
2448: 't'
2449: 'T'
2450: 'a'
2451: 'A'
2452: 'm'
2453: 'M'
2454: 'p'
2455: 'P'
2456: 't'
2457: 'T'
2458: 'o'
2459: 'O'
2460: 't'
2461: 'T'
2462: 'r'
2463: 'R'
2464: 'a'
2465: 'A'
2466: 'n'
2467: 'N'
2468: 's'
2469: 'S'
2470: 'a'
2471: 'A'
2472: 'c'
2473: 'C'
2474: 't'
2475: 'T'
2476: 'i'
2477: 'I'
2478: 'o'
2479: 'O'
2480: 'n'
2481: 'N'
2482: 't'
2483: 'T'
2484: 'r'
2485: 'R'
2486: 'a'
2487: 'A'
2488: 'n'
2489: 'N'
2490: 's'
2491: 'S'
2492: 'f'
2493: 'F'
2494: 'o'
2495: 'O'
2496: 'r'
2497: 'R'
2498: 'm'
2499: 'M'
2500: 't'
2501: 'T'
2502: 'r'
2503: 'R'
2504: 'u'
2505: 'U'
2506: 'n'
2507: 'N'
2508: 'c'
2509: 'C'
2510: 'a'
2511: 'A'
2512: 't'
2513: 'T'
2514: 'e'
2515: 'E'
2516: 't'
2517: 'T'
2518: 'y'
2519: 'Y'
2520: 'p'
2521: 'P'
2522: 'e'
2523: 'E'
2524: 'u'
2525: 'U'
2526: 'n'
2527: 'N'
2528: 'b'
2529: 'B'
2530: 'o'
2531: 'O'
2532: 'u'
2533: 'U'
2534: 'n'
2535: 'N'
2536: 'd'
2537: 'D'
2538: 'e'
2539: 'E'
2540: 'd'
2541: 'D'
2542: 'u'
2543: 'U'
2544: 'n'
2545: 'N'
2546: 'i'
2547: 'I'
2548: 'o'
2549: 'O'
2550: 'n'
2551: 'N'
2552: 'u'
2553: 'U'
2554: 'n'
2555: 'N'
2556: 'i'
2557: 'I'
2558: 'q'
2559: 'Q'
2560: 'u'
2561: 'U'
2562: 'e'
2563: 'E'
2564: 'u'
2565: 'U'
2566: 'n'
2567: 'N'
2568: 'n'
2569: 'N'
2570: 'e'
2571: 'E'
2572: 's'
2573: 'S'
2574: 't'
2575: 'T'
2576: 'u'
2577: 'U'
2578: 'n'
2579: 'N'
2580: 'p'
2581: 'P'
2582: 'i'
2583: 'I'
2584: 'v'
2585: 'V'
2586: 'o'
2587: 'O'
2588: 't'
2589: 'T'
2590: 'u'
2591: 'U'
2592: 'n'
2593: 'N'
2594: 't'
2595: 'T'
2596: 'i'
2597: 'I'
2598: 'l'
2599: 'L'
2600: 'u'
2601: 'U'
2602: 'p'
2603: 'P'
2604: 'd'
2605: 'D'
2606: 'a'
2607: 'A'
2608: 't'
2609: 'T'
2610: 'e'
2611: 'E'
2612: 'u'
2613: 'U'
2614: 's'
2615: 'S'
2616: 'i'
2617: 'I'
2618: 'n'
2619: 'N'
2620: 'g'
2621: 'G'
2622: 'v'
2623: 'V'
2624: 'a'
2625: 'A'
2626: 'l'
2627: 'L'
2628: 'u'
2629: 'U'
2630: 'e'
2631: 'E'
2632: 'v'
2633: 'V'
2634: 'a'
2635: 'A'
2636: 'l'
2637: 'L'
2638: 'u'
2639: 'U'
2640: 'e'
2641: 'E'
2642: 's'
2643: 'S'
2644: 'v'
2645: 'V'
2646: 'i'
2647: 'I'
2648: 'e'
2649: 'E'
2650: 'w'
2651: 'W'
2652: 'v'
2653: 'V'
2654: 'i'
2655: 'I'
2656: 'e'
2657: 'E'
2658: 'w'
2659: 'W'
2660: 's'
2661: 'S'
2662: 'v'
2663: 'V'
2664: 'o'
2665: 'O'
2666: 'l'
2667: 'L'
2668: 'a'
2669: 'A'
2670: 't'
2671: 'T'
2672: 'i'
2673: 'I'
2674: 'l'
2675: 'L'
2676: 'e'
2677: 'E'
2678: 'w'
2679: 'W'
2680: 'e'
2681: 'E'
2682: 'i'
2683: 'I'
2684: 'g'
2685: 'G'
2686: 'h'
2687: 'H'
2688: 't'
2689: 'T'
2690: 'w'
2691: 'W'
2692: 'h'
2693: 'H'
2694: 'e'
2695: 'E'
2696: 'n'
2697: 'N'
2698: 'w'
2699: 'W'
2700: 'h'
2701: 'H'
2702: 'e'
2703: 'E'
2704: 'r'
2705: 'R'
2706: 'e'
2707: 'E'
2708: 'w'
2709: 'W'
2710: 'h'
2711: 'H'
2712: 'i'
2713: 'I'
2714: 'l'
2715: 'L'
2716: 'e'
2717: 'E'
2718: 'w'
2719: 'W'
2720: 'i'
2721: 'I'
2722: 'n'
2723: 'N'
2724: 'd'
2725: 'D'
2726: 'o'
2727: 'O'
2728: 'w'
2729: 'W'
2730: 'w'
2731: 'W'
2732: 'i'
2733: 'I'
2734: 't'
2735: 'T'
2736: 'h'
2737: 'H'
2738: 'w'
2739: 'W'
2740: 'r'
2741: 'R'
2742: 'i'
2743: 'I'
2744: 't'
2745: 'T'
2746: 'e'
2747: 'E'
2748: 'z'
2749: 'Z'
2750: 'o'
2751: 'O'
2752: 'n'
2753: 'N'
2754: 'e'
2755: 'E'
2756: '-'
2757: '-'
2758: '\n'
2759: '#'
2760: '\n'
2761: ';'
2762: ','
2763: '('
2764: ')'
2765: '.'
2766: '*'
2767: '*'
2768: '<'
2769: '>'
2770: '|'
2771: '^'
2772: '&'
2773: '['
2774: ']'
2775: '.'
2776: '+'
2777: '-'
2778: '~'
2779: '/'
2780: '|'
2781: '|'
2782: '<'
2783: '<'
2784: '>'
2785: '>'
2786: '@'
2787: '='
2788: '!'
2789: '='
2790: '<'
2791: '>'
2792: '<'
2793: '='
2794: '>'
2795: '='
2796: '\'
2797: '"'
2798: '''
2799: '\'
2800: '\n'
2801: '\r'
2802: 'r'
2803: 'R'
2804: 'r'
2805: 'R'
2806: 'b'
2807: 'B'
2808: 'b'
2809: 'B'
2810: 'r'
2811: 'R'
2812: 'b'
2813: 'B'
2814: '"'
2815: '''
2816: '\n'
2817: '\n'
2818: '''
2819: '"'
2820: '.'
2821: '.'
2822: 'e'
2823: 'E'
2824: '+'
2825: '-'
2826: '0'
2827: 'x'
2828: 'X'
2829: '*'
2830: '*'
2831: '/'
2832: '/'
2833: '*'
2834: '*'
2835: '_'
2836: '`'
2837: ' '
2838: '\t'
2839: '\n'
2840: '\r'
2841: \u00a0
2842: 'a'-'z'
2843: 'A'-'Z'
2844: 'g'-'z'
2845: 'G'-'Z'
2846: 'a'-'z'
2847: 'A'-'Z'
2848: \u0001-'\t'
2849: '\v'-'!'
2850: '#'-'&'
2851: '('-'['
2852: ']'-\u007f
2853: \u0080-\ufffc
2854: \ufffe-\U0010ffff
2855: '0'-'9'
2856: 'a'-'f'
2857: 'A'-'F'
2858: \u0001-'\t'
2859: '\v'-\u007f
2860: \u0080-\ufffc
2861: \ufffe-\U0010ffff
2862: '0'-'9'
2863: 'a'-'z'
2864: 'A'-'Z'
2865: \u0001-'\t'
2866: '\v'-'['
2867: ']'-'_'
2868: 'a'-\u007f
2869: .
*/
