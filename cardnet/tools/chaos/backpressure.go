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
	concurrency = flag.Int("concurrency", 100, "Number of concurrent workers")
	duration    = flag.Duration("duration", 60*time.Second, "Test duration")
	rate        = flag.Int("rate", 1000, "Requests per second per worker (high rate to test backpressure)")
)

// Backpressure test: Send requests at a rate higher than system capacity
// to verify the system handles overload gracefully
func main() {
	flag.Parse()

	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := grpcapi.NewAuthGatewayClient(conn)

	var (
		totalReqs      int64
		successfulReqs int64
		failedReqs     int64
		rateLimited    int64
		timeoutReqs    int64
		rejectedReqs   int64
		totalLatency   int64
		maxLatency     int64
	)

	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), *duration)
	defer cancel()

	startTime := time.Now()

	// Start workers at high rate
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
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
					reqID := fmt.Sprintf("bp-req-%d-%d", workerID, requestID)

					reqCtx, cancelReq := context.WithTimeout(ctx, 2*time.Second)
					start := time.Now()

					resp, err := client.Authorize(reqCtx, &grpcapi.AuthRequest{
						RequestId:  reqID,
						MerchantId: fmt.Sprintf("merchant_%03d", workerID%10),
						CardToken:  fmt.Sprintf("tok_visa_%04d", requestID%10000),
						Amount:     1000 + int64(requestID%10000),
						Currency:   "INR",
						Mcc:        "5411",
					})

					latency := time.Since(start).Milliseconds()
					cancelReq()

					atomic.AddInt64(&totalReqs, 1)
					atomic.AddInt64(&totalLatency, latency)

					// Update max latency
					for {
						oldMax := atomic.LoadInt64(&maxLatency)
						if latency <= oldMax {
							break
						}
						if atomic.CompareAndSwapInt64(&maxLatency, oldMax, latency) {
							break
						}
					}

					if err != nil {
						atomic.AddInt64(&failedReqs, 1)
						if reqCtx.Err() == context.DeadlineExceeded {
							atomic.AddInt64(&timeoutReqs, 1)
						}
					} else {
						atomic.AddInt64(&successfulReqs, 1)
						switch resp.Status {
						case "RATE_LIMITED":
							atomic.AddInt64(&rateLimited, 1)
						case "REJECTED", "DECLINED":
							atomic.AddInt64(&rejectedReqs, 1)
						}
					}
				}
			}
		}(i)
	}

	// Print periodic stats
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				total := atomic.LoadInt64(&totalReqs)
				successful := atomic.LoadInt64(&successfulReqs)
				failed := atomic.LoadInt64(&failedReqs)
				rlCount := atomic.LoadInt64(&rateLimited)
				elapsed := time.Since(startTime).Seconds()

				if elapsed > 0 {
					tps := float64(total) / elapsed
					successRate := float64(0)
					if total > 0 {
						successRate = float64(successful) / float64(total) * 100
					}

					fmt.Printf("[%s] Total: %d | Success: %d (%.1f%%) | Failed: %d | Rate Limited: %d | TPS: %.1f\n",
						time.Now().Format("15:04:05"),
						total, successful, successRate, failed, rlCount, tps)
				}
			}
		}
	}()

	wg.Wait()
	testDuration := time.Since(startTime)

	// Final results
	total := atomic.LoadInt64(&totalReqs)
	successful := atomic.LoadInt64(&successfulReqs)
	failed := atomic.LoadInt64(&failedReqs)
	rateLimitedCount := atomic.LoadInt64(&rateLimited)
	timeouts := atomic.LoadInt64(&timeoutReqs)
	rejected := atomic.LoadInt64(&rejectedReqs)
	totalLat := atomic.LoadInt64(&totalLatency)
	maxLat := atomic.LoadInt64(&maxLatency)

	avgLatency := float64(0)
	if total > 0 {
		avgLatency = float64(totalLat) / float64(total)
	}

	tps := float64(total) / testDuration.Seconds()
	successRate := float64(0)
	if total > 0 {
		successRate = float64(successful) / float64(total) * 100
	}

	fmt.Println("\n" + repeat("=", 70))
	fmt.Println("BACKPRESSURE TEST RESULTS")
	fmt.Println(repeat("=", 70))
	fmt.Printf("Duration:           %v\n", testDuration)
	fmt.Printf("Concurrency:        %d workers\n", *concurrency)
	fmt.Printf("Request Rate:        %d req/s per worker\n", *rate)
	fmt.Printf("Total Requests:     %d\n", total)
	fmt.Printf("Successful:         %d (%.2f%%)\n", successful, successRate)
	fmt.Printf("Failed:             %d\n", failed)
	fmt.Printf("Timeouts:           %d\n", timeouts)
	fmt.Printf("Rate Limited:       %d\n", rateLimitedCount)
	fmt.Printf("Rejected:           %d\n", rejected)
	fmt.Printf("\nThroughput:         %.2f req/s\n", tps)
	fmt.Printf("Avg Latency:        %.2f ms\n", avgLatency)
	fmt.Printf("Max Latency:        %.2f ms\n", float64(maxLat))
	fmt.Println(repeat("=", 70))

	// Backpressure indicators
	if successRate < 95 {
		fmt.Printf("\n⚠️  WARNING: Low success rate (%.2f%%) indicates system under stress\n", successRate)
	}

	if avgLatency > 100 {
		fmt.Printf("\n⚠️  WARNING: High average latency (%.2f ms) indicates backpressure\n", avgLatency)
	}

	if float64(rateLimitedCount)/float64(total) > 0.1 {
		fmt.Printf("\n✅ Rate limiting is working (%.1f%% rate limited)\n", float64(rateLimitedCount)/float64(total)*100)
	}

	if successRate >= 95 && avgLatency < 100 {
		fmt.Printf("\n✅ System handled backpressure gracefully\n")
	}
}

func repeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
