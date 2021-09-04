Dependencies:

    - github.com/jmoiron/sqlx

Example for connection:

	connString := fmt.Sprintf("user=%s dbname=%s sslmode=disable password=%s", config.Config.Postgres.Username, config.Config.Postgres.Database, config.Config.Postgres.Password)
	db := dbx.MustConnect()
	db.QuoteIdentifier = pq.QuoteIdentifier


Example for current methods:

    db.Create(&model)
    db.Get(&model, id)
    db.Select(&[]model, condition)
    
They all use the Tabler interface, so if the name of the table is not the plural of the struct name, fullfil the interface.