Clear-Host
wt -w 0 -d .\order\ go run .\cmd\order\order.go nt
wt -w 0 -d .\orchestrator\ go run .\cmd\orchestrator\orchestrator.go nt
wt -w 0 -d .\user\ go run .\cmd\user\user.go nt
wt -w 0 -d .\book\ go run .\cmd\book\book.go nt
wt -w 0 -d .\payment\ go run .\cmd\payment\payment.go nt