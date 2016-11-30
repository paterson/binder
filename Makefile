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
