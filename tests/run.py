#!/usr/bin/env python

from __future__ import print_function
from hashlib import sha512
from subprocess import Popen, PIPE, STDOUT

data = open("index.html", "r").read()

for line in open("cmds.txt", "r"):
    line = line.strip()
    p = Popen(['pup', line], stdout=PIPE, stdin=PIPE, stderr=PIPE)
    h = sha512()
    h.update(p.communicate(input=data)[0])
    print("%s %s" % (h.hexdigest(), line))
