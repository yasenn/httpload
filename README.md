# HTTPload - tool for test webserver availability under highload

Upd: use `golang` implementation instead of `httload.c`

`httpload` tries to connect to web-server with provided verb (`GET`|`POST`) 

You can modify:

* concurrent threads count
* requests count

# Usage
```
httpload https://google.com -b POST -c 100 -n 1000
```

# Build go module

```
make go
./httpload https://google.com
```

## Clean build

```
make clean
```

# Legacy `httpload`

## Build 

```
make
./httpload <count>
```

## Clean build

```
make clean
```

## Build settings

| Const      | Value    | Description            |
| :---       | ---:     | ---                    |
| TH_NUM     | 64       | Threads limit          |
| DEST_PORT  | 80       | Destination port       |
| DEST_IP    | 1.1.1.1  | Description IP         |
| TO_ENABLED | 1        | Timeout enabled        |
| TO_SEC     | 0        | Timeout, seconds       |
| TO_NSEC    | 50000000 | Timeout, ns	/* 50ms */ |
