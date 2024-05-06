invest:
	go run . -USD 100

investbig:
	go run . -USD 1000000000.50

testv:
	go test . -v	

test:
	go test .
	
tidy:
	go mod tidy

.PHONY: invest,investbig,test,testv,tidy