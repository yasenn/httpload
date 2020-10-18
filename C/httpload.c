#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <pthread.h>
#include <arpa/inet.h>

/* threads count */
#define TH_NUM    64

/* destination ip and port */
#define DEST_PORT 80
#ifdef DEST_IP
#else
  #define DEST_IP "1.1.1.1"	/* cloudflare */
#endif

/* timeout settings */
#define TO_ENABLED 1
#define TO_SEC    0
#define TO_NSEC   50000000		/* 50ms */

int usage(char *p){
        fprintf(stderr, "usage: %s [count]\n", p);
        exit(1);
}

void *pworker(void *arg){
        int *tid = (int *)arg;
        int sd;
        struct sockaddr_in servaddr;

        if ((sd = socket(PF_INET, SOCK_STREAM, IPPROTO_TCP)) == -1){
                printf("s(%d)", *tid);
                return NULL;
        }

        memset(&servaddr, 0, sizeof(servaddr));
        servaddr.sin_family = AF_INET;
        servaddr.sin_port = htons(DEST_PORT);
        inet_pton(AF_INET, DEST_IP, &servaddr.sin_addr);

        if (connect(sd, (struct sockaddr *)&servaddr, sizeof(servaddr)) < 0){
                printf("x(%d)", *tid);
        } else {
                printf(".");
        }

        close(sd);
}

int main(int argc, char **argv)
{
        int i, j, count = 2;
        /* threads */
        int ids[TH_NUM];
        pthread_t tid[TH_NUM];
        /* current time */
        time_t t;
        struct tm tm;
        /* timeout */
        struct timespec timeo = {
                .tv_sec  = TO_SEC,
                .tv_nsec = TO_NSEC
        };

        if (argc > 2) {
                usage(argv[0]);         /* usage */
        } else if (argc == 2) {
                count = atoi(argv[1]) + 1;
        }

        i = 1;

        do {
                t = time(NULL);
                tm = *localtime(&t);

                printf("[%d-%02d-%02d %02d:%02d:%02d] started(%d): ",
                        tm.tm_year + 1900,
                        tm.tm_mon + 1,
                        tm.tm_mday,
                        tm.tm_hour,
                        tm.tm_min,
                        tm.tm_sec,
                        i);

                for(j = 0; j < TH_NUM; j++) {
                        ids[j] = j;
                        pthread_create(&tid[j], NULL, pworker, (void*)(ids + j)); 
                }
        
                for(j = 0; j < TH_NUM; j++)
                        pthread_join(tid[j], NULL);

                if ((i != (count - 1)) && TO_ENABLED) {
                        printf("\nsleep...\n");
                        nanosleep(&timeo, NULL);
                } else if ((i != (count - 1)) && !TO_ENABLED) {
                        printf("\n");
                } else {
                        printf("\nfinished\n");
                }

        } while (++i != count);

        return 0;
}
