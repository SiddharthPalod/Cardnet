package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	grpcapi "cardnet/internal/api/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr        = flag.String("addr", "localhost:50051", "Auth gateway address")
	concurrency = flag.Int("concurrency", 10, "Number of concurrent workers")
	duration    = flag.Duration("duration", 30*time.Second, "Test duration")
	rate        = flag.Int("rate", 100, "Requests per second per worker")
	merchantID  = flag.String("merchant", "merchant_001", "Merchant ID")
)

type Stats struct {
	TotalRequests    int64
	SuccessfulReqs   int64
	FailedReqs       int64
	RateLimitedReqs  int64
	TimeoutReqs      int64
	TotalLatency     int64 // microseconds
	MinLatency       int64
	MaxLatency       int64
}

func main() {
	flag.Parse()

	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := grpcapi.NewAuthGatewayClient(conn)

	stats := &Stats{
		MinLatency: int64(^uint64(0) >> 1), // Max int64
	}

	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), *duration)
	defer cancel()

	startTime := time.Now()

	// Start workers
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go worker(ctx, client, stats, &wg, i)
	}

	// Wait for all workers
	wg.Wait()

	duration := time.Since(startTime)

	// Print results
	printResults(stats, duration)
}

func worker(ctx context.Context, client grpcapi.AuthGatewayClient, stats *Stats, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	ticker := time.NewTicker(time.Second / time.Duration(*rate))
	defer ticker.Stop()

	requestID := 0

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			requestID++
			reqID := fmt.Sprintf("req-%d-%d", workerID, requestID)

			start := time.Now()
			resp, err := client.Authorize(ctx, &grpcapi.AuthRequest{
				RequestId:  reqID,
				MerchantId: *merchantID,
				CardToken:  fmt.Sprintf("tok_visa_%04d", requestID%10000),
				Amount:     1000 + int64(requestID%10000),
				Currency:   "INR",
				Mcc:        "5411",
			})

			latency := time.Since(start).Microseconds()
			atomic.AddInt64(&stats.TotalRequests, 1)
			atomic.AddInt64(&stats.TotalLatency, latency)

			// Update min/max latency
			for {
				oldMin := atomic.LoadInt64(&stats.MinLatency)
				if latency >= oldMin {
					break
				}
				if atomic.CompareAndSwapInt64(&stats.MinLatency, oldMin, latency) {
					break
				}
			}

			for {
				oldMax := atomic.LoadInt64(&stats.MaxLatency)
				if latency <= oldMax {
					break
				}
				if atomic.CompareAndSwapInt64(&stats.MaxLatency, oldMax, latency) {
					break
				}
			}

			if err != nil {
				atomic.AddInt64(&stats.FailedReqs, 1)
				if ctx.Err() == context.DeadlineExceeded {
					atomic.AddInt64(&stats.TimeoutReqs, 1)
				}
			} else {
				atomic.AddInt64(&stats.SuccessfulReqs, 1)
				if resp.Status == "RATE_LIMITED" {
					atomic.AddInt64(&stats.RateLimitedReqs, 1)
				}
			}
		}
	}
}

func printResults(stats *Stats, duration time.Duration) {
	total := atomic.LoadInt64(&stats.TotalRequests)
	successful := atomic.LoadInt64(&stats.SuccessfulReqs)
	failed := atomic.LoadInt64(&stats.FailedReqs)
	rateLimited := atomic.LoadInt64(&stats.RateLimitedReqs)
	timeouts := atomic.LoadInt64(&stats.TimeoutReqs)
	totalLatency := atomic.LoadInt64(&stats.TotalLatency)
	minLatency := atomic.LoadInt64(&stats.MinLatency)
	maxLatency := atomic.LoadInt64(&stats.MaxLatency)

	avgLatency := float64(0)
	if total > 0 {
		avgLatency = float64(totalLatency) / float64(total) / 1000.0 // Convert to ms
	}

	tps := float64(total) / duration.Seconds()
	successRate := float64(0)
	if total > 0 {
		successRate = float64(successful) / float64(total) * 100
	}

	fmt.Println("\n" + repeat("=", 60))
	fmt.Println("LOAD TEST RESULTS")
	fmt.Println(repeat("=", 60))
	fmt.Printf("Duration:           %v\n", duration)
	fmt.Printf("Total Requests:     %d\n", total)
	fmt.Printf("Successful:         %d (%.2f%%)\n", successful, successRate)
	fmt.Printf("Failed:             %d\n", failed)
	fmt.Printf("Rate Limited:       %d\n", rateLimited)
	fmt.Printf("Timeouts:           %d\n", timeouts)
	fmt.Printf("\nThroughput:         %.2f req/s\n", tps)
	fmt.Printf("Avg Latency:        %.2f ms\n", avgLatency)
	if minLatency != int64(^uint64(0)>>1) {
		fmt.Printf("Min Latency:        %.2f ms\n", float64(minLatency)/1000.0)
	}
	fmt.Printf("Max Latency:        %.2f ms\n", float64(maxLatency)/1000.0)
	fmt.Println(repeat("=", 60))

	// Check SLA (<50ms target)
	if avgLatency > 50 {
		fmt.Printf("\n⚠️  WARNING: Average latency (%.2f ms) exceeds SLA target (50 ms)\n", avgLatency)
	} else {
		fmt.Printf("\n✅ Average latency (%.2f ms) meets SLA target (50 ms)\n", avgLatency)
	}
}

// Helper function to repeat strings
func repeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
