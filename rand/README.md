generates random numbers using wincrypt api

Compile:
```
gcc -Wall -std=gnu11 -pedantic -o rand rand.c -s -O3 -nostdlib -static -fno-asynchronous-unwind-tables -fno-ident -ffunction-sections -lkernel32 -ladvapi32
```

```
Usage of rand:
  -s:  lower range end
  -e:  upper range end
  -c:  coin mode, only prints yes or no, and returns 1 or 0
  
$ rand -c
yes

$ rand -s 1 -e 10
4
```
  