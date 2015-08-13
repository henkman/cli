Searches files by filename and content using regular expressions
(port of search to C)

Build with 
$ gcc -Wall -o pose pose.c -O3 -s -nostdlib -fno-asynchronous-unwind-tables -fno-ident -ffunction-sections -static -Wl,-e,__main -lkernel32 -lshlwapi

```
Usage of pose:
  -a=false: only show absolute paths
  -b=".": the base directory
  -c="": regex of thing(s) that should be in file(s)
  -f=false: stop on the first occurance
  -ic=false: ignore case in contains regex
  -is=false: ignore case in file regex
  -s="": regex of file(s) to search for
```
