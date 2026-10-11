package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// 1. 使用 Struct 定義 Handler (傳統方式)
type aboutHandler struct {}

func (h *aboutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "這是工程變更影響分析系統 (V0) - 關於頁面\n時間: %s",time.Now().Format(time.RFC3339))
}

func main() {
	// 初始化 ServeMux (多路復用器)
	mux := http.NewServeMux()

	// 2. 使用 HandlerFunc 適配器 (簡潔方式)
	// 直接傳入一個匿名函式，並利用 http.HandlerFunc 轉換
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 簡單的路徑檢查，因為 "/" 是前綴匹配，會攔截所有路徑
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "歡迎來到 Engineering Change Impact Service API"}`))
	})

	// 註冊使用 Struct 實作的 Handler
	// 建立一個新的 aboutHandler 實例，並把它的記憶體位址存入變數 about 中
	about := &aboutHandler{}
	mux.Handle("/about", about)

	// 設定 Server 參數
	// 在後端工程中，明確設定 Timeout 是非常重要的習慣
	// http.Server 是一個包含非常多配置欄位的大型 Struct，如果不用 &，Go 就會建立一個該結構體的完整複本
	// 這在處理這種「重量級」物件時是不必要的資源浪費
	server := &http.Server{
		Addr: ":8080",
		Handler: mux,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 15 * time.Second,
	}

	fmt.Println("Server 正在啟動，監聽埠號 :8080...")
	// 啟動 Server
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server 啟動失敗: %v", err)
	}
}