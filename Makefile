GOCCFLAGS = -a # -debug_parser

.PHONY: build
build: generated_types generated_parser
	go build cmd/bqfmt/bqfmt.go

.PHONY: test
test: build
	gotestsum --format dots ./zetasql/...

.PHONY: test_lexer
test_lexer: build
	go test -v ./zetasql/lexer_test.go

.PHONY: debug_conflicts
debug_conflicts:
	cd zetasql && gocc -v zetasql.bnf

debug_test: GOCCFLAGS += -debug_parser
debug_test: test

generated_types: zetasql/ast/types_generated.go

generated_parser: zetasql/parser/productionstable.go

zetasql/ast/types_generated.go: zetasql/ast/types_generated.go.j2 zetasql/ast/gen_types.py
	cd zetasql/ast && python gen_types.py

zetasql/parser/productionstable.go: zetasql/zetasql.bnf
	cd zetasql && gocc $(GOCCFLAGS) zetasql.bnf
