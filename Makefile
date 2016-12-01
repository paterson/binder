build:
	$(MAKE) -C authservice install
	$(MAKE) -C client install
	$(MAKE) -C directoryservice install
	$(MAKE) -C fileserver install

run:
	$(MAKE) -C authservice run
	$(MAKE) -C client run
	$(MAKE) -C directoryservice run
	$(MAKE) -C fileserver run

kill:
	lsof -i tcp:3000 | awk 'NR!=1 {print $2}' | xargs kill 
	lsof -i tcp:3001 | awk 'NR!=1 {print $2}' | xargs kill 
	lsof -i tcp:3002 | awk 'NR!=1 {print $2}' | xargs kill