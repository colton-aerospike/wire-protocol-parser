package main

import (
        "encoding/binary"
        "encoding/hex"
        "fmt"
        "io"
        "log"
        "net/http"
        "os"
)

var fieldTypes = []string{
        "AS_MSG_FIELD_TYPE_NAMESPACE",
        "AS_MSG_FIELD_TYPE_SET",
        "AS_MSG_FIELD_TYPE_KEY",
        "AS_MSG_FIELD_TYPE_BIN",
        "AS_MSG_FIELD_TYPE_DIGEST_RIPE",
        "AS_MSG_FIELD_TYPE_GU_TID",
        "AS_MSG_FIELD_TYPE_DIGEST_RIPE_ARRAY",
        "AS_MSG_FIELD_TYPE_TRID",
        "AS_MSG_FIELD_TYPE_SCAN_OPTIONS",
        "AS_MSG_FIELD_TYPE_SOCKET_TIMEOUT",
        "AS_MSG_FIELD_TYPE_RECS_PER_SEC",
        "AS_MSG_FIELD_TYPE_PID_ARRAY",
        "AS_MSG_FIELD_TYPE_DIGEST_ARRAY",
        "AS_MSG_FIELD_TYPE_SAMPLE_MAX",
        "AS_MSG_FIELD_TYPE_LUT",
        "AS_MSG_FIELD_TYPE_BVAL_ARRAY",
        "",
        "",
        "",
        "",
        "",
        "AS_MSG_FIELD_TYPE_INDEX_NAME",
        "AS_MSG_FIELD_TYPE_INDEX_RANGE",
        "AS_MSG_FIELD_TYPE_INDEX_CONTEXT",
        "",
        "",
        "AS_MSG_FIELD_TYPE_INDEX_TYPE",
        "",
        "",
        "",
        "AS_MSG_FIELD_TYPE_UDF_FILENAME",
        "AS_MSG_FIELD_TYPE_UDF_FUNCTION",
        "AS_MSG_FIELD_TYPE_UDF_ARGLIST",
        "AS_MSG_FIELD_TYPE_UDF_OP",
        "",
        "",
        "",
        "",
        "",
        "",
        "AS_MSG_FIELD_TYPE_QUERY_BINLIST",
        "AS_MSG_FIELD_TYPE_BATCH",
        "AS_MSG_FIELD_TYPE_BATCH_WITH_SET",
        "AS_MSG_FIELD_TYPE_PREDEXP",
}

var opTypes = []string{
        "",
        "AS_MSG_OP_READ",
        "AS_MSG_OP_WRITE",
        "AS_MSG_OP_CDT_READ",
        "AS_MSG_OP_CDT_MODIFY",
        "AS_MSG_OP_INCR",
        "",
        "AS_MSG_OP_EXP_READ",
        "AS_MSG_OP_EXP_MODIFY",
        "AS_MSG_OP_APPEND",
        "AS_MSG_OP_PREPEND",
        "AS_MSG_OP_TOUCH",
        "AS_MSG_OP_BITS_READ",
        "AS_MSG_OP_BITS_MODIFY",
        "AS_MSG_OP_DELETE_ALL",
        "AS_MSG_OP_HLL_READ",
        "AS_MSG_OP_HLL_MODIFY",
}

var binTypes = []string{
        "",
        "AS_PARTICLE_TYPE_INTEGER",
        "AS_PARTICLE_TYPE_FLOAT",
        "AS_PARTICLE_TYPE_STRING",
        "AS_PARTICLE_TYPE_BLOB",
        "",
        "",
        "AS_PARTICLE_TYPE_JAVA_BLOB",
        "AS_PARTICLE_TYPE_CSHARP_BLOB",
        "AS_PARTICLE_TYPE_PYTHON_BLOB",
        "AS_PARTICLE_TYPE_RUBY_BLOB",
        "AS_PARTICLE_TYPE_PHP_BLOB",
        "AS_PARTICLE_TYPE_ERLANG_BLOB",
        "",
        "",
        "",
        "",
        "AS_PARTICLE_TYPE_BOOL",
        "AS_PARTICLE_TYPE_HLL",
        "AS_PARTICLE_TYPE_MAP",
        "AS_PARTICLE_TYPE_LIST",
        "",
        "",
        "AS_PARTICLE_TYPE_GEOJSON",
}

