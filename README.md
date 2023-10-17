# Case Study Back End Engineer

docker-compose.yaml -> untuk menjalanlan depedencynya


## Run Auth Service
1. cd auth-se 
2. go run main.go db:migrate up
3. go run main.go http


## Run Wallet Service
1. cd wallet-se  
2. go run main.go db:migrate up
3. go run main.go http


## Run Manage Service
1. cd manage-se  
2. go run main.go db:migrate up
3. go run main.go http


for detail explanation, read read.me inside each service