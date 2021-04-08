package depantmodel

import (
	"database/sql"
	"runtime"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type EmployeeAD struct {
	ADID              uint32    `json:"ad_id"`
	EmployeeID        string    `json:"Employee_ID"`
	Description       string    `json:"description"`
	Department        string    `json:"department"`
	Company           string    `json:"company"`
	DN                string    `json:"dn"`
	CN                string    `json:"cn"`
	Mail              string    `json:"mail"`
	DistinguishedName string    `json:"distinguishedName"`
	SAMAccountName    string    `json:"sAMAccountName"`
	Createdate        time.Time `json:"createdate"`
	Updatedate        time.Time `json:"updatedate"`
	whenCreated       string    `json:"whenCreated"`
	whenChanged       string    `json:"whenChanged"`
}
type ParamSetup struct {
	ParamID       int       `json:"param_id"`
	ParamCode     string    `json:"param_code"`
	ParamVal      string    `json:"param_val"`
	ParamCategory string    `json:"param_category"`
	Note          string    `json:"note"`
	CreateDate    time.Time `json:"createdate"`
	UpdateDate    time.Time `json:"updatedate"`
}
type Employee struct {
	UsrID          uint32 `gorm:"primary_key;auto_increment" json:"usr_id"`
	UserLevelID    int    `gorm:"type:int;" json:"userlevelid"`
	EmployeeID     string `gorm:"type:varchar(15);" json:"Employee_ID"`
	LDEPEmployeeID string `gorm:"type:varchar(15);" json:"LDAP_Employee_ID"`
	EmployeeName   string `gorm:"type:varchar(255);" json:"Employee_Name"`
	Organization   string `gorm:"type:varchar(100);" json:"Organization"`
	BusinessUnit   string `gorm:"type:varchar(200);" json:"Business_Unit"`
	AppGroupID     string `gorm:"type:varchar(255);" json:"AppGroupID"`
	MD             string `gorm:"type:varchar(255);" json:"MD"`
	DSOrganize     string `gorm:"type:varchar(255);" json:"DS_Organize"`
	DMOrganize     string `gorm:"type:varchar(255);" json:"DM_Organize"`

	DSProject  string `gorm:"type:varchar(255);" json:"DS_Project"`
	DMProject  string `gorm:"type:varchar(255);" json:"DM_Project"`
	OnOffBoard string `gorm:"type:varchar(255);" json:"On_Off_Board"`
	Remarks    string `gorm:"type:varchar(255);" json:"Remarks"`
	Password   string `gorm:"type:varchar(255);" json:"password"`
	Email      string `gorm:"type:varchar(255);" json:"email"`

	Company  string `json:"Company"`
	Mobile   string `json:"Mobile"`
	Nickname string `json:"nickname"`
	NameTH   string `json:"name_th"`

	StsID      string    `json:sts_id`
	CreateDate time.Time `json:"createdate"`
	UpdateDate time.Time `json:"updatedate"`
}

func (e *Employee) IsEmployeeExist(employeeID string, dbm sql.DB) (bool, error) {
	sqlFind := "SELECT count(usr_id) as c FROM employee WHERE Employee_ID = ? ;"

	pstmt, err := dbm.Prepare(sqlFind)
	defer pstmt.Close()
	if err != nil {
		return false, err
	}

	rows, err := pstmt.Query(employeeID)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	cc := 0
	if rows.Next() {
		if err := rows.Scan(cc); err != nil {
			return false, err
		}
	}

	result := false
	if cc > 0 {
		result = true
	}
	return result, nil
}

func (e *Employee) GetEmployee(employeeID string, dburl string) error {

	sqlemp := "SELECT usr_id, Employee_ID, Employee_Name, Organization, Business_Unit, AppGroupID, Company, Mobile, Nickname, Name_Thai, email, sts_id, createdate, updatedate from employee where Employee_ID = ?"
	dbm, err := sql.Open("mysql", dburl)
	defer dbm.Close()
	if err != nil {

		return err
	}

	pstmt, err := dbm.Prepare(sqlemp)
	defer pstmt.Close()
	if err != nil {

		return err
	}

	rows, err := pstmt.Query(employeeID)
	if err != nil {

		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&e.UsrID, &e.EmployeeID, &e.EmployeeName, &e.Organization, &e.BusinessUnit, &e.AppGroupID, &e.Company, &e.Mobile, &e.Nickname, e.NameTH, e.Email, e.StsID, e.CreateDate, e.UpdateDate); err != nil {

			return err
		}
	}

	return nil
}

