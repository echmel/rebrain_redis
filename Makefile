build:
	go mod tidy
	go mod vendor
	#go build -o ./bin/task6 ./task6/task6.go
	#go build -o ./bin/task6_2 ./task6_2/task6_2.go
	#go build -o ./bin/task7_1 ./task7_1/task7_1.go
	#go build -o ./bin/task7_2 ./task7_2/task7_2.go
	#go build -o ./bin/task8 ./task8/task8.go
	#go build -o ./bin/task9 ./task9/task9.go
	#go build -o ./bin/task9_2 ./task9_2/task9_2.go
	go build -o ./bin/task10 ./task10/task10.go
run_task6:
	go run ./task6/task6.go $(var)

go_install:
	sudo snap install go --classic

init_redis:
	#redis-cli SET cmd-shuf "/home/user/rebrain_redis/bin/task9"
	#redis-cli SET cmd-7777 "/home/user/rebrain_redis/bin/task9_2"
	redis-cli SET cmd-wait "/home/user/rebrain_redis/bin/task10"

#redis-cli SET cmd-lenq "/home/user/rebrain_redis/bin/task6"
#redis-cli SET cmd-ldeq "/home/user/rebrain_redis/bin/task6_2"
