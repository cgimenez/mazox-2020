#
# find . -name '*.go' | ./entr -r ./t.sh
#
find . -name '*.go' | ./entr -r go test -v -tags runtime_test game/*.go
