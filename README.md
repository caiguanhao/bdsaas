# bdsaas

```go
client := bdsaas.NewClient("XXXXXXXXXX")

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

seats, err := client.GetSeats(ctx)

sessionId, err := client.Call(ctx, "13812345678", "13912345678", "127.0.0.1", "TEST")

records, err := client.Query(ctx, sessionId)
```
