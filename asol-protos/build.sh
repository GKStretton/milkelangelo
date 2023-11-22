NANOPB_VERSION=0.4.6
NANOPB_ROOT_DIR=/opt/nanopb-$NANOPB_VERSION-linux-x86

PROTO_PATH=./proto
C_OUT_DIR=../c
GO_OUT_DIR=../go
TYPESCRIPT_OUT_DIR=../typescript
PYTHON_OUT_DIR=../python
GO_MODULE=github.com/gkstretton/asol-protos

cd $PROTO_PATH

mkdir -p $C_OUT_DIR
mkdir -p $GO_OUT_DIR
mkdir -p $PYTHON_OUT_DIR
mkdir -p $TYPESCRIPT_OUT_DIR

###### C  ######

# Copy nanopb library files 
mkdir -p $C_OUT_DIR/nanopb
cp $NANOPB_ROOT_DIR/pb* $C_OUT_DIR/nanopb

# Build protos for nanopb
mkdir -p $C_OUT_DIR/machinepb
$NANOPB_ROOT_DIR/generator-bin/protoc \
	--nanopb_out $C_OUT_DIR/machinepb \
	machine.proto

###### GO ######

# Build protos for go
mkdir -p $GO_OUT_DIR/machinepb
protoc --go_out $GO_OUT_DIR/machinepb --go_opt=module=$GO_MODULE/go/machinepb ./machine.proto

cd ..

go mod init $GO_MODULE
go mod tidy

cd $PROTO_PATH

###### PYTHON ######
mkdir -p $PYTHON_OUT_DIR/machinepb
protoc --python_betterproto_out $PYTHON_OUT_DIR/machinepb ./machine.proto


###### TYPESCRIPT ######
mkdir -p $TYPESCRIPT_OUT_DIR/machinepb

# Hard to do JSON conversions with this generator
# protoc --js_out=import_style=commonjs,binary:$TYPESCRIPT_OUT_DIR/machinepb ./machine.proto 
# protoc --ts_out $TYPESCRIPT_OUT_DIR/machinepb ./machine.proto

# snakeToCamel: json keys are snake, code keys camel, consistent with Go default
protoc \
	--ts_proto_opt=esModuleInterop=true \
	--ts_proto_opt=snakeToCamel=keys \
	--ts_proto_out $TYPESCRIPT_OUT_DIR/machinepb ./machine.proto

