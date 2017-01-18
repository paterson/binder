install:
	go install

run:
	authservice 3000 &

clean:
	rm -f auth.db

docker:
	docker build -t paterson/binder-authservice .
