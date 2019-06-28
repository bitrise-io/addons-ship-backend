release: go run db/main.go -dir db up
web: addons-ship-backend -port=$PORT
worker: export WORKER='true' && addons-ship-backend