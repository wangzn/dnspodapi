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

# get record info of domain ashtray.wang, with record proxy,proxy11,proxy12
dpctl record -a info -d ashtray.wang -r proxy,proxy11,proxy12

# disalbe these records and enable them again
dpctl record -a disable -d ashtray.wang -r proxy,proxy11,proxy12
dpctl record -a eanble -d ashtray.wang -r proxy,proxy11,proxy12

# create a record
dpctl  record -a create -d ashtray.wang -r testcreate -t A -v 1.1.1.1

## import records from local file but DO NOT CLEAR conflict records
dpctl record -a import -d ashtray.wang -f testdata/records.lst

# import records from local file and CLEAR conflict ones
dpctl record -a import -d ashtray.wang -f testdata/records.lst -c

# import record and force create domain if not exist
dpctl record -a import -d ashtray.wang,abcde.xyz -f testdata/records.lst --force-domain

# export records to file, append if file exist
dpctl record -d ashtray.wang -a export -f /dev/stdout --export-file-mode append

# export records to file, with records filter www21,www22, and overwrite local file if exist
dpctl record -d ashtray.wang -a export -f testdata/export.lst --export-file-mode overwrite -r www11,www12

```

### Playbook utils
We can predefined some common used cmds into playbook (named as __scene__ ), and reuse them in short command.

```bash
dpctl help playbook
```

### Scene definition

__scene__ is a list of __actions__ , e.g.:

```yaml
  scene1:
    - auth: "default" # could omit to use default
      category: "domain"
      action: "create"
      subject: "abc.xyz,abcde.xyz"
      params:

    - auth: "default" # could omit to use default
      category: "record"
      action: "import"
      subject: ""
      params:
        domain: "abc.xyz,abcde.xyz"
        clear: "on"
        force_domain: "off"
        record_file: "testdata/record.lst"

    - auth: "default" # could omit to use default
      category: "record"
      action: "export"
      subject: "www,proxy"
      params:
        domain: "abc.xyz,abcde.xyz"
        file_mode: "append" # overwrite, default is exit for no damage
        record_file: "testdata/record.lst"

```

In this scene above, we define three actions:
  
  * create domain `abc.xyz,abcde.xyz`
  * import record from local file `testdata/record.lst` into domain `abc.xyz,abcde.xyz`
  * export record into local file `testdata/record.lst` with record `www` and `proxy` in domain `abc.xyz,abcde.xyz`

### Preview scene
then we could preview this scene in case of any unexpected cmds

```bash
dpctl playbook -a preview -s scene1
```

### Run scene

finally, we could run this scene

```bash
dpctl playbook -a run -s scene1
```


### TODO
 
- [x] operate domains
- [x] add more info commands
- [x] more config entry for config file
- [x] add playbook cmd and modify config file structure



[DNSPOD API]: https://www.dnspod.cn/docs/index.html
[demo config]: https://raw.githubusercontent.com/wangzn/dnspodapi/master/dpctl/testdata/dnspod.yaml

