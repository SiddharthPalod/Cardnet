## Run Services
- go run cmd/auth-gateway/main.go
- go run cmd/rate-limiter/main.go 

## Build Project
- go build ./...
- go test ./...

## Proto File Auto Gen
-  protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. internal/api/grpc/<filename.proto>

## GRPCURL For testing
- for ($i = 1; $i -le 50; $i++) { grpcurl -plaintext -d '{\"request_id\":\"req-011\",\"merchant_id\":\"merchant_124\",\"card_token\":\"tok_visa_133\",\"amount\":1000,\"currency\":\"INR\"}' localhost:50051 auth.AuthGateway/Authorize }
- Get-Content risk.json | grpcurl --% -plaintext -emit-defaults -d @ localhost:60052 risk.RiskEngine.Evaluate

## Docker Commands
- docker-compose up --build
- docker-compose down

## Frontend 
- npm run dev -- -p 3001 
- npm run dev