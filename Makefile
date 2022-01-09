GOCCFLAGS = -a

build: zetasql/ast/types_generated.go zetasql/parser/productionstable.go
	go build

test: build
	go test -v ./zetasql/

debug_test: GOCCFLAGS += -debug_parser
debug_test: test

zetasql/ast/types_generated.go: zetasql/ast/types_generated.go.j2 zetasql/ast/gen_types.py
	cd zetasql/ast && python gen_types.py

zetasql/parser/productionstable.go: zetasql/zetasql.bnf
	cd zetasql && gocc $(GOCCFLAGS) zetasql.bnf
