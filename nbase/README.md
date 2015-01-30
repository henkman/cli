converts decimal numbers to other number systems and vice versa

Usage of nbase:
  -b=0: base <= length of digits string
  -d="0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ": digits to be used
  -n=0: number
  -s="": string

Sample usage:
```
$ nbase -n 255 -b 16
ff

$ nbase -s ff -b 16
255

$ nbase -s hello -b 60
223420884
```
