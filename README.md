# wire protocol parser
Simple application that can parse basic wire protocol from Aerospike server.

This can be ran in a web server mode where it prints out parsed incoming transactions (i.e. from ESP connector), or it can be manually fed hexadecimal strings to parse out (i.e from tcpdump)

## Installation
```bash
git clone https://github.com/colton-aerospike/wire-protocol-parser && cd wire-protocol-parser
```

## Running as web server
```bash
go run main.go
2024/12/16 23:10:54 Running webserver
```

## Ingest a hexadecimal string
```bash
go run main.go 01020300000000005d161003000000000000020000000000000000000500000000000500746573740000001504f9bb902ddd1874858475b66c8f36b950badb44f700000006016d797365740000000a0201000000000000000e000000090e0000016dea751d58
========== 0 ==========
2024/12/16 23:07:50 Protocol: 2
2024/12/16 23:07:50 Message type: 3 (message)
2024/12/16 23:07:50 Message size following this header: 93
2024/12/16 23:07:50 Header size: 22
2024/12/16 23:07:50 info1: AS_MSG_INFO1_XDR
2024/12/16 23:07:50 info2: AS_MSG_INFO2_WRITE
2024/12/16 23:07:50 info2: AS_MSG_INFO2_DELETE
2024/12/16 23:07:50 Result Code: 0
2024/12/16 23:07:50 Generation: 2
2024/12/16 23:07:50 RecTtl: 0
2024/12/16 23:07:50 TransactionTtl: 0
2024/12/16 23:07:50 n_fields: 5
2024/12/16 23:07:50 n_ops: 0
2024/12/16 23:07:50 Field 0: msgSz=5 type=AS_MSG_FIELD_TYPE_NAMESPACE(0) data=test
2024/12/16 23:07:50 Field 1: msgSz=21 type=AS_MSG_FIELD_TYPE_DIGEST_RIPE(4) data=f9bb902ddd1874858475b66c8f36b950badb44f7
2024/12/16 23:07:50 Field 2: msgSz=6 type=AS_MSG_FIELD_TYPE_SET(1) data=myset
2024/12/16 23:07:50 Field 3: msgSz=10 type=AS_MSG_FIELD_TYPE_KEY(2) data=01000000000000000e
2024/12/16 23:07:50 Field 4: msgSz=9 type=AS_MSG_FIELD_TYPE_LUT(14) data=0000016dea751d58
```

## Get hexadecimal string from tcpdump
```bash
[root@aerolab4-src-1 ~]# tcpdump -x -vvv -i lo port 8901 and greater 100
dropped privs to tcpdump
tcpdump: listening on lo, link-type EN10MB (Ethernet), snapshot length 262144 bytes
23:11:07.138925 IP6 (flowlabel 0x9e580, hlim 64, next-header TCP (6) payload length: 133) localhost.40176 > localhost.jmb-cds2: Flags [P.], cksum 0x008d (incorrect -> 0xcb01), seq 2746282963:2746283064, ack 1806250174, win 512, options [nop,nop,TS val 4021083318 ecr 4021083318], length 101
        0x0000:  6009 e580 0085 0640 0000 0000 0000 0000
        0x0010:  0000 0000 0000 0001 0000 0000 0000 0000
        0x0020:  0000 0000 0000 0001 9cf0 22c5 a3b0 f3d3
        0x0030:  6ba9 30be 8018 0200 008d 0000 0101 080a
        0x0040:  efac dcb6 efac dcb6 0203 0000 0000 005d
        0x0050:  1610 0110 0000 0000 0001 ffff ffff 0000
        0x0060:  0000 0005 0000 0000 0005 0074 6573 7400
        0x0070:  0000 1504 035a f20d be55 f811 624e 5b65
        0x0080:  4cff d520 1357 2b28 0000 0006 016d 7973
        0x0090:  6574 0000 000a 0201 0000 0000 0000 0010
        0x00a0:  0000 0009 0e00 0001 6dea 8edf 5d
```

