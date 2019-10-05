mpdf : linux darwin

linux darwin :
	@echo "building for platform $@"
	mkdir -p build/$@ && env GOOS=$@ go build -o build/$@/mpdf . && zip -j build/$@/mpdf-$@.zip build/$@/mpdf
	mv build/$@/mpdf-$@.zip build && rm -rf build/$@
	@echo "done building for platform $@"

clean :
	@rm -rf build 
