# dCryptClient
Simple command line utility for uploading Download Link Container (DLC) to dcrypt.it and write the decrypted links to terminal or file.

## Install
    go get github.com/Dome91/dCryptClient

## Usage
    ./dCryptClient -h

    Usage of ./dCryptClient:
    -dlc string
            DLC file to decrypt
    -f	Decrypted links should be written to file, not to terminal

Uploading DLC and writing decrypted links to file (links.txt):

    dCryptClient -dlc test.dlc -f

Uploading DLC and writing decrypted links to command line:

    dCryptClient -dlc test.dlc