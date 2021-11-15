## go-send-to

There's no such an api in ethereum which can get all the transactions send to a specific address, so I build this project to achieve it.

There's no magic way to do it, the only way is to loop the blocks, but take advantage of golang's concurrency & performance, it should be faster.

## get started

```bash
git clone https://github.com/huahuayu/go-send-to.git
cd go-send-to
go run main.go --address=0x10ED43C718714eb63d5aA57B78B54704E256024E --from=12676566 --to=12676576 --workers=10
```

supported flags

```text
  -address string
        address (default "0x10ED43C718714eb63d5aA57B78B54704E256024E")
  -from uint
        from block (default 12676576)
  -node string
        node url (default "https://bsc-dataseed.binance.org/")
  -to uint
        to block (default 12676586)
  -workers int
        workers to deal block concurrently (default 10)
```
