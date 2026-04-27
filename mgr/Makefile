TOPTARGETS := all build clean dist test mod lint wire swagger
SUBDIRS := server webui

$(TOPTARGETS): $(SUBDIRS)

$(SUBDIRS):
	$(MAKE) -C $@ $(MAKECMDGOALS)

.PHONY: $(TOPTARGETS) $(SUBDIRS) server-start webui-dev

server-start:
	$(MAKE) -C server start

webui-dev:
	$(MAKE) -C webui dev
