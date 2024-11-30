#!/bin/bash

/home/keith/restler_bin/restler/Restler fuzz-lean \
              --grammar_file Compile/grammar.py \
              --dictionary_file Compile/dict.json \
              --settings Compile/engine_settings.json \
              --target_ip 127.0.0.1 \
              --target_port 80 \
              --no_ssl
