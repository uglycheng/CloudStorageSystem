cd ./test
rm -r user1;
rm -r user2;
mkdir user1;
mkdir user2;
cp ./resources/1.jpg ./user1/;
cd ../;
go run cmd/SurfstoreClientExec/main.go localhost:8081 ./test/user1 4096;
cp ./test/resources/2.jpg ./test/user1/1.jpg;
go run cmd/SurfstoreClientExec/main.go localhost:8081 ./test/user1 4096;
cp ./test/resources/2.jpg ./test/user1/1.jpg;
go run cmd/SurfstoreClientExec/main.go localhost:8081 ./test/user1 4096;
cp ./test/resources/2.jpg ./test/user1/1.jpg;
go run cmd/SurfstoreClientExec/main.go localhost:8081 ./test/user1 4096;
cp ./test/resources/2.jpg ./test/user1/1.jpg;
go run cmd/SurfstoreClientExec/main.go localhost:8081 ./test/user1 4096;
cp ./test/resources/3.jpg ./test/user1/1.jpg;
go run cmd/SurfstoreClientExec/main.go localhost:8081 ./test/user1 4096;

go run cmd/SurfstoreClientExec/main.go localhost:8081 ./test/user2 4096;
cp ./test/resources/4.jpg ./test/user2/1.jpg;
go run cmd/SurfstoreClientExec/main.go localhost:8081 ./test/user2 4096;

go run cmd/SurfstoreClientExec/main.go localhost:8081 ./test/user1 4096;


