package user

import (
	"database/sql"
	"fmt"
	"quarto/models/postgresql"
	"time"

	"github.com/charmbracelet/log"

	"github.com/jackc/pgx/v4"
)

func ScanUser(row pgx.Row) (u User, err error) {

	var (
		id                                      sql.NullInt64
		email, username, password, recoverToken sql.NullString
		admin, enable                           sql.NullBool
	)

	err = row.Scan(
		&id,
		&email,
		&username,
		&password,
		&recoverToken,
		&admin,
		&enable,
	)

	if err != nil {
		return
	}

	u = User{
		ID:           id.Int64,
		Email:        email.String,
		Username:     username.String,
		Password:     password.String,
		RecoverToken: recoverToken.String,
		Enable:       enable.Bool,
		Admin:        admin.Bool,
	}

	return
}

func GetSQLUserToken(email, password string) (token UserToken, err error) {

	query := "select * from account " +
		"where enable=true and email=$1 and password=crypt($2, password)"

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	row := sqlCo.QueryRow(postgresql.SQLCtx, query, email, password)
	u, err := ScanUser(row)

	if err == pgx.ErrNoRows {
		return
	} else if err != nil {
		log.Error("During GetSQLUserToken query", "error", err)
		return
	}

	token = UserToken{
		User:      u,
		CreatedAt: time.Now(),
	}

	return
}

func CreateAccount(email, username, password string) (id int64) {

	query := "insert into account (email, username, password) " +
		"VALUES ($1,$2,crypt($3, gen_salt('bf'))) RETURNING id"

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return -1
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	id = time.Now().UnixNano()

	err = sqlCo.QueryRow(postgresql.SQLCtx, query, email, username, password).Scan(&id)
	if err != nil {
		log.Warn("During CreateAccount query", "error", err)
		return -1
	}
	return id
}

func DeleteAccount(id int64) (msg string) {

	query := "update account set enable=false where id=$1"

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		msg = "Internal server error"
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	cmd, err := sqlCo.Exec(postgresql.SQLCtx, query, id)
	if err != nil {
		return "Internal server error"
	} else if cmd.RowsAffected() == 0 {
		return "Account not found"
	}
	return
}

func PasswordCheck(id int64, password string) (checked bool) {
	query := "select id from account where enable=true and id=$1 and password=crypt($2, password)"

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	var id_ int64
	err = sqlCo.QueryRow(postgresql.SQLCtx, query, id, password).Scan(&id_)
	if err == nil {
		return id_ == id
	}
	return
}

func CheckEmailAvailability(email string) (available bool) {
	query := "select id from account where email=$1"

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	var id int64
	err = sqlCo.QueryRow(postgresql.SQLCtx, query, email).Scan(&id)
	if err == pgx.ErrNoRows {
		return true
	}
	return
}

func GetUserById(UserId int64) (u User, err error) {

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	var query = "SELECT * FROM account WHERE enable=true and id=$1"

	row := sqlCo.QueryRow(postgresql.SQLCtx, query, UserId)
	u, err = ScanUser(row)

	return
}

func GetAllUsers() (users []User) {

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}

	defer sqlCo.Close(postgresql.SQLCtx)

	query := "SELECT * FROM account"

	rows, err := sqlCo.Query(postgresql.SQLCtx, query)
	if err != nil {
		log.Error("During GetAllUsers query", "error", err)
		return
	}
	defer rows.Close()

	for rows.Next() {

		u, err := ScanUser(rows)

		if err == pgx.ErrNoRows {
			return
		} else if err != nil {
			log.Error("During GetAllUsers scan", "error", err)
			return
		}

		users = append(users, u)
	}

	return
}

func GetUserByEmail(userEmail string) (u User, msg string) {

	if userEmail == "" {
		msg = "Empty email"
		return
	}

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		msg = "Internal server error"
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	query := "SELECT * FROM account WHERE email=$1"

	row := sqlCo.QueryRow(postgresql.SQLCtx, query, userEmail)
	u, err = ScanUser(row)

	if err == pgx.ErrNoRows {
		msg = "Username not found"
	} else if err != nil {
		log.Error("During GetUserByEmail query", "error", err)
		msg = "Internal server error"
	}

	return
}

func UpdateUser(user User, password bool) (ok bool) {

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	var query string
	var args []any
	if password {
		query = "UPDATE account set (email, username, password) = ($1,$2,crypt($3, gen_salt('bf'))) " +
			"WHERE id=$4"
		args = []any{
			user.Email,
			user.Username,
			user.Password,
			user.ID,
		}
	} else {
		query = "UPDATE account set (email, username) = ($1,$2) " +
			"WHERE id=$3"
		args = []any{
			user.Email,
			user.Username,
			user.ID,
		}
	}

	cmd, err := sqlCo.Exec(postgresql.SQLCtx, query, args...)
	ok = cmd.RowsAffected() == 1 && err == nil
	return
}

func CreateRecoverToken(email string) (string, error) {

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return "", err
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	var username, recoverToken sql.NullString

	query := "UPDATE account set recover_token=gen_random_uuid() where email=$1 RETURNING username, recover_token"
	err = sqlCo.QueryRow(postgresql.SQLCtx, query, email).Scan(&username, &recoverToken)
	if err == pgx.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return recoverToken.String, nil
}

