hold
==============
高性能日志库

### go get

```bash
go get github.com/Bestfeel/hole
```

### go test 

```bash
INFO[2018-03-02T19:06:03+08:00] log info                                     
INFO[2018-03-02T19:06:03+08:00] log info                                      name=hole
ERRO[2018-03-02T19:06:03+08:00] log error                                    
ERRO[2018-03-02T19:06:03+08:00] log error                                     name=hole
WARN[2018-03-02T19:06:03+08:00] log warn                                     
WARN[2018-03-02T19:06:03+08:00] log warn                                      name=hole
```


### hole

```bash
cat /tmp/hole.log|hole
```