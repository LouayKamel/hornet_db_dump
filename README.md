# Hornet DB Dump
Dumps hornet txs to a file

## Usage
```
Usage of ./main:
  -dbPath string
        directory that contains the tangle.db. e.g mainnetdb
  -output string
        output file to store the dump (default "hornet_tx_dump.txt")
```

You should see lines in the console like.
```
....
....
2020/09/09 00:42:28 Z9ZFQMJVVL9HXNANMKIOHRGTWXEQROHHVMIGEYRSHF9JRFXAADGGYAADJVVC9VTCJYCUQNDAHUGXA9999  done...
2020/09/09 00:42:28 ZIZOINCCSSVANAHSUB9WYFXSZFXQLAYITLIKFPVO9H9HJ9UJNKTWO9MJLMLOIBXMZWHHSCPBBZPGZ9999  done...
2020/09/09 00:42:28 ZR99HM9TIXC9STCLKMQUXXPOFZ9CGSQY9A9ZWBFYUJXLKBCYQHZFQBIQCSDUCYWASZKNKQOZFAJX99999  done...
2020/09/09 00:42:28 ZR9CLKSOTPM9NACUGBXJRVAWY9DLWIDJI9BPPHTZPAPW9DFCBYSMKPRPMCCGZ9NAIHWQCZEDSA9UA9999  done...
2020/09/09 00:42:28 Total txs:  56746
2020/09/09 00:42:28 Success:  56746
```

If an error occurs parsing a tx, an error will be printed on the console. The `successCount` will not be incremented.