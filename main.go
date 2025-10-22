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

// 定义请求的结构体
type SwapRequest struct {
	From        string  `json:"from"`
	To          string  `json:"to"`
	FromAmount  string  `json:"fromAmount"`
	Slippage    float64 `json:"slippage"`
	Payer       string  `json:"payer"`
	PriorityFee float64 `json:"priorityFee"`
	ForceLegacy bool    `json:"forceLegacy"`
}

// 发起POST请求的函数
func sendPostRequest(url string, data SwapRequest) (time.Duration, error) {
	// 将请求数据转换为JSON
	requestBody, err := json.Marshal(data)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal request body: %v", err)
	}

	// 创建POST请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 记录请求开始时间
	startTime := time.Now()

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read response body: %v", err)
	}

	// 打印响应状态和响应体
	fmt.Printf("Response Status: %s\n", resp.Status)
	fmt.Printf("Response Body: %s\n", string(respBody))

	// 计算请求延迟
	duration := time.Since(startTime)
	return duration, nil
}

// 批量发送请求的函数
func sendMultipleRequests(url string, data SwapRequest, count int) {
	for i := 0; i < count; i++ {
		// 发送请求并记录延迟
		duration, err := sendPostRequest(url, data)
		if err != nil {
			log.Printf("Request %d failed: %v\n", i+1, err)
			continue
		}

		// 打印请求延迟
		log.Printf("Request %d took %v\n", i+1, duration)
	}
}

func main() {
	url := "https://swap-v2.solanatracker.io/swap"

	// 请求体数据
	data := SwapRequest{
		From:        "EKpQGSJtjMFqKZ9KQanSqYXRcF8fBopzLHYxdM65zcjm",
		To:          "So11111111111111111111111111111111111111112",
		FromAmount:  "50%",
		Slippage:    0.5,
		Payer:       "Ef3GGMdwmgimvo2xP6MXrx3y15R2sjJzNdDLc3Nz8kVq",
		PriorityFee: 0.0001,
		ForceLegacy: false,
	}

	// 选择功能
	fmt.Println("选择测试功能: ")
	fmt.Println("1. 测试单个请求延迟")
	fmt.Println("2. 连续发起多个请求")
	var choice int
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		// 测试单个请求延迟
		duration, err := sendPostRequest(url, data)
		if err != nil {
			log.Fatalf("请求失败: %v", err)
		}
		fmt.Printf("请求延迟: %v\n", duration)

	case 2:
		// 测试多个请求
		var count int
		fmt.Println("请输入发起请求的数量:")
		fmt.Scanln(&count)
		sendMultipleRequests(url, data, count)

	default:
		fmt.Println("无效选择!")
	}
}
