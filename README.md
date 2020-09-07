# HTTPload - tool for test webserver availability under highload

# Usage
```
httpload <count>
```

# Build

```
make 
./httpload 100
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
