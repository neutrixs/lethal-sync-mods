package util

import "fmt"

func FormatByteSize(size int64) string {
    // we won't need size bigger than TB, won't we?
    format := []string{"B","KB","MB","GB","TB"}
    tempsize := float64(size)

    i := 0
    for ; i < 5; i++ {
        if tempsize < 1024 {
            break
        }
        
        tempsize = tempsize / (1 << 10)
    }

    return fmt.Sprintf("%.2f%s", tempsize, format[i])
}