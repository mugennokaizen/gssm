mup:
	migrate -source file://db/migrations -database postgres://user:password@localhost:5601/gssn?sslmode=disable up 1
mdown:
	migrate -source file://db/migrations -database postgres://user:password@localhost:5601/gssn?sslmode=disable down 1