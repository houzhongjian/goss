cd ./app/api/
go build -o ../../bin/api.bin
cd ../store/
go build -o ../../bin/store.bin
cd ../master/
go build -o ../../bin/master.bin
cp -r ./admin/ ../../bin/admin/