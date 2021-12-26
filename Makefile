BINARY=engine
engine:
	go build -o ${BINARY} .

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

docker:
	docker build -t stripe-payment .

run:
	docker-compose up --build -d

stop:
	docker-compose down
