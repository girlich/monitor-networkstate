B=networkstate
ARCHS=amd64 arm
ADD_arm=GOARM=5
BINS=$(addprefix $(B)_,$(ARCHS))

all: $(BINS)

$(BINS) : $(B)_% : $(B).go
	GOARCH=$* $(ADD_$*) go build -o $@ $<

clean:
	rm -f $(BINS)

