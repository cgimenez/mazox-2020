#
# find . -name '*.go' | ./entr -r ./b.sh
#
printf "Building - "
date
OUTPUT=$(GOOS=js GOARCH=wasm go build -tags runtime_js -o www/mazox.wasm 2>&1)
#OUTPUT=$(go build -tags runtime_sdl -o mazox 2>&1)
#OUTPUT=$(GOOS=linux GOARCH=adm64 go build -tags runtime_sdl -o mazox.linux 2>&1)
if [ ${#OUTPUT} -ge 1 ];
then
  printf "Failed - $OUTPUT\n"
  echo 'display notification "'"$OUTPUT"'"' | osascript
else
  go build -tags runtime_sdl -o mazox
fi
printf "Succeed\n"
