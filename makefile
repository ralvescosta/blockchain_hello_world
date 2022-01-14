start blockchain:
	DB_PATH=./temp/blocks/ go run .

add block:
	DB_PATH=./temp/blocks/ go run . add -block "YOUR BLOCK DATA HERE"