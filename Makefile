# Paths to tools needed in dependencies
GO := $(shell which go)
TINYGO := $(shell which tinygo)

# Paths to locations, etc
BUILD_DIR := "build"
PICO_CMD_DIR := $(filter-out cmd/pico/README.md, $(wildcard cmd/pico/*))
PICO_BUILD_FLAGS := -target pico
RPI_CMD_DIR := $(filter-out cmd/rpi/README.md, $(wildcard cmd/rpi/*))
RPI_BUILD_FLAGS := -tags rpi

# Targets
all: pico rpi

pico: $(PICO_CMD_DIR)

rpi: $(RPI_CMD_DIR)

$(PICO_CMD_DIR): dependencies mkdir FORCE
	@echo Build pico $(notdir $@)
	@${TINYGO} build ${PICO_BUILD_FLAGS} -o ${BUILD_DIR}/$(notdir $@).uf2 ./$@

$(RPI_CMD_DIR): dependencies mkdir FORCE
	@echo Build rpi $(notdir $@)
	@${GO} build ${RPI_BUILD_FLAGS} -o ${BUILD_DIR}/$(notdir $@).uf2 ./$@

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
