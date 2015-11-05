include ../includes.mk

MONITOR_IMAGES=alertmanager daemon

ALERTMANAGER_IMAGE = $(IMAGE_PREFIX)monitor-alertmanager:$(BUILD_TAG)
ALERTMANAGER_DEV_IMAGE = $(REGISTRY)$(ALERTMANAGER_IMAGE)
DAEMON_IMAGE = $(IMAGE_PREFIX)monitor-daemon:$(BUILD_TAG)
DAEMON_DEV_IMAGE = $(REGISTRY)$(DAEMON_IMAGE)

build: check-docker
	@# Build base as normal
	$(foreach I, $(MONITOR_IMAGES), \
		docker build -t $(IMAGE_PREFIX)monitor-$(I):$(BUILD_TAG) $(I)/ || exit 1; \
	)

clean: check-docker check-registry
	$(foreach I, $(MONITOR_IMAGES), \
		docker rmi $(IMAGE_PREFIX)monitor-$(I):$(BUILD_TAG) ; \
		docker rmi $(REGISTRY)/$(IMAGE_PREFIX)monitor-$(I):$(BUILD_TAG) ; \
	)

full-clean: check-docker check-registry
	$(foreach I, $(MONITOR_IMAGES), \
		docker images -q $(IMAGE_PREFIX)monitor-$(I) | xargs docker rmi -f ; \
		docker images -q $(REGISTRY)/$(IMAGE_PREFIX)monitor-$(I) | xargs docker rmi -f ; \
	)

install: check-deisctl
	deisctl install monitor-alertmanager
	deisctl install monitor-daemon # @fixme: both should be scalable
	#deisctl scale monitor-daemon=1

uninstall: check-deisctl
	#deisctl scale smonitor-daemon=0
	deisctl uninstall monitor-alertmanager
	deisctl uninstall monitor-daemon

start: check-deisctl
	deisctl start monitor-alertmanager
	deisctl start monitor-daemon

stop: check-deisctl
	deisctl stop monitor-alertmanager
	deisctl stop monitor-daemon

restart: stop start

run: install start

dev-release: push set-image

push: check-registry
	docker tag -f $(ALERTMANAGER_IMAGE) $(ALERTMANAGER_DEV_IMAGE)
	docker push $(ALERTMANAGER_DEV_IMAGE)
	docker tag -f $(DAEMON_IMAGE) $(DAEMON_DEV_IMAGE)
	docker push $(DAEMON_DEV_IMAGE)

set-image: check-deisctl
	deisctl config store-admin set image=$(ALERTMANAGER_DEV_IMAGE)
	deisctl config store-daemon set image=$(DAEMON_DEV_IMAGE)

release:
#	docker push $(ALERTMANAGER_IMAGE) # fixme: uncomment
#	docker push $(DAEMON_IMAGE)

deploy: build dev-release restart

#$test: test-style test-unit test-functional # fixme: no tests

test-functional:
	@echo no functional tests # fixme: yet
#	@$(MAKE) -C ../tests/ test-etcd
#	GOPATH=`cd ../tests/ && godep path`:$(GOPATH) go test -v ./tests/...

test-style:
	@echo no style tests

test-unit:
	@echo no unit tests

