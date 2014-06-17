### govendor

Minimalistic golang's package management. Our approach is:
- clones dependencies listed in `deps.json` into `_vendor` directory
- output `export GOPATH=$(pwd)/_vendor:$GOPATH` into `.env` file
- use https://github.com/kennethreitz/autoenv to activate our GOPATH

### How to use

- Create a `deps.json` file that contains your third-party dependency such as

```json
[
  {
    "vcs": "hg",
    "repo": "https://code.google.com/p/go.crypto",
    "rev": "fe6c00a82e55",
    "path": "code.google.com/p/go.crypto"
  }
]
```
