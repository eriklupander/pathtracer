all: fmt build test

build:
	go build -o bin/pt cmd/pt/main.go

run: build
	./bin/pt

test:
	go test ./... -count=1 -v

fmt:
	go fmt ./...

asm: asm-cross asm-mulmat

asm-cross:
	clang -S -mavx2 -mfma -masm=intel -mno-red-zone -mstackrealign -mllvm -inline-threshold=1000 -fno-asynchronous-unwind-tables -fno-exceptions -fno-rtti -c -O3 cfiles/CrossProduct.c
	mv CrossProduct.s cfiles/
	c2goasm -a -f cfiles/CrossProduct.s internal/app/geom/CrossProduct_amd64.s
	rm cfiles/CrossProduct.s
asm-dot:
	clang -S -mavx2 -mfma -masm=intel -mno-red-zone -mstackrealign -mllvm -inline-threshold=1000 -fno-asynchronous-unwind-tables -fno-exceptions -fno-rtti -c -O3 cfiles/DotProduct.c
	mv DotProduct.s cfiles/
	c2goasm -a -f cfiles/DotProduct.s internal/app/geom/DotProduct_amd64.s
	rm cfiles/DotProduct.s
asm-dp:
	clang -S -mavx2 -mfma -masm=intel -mno-red-zone -mstackrealign -mllvm -inline-threshold=1000 -fno-asynchronous-unwind-tables -fno-exceptions -fno-rtti -c -O3 cfiles/DP.c
	mv DP.s cfiles/
	c2goasm -a -f cfiles/DP.s internal/app/geom/DP_amd64.s
	rm cfiles/DP.s
asm-mulmat:
	clang -S -mavx2 -mfma -masm=intel -mno-red-zone -mstackrealign -mllvm -inline-threshold=1000 -fno-asynchronous-unwind-tables -fno-exceptions -fno-rtti -c -O3 cfiles/MultiplyMatrixByVec64.c
	mv MultiplyMatrixByVec64.s cfiles/
	c2goasm -a -f cfiles/MultiplyMatrixByVec64.s internal/app/geom/MultiplyMatrixByVec64_amd64.s
	rm cfiles/MultiplyMatrixByVec64.s