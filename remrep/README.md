renames files using regular expressions

Usage of remrep:
  -p="": replacement
  -r="": regex
  -t=false: do not move for realz, only print
  
Sample usage:
```
$ ls
wat_1_thingidontwantinfilename.ogg
wat_2_thingidontwantinfilename.ogg
...
wat_1231_thingidontwantinfilename.ogg

$ remrep -r (wat_\d+)_thingidontwantinfilename\.ogg -p $1.ogg -t
wat_1_thingidontwantinfilename.ogg -> wat_1.ogg
wat_2_thingidontwantinfilename.ogg -> wat_2.ogg
...
wat_1231_thingidontwantinfilename.ogg -> wat_1231.ogg

$ remrep -r (wat_\d+)_thingidontwantinfilename\.ogg -p $1.ogg

$ ls
wat_1.ogg
wat_2.ogg
...
wat_1231.ogg
```