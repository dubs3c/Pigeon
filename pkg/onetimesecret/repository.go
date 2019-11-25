package onetimesecret

// Secret model
type Secret struct {
	token    string
	message  string
	password string
	expire   string
	maxviews int
	views    int
}

// GetSecretByToken : Fetches a secret by token
func (db *DB) GetSecretByToken(token string) (*Secret, error) {
	row := db.QueryRow("SELECT * FROM Secrets WHERE token=?", token)

	sec := new(Secret)
	if err := row.Scan(&sec.token, &sec.message, &sec.password, &sec.expire, &sec.maxviews, &sec.views); err != nil {
		return nil, err
	}

	return sec, nil
}

// GetSecretByTokenAndPassword : Fetch a secret by token and its password
func (db *DB) GetSecretByTokenAndPassword(token string, password string) (*Secret, error) {
	row := db.QueryRow("SELECT * FROM Secrets WHERE token=? AND password=?", token, password)

	sec := new(Secret)
	if err := row.Scan(&sec.token, &sec.message, &sec.password, &sec.expire, &sec.maxviews, &sec.views); err != nil {
		return nil, err
	}

	return sec, nil
}

// CreateSecret : Create a one-time secret
func (db *DB) CreateSecret(secret Secret) error {
	_, err := db.Exec("INSERT INTO Secrets(token, secret, password, expire, maxviews, views) VALUES(?, ?, ?, datetime('now', ?, 'localtime'), ?, ?)", secret.token, secret.message, secret.password, secret.expire, secret.maxviews, secret.views)
	return err
}

// DeleteSecret : Delete a given secret
func (db *DB) DeleteSecret(token string) error {
	_, err := db.Exec("DELETE FROM Secrets WHERE token=?", token)
	return err
}

// IncrementViews : Increment amount of views for a given secret
func (db *DB) IncrementViews(token string) error {
	tx, err := db.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	stmt, err := tx.Prepare("UPDATE Secrets SET views=views + 1 WHERE token=?")

	if err != nil {
		return err
	}

	defer stmt.Close() // danger!

	_, err = stmt.Exec(token)

	if err != nil {
		return err
	}

	err = tx.Commit()

	return err
}
