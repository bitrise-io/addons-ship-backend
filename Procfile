release: go run db/main.go -dir db up && go run tools/migration/main.go
web: addons-ship-backend -port=$PORT
worker: export WORKER='true' && addons-ship-backend
