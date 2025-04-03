package main

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
)

var existingClientError = errors.New("client already exists")

func createConnection(databaseName string) (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", databaseName)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func closeConnection(db *sql.DB) {
	db.Close()
}

func checkPrivilege(db *sql.DB, tgID int64) (prev string, err error) {
	err = db.QueryRow("select p.P_VALUE from \"Privileges\" p \ninner join Clients c on c.C_PREV = p.P_ID \nwhere c.C_TID = ?", tgID).Scan(&prev)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("no client with provided tid found")
		}
		return "", err
	}
	return prev, nil
}

func addNewClient(db *sql.DB, tgID int64, nickName string, firstName string, lastName string) (err error) {
	exists, err := clientExists(db, tgID)
	if err != nil {
		return err
	}

	if exists {
		return existingClientError
	}

	stmt, err := db.Prepare("insert into Clients (C_TID, C_NICKNAME, C_FIRSTNAME, C_LASTNAME) values (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(tgID, nickName, firstName, lastName)
	if err != nil {
		return err
	}
	return nil
}

func clientExists(db *sql.DB, tgID int64) (exists bool, err error) {
	var id int
	err = db.QueryRow("select C_ID from Clients where C_TID = ?", tgID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func authorizeClient(db *sql.DB, tgID int64) (err error) {
	exists, err := clientExists(db, tgID)
	if err != nil {
		return err
	}
	if !exists {
		err = addNewClient(db, tgID, "", "", "")
		if err != nil {
			return err
		}
	}

	stmt, err := db.Prepare("update Clients set C_PREV = 1 where C_TID = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(tgID)
	if err != nil {
		return err
	}
	return nil
}

func deauthorizeClient(db *sql.DB, tgID int64) (err error) {
	exists, err := clientExists(db, tgID)
	if err != nil {
		return err
	}
	if !exists {
		err = addNewClient(db, tgID, "", "", "")
		if err != nil {
			return err
		}
	}

	stmt, err := db.Prepare("update Clients set C_PREV = 0 where C_TID = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(tgID)
	if err != nil {
		return err
	}
	return nil
}

func setClientNames(db *sql.DB, tgID int64, nickName string, firstName string, lastName string) (err error) {
	var (
		dbNick      string
		dbFirstName string
		dbLastName  string
	)

	err = db.QueryRow("select C_NICKNAME, C_FIRSTNAME, C_LASTNAME from Clients where C_TID = ?", tgID).Scan(&dbNick, &dbFirstName, &dbLastName)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("no client found")
		}
		return err
	}
	if dbNick != nickName {
		stmt, err := db.Prepare("update Clients set C_NICKNAME = ? where C_TID = ?")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(nickName, tgID)
		if err != nil {
			return err
		}
	}

	if dbFirstName != firstName {
		stmt, err := db.Prepare("update Clients set C_FIRSTNAME = ? where C_TID = ?")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(firstName, tgID)
		if err != nil {
			return err
		}
	}

	if dbLastName != lastName {
		stmt, err := db.Prepare("update Clients set C_LASTNAME = ? where C_TID = ?")
		if err != nil {
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(lastName, tgID)
		if err != nil {
			return err
		}
	}
	return nil
}

func getClientId(db *sql.DB, tgID int64) (clientId int, err error) {
	err = db.QueryRow("select C_ID from Clients where C_TID = ?", tgID).Scan(&clientId)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, errors.New("client does not exist")
		}
		return -1, err
	}
	return clientId, nil
}

func writeLog(db *sql.DB, clientId int, command string, request string, response string, output string, result int) (err error) {
	stmt, err := db.Prepare("insert into Log (L_C_ID, L_COMMAND, L_REQUEST, L_RESPONSE, L_OUTPUT, L_RESULT) values (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(clientId, command, request, response, output, result)
	if err != nil {
		return err
	}
	return nil
}

func getAuthorizedClients(db *sql.DB) (clients []int64, err error) {
	rows, err := db.Query("select C_TID from Clients where C_PREV != 0")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var tid int64
		err = rows.Scan(&tid)
		if err != nil {
			return nil, err
		}
		clients = append(clients, tid)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return clients, nil
}
