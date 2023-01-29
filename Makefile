build_task6:
	go build -o ./bin/task6 ./task6/task6.go
	go build -o ./bin/task6_2 ./task6_2/task6_2.go
run_task6:
	go run ./task6/task6.go $(var)
