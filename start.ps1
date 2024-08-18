Clear-Host
wt -w 0 -d .\orchestrator\ go run .\cmd\orchestrator\orchestrator.go nt
wt -w 1 -d .\order\ go run .\cmd\order\order.go nt
wt -w 1 -d .\user\ go run .\cmd\user\user.go nt
wt -w 1 -d .\book\ go run .\cmd\book\book.go nt
wt -w 1 -d .\payment\ go run .\cmd\payment\payment.go nt