# Paths to tools needed in dependencies
GO := $(shell which go)
TINYGO := $(shell which tinygo)

# Paths to locations, etc
BUILD_DIR := "build"
CMD_DIR := $(filter-out cmd/_old, $(wildcard cmd/*))
PICO_BUILD_FLAGS := -target pico -tags "rp2040 debug"

# Targets
all: $(CMD_DIR)

$(CMD_DIR): dependencies mkdir FORCE
	@echo Build $(notdir $@)
	${TINYGO} build ${PICO_BUILD_FLAGS} -o ${BUILD_DIR}/$(notdir $@).uf2 ./$@

FORCE:

dependencies:
ifeq (,${GO})
        $(error "Missing go binary")
endif
ifeq (,${TINYGO})
        $(error "Missing tinygo binary")
endif

mkdir:
	@echo Mkdir ${BUILD_DIR}
	@install -d ${BUILD_DIR}

clean:
	@echo Clean
	@rm -fr $(BUILD_DIR)
