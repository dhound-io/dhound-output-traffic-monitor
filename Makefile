VERSION=1.1.3
PACKAGE_NAME=dhound-output-traffic-monitor

.PHONY: default
default: compile

PROJECT=dhound-output-traffic-monitor

.PHONY: compile
compile: $(PROJECT)

dhound-output-traffic-monitor:
	go build --ldflags '-extldflags "-static"' -o $(PACKAGE_NAME)

.PHONY: clean
clean:
	#-rm $(PROJECT)
	-rm -rf build
	-rm -rf publish/empty
	#-rm -rf publish/packages

.PHONY: buildempty
buildempty:
	-mkdir publish/empty
	#-mkdir publish/packages


.PHONY: rpm deb
deb: BEFORE_INSTALL=publish/pkg/ubuntu/before-install.sh
deb: AFTER_INSTALL=publish/pkg/ubuntu/after-install.sh
deb: BEFORE_REMOVE=publish/pkg/ubuntu/before-remove.sh

rpm: AFTER_INSTALL=publish/pkg/centos/after-install.sh
rpm: BEFORE_INSTALL=publish/pkg/centos/before-install.sh
rpm: BEFORE_REMOVE=publish/pkg/centos/before-remove.sh

rpm deb: PREFIX=/opt/dhound-output-traffic-monitor
rpm deb: clean compile buildempty
	fpm  -f -s dir -t $@ -n $(PACKAGE_NAME) -v $(VERSION) \
		--deb-no-default-config-files \
		--architecture $(ARCHITECTURE) \
		--replaces $(PACKAGE_NAME) \
		--description "$(PACKAGE_NAME) tool to track and log output traffic" \
		--after-install $(AFTER_INSTALL) \
		--before-install $(BEFORE_INSTALL) \
		--before-remove $(BEFORE_REMOVE) \
		./$(PACKAGE_NAME)=$(PREFIX)/bin/ \
		./publish/etc=/ \
		./publish/empty/=/var/lib/$(PACKAGE_NAME)/ \
		./publish/empty/=/var/log/dhound/