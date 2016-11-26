build:
	$(MAKE) -C authservice install
	$(MAKE) -C client install
	$(MAKE) -C clientproxy install
	$(MAKE) -C directoryservice install
	$(MAKE) -C fileserver install

run:
	$(MAKE) -C authservice run
	$(MAKE) -C client run
	$(MAKE) -C clientproxy run
	$(MAKE) -C directoryservice run
	$(MAKE) -C fileserver run