func ResetPassword(token, password string) (ok bool) {

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return false
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	query := "UPDATE account set password=crypt($2, gen_salt('bf')), recover_token=null where recover_token=$1"
	cmd, err := sqlCo.Exec(postgresql.SQLCtx, query, token, password)
	return cmd.RowsAffected() == 1 && err == nil
}

func IsInOrganization(userId, orgId int64) (ok bool) {

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return false
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	query := "SELECT account FROM account_in_organization WHERE account=$1 and organization=$2"
	row := sqlCo.QueryRow(postgresql.SQLCtx, query, userId, orgId)

	var id int64
	err = row.Scan(&id)
	return err == nil
}

func ListOrganizationMembers(orgID int64) (users UserList, err error) {

	query := `SELECT a.* 
				FROM account a 
				JOIN account_in_organization aio ON a.id = aio.account 
				WHERE aio.organization = $1`

	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	rows, err := sqlCo.Query(postgresql.SQLCtx, query, orgID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		user, err = ScanUser(rows)
		if err != nil {
			return
		}

		users = append(users, user)
	}

	return
}

// GetUserByID récupère un utilisateur par son ID
func GetUserByID(userID int64) (*User, error) {
	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return nil, err
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	query := "SELECT id, email, username, password, recover_token, admin, enable FROM account WHERE id = $1 AND enable = true"
	row := sqlCo.QueryRow(postgresql.SQLCtx, query, userID)

	user, err := ScanUser(row)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Utilisateur non trouvé
		}
		return nil, err
	}

	return &user, nil
}

// GetUsersWithPagination récupère la liste des utilisateurs avec pagination
func GetUsersWithPagination(page, pageSize int, search string) (*UserPaginationResponse, error) {
	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return nil, err
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	// Calculer l'offset
	offset := (page - 1) * pageSize

	// Construire la requête avec recherche optionnelle
	whereClause := "WHERE enable = true"
	args := []interface{}{pageSize, offset}
	argIndex := 3

	if search != "" {
		whereClause += " AND username ILIKE $" + fmt.Sprintf("%d", argIndex)
		args = append(args, "%"+search+"%")
		argIndex++
	}

	// Compter le total
	countQuery := "SELECT COUNT(*) FROM account " + whereClause[6:] // Enlever "WHERE " du début
	var total int64
	err = sqlCo.QueryRow(postgresql.SQLCtx, countQuery, args[2:]...).Scan(&total)
	if err != nil {
		return nil, err
	}

	// Récupérer les utilisateurs
	query := "SELECT id, email, username, password, recover_token, admin, enable FROM account " +
		whereClause + " ORDER BY username ASC LIMIT $1 OFFSET $2"

	rows, err := sqlCo.Query(postgresql.SQLCtx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []UserPublic
	for rows.Next() {
		user, err := ScanUser(rows)
		if err != nil {
			continue
		}
		users = append(users, user.ToPublic())
	}

	// Calculer le nombre total de pages
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))

	return &UserPaginationResponse{
		Users:      users,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

// GetUsersPaginated récupère une liste paginée des utilisateurs
func GetUsersPaginated(page, pageSize int) ([]UserPublic, int64, error) {
	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return nil, 0, fmt.Errorf("erreur de connexion à la base de données: %v", err)
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	// Calculer l'offset
	offset := (page - 1) * pageSize

	// Compter le total d'utilisateurs actifs
	var total int64
	err = sqlCo.QueryRow(postgresql.SQLCtx, "SELECT COUNT(*) FROM account WHERE enable = TRUE").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("erreur lors du comptage des utilisateurs: %v", err)
	}

	// Récupérer les utilisateurs paginés
	query := `
		SELECT id, username 
		FROM account 
		WHERE enable = TRUE 
		ORDER BY username ASC 
		LIMIT $1 OFFSET $2`

	rows, err := sqlCo.Query(postgresql.SQLCtx, query, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("erreur lors de la récupération des utilisateurs: %v", err)
	}
	defer rows.Close()

	var users []UserPublic
	for rows.Next() {
		var user UserPublic
		err := rows.Scan(&user.ID, &user.Username)
		if err != nil {
			log.Error("Erreur lors du scan d'un utilisateur", "error", err)
			continue
		}
		users = append(users, user)
	}

	return users, total, nil
}

// GetUserPublicByID récupère un utilisateur par son ID (version publique)
func GetUserPublicByID(userID int64) (*UserPublic, error) {
	sqlCo, err := pgx.ConnectConfig(postgresql.SQLCtx, postgresql.SQLConn)
	if err != nil {
		return nil, fmt.Errorf("erreur de connexion à la base de données: %v", err)
	}
	defer sqlCo.Close(postgresql.SQLCtx)

	query := "SELECT id, username FROM account WHERE id = $1 AND enable = TRUE"
	row := sqlCo.QueryRow(postgresql.SQLCtx, query, userID)

	var user UserPublic
	err = row.Scan(&user.ID, &user.Username)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("utilisateur non trouvé")
		}
		return nil, fmt.Errorf("erreur lors de la récupération de l'utilisateur: %v", err)
	}

	return &user, nil
}
