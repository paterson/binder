build:
	$(MAKE) -C authservice install
	$(MAKE) -C directoryserver install
	$(MAKE) -C fileserver install
	$(MAKE) -C client install

run:
	$(MAKE) -C authservice run
	$(MAKE) -C directoryserver run
	$(MAKE) -C fileserver run
	fileserver 3003

kill:
	lsof -ti:3000 | xargs kill 
	lsof -ti:3001 | xargs kill
	lsof -ti:3002 | xargs kill