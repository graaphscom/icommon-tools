.PHONY: start-redis-stack

start-redis-stack:
	docker run -v icommon-tools-redis:/data --rm -p 6379:6379 -p 8001:8001 redis/redis-stack
