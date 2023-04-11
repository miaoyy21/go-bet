#!/bin/sh

make release

mv go-bet_linux_386 $DOCUMENTS/start/bet_150.158.171.47/os/bet/go-bet_linux_386


# cd ~/start/bet_150.158.171.47
# rm -rf os.zip
# zip -qr os.zip ~/start/bet_150.158.171.47/os
# scp -i 20220427 os.zip root@150.158.171.47:~

# ssh -i 20220427 root@150.158.171.47
