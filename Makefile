GO=go
GOFLAGS = -v -buildmode=pie
BIN_OUT=$(GOBIN)/goCarousel
BIN_OUT_UTIL=$(GOBIN)/goUnixStyle
MAIN=cmd/*.go
# Packagers only
PKG_NAME=go-carousel
PKG_VERSION=1.0.0
PKG_REVISION=1
PKG_ARCH=amd64
PKG_FULLNAME=${PKG_NAME}_${PKG_VERSION}-${PKG_REVISION}_${PKG_ARCH}
PKG_BUILD_DIR=${HOME}/Develop/Distrib/Build/${PKG_NAME}
PKG_PPA_DIR=${HOME}/Develop/Distrib/PPA

.PHONY: clean build

build:
	$(GO) build -tags logx $(GOFLAGS) -o ${BIN_OUT} ${MAIN}

buildwin:
	$(GO) build -tags logx $(GOFLAGS) -o ${BIN_OUT}.exe ${MAIN}

release:
	$(GO) build $(GOFLAGS) -o ${BIN_OUT} ${MAIN}

clean:
	go clean

util:
	$(GO) build $(GOFLAGS) -o ${BIN_OUT_UTIL} cmd/util/*go

run:
	go run -race  $MAIN

lint: 
	@gofmt -l . | grep ".*\.go"

test:
	go test tests/*test.go	

debian:
	rm -fR ${PKG_BUILD_DIR}
	mkdir -p ${PKG_BUILD_DIR}/DEBIAN
	ln -s ${PKG_BUILD_DIR}/DEBIAN ${PKG_BUILD_DIR}/debian
	cp -R distrib/DEBIAN/* ${PKG_BUILD_DIR}/DEBIAN
	mkdir -p ${PKG_BUILD_DIR}/usr/bin
	mkdir -p ${PKG_BUILD_DIR}/usr/share/doc/${PKG_NAME}/assets
	mkdir -p ${PKG_BUILD_DIR}/usr/share/man/man1
	gzip -n -9 -c distrib/manpages/man1/goCarousel.1 > ${PKG_BUILD_DIR}/usr/share/man/man1/goCarousel.1.gz
	gzip -n -9 -c distrib/manpages/man1/goUnixStyle.1 > ${PKG_BUILD_DIR}/usr/share/man/man1/goUnixStyle.1.gz
	mkdir -p ${PKG_BUILD_DIR}/usr/share/man/man5
	gzip -n -9 -c distrib/manpages/man5/goCarousel.5 > ${PKG_BUILD_DIR}/usr/share/man/man1/goCarousel.5.gz
	strip --strip-unneeded ${BIN_OUT}
	cp ${BIN_OUT} ${PKG_BUILD_DIR}/usr/bin
	strip --strip-unneeded ${BIN_OUT_UTIL}
	cp ${BIN_OUT_UTIL} ${PKG_BUILD_DIR}/usr/bin
	cp distrib/DEBIAN/copyright ${PKG_BUILD_DIR}/usr/share/doc/${PKG_NAME}
	cp docs/README.md ${PKG_BUILD_DIR}/usr/share/doc/${PKG_NAME}
	cp docs/assets/* ${PKG_BUILD_DIR}/usr/share/doc/${PKG_NAME}/assets
	gzip -n -9 -c distrib/DEBIAN/changelog > ${PKG_BUILD_DIR}/usr/share/doc/${PKG_NAME}/changelog.gz
	(cd ${PKG_BUILD_DIR} && dpkg-deb --root-owner-group -b ./ ${PKG_FULLNAME}.deb)
	#(cd ${PKG_BUILD_DIR} && fakeroot /usr/bin/dpkg-buildpackage --build=binary -us -uc -b ./ ${PKG_FULLNAME})
	#@mv /tmp/${PKG_FULLNAME}.deb ${DEST_REPOSITORY}
