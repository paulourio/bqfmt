GOCCFLAGS = -a

build: generated_types generated_parser
	go build

test: build
	go test -v ./zetasql/

debug_test: GOCCFLAGS += -debug_parser
debug_test: test

generated_types: zetasql/ast/types_generated.go

generated_parser: zetasql/parser/productionstable.go

zetasql/ast/types_generated.go: zetasql/ast/types_generated.go.j2 zetasql/ast/gen_types.py
	cd zetasql/ast && python gen_types.py

zetasql/parser/productionstable.go: zetasql/zetasql.bnf
	cd zetasql && gocc $(GOCCFLAGS) zetasql.bnf

.PHONY: debug_conflicts
debug_conflicts:
	cd zetasql && gocc -v zetasql.bnf