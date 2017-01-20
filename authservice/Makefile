install:
	go install

run:
	PORT=3000 authservice &

clean:
	rm -f auth.db

docker:
	docker build -t paterson/binder-authservice .
