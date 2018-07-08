# You need the following packages for this to work:
# gofmt
# golang.org/x/tools/cmd/goimports
# github.com/droptheplot/abcgo
# github.com/fzipp/gocyclo
# github.com/alexkohler/nakedret
# stathat.com/c/splint
# github.com/mibk/dupl
# github.com/qiniu/checkstyle/gocheckstyle
# github.com/jgautheron/goconst/cmd/goconst
# golang.org/x/lint/golint
# github.com/fsamin/nofuncflags
# github.com/jgautheron/usedexports
# github.com/kyoh86/scopelint
# github.com/tsenart/deadcode
# github.com/gordonklaus/ineffassign
# github.com/alexkohler/prealloc


echo ----gofmt----
echo ""
gofmt -s -w -e ./
echo ""
echo ----END----

echo ""

echo ----goimports----
echo ""
goimports -e -w ./
echo ""
echo ----END----

echo ""

echo ----abcgo----
echo ""
echo If these are above 20, please reduce the complexity
echo ""
abcgo -path ./ -sort 2>&1 | head -n 4
echo ""
echo ----END----

echo ""

echo ----gocyclo----
echo ""
echo If any functions appear below, please reduce their complexity
echo ""
gocyclo -over 10 ./
echo ""
echo ----END----

echo ""

echo ----nakedret----
echo ""
echo If any functions appear below, please replace the naked returns
echo ""
nakedret ./...
echo ""
echo ----END----

echo ""

echo ----splint----
echo ""
splint **/*.go
echo ""
echo ----END----
echo ""

echo ""

echo ----dupl----
echo ""
echo This can be erroneous in certain cases. If there is code duplication, please consolidate it into one function
echo ""
dupl -t 30
echo ""
echo ----END----
echo ""

echo ""

echo ----goconst----
echo ""
goconst -numbers ./...
echo ""
echo ----END----
echo ""

echo ""

echo ----golint----
echo ""
golint -min_confidence 0.0
echo ""
echo ----END----
echo ""

echo ""

echo ----nofuncflags----
echo ""
nofuncflags
echo ""
echo ----END----
echo ""

echo ""

echo ----usedexports----
echo ""
usedexports ./...
echo ""
echo ----END----
echo ""

echo ""

echo ----govet----
echo ""
go vet
echo ""
echo ----END----
echo ""

echo ""

echo ----scopelint----
echo ""
scopelint ./...
echo ""
echo ----END----
echo ""

echo ""

echo ----deadcode----
echo ""
deadcode
echo ""
echo ----END----
echo ""

echo ""

echo ----ineffassign----
echo ""
ineffassign ./
echo ""
echo ----END----
echo ""

echo ""

echo ----prealloc----
echo ""
prealloc -simple=false -forloops  ./...
echo ""
echo ----END----
echo ""

echo ""

echo ----go test----
echo ""
echo Please keep code coverage above 90%
echo ""
go test -cover
echo ""
echo ----END----
echo ""

echo ""