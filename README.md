# teller

Tool for getting basic isights on spending based on manually downloaded bank/credit card statements. Banks, generally, don't allow users to get spending data from their accounts in an automated way without an API key. There are some apps which provide have keys to major banks and provide spending insights (ie Mint), but not all banks are supported, in particular non-US banks. This tool supports parsing files from some banks which are not supported on these platforms.

## Running

```bash
$ git clone https://github.com/yobrosoft/teller.git
$ cd teller
$ make build
$ ./out/teller -h
```