func handler(w http.ResponseWriter, r *http.Request) {
        // Read the body of the request
        body, err := io.ReadAll(r.Body)
        if err != nil {
                log.Printf("Error reading body: %v", err)
                http.Error(w, "can't read body", http.StatusBadRequest)
                return
        }

        // It's a good practice to close the body when you're done with it
        defer r.Body.Close()
        d := hex.EncodeToString(body)
        log.Printf("Req body (hex): %s", d)
        handleProtocol(body)
}

func handleProtocol(body []byte) {
        // Convert the body to a string and log it
        protocol := body[1]
        log.Printf("Protocol: %d", int(protocol))

        msgType := body[2]
        switch msgType {
        case 1:
                log.Print("Message type: 1 (info)")
        case 3:
                log.Print("Message type: 3 (message)")
        default:
                log.Printf("Unknown message type: %d, skipping", int(msgType))
                return
        }

        // uint64 - needs 8 bytes, not 6
        szBytes := append([]byte{0, 0}, body[3:9]...)
        sz := binary.BigEndian.Uint64(szBytes)
        log.Printf("Message size following this header: %d", sz)

        if msgType == 1 {
                log.Print("Not handling info, skipping")
                return
        }

        // must be 22 bytes in size in protocol version 2
        if body[9] != 22 {
                log.Printf("Wrong header size: %d, skipping", int(body[9]))
                return
        }
        log.Print("Header size: 22")

        info1 := body[10]
        if info1&(1<<0) != 0 {
                log.Print("info1: AS_MSG_INFO1_READ")
        }
        if info1&(1<<1) != 0 {
                log.Print("info1: AS_MSG_INFO1_GET_ALL")
        }
        if info1&(1<<2) != 0 {
                log.Print("info1: AS_MSG_INFO1_SHORT_QUERY")
        }
        if info1&(1<<3) != 0 {
                log.Print("info1: AS_MSG_INFO1_BATCH")
        }
        if info1&(1<<4) != 0 {
                log.Print("info1: AS_MSG_INFO1_XDR")
        }
        if info1&(1<<5) != 0 {
                log.Print("info1: AS_MSG_INFO1_GET_NO_BINS")
        }
        if info1&(1<<6) != 0 {
                log.Print("info1: AS_MSG_INFO1_CONSISTENCY_LEVEL_ALL")
        }
        if info1&(1<<7) != 0 {
                log.Print("info1: AS_MSG_INFO1_COMPRESS_RESPONSE")
        }

        info2 := body[11]
        if info2&(1<<0) != 0 {
                log.Print("info2: AS_MSG_INFO2_WRITE")
        }
        if info2&(1<<1) != 0 {
                log.Print("info2: AS_MSG_INFO2_DELETE")
        }
        if info2&(1<<2) != 0 {
                log.Print("info2: AS_MSG_INFO2_GENERATION")
        }
        if info2&(1<<3) != 0 {
                log.Print("info2: AS_MSG_INFO2_GENERATION_GT")
        }
        if info2&(1<<4) != 0 {
                log.Print("info2: AS_MSG_INFO2_DURABLE_DELETE")
        }
        if info2&(1<<5) != 0 {
                log.Print("info2: AS_MSG_INFO2_CREATE_ONLY")
        }
        if info2&(1<<6) != 0 {
                log.Print("info2: unused byte set")
        }
        if info2&(1<<7) != 0 {
                log.Print("info2: AS_MSG_INFO2_RESPOND_ALL_OPS")
        }

        info3 := body[12]
        if info3&(1<<0) != 0 {
                log.Print("info3: AS_MSG_INFO3_LAST")
        }
        if info3&(1<<1) != 0 {
                log.Print("info3: AS_MSG_INFO3_COMMIT_LEVEL_MASTER")
        }
        if info3&(1<<2) != 0 {
                log.Print("info3: AS_MSG_INFO3_PARTITION_DONE")
        }
        if info3&(1<<3) != 0 {
                log.Print("info3: AS_MSG_INFO3_UPDATE_ONLY")
        }
        if info3&(1<<4) != 0 {
                log.Print("info3: AS_MSG_INFO3_CREATE_OR_REPLACE")
        }
        if info3&(1<<5) != 0 {
                log.Print("info3: AS_MSG_INFO3_REPLACE_ONLY")
        }
        if info3&(1<<6) != 0 {
                log.Print("info3: AS_MSG_INFO3_SC_READ_TYPE")
        }
        if info3&(1<<7) != 0 {
                log.Print("info3: AS_MSG_INFO3_SC_READ_RELAX")
        }

        resultCode := body[14]
        log.Printf("Result Code: %d", int(resultCode))

        // uint32 - 4 bytes
        generationBytes := body[15:19]
        generation := binary.BigEndian.Uint32(generationBytes)
        log.Printf("Generation: %d", generation)

        recTtlBytes := body[19:23]
        recTtl := int32(binary.BigEndian.Uint32(recTtlBytes))
        log.Printf("RecTtl: %d", recTtl)

        transactionTtlBytes := body[23:27]
        transactionTtl := int32(binary.BigEndian.Uint32(transactionTtlBytes))
        log.Printf("TransactionTtl: %d", transactionTtl)

        nFieldsBytes := body[27:29]
        nFields := binary.BigEndian.Uint16(nFieldsBytes)
        log.Printf("n_fields: %d", nFields)

        nOpsBytes := body[29:31]
        nOps := binary.BigEndian.Uint16(nOpsBytes)
        log.Printf("n_ops: %d", nOps)

        // data starts at byte 31
        offset := 31
        for i := uint16(0); i < nFields; i++ {
                msgSzBytes := body[offset : offset+4]
                msgSz := int32(binary.BigEndian.Uint32(msgSzBytes))
                offset += 4

                fieldType := body[offset]
                offset++

                data := body[offset : offset+int(msgSz)-1]
                var d string
                if fieldType == 0 || fieldType == 1 || fieldType == 21 || fieldType == 31 || fieldType == 30 {
                        d = string(data)
                } else {
                        d = hex.EncodeToString(data)
                }
                offset += (int(msgSz) - 1)
                log.Printf("Field %d: msgSz=%d type=%s(%d) data=%s", i, msgSz, fieldTypes[fieldType], fieldType, d)
        }

        for i := uint16(0); i < nOps; i++ {
                msgSzBytes := body[offset : offset+4]
                msgSz := int32(binary.BigEndian.Uint32(msgSzBytes))
                offset += 4

                opType := body[offset]
                offset++

                dataType := body[offset]
                offset++

                // unused
                offset++

                binNameLen := body[offset]
                offset++
                binName := body[offset : offset+int(binNameLen)]
                offset += int(binNameLen)

                data := body[offset : offset+int(msgSz)-4-int(binNameLen)]
                // handle data for all types
                d := hex.EncodeToString(data)
                offset += (int(msgSz) - 4 - int(binNameLen))
                if dataType == 3 {
                        // it's a string special handling
                        d = string(data)
                } else if dataType == 1 {
                        // it's an int, special handling
                        for len(data) < 8 {
                                data = append([]byte{0}, data...)
                        }
                        d = fmt.Sprintf("%d", int(binary.BigEndian.Uint64(data)))
                }

                log.Printf("Op %d: msgSz=%d opType=%s(%d) dataType=%s(%d) binName=%s data=%s", i, msgSz, opTypes[opType], opType, binTypes[dataType], dataType, binName, d)
        }

        if offset != len(body) {
                log.Printf("WARNING: Extra data at end of body (sz: %d offset: %d)", len(body), offset)
                data := body[offset:]
                d := hex.EncodeToString(data)
                log.Print(d)
        }
}

func main() {
        if len(os.Args) == 1 {
                http.HandleFunc("/", handler)
                log.Print("Running webserver")
                log.Fatal(http.ListenAndServe(":8080", nil))
        }
        for i, body := range os.Args[1:] {
                fmt.Printf("========== %d ==========\n", i)
                byteArray, err := hex.DecodeString(body)
                if err != nil {
                        log.Fatal(err)
                }
                handleProtocol(byteArray)
                fmt.Println("")
        }
}