func (e *Employee) AddEmployee(dburl string) error {
	sqlemp := "insert into employee (Employee_ID, Employee_Name, Organization, Business_Unit, AppGroupID, Company, Mobile, Nickname, Name_Thai, email, sts_id, createdate, updatedate) values (?, ?, ?, ?)"
	dbm, err := sql.Open("mysql", dburl)
	defer dbm.Close()
	if err != nil {
		return err
	}
	pstmt, err := dbm.Prepare(sqlemp)
	defer pstmt.Close()
	if err != nil {
		return err
	}

	_, err = pstmt.Exec(e.EmployeeID, e.EmployeeName, e.Organization, e.BusinessUnit, e.AppGroupID, e.Company, e.Mobile, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return err
	}

	return nil
}

func (e *Employee) UpdateEmployee(dburl string) error {
	sqlemp := "update employee set Employee_Name=?, Organization=?, Business_Unit=?, AppGroupID=?, Company=?, Mobile=?, Nickname=?, Name_Thai=?, email=?, sts_id=?, updatedate=?) where Employee_ID = ?"
	dbm, err := sql.Open("mysql", dburl)
	defer dbm.Close()
	if err != nil {
		return err
	}
	pstmt, err := dbm.Prepare(sqlemp)
	defer pstmt.Close()
	if err != nil {
		return err
	}

	_, err = pstmt.Exec(e.EmployeeName, e.Organization, e.BusinessUnit, e.AppGroupID, e.Company, e.Mobile, e.Nickname, e.NameTH, e.Email, e.StsID, time.Now().Format("2006-01-02 15:04:05"), e.EmployeeID)
	if err != nil {
		return err
	}

	return nil
}

func (e *EmployeeAD) IsEmployeeAdChanged(sAMAcName string, compare string, dburl string) (bool, error) {
	sqlFind := "SELECT count(Employee_ID) from employee_ad where sAMAccountName = ? and md5(Employee_ID, sAMAccountName, dn, cn, company, department) <> ?"
	dbm, err := sql.Open("mysql", dburl)
	defer dbm.Close()
	if err != nil {
		return false, err
	}
	pstmt, err := dbm.Prepare(sqlFind)
	defer pstmt.Close()
	if err != nil {
		return false, err
	}

	rows, err := pstmt.Query(sAMAcName, compare)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	cc := 0
	if rows.Next() {
		if err := rows.Scan(cc); err != nil {
			return false, err
		}
	}

	result := false
	if cc > 0 {
		result = true
	}
	return result, nil

}

func (e *EmployeeAD) GetEmployeeAd(sAMAcName string, dburl string) error {

	sqlemp := "SELECT ad_id,Employee_ID, sAMAccountName, dn, cn, company, department from employee_ad where sAMAccountName = ?"
	dbm, err := sql.Open("mysql", dburl)
	defer dbm.Close()
	if err != nil {

		return err
	}

	pstmt, err := dbm.Prepare(sqlemp)
	defer pstmt.Close()
	if err != nil {

		return err
	}

	rows, err := pstmt.Query(sAMAcName)
	if err != nil {

		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&e.ADID, &e.EmployeeID, e.SAMAccountName, e.DN, e.CN, e.Company, e.Department); err != nil {
			return err
		}
	}

	return nil
}

func (e *EmployeeAD) AddEmployeeAd(dburl string) error {
	sqlemp := "insert into employee_ad (Employee_ID, sAMAccountName, company, department, mail, whenCreated, whenChanged, createdate, updatedate) values (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	dbm, err := sql.Open("mysql", dburl)
	defer dbm.Close()
	if err != nil {
		return err
	}
	pstmt, err := dbm.Prepare(sqlemp)
	defer pstmt.Close()
	if err != nil {
		return err
	}

	_, err = pstmt.Exec(e.EmployeeID, e.SAMAccountName, e.Company, e.Department, e.Mail, e.whenCreated, e.whenChanged, time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return err
	}

	return nil
}

func (e *EmployeeAD) UpdateEmployeeAd(dburl string) error {
	sqlemp := "update employee_ad set Employee_ID=?, cn=?, dn=?, mail=?, displayName=?, whenChanged=?, updatedate=?  where sAMAccountName = ?"
	dbm, err := sql.Open("mysql", dburl)
	defer dbm.Close()
	if err != nil {
		return err
	}
	pstmt, err := dbm.Prepare(sqlemp)
	defer pstmt.Close()
	if err != nil {
		return err
	}

	_, err = pstmt.Exec(e.EmployeeID, e.SAMAccountName, e.Company, e.Department, e.Mail, e.whenCreated, e.whenChanged, time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		return err
	}

	return nil
}

func funcName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}
