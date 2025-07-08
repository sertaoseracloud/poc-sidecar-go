package main

import (
    "bytes"
    "context"
    "encoding/json"
    "io"
    "log"
    "net/http"
    "os"
    "time"

    "google.golang.org/grpc"
    pb "sidecar/sidecar/proto"
)

func httpRequest(provider string) {
    b, _ := json.Marshal(map[string]string{"provider": provider})
    resp, err := http.Post("http://sidecar-go:3000/auth", "application/json", bytes.NewReader(b))
    if err != nil {
        log.Println("http error:", err)
        return
    }
    defer resp.Body.Close()
    body, _ := io.ReadAll(resp.Body)
    log.Println("HTTP credentials", string(body))
}

func grpcRequest(provider string) {
    conn, err := grpc.Dial("sidecar-go:50051", grpc.WithInsecure())
    if err != nil {
        log.Println("grpc dial:", err)
        return
    }
    defer conn.Close()
    client := pb.NewAuthServiceClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    res, err := client.RequestAuth(ctx, &pb.AuthRequest{Provider: provider})
    if err != nil {
        log.Println("grpc call:", err)
        return
    }
    log.Println("gRPC credentials", res.Json)
}

func main() {
    provider := os.Getenv("PROVIDER")
    if provider == "" {
        provider = "aws"
    }
    httpRequest(provider)
    grpcRequest(provider)
}
