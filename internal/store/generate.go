package store

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate ./schema
//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/upsert --template ./templates ./schema