### Grab just the hexadecimal starting from 0203 (0x0040 in this example above) and paste it into a file
```bash
cat tcpdump.text
        0x0040:  ef7c c679 ef7c c679 0203 0000 0000 009a
        0x0050:  1610 1100 0000 0000 0001 ffff ffff 0000
        0x0060:  0000 0005 0003 0000 0005 0074 6573 7400
        0x0070:  0000 1504 78f2 6d30 728c 436a c5b7 e2b7
        0x0080:  9f7f aa4b 1c01 1576 0000 0006 016d 7973
        0x0090:  6574 0000 000a 0201 0000 0000 0000 0004
        0x00a0:  0000 0009 0e00 0001 6dea 5ec9 2d00 0000
        0x00b0:  0e02 0300 046e 616d 6563 6f6c 746f 6e00
        0x00c0:  0000 0f02 0100 0361 6765 0000 0000 0000
        0x00d0:  001b 0000 0014 0213 0007 6d61 7054 6573
        0x00e0:  7482 a203 6101 a203 6202
```

### Now extract and refactor the data using awk to get a single string
```bash
egrep -o '0x.*' tcpdump.text |awk -F':' '{print $2}' |egrep -o '[a-f0-9]' |while read char; do echo -n $char; done
ef7cc679ef7cc679020300000000009a16101100000000000001ffffffff0000000000050003000000050074657374000000150478f26d30728c436ac5b7e2b79f7faa4b1c01157600000006016d797365740000000a02010000000000000004000000090e0000016dea5ec92d0000000e020300046e616d65636f6c746f6e0000000f02010003616765000000000000001b00000014021300076d61705465737482a2036101a2036202
```

# == IMPORTANT! == 
We only care about the part starting with `0203` BUT we need one more byte prior for it to parse. (e.g. `790203...`)
### 
### Paste it into the application
```bash
go run main.go 79020300000000009a16101100000000000001ffffffff0000000000050003000000050074657374000000150478f26d30728c436ac5b7e2b79f7faa4b1c01157600000006016d797365740000000a02010000000000000004000000090e0000016dea5ec92d0000000e020300046e616d65636f6c746f6e0000000f02010003616765000000000000001b00000014021300076d61705465737482a2036101a2036202
========== 0 ==========
2024/12/16 23:16:53 Protocol: 2
2024/12/16 23:16:53 Message type: 3 (message)
2024/12/16 23:16:53 Message size following this header: 154
2024/12/16 23:16:53 Header size: 22
2024/12/16 23:16:53 info1: AS_MSG_INFO1_XDR
2024/12/16 23:16:53 info2: AS_MSG_INFO2_WRITE
2024/12/16 23:16:53 info2: AS_MSG_INFO2_DURABLE_DELETE
2024/12/16 23:16:53 Result Code: 0
2024/12/16 23:16:53 Generation: 1
2024/12/16 23:16:53 RecTtl: -1
2024/12/16 23:16:53 TransactionTtl: 0
2024/12/16 23:16:53 n_fields: 5
2024/12/16 23:16:53 n_ops: 3
2024/12/16 23:16:53 Field 0: msgSz=5 type=AS_MSG_FIELD_TYPE_NAMESPACE(0) data=test
2024/12/16 23:16:53 Field 1: msgSz=21 type=AS_MSG_FIELD_TYPE_DIGEST_RIPE(4) data=78f26d30728c436ac5b7e2b79f7faa4b1c011576
2024/12/16 23:16:53 Field 2: msgSz=6 type=AS_MSG_FIELD_TYPE_SET(1) data=myset
2024/12/16 23:16:53 Field 3: msgSz=10 type=AS_MSG_FIELD_TYPE_KEY(2) data=010000000000000004
2024/12/16 23:16:53 Field 4: msgSz=9 type=AS_MSG_FIELD_TYPE_LUT(14) data=0000016dea5ec92d
2024/12/16 23:16:53 Op 0: msgSz=14 opType=AS_MSG_OP_WRITE(2) dataType=AS_PARTICLE_TYPE_STRING(3) binName=name data=colton
2024/12/16 23:16:53 Op 1: msgSz=15 opType=AS_MSG_OP_WRITE(2) dataType=AS_PARTICLE_TYPE_INTEGER(1) binName=age data=27
2024/12/16 23:16:53 Op 2: msgSz=20 opType=AS_MSG_OP_WRITE(2) dataType=AS_PARTICLE_TYPE_MAP(19) binName=mapTest data=82a2036101a2036202
```
