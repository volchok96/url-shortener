package random

import (
    "crypto/rand"
    "encoding/binary"
)

func NewRandomString(size int) string {
    chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
        "abcdefghijklmnopqrstuvwxyz" +
        "0123456789")

    b := make([]rune, size)
    for i := range b {
        var index uint32
        _ = binary.Read(rand.Reader, binary.LittleEndian, &index) // Ignore error
        b[i] = chars[int(index)%len(chars)]
    }

    return string(b)
}