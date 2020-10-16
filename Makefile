default: httpload

httpload: httpload.c
	gcc -pthread -o httpload httpload.c

go: main.go
	go build  -o httpload main.go

clean:
	rm -f *.o
	rm -f httpload
