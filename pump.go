package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// 请求体结构体
type TradeRequest struct {
	Action          string  `json:"action"`
	Amount          string  `json:"amount"`
	DenominatedInSol string `json:"denominatedInSol"`
	Mint            string  `json:"mint"`
	Pool            string  `json:"pool"`
	PriorityFee     float64 `json:"priorityFee"`
	PublicKey       string  `json:"publicKey"`
	Slippage        int     `json:"slippage"`
}

// 发送POST请求函数
func sendPostRequest(url string, data TradeRequest) (time.Duration, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return 0, fmt.Errorf("marshal error: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return 0, fmt.Errorf("create request error: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	start := time.Now()
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	duration := time.Since(start)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return duration, fmt.Errorf("read body error: %v", err)
	}

	fmt.Printf("=== Response ===\nStatus: %s\nBody: %s\nDuration: %v\n\n", resp.Status, string(respBody), duration)
	return duration, nil
}

// 连续发起多个请求测试速率限制
func sendMultipleRequests(url string, data TradeRequest, count int, interval time.Duration) {
	for i := 0; i < count; i++ {
		fmt.Printf("----- Request #%d -----\n", i+1)
		duration, err := sendPostRequest(url, data)
		if err != nil {
			log.Printf("Request %d failed: %v\n", i+1, err)
		} else {
			log.Printf("Request %d completed in %v\n", i+1, duration)
		}

		// 等待间隔（避免太快触发速率限制）
		time.Sleep(interval)
	}
}

func main() {
	url := "https://pumpportal.fun/api/trade-local"

	// 构造请求体
	data := TradeRequest{
		Action:          "sell",
		Amount:          "10%",
		DenominatedInSol: "false",
		Mint:            "EKpQGSJtjMFqKZ9KQanSqYXRcF8fBopzLHYxdM65zcjm",
		Pool:            "auto",
		PriorityFee:     0.005,
		PublicKey:       "Ef3GGMdwmgimvo2xP6MXrx3y15R2sjJzNdDLc3Nz8kVq",
		Slippage:        10,
	}

	fmt.Println("选择测试模式:")
	fmt.Println("1. 单次请求延迟测试")
	fmt.Println("2. 多次请求速率测试")

	var choice int
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		// 单次请求
		duration, err := sendPostRequest(url, data)
		if err != nil {
			log.Fatalf("请求失败: %v", err)
		}
		fmt.Printf("请求总耗时: %v\n", duration)

	case 2:
		// 多次请求
		var count int
		fmt.Println("请输入要发送的请求数量:")
		fmt.Scanln(&count)

		var intervalMs int
		fmt.Println("请输入每次请求间隔 (毫秒):")
		fmt.Scanln(&intervalMs)

		sendMultipleRequests(url, data, count, time.Duration(intervalMs)*time.Millisecond)

	default:
		fmt.Println("无效选项。")
	}
}
