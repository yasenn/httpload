default: httpload

httpload: httpload.c
	gcc -pthread -o httpload httpload.c

clean:
	rm -f *.o
	rm -f httpload
