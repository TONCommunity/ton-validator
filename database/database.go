package database

import (
	"database/sql"
	"fmt"
	"log"

	// Init
	_ "github.com/mattn/go-sqlite3"
	tonlib "github.com/mercuryoio/tonlib-go/v2"
)

type store struct {
	db *sql.DB
}

//NewClient init new connection to the database
func NewClient(pathDB string) (*store, error) {
	db, err := sql.Open("sqlite3", pathDB)
	if err != nil {
		return &store{}, err
	}
	return &store{db: db}, nil
}

//Node info
type Node struct {
	ID         int
	HostPort   string
	ServerPub  string
	ClientCert string
	Enabled    int
}

//Wallet info
type Wallet struct {
	ID       int
	FilePath string
	Addr     string
	Balance  int64
	Enabled  int
}

//Election info
type Election struct {
	ID              int
	ElectionID      int64
	StartAt         int64
	CloseAt         int64
	NextElectionsAt int64
}

//Participate record
type Participate struct {
	NodeID      int
	ElectionID  int64
	StakeAmount int64
	MaxFactor   string
}

//Key validator keys
type Key struct {
	Key        string
	ElectionID int64
	NodeID     int
	Type       string
}

//AddNode Add a node info to database
func (store *store) AddNode(hostPort, serverPub, clientCert string, walletID int) (int64, error) {
	stmt, err := store.db.Prepare("INSERT INTO nodes(host_port, server_pub, client_cert, wallet_id, enabled) values(?,?,?,?,?)")
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(hostPort, serverPub, clientCert, walletID, "1")
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

//DelNode Del node from database
func (store *store) DelNode(id int) error {
	query := fmt.Sprintf("delete from nodes where id = %d", id)
	_, err := store.db.Exec(query)
	if err != nil {
		fmt.Printf("Failed to delete node with ID: %d: %s\n", id, err)
		return err
	}
	return nil
}

//AddWallet Add wallet info to database
func (store *store) AddWallet(walletFile, walletAddr string) (int64, error) {
	stmt, err := store.db.Prepare("INSERT INTO wallets(wallet_file, wallet_addr, balance, enabled) values(?,?,?,?)")
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(walletFile, walletAddr, "0", "1")
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

//DelWallet Delete wallet by ID
func (store *store) DelWallet(id int) error {
	query := fmt.Sprintf("delete from wallets where id = %d", id)
	_, err := store.db.Exec(query)
	if err != nil {
		fmt.Printf("Failed to delete wallet with ID: %d: %s\n", id, err)
		return err
	}
	return nil
}

//GetWallets Get wallet info
func (store *store) GetWallets(enabled int) ([]Wallet, error) {
	var query string
	if enabled > 1 {
		query = "select id,wallet_file,wallet_addr,balance,enabled from wallets"
	} else {
		query = fmt.Sprintf("select id,wallet_file,wallet_addr,balance,enabled from wallets where enabled=%d", enabled)
	}
	rows, err := store.db.Query(query)
	if err != nil {
		//fmt.Println(err)
		return []Wallet{}, err
	}
	defer rows.Close()
	var wallets []Wallet
	for rows.Next() {
		var wallet Wallet
		err = rows.Scan(&wallet.ID, &wallet.FilePath, &wallet.Addr, &wallet.Balance, &wallet.Enabled)
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Println(walletFile, walletAddr)
		wallets = append(wallets, wallet)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
	}

	return wallets, nil
}

//UpdateWalletBalance update wallet balance by id
func (store *store) UpdateWalletBalance(walletID int, newBalance int64) error {
	stmt, err := store.db.Prepare("update wallets set balance=? where id=?")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(newBalance, walletID)
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}
	return nil
	//fmt.Println(affect)
}

//GetNodes Get info about nodes
func (store *store) GetNodes(walletID, enabled int) ([]Node, error) {
	var query string
	if enabled > 1 {
		query = fmt.Sprintf("select id,host_port,server_pub,client_cert,enabled from nodes where wallet_id=%d", walletID)
	} else {
		query = fmt.Sprintf("select id,host_port,server_pub,client_cert,enabled from nodes where enabled=%d and wallet_id=%d", enabled, walletID)
	}
	rows, err := store.db.Query(query)
	if err != nil {
		return []Node{}, err
	}
	defer rows.Close()
	var nodes []Node
	for rows.Next() {
		var node Node
		err = rows.Scan(&node.ID, &node.HostPort, &node.ServerPub, &node.ClientCert, &node.Enabled)
		if err != nil {
			return nodes, err
		}
		//fmt.Println(hostPort, serverPub, clientCert)
		nodes = append(nodes, node)
	}
	err = rows.Err()
	if err != nil {
		return []Node{}, err

	}
	return nodes, nil
}

//GetElection Check if election exists
func (store *store) GetElection(electionID int64) (Election, error) {
	sqlStmt := "select id,election_id,start_at,close_at,next_elections_at from elections where election_id=?"
	var election Election
	err := store.db.QueryRow(sqlStmt, electionID).Scan(&election.ID, &election.ElectionID, &election.StartAt, &election.CloseAt, &election.NextElectionsAt)
	if err != nil {
		if err != sql.ErrNoRows {
			// a real error happened! you should change your function return
			// to "(bool, error)" and return "false, err" here
			log.Print(err)
		}

		return Election{}, err
	}

	return election, nil
}

//GetParticipates log
func (store *store) GetParticipates(nodeID int, electionID int64) []Participate {

	query := fmt.Sprintf("select participate.node_id,participate.election_id,participate.stake_amount,participate.max_factor from participate inner join elections on elections.election_id = participate.election_id where elections.election_id=%d and participate.node_id=%d", electionID, nodeID)

	rows, err := store.db.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	var participates []Participate
	for rows.Next() {
		var p Participate
		err = rows.Scan(&p.NodeID, &p.ElectionID, &p.StakeAmount, &p.MaxFactor)
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Println(walletFile, walletAddr)
		participates = append(participates, p)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
	}

	return participates
}

//GetKey from db
func (store *store) GetKey(keyType string, nodeID int, electionID int64) (Key, error) {
	sqlStmt := "select key,election_id,node_id,type from keys where election_id=? and node_id=? and type=?"
	var key Key
	err := store.db.QueryRow(sqlStmt, electionID, nodeID, keyType).Scan(&key.Key, &key.ElectionID, &key.NodeID, &key.Type)
	if err != nil {
		if err != sql.ErrNoRows {
			// a real error happened! you should change your function return
			// to "(bool, error)" and return "false, err" here
			log.Fatalln(err)
		}

		return key, err
	}

	return key, nil
}

//AddKey add key to db
func (store *store) AddKey(key Key) (int64, error) {
	stmt, err := store.db.Prepare("INSERT INTO keys(election_id,key,type,node_id) values(?,?,?,?)")
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(key.ElectionID, key.Key, key.Type, key.NodeID)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

//AddElection Add new election id
func (store *store) AddElection(election Election) (int64, error) {
	stmt, err := store.db.Prepare("INSERT INTO elections(election_id,start_at,close_at,next_elections_at) values(?,?,?,?)")
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(election.ElectionID, election.StartAt, election.CloseAt, election.NextElectionsAt)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

//AddParticipate log
func (store *store) AddParticipate(p Participate) (int64, error) {
	stmt, err := store.db.Prepare("INSERT INTO participate(node_id,election_id,stake_amount,max_factor) values(?,?,?,?)")
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(p.NodeID, p.ElectionID, p.StakeAmount, p.MaxFactor)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

//SyncWalletsBalance sync wallets balance to db
func (store *store) SyncWalletsBalance(cln *tonlib.Client) error {
	wallets, err := store.GetWallets(1)
	if err != nil {
		return err
	}
	if len(wallets) == 0 {
		return fmt.Errorf("No wallets found")
	}
	for _, wallet := range wallets {
		var walletBalance int64
		AccountState, err := cln.GetAccountState(*tonlib.NewAccountAddress(wallet.Addr))
		if err != nil {
			return err
		}
		walletBalance = int64(AccountState.Balance)
		err = store.UpdateWalletBalance(wallet.ID, walletBalance)
		if err != nil {
			return err
		}

	}
	return nil
}
