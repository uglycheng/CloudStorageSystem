cd ./test
rm -r user1;
rm -r user2;
mkdir user1;
mkdir user2;
cp ./resources/empty.txt ./user1/;
cd ../;
go run cmd/SurfstoreClientExec/main.go localhost:8081 ./test/user1 4096;
go run cmd/SurfstoreClientExec/main.go localhost:8081 ./test/user2 4096;

