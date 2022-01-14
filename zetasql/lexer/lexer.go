// Code generated by gocc; DO NOT EDIT.

package lexer

import (
	"io/ioutil"
	"unicode/utf8"

	"github.com/paulourio/bqfmt/zetasql/token"
)

const (
	NoState    = -1
	NumStates  = 562
	NumSymbols = 999
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
2: 'f'
3: 'F'
4: 'a'
5: 'A'
6: 'l'
7: 'L'
8: 's'
9: 'S'
10: 'e'
11: 'E'
12: 't'
13: 'T'
14: 'r'
15: 'R'
16: 'u'
17: 'U'
18: 'e'
19: 'E'
20: 'n'
21: 'N'
22: 'u'
23: 'U'
24: 'm'
25: 'M'
26: 'e'
27: 'E'
28: 'r'
29: 'R'
30: 'i'
31: 'I'
32: 'c'
33: 'C'
34: 'd'
35: 'D'
36: 'e'
37: 'E'
38: 'c'
39: 'C'
40: 'i'
41: 'I'
42: 'm'
43: 'M'
44: 'a'
45: 'A'
46: 'l'
47: 'L'
48: 'b'
49: 'B'
50: 'i'
51: 'I'
52: 'g'
53: 'G'
54: 'n'
55: 'N'
56: 'u'
57: 'U'
58: 'm'
59: 'M'
60: 'e'
61: 'E'
62: 'r'
63: 'R'
64: 'i'
65: 'I'
66: 'c'
67: 'C'
68: 'b'
69: 'B'
70: 'i'
71: 'I'
72: 'g'
73: 'G'
74: 'd'
75: 'D'
76: 'e'
77: 'E'
78: 'c'
79: 'C'
80: 'i'
81: 'I'
82: 'm'
83: 'M'
84: 'a'
85: 'A'
86: 'l'
87: 'L'
88: 'a'
89: 'A'
90: 'l'
91: 'L'
92: 'l'
93: 'L'
94: 'a'
95: 'A'
96: 'n'
97: 'N'
98: 'd'
99: 'D'
100: 'a'
101: 'A'
102: 'r'
103: 'R'
104: 'r'
105: 'R'
106: 'a'
107: 'A'
108: 'y'
109: 'Y'
110: 'a'
111: 'A'
112: 's'
113: 'S'
114: 'a'
115: 'A'
116: 's'
117: 'S'
118: 'c'
119: 'C'
120: 'a'
121: 'A'
122: 't'
123: 'T'
124: 'b'
125: 'B'
126: 'e'
127: 'E'
128: 't'
129: 'T'
130: 'w'
131: 'W'
132: 'e'
133: 'E'
134: 'e'
135: 'E'
136: 'n'
137: 'N'
138: 'b'
139: 'B'
140: 'y'
141: 'Y'
142: 'c'
143: 'C'
144: 'a'
145: 'A'
146: 's'
147: 'S'
148: 'e'
149: 'E'
150: 'c'
151: 'C'
152: 'a'
153: 'A'
154: 's'
155: 'S'
156: 't'
157: 'T'
158: 'c'
159: 'C'
160: 'r'
161: 'R'
162: 'o'
163: 'O'
164: 's'
165: 'S'
166: 's'
167: 'S'
168: 'c'
169: 'C'
170: 'u'
171: 'U'
172: 'r'
173: 'R'
174: 'r'
175: 'R'
176: 'e'
177: 'E'
178: 'n'
179: 'N'
180: 't'
181: 'T'
182: 'd'
183: 'D'
184: 'a'
185: 'A'
186: 't'
187: 'T'
188: 'e'
189: 'E'
190: 'd'
191: 'D'
192: 'a'
193: 'A'
194: 't'
195: 'T'
196: 'e'
197: 'E'
198: 't'
199: 'T'
200: 'i'
201: 'I'
202: 'm'
203: 'M'
204: 'e'
205: 'E'
206: 'd'
207: 'D'
208: 'e'
209: 'E'
210: 's'
211: 'S'
212: 'c'
213: 'C'
214: 'd'
215: 'D'
216: 'i'
217: 'I'
218: 's'
219: 'S'
220: 't'
221: 'T'
222: 'i'
223: 'I'
224: 'n'
225: 'N'
226: 'c'
227: 'C'
228: 't'
229: 'T'
230: 'e'
231: 'E'
232: 'l'
233: 'L'
234: 's'
235: 'S'
236: 'e'
237: 'E'
238: 'e'
239: 'E'
240: 'n'
241: 'N'
242: 'd'
243: 'D'
244: 'e'
245: 'E'
246: 'x'
247: 'X'
248: 'c'
249: 'C'
250: 'e'
251: 'E'
252: 'p'
253: 'P'
254: 't'
255: 'T'
256: 'e'
257: 'E'
258: 'x'
259: 'X'
260: 't'
261: 'T'
262: 'r'
263: 'R'
264: 'a'
265: 'A'
266: 'c'
267: 'C'
268: 't'
269: 'T'
270: 'f'
271: 'F'
272: 'i'
273: 'I'
274: 'r'
275: 'R'
276: 's'
277: 'S'
278: 't'
279: 'T'
280: 'f'
281: 'F'
282: 'o'
283: 'O'
284: 'l'
285: 'L'
286: 'l'
287: 'L'
288: 'o'
289: 'O'
290: 'w'
291: 'W'
292: 'i'
293: 'I'
294: 'n'
295: 'N'
296: 'g'
297: 'G'
298: 'f'
299: 'F'
300: 'o'
301: 'O'
302: 'r'
303: 'R'
304: 'm'
305: 'M'
306: 'a'
307: 'A'
308: 't'
309: 'T'
310: 'f'
311: 'F'
312: 'r'
313: 'R'
314: 'o'
315: 'O'
316: 'm'
317: 'M'
318: 'f'
319: 'F'
320: 'u'
321: 'U'
322: 'l'
323: 'L'
324: 'l'
325: 'L'
326: 'g'
327: 'G'
328: 'r'
329: 'R'
330: 'o'
331: 'O'
332: 'u'
333: 'U'
334: 'p'
335: 'P'
336: 'h'
337: 'H'
338: 'a'
339: 'A'
340: 'v'
341: 'V'
342: 'i'
343: 'I'
344: 'n'
345: 'N'
346: 'g'
347: 'G'
348: 'i'
349: 'I'
350: 'g'
351: 'G'
352: 'n'
353: 'N'
354: 'o'
355: 'O'
356: 'r'
357: 'R'
358: 'e'
359: 'E'
360: 'i'
361: 'I'
362: 'n'
363: 'N'
364: 'i'
365: 'I'
366: 'n'
367: 'N'
368: 'n'
369: 'N'
370: 'e'
371: 'E'
372: 'r'
373: 'R'
374: 'i'
375: 'I'
376: 'n'
377: 'N'
378: 't'
379: 'T'
380: 'e'
381: 'E'
382: 'r'
383: 'R'
384: 's'
385: 'S'
386: 'e'
387: 'E'
388: 'c'
389: 'C'
390: 't'
391: 'T'
392: 'i'
393: 'I'
394: 'n'
395: 'N'
396: 't'
397: 'T'
398: 'e'
399: 'E'
400: 'r'
401: 'R'
402: 'v'
403: 'V'
404: 'a'
405: 'A'
406: 'l'
407: 'L'
408: 'i'
409: 'I'
410: 's'
411: 'S'
412: 'j'
413: 'J'
414: 'o'
415: 'O'
416: 'i'
417: 'I'
418: 'n'
419: 'N'
420: 'j'
421: 'J'
422: 's'
423: 'S'
424: 'o'
425: 'O'
426: 'n'
427: 'N'
428: 'l'
429: 'L'
430: 'a'
431: 'A'
432: 's'
433: 'S'
434: 't'
435: 'T'
436: 'l'
437: 'L'
438: 'e'
439: 'E'
440: 'f'
441: 'F'
442: 't'
443: 'T'
444: 'l'
445: 'L'
446: 'i'
447: 'I'
448: 'k'
449: 'K'
450: 'e'
451: 'E'
452: 'l'
453: 'L'
454: 'i'
455: 'I'
456: 'm'
457: 'M'
458: 'i'
459: 'I'
460: 't'
461: 'T'
462: 'n'
463: 'N'
464: 'o'
465: 'O'
466: 't'
467: 'T'
468: 'n'
469: 'N'
470: 'u'
471: 'U'
472: 'l'
473: 'L'
474: 'l'
475: 'L'
476: 'n'
477: 'N'
478: 'u'
479: 'U'
480: 'l'
481: 'L'
482: 'l'
483: 'L'
484: 's'
485: 'S'
486: 'o'
487: 'O'
488: 'f'
489: 'F'
490: 'f'
491: 'F'
492: 's'
493: 'S'
494: 'e'
495: 'E'
496: 't'
497: 'T'
498: 'o'
499: 'O'
500: 'n'
501: 'N'
502: 'o'
503: 'O'
504: 'r'
505: 'R'
506: 'o'
507: 'O'
508: 'r'
509: 'R'
510: 'd'
511: 'D'
512: 'e'
513: 'E'
514: 'r'
515: 'R'
516: 'o'
517: 'O'
518: 'u'
519: 'U'
520: 't'
521: 'T'
522: 'e'
523: 'E'
524: 'r'
525: 'R'
526: 'o'
527: 'O'
528: 'v'
529: 'V'
530: 'e'
531: 'E'
532: 'r'
533: 'R'
534: 'p'
535: 'P'
536: 'a'
537: 'A'
538: 'r'
539: 'R'
540: 't'
541: 'T'
542: 'i'
543: 'I'
544: 't'
545: 'T'
546: 'i'
547: 'I'
548: 'o'
549: 'O'
550: 'n'
551: 'N'
552: 'p'
553: 'P'
554: 'e'
555: 'E'
556: 'r'
557: 'R'
558: 'c'
559: 'C'
560: 'e'
561: 'E'
562: 'n'
563: 'N'
564: 't'
565: 'T'
566: 'p'
567: 'P'
568: 'r'
569: 'R'
570: 'e'
571: 'E'
572: 'c'
573: 'C'
574: 'e'
575: 'E'
576: 'd'
577: 'D'
578: 'i'
579: 'I'
580: 'n'
581: 'N'
582: 'g'
583: 'G'
584: 'q'
585: 'Q'
586: 'u'
587: 'U'
588: 'a'
589: 'A'
590: 'l'
591: 'L'
592: 'i'
593: 'I'
594: 'f'
595: 'F'
596: 'y'
597: 'Y'
598: 'r'
599: 'R'
600: 'a'
601: 'A'
602: 'n'
603: 'N'
604: 'g'
605: 'G'
606: 'e'
607: 'E'
608: 'r'
609: 'R'
610: 'e'
611: 'E'
612: 'p'
613: 'P'
614: 'e'
615: 'E'
616: 'a'
617: 'A'
618: 't'
619: 'T'
620: 'a'
621: 'A'
622: 'b'
623: 'B'
624: 'l'
625: 'L'
626: 'e'
627: 'E'
628: 'r'
629: 'R'
630: 'e'
631: 'E'
632: 'p'
633: 'P'
634: 'l'
635: 'L'
636: 'a'
637: 'A'
638: 'c'
639: 'C'
640: 'e'
641: 'E'
642: 'r'
643: 'R'
644: 'e'
645: 'E'
646: 's'
647: 'S'
648: 'p'
649: 'P'
650: 'e'
651: 'E'
652: 'c'
653: 'C'
654: 't'
655: 'T'
656: 'r'
657: 'R'
658: 'i'
659: 'I'
660: 'g'
661: 'G'
662: 'h'
663: 'H'
664: 't'
665: 'T'
666: 'r'
667: 'R'
668: 'o'
669: 'O'
670: 'w'
671: 'W'
672: 'r'
673: 'R'
674: 'o'
675: 'O'
676: 'w'
677: 'W'
678: 's'
679: 'S'
680: 's'
681: 'S'
682: 'a'
683: 'A'
684: 'f'
685: 'F'
686: 'e'
687: 'E'
688: '_'
689: 'c'
690: 'C'
691: 'a'
692: 'A'
693: 's'
694: 'S'
695: 't'
696: 'T'
697: 's'
698: 'S'
699: 'e'
700: 'E'
701: 'l'
702: 'L'
703: 'e'
704: 'E'
705: 'c'
706: 'C'
707: 't'
708: 'T'
709: 's'
710: 'S'
711: 't'
712: 'T'
713: 'r'
714: 'R'
715: 'u'
716: 'U'
717: 'c'
718: 'C'
719: 't'
720: 'T'
721: 't'
722: 'T'
723: 'a'
724: 'A'
725: 'b'
726: 'B'
727: 'l'
728: 'L'
729: 'e'
730: 'E'
731: 's'
732: 'S'
733: 'a'
734: 'A'
735: 'm'
736: 'M'
737: 'p'
738: 'P'
739: 'l'
740: 'L'
741: 'e'
742: 'E'
743: 't'
744: 'T'
745: 'h'
746: 'H'
747: 'e'
748: 'E'
749: 'n'
750: 'N'
751: 't'
752: 'T'
753: 'i'
754: 'I'
755: 'm'
756: 'M'
757: 'e'
758: 'E'
759: 't'
760: 'T'
761: 'i'
762: 'I'
763: 'm'
764: 'M'
765: 'e'
766: 'E'
767: 's'
768: 'S'
769: 't'
770: 'T'
771: 'a'
772: 'A'
773: 'm'
774: 'M'
775: 'p'
776: 'P'
777: 't'
778: 'T'
779: 'o'
780: 'O'
781: 'u'
782: 'U'
783: 'n'
784: 'N'
785: 'b'
786: 'B'
787: 'o'
788: 'O'
789: 'u'
790: 'U'
791: 'n'
792: 'N'
793: 'd'
794: 'D'
795: 'e'
796: 'E'
797: 'd'
798: 'D'
799: 'u'
800: 'U'
801: 'n'
802: 'N'
803: 'i'
804: 'I'
805: 'o'
806: 'O'
807: 'n'
808: 'N'
809: 'u'
810: 'U'
811: 'n'
812: 'N'
813: 'n'
814: 'N'
815: 'e'
816: 'E'
817: 's'
818: 'S'
819: 't'
820: 'T'
821: 'u'
822: 'U'
823: 's'
824: 'S'
825: 'i'
826: 'I'
827: 'n'
828: 'N'
829: 'g'
830: 'G'
831: 'w'
832: 'W'
833: 'e'
834: 'E'
835: 'i'
836: 'I'
837: 'g'
838: 'G'
839: 'h'
840: 'H'
841: 't'
842: 'T'
843: 'w'
844: 'W'
845: 'h'
846: 'H'
847: 'e'
848: 'E'
849: 'n'
850: 'N'
851: 'w'
852: 'W'
853: 'h'
854: 'H'
855: 'e'
856: 'E'
857: 'r'
858: 'R'
859: 'e'
860: 'E'
861: 'w'
862: 'W'
863: 'i'
864: 'I'
865: 'n'
866: 'N'
867: 'd'
868: 'D'
869: 'o'
870: 'O'
871: 'w'
872: 'W'
873: 'w'
874: 'W'
875: 'i'
876: 'I'
877: 't'
878: 'T'
879: 'h'
880: 'H'
881: 'z'
882: 'Z'
883: 'o'
884: 'O'
885: 'n'
886: 'N'
887: 'e'
888: 'E'
889: '-'
890: '-'
891: '\n'
892: '#'
893: '\n'
894: ';'
895: ','
896: '*'
897: '('
898: ')'
899: '<'
900: '>'
901: '|'
902: '^'
903: '&'
904: '['
905: ']'
906: '.'
907: '+'
908: '-'
909: '~'
910: '/'
911: '|'
912: '|'
913: '<'
914: '<'
915: '>'
916: '>'
917: '@'
918: '='
919: '!'
920: '='
921: '<'
922: '>'
923: '<'
924: '='
925: '>'
926: '='
927: '\'
928: '"'
929: '''
930: '\'
931: '\n'
932: '\r'
933: 'r'
934: 'R'
935: 'r'
936: 'R'
937: 'b'
938: 'B'
939: 'b'
940: 'B'
941: 'r'
942: 'R'
943: 'b'
944: 'B'
945: '"'
946: '''
947: '\n'
948: '\n'
949: '''
950: '"'
951: '.'
952: '.'
953: 'e'
954: 'E'
955: '+'
956: '-'
957: '0'
958: 'x'
959: 'X'
960: '*'
961: '*'
962: '/'
963: '/'
964: '*'
965: '*'
966: '_'
967: '`'
968: ' '
969: '\t'
970: '\n'
971: '\r'
972: \u00a0
973: 'a'-'z'
974: 'A'-'Z'
975: 'a'-'z'
976: 'A'-'Z'
977: \u0001-'\t'
978: '\v'-'!'
979: '#'-'&'
980: '('-'['
981: ']'-\u007f
982: \u0080-\ufffc
983: \ufffe-\U0010ffff
984: '0'-'9'
985: 'a'-'f'
986: 'A'-'F'
987: \u0001-'\t'
988: '\v'-\u007f
989: \u0080-\ufffc
990: \ufffe-\U0010ffff
991: '0'-'9'
992: 'a'-'z'
993: 'A'-'Z'
994: \u0001-'\t'
995: '\v'-'['
996: ']'-'_'
997: 'a'-\u007f
998: .
*/
