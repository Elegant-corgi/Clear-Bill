# note: invoke sub-project targets from the repository root

TOPTARGETS := all build clean dist test mod lint wire swagger
SUBDIRS := mgr scripts tools
VERSION := $(shell cat VERSION)

$(TOPTARGETS): $(SUBDIRS)

$(SUBDIRS):
	$(MAKE) -C $@ $(MAKECMDGOALS)

.PHONY: $(TOPTARGETS) $(SUBDIRS) release server-start webui-dev

server-start:
	$(MAKE) -C mgr server-start

webui-dev:
	$(MAKE) -C mgr webui-dev

release:
	@mkdir -p release/clear-bill-$(VERSION)/server/website
	@mkdir -p release/clear-bill-$(VERSION)/server/configs
	@mkdir -p release/clear-bill-$(VERSION)/scripts
	@mkdir -p release/clear-bill-$(VERSION)/tools
	@cp -f VERSION release/clear-bill-$(VERSION)/
	@cp -f README.md build.sh deploy.sh release/clear-bill-$(VERSION)/
	@if [ -d mgr/server/dist ]; then cp -rf mgr/server/dist/* release/clear-bill-$(VERSION)/server/; fi
	@if [ -d mgr/webui/dist ]; then cp -rf mgr/webui/dist/* release/clear-bill-$(VERSION)/server/website/; fi
	@if [ -d mgr/server/configs ]; then cp -rf mgr/server/configs/* release/clear-bill-$(VERSION)/server/configs/; fi
	@if [ -d scripts/dist ]; then cp -rf scripts/dist/* release/clear-bill-$(VERSION)/scripts/; fi
	@if [ -d tools/dist ]; then cp -rf tools/dist/* release/clear-bill-$(VERSION)/tools/; fi
	@echo "release package created at release/clear-bill-$(VERSION)"
