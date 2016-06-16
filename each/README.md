# contains
utility to execute a command for each line coming from stdin

```
Usage of each:
  -do string
        command for each element, $0 is complete match, $1..$n are groups
  -e    ignore errors
  -re string
        regex to use for partitioning, if empty whole string is used
  -t    don't do, just print
```

Sample:
```
$ ls -l | each -do "echo $0"
total 2265
-rw-r--r-- 1 user group group 638 Jun 16 15:04 README.md
-rwxr-xr-x 1 user group group 2310656 Jun 16 15:03 each.exe
-rw-r--r-- 1 user group group 1585 Jun 16 15:03 main.go

$ ls -l | each -re "^[rwx-]+" -do "echo $0" -t
echo -rw-r--r--
echo -rwxr-xr-x
echo -rw-r--r--

$ ls -l | each -re "^[d-]([rwx-]{3})([rwx-]{3})([rwx-]{3})" -do "echo $1 $2 $3"
rw- r-- r--
rwx r-x r-x
rw- r-- r--

```