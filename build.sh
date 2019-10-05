cd ./app/master/
go build -o ../../bin/master.bin
cd ../store/
go build -o ../../bin/store.bin
cd ../zoo/
go build -o ../../bin/zoo.bin
cp -r ./admin/ ../../bin/admin/