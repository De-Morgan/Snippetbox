
# mysql:
# 	docker run --name mysql -p 3306:3306 -e MYSQL_DATABASE=snippetbox -e MYSQL_USER=user -e MYSQL_PASSWORD=secret  -e MYSQL_ROOT_PASSWORD=rootsecret -d mysql:latest

migrateup:
	 migrate -path pkg/migration -database "mysql://root:secret@/snippets?parseTime=true&multiStatements=true" --verbose up