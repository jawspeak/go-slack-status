#!/usr/bin/python

# This is only to be called from the Makefile in the deploy process

import json
import sys
import fileinput

if __name__ == "__main__":
  config_base = sys.argv[1]
  env = sys.argv[2]
  stash_password = None

  with open(env, 'r') as j:
    stash_password = json.load(j)["stash_password"]
  if not stash_password:
    raise "missing stash_password in %s" % env

  with open(config_base, 'r') as f:
    print(f.read().replace('<<STASH_PASSWORD>>', stash_password))
