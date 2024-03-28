.PHONY: start-redis-stack

start-redis-stack:
	docker run --rm -p 6379:6379 -p 8001:8001 redis/redis-stack
