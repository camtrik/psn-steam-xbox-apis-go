# linux
# run:
# 	sudo service redis-server start && go run cmd/main.go

# mac 
run:
	brew services start redis || echo "Redis may already be running" && go run cmd/main.go
