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

### Test API token and print API version

```bash
# test connection and print info
dpctl info
```

### Domain utils

```bash
# list all domains
dpctl domain -a list

# get domain info
dpctl domain -d ashtray.wang,ashtray.name -a info

# create domain and add default record
dpctl domain -d abc.xyz,abcde.xyz -a create

# remove domain
dpctl domain -d abc.xyz,abcde.xyz -a remove

# enable domain
dpctl domain -d abc.xyz,abcde.xyz -a enable

# disable domain
dpctl domain -d abc.xyz,abcde.xyz -a diable

```

### Record utils

```bash
# list record of domain ashtray.wang and ashtray.name
dpctl record -a list -d ashtray.wang,ashtray.name

# create record but DO NOT CLEAR conflict records
dpctl record -a create -d ashtray.wang -f testdata/records.lst

# create record and CLEAR conflict ones
dpctl record -a create -d ashtray.wang -f testdata/records.lst -c

# create record and force create domain if not exist
dpctl record -a create -d ashtray.wang,abcde.xyz -f testdata/records.lst --force-domain

```

### TODO
 
- [x] operate domains
- [x] add more info commands
- [ ] more config entry for config file



[DNSPOD API]: https://www.dnspod.cn/docs/index.html
[demo config]: https://raw.githubusercontent.com/wangzn/dnspodapi/master/dpctl/testdata/dnspod.yaml

