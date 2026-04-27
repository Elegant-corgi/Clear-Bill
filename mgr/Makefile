TOPTARGETS := all build clean dist test
SUBDIRS := server webui

$(TOPTARGETS): $(SUBDIRS)

$(SUBDIRS):
	$(MAKE) -C $@ $(MAKECMDGOALS)

.PHONY: $(TOPTARGETS) $(SUBDIRS) server-start webui-dev

server-start:
	$(MAKE) -C server start

webui-dev:
	$(MAKE) -C webui dev
