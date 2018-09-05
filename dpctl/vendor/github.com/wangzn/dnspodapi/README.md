# dnspodapi

Go client wrapper and CTL tool for [DNSPOD API]

## INSTALL

Make sure you have go installed in your env.

```bash
go get -u github.com/wangzn/dnspodapi/dpctl
```

## Usage

```bash
dpctl help
```

### Config

Default config file is at $HOME/.dnspod.yaml

Demo config file could be found at [demo config]

### operate record

```bash
# list record of domain ashtray.wang and ashtray.name
dpctl record -a list -z ashtray.wang,ashtray.name

# create record but DO NOT CLEAR conflict records
dpctl record -a create -z ashtray.wang -f testdata/records.lst

# create record and CLEAR conflict ones
dpctl record -a create -z ashtray.wang -f testdata/records.lst -c
```

### TODO
 
- [ ] operate domains
- [ ] add more info commands




[DNSPOD API]: https://www.dnspod.cn/docs/index.html
[demo config]: https://raw.githubusercontent.com/wangzn/dnspodapi/master/dpctl/testdata/dnspod.yaml

