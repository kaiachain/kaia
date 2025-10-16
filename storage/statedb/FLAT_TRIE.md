# FlatTrie development guide

## How to import erigon repo

Importing the erigon library, especially importing a fork repo isn't trivial. Precisely follow below steps to import your erigon-lib, or bump the imported version.

1. Choose the repository url and commit hash. You can use your own fork. Do not use tags (like v3.0.0). Use commit hash or branch name. Commit hash is the most reliable way though.
```
url: github.com/kaiachain/kaia-erigon
commit: 1a610b3f57
```
2. Go get your repo. It will fail, but we must harvest the error message.
```
$ go get github.com/kaiachain/kaia-erigon@1a610b3f57
go: github.com/kaiachain/kaia-erigon@1a610b3f57 (v1.9.7-0.20250723081440-1a610b3f574d) requires github.com/kaiachain/kaia-erigon@v1.9.7-0.20250723081440-1a610b3f574d: parsing go.mod:
        module declares its path as: github.com/erigontech/erigon
                but was required as: github.com/kaiachain/kaia-erigon
```
3. Insert the following lines in the `go.mod` file, if not exists. For the version string, copy from the above `go get` error message. But change the version number to `v0.0.0-` despite the error message.
```
- Good example: v0.0.0-20250723081440-1a610b3f574d
- Bad example:  v1.9.7-20250723081440-1a610b3f574d    (nonzero version number)
- Bad example:  v0.0.0-0.20250723081440-1a610b3f574d  (extra fourth version number)
```
```
replace (
	github.com/erigontech/erigon => github.com/kaiachain/kaia-erigon v0.0.0-20250723081440-1a610b3f574d
	github.com/erigontech/erigon-lib => github.com/kaiachain/kaia-erigon/erigon-lib v0.0.0-20250723081440-1a610b3f574d
	github.com/holiman/bloomfilter/v2 => github.com/AskAlexSharov/bloomfilter/v2 v2.0.9
)
```
4. Import the package under the name of original repo. e.g. in flat_trie_test.go
```go
import (
	"testing"

	"github.com/erigontech/erigon-lib/commitment"
)

func Test_FlatTrie_Import(t *testing.T) {
	t.Log(commitment.ModeDirect)
}
```
5. If the compiler complains, run `go mod tidy`.
```
$ go mod tidy
```
