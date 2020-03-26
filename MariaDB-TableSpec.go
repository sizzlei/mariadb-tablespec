package main
import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"flag"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
)


func main() {
	
	srvhost := flag.String("host","localhost","Target Database Host.")
	srvport := flag.String("port","3306","Target Database Port.")
	srvuser := flag.String("user","root","Database Connect User.")
	srvpass := flag.String("password","","Database Connect Password.")
	xlsfile := flag.String("file","default.xlsx","Export File Path and Name.")
	srvdb := flag.String("database","mysql","Export for Database.")

	flag.Parse()

	// Declare Config
	var srvHost string = *srvhost
	var srvPort string = *srvport
	var srvUser string = *srvuser
	var srvPass string = *srvpass
	var xlsFile string = *xlsfile
	var srvDB string = *srvdb

	// Declare Table
	type dbTable struct {
		schema 		string
		table_nm 	string
		table_type 	string
		engine		string
		row_format 	string
		collation 	string
		comment 	sql.NullString
	}

	var dvTB dbTable

	type tbColumn struct {
		col_nm 		string
		col_default sql.NullString
		nullable	string
		col_type	string
		col_length	sql.NullInt64
		charset		sql.NullString
		collation	sql.NullString
		colkey		sql.NullString
		extra		sql.NullString
		comment		sql.NullString
	}

	var dvCol tbColumn

	type colConst struct {
		const_key	string
		const_col	string
		const_refer	sql.NullString
		const_del	string
		const_upd	string
	}

	var dvConst colConst

	type colIndex struct {
		idx_nm		string
		idx_unique	int
		idx_col		string
	}

	var dvIdx colIndex
	// database Connection
	var dnsFormat string = "%s:%s@tcp(%s:%s)/information_schema"
	var DNS string = fmt.Sprintf(dnsFormat,srvUser,srvPass,srvHost,srvPort)
	dbCon, err := sql.Open("mysql",DNS)
	if err != nil {
		panic(err)
	}
	defer dbCon.Close()

	// Create Excel File
	exFile := excelize.NewFile()
	exSheet := exFile.NewSheet(srvDB)
	exFile.DeleteSheet("Sheet1")
	exFile.SetActiveSheet(exSheet)

	// Style Set
	exCont, err := exFile.NewStyle(`{"border":[{"type":"left","color":"00d3d3d3","style":1},{"type":"right","color":"00d3d3d3","style":1},{"type":"top","color":"00d3d3d3","style":1},{"type":"bottom","color":"00d3d3d3","style":1}]}`)
	if err != nil {
		panic(err)
	}

	exTitle, err := exFile.NewStyle(`{"font":{"bold":true,"color":"#000000"},"fill":{"type":"pattern","color":["#E0EBF5"],"pattern":1},"border":[{"type":"left","color":"00d3d3d3","style":1},{"type":"right","color":"00d3d3d3","style":1},{"type":"top","color":"00d3d3d3","style":1},{"type":"bottom","color":"00d3d3d3","style":1}]}`)
	if err != nil {
		panic(err)
	}

	exTitleGr, err := exFile.NewStyle(`{"font":{"bold":true,"color":"#ffffff"},"fill":{"type":"pattern","color":["#848484"],"pattern":1},"border":[{"type":"left","color":"00d3d3d3","style":1},{"type":"right","color":"00d3d3d3","style":1},{"type":"top","color":"00d3d3d3","style":1},{"type":"bottom","color":"00d3d3d3","style":1}]}`)
	if err != nil {
		panic(err)
	}

	// Get Database Info
	var getInfoFormat string = "SELECT table_schema,table_name,table_type,ENGINE,row_format,table_collation,table_comment FROM information_schema.TABLES WHERE table_schema='%s'"
	var getTB string = fmt.Sprintf(getInfoFormat,srvDB)
	tbRows, err:= dbCon.Query(getTB)
	if err != nil {
		fmt.Printf("Not Exists Table.")
		panic(err)
	}
	defer tbRows.Close()

	// declare
	colRownum := 1

	for tbRows.Next() {
		err := tbRows.Scan(
			&dvTB.schema,
			&dvTB.table_nm,
			&dvTB.table_type,
			&dvTB.engine,
			&dvTB.row_format,
			&dvTB.collation,
			&dvTB.comment)
		if err != nil {
			panic(err)
		}
		
		/*
		
		--------------------- Table Name
		
		*/
		
		// Title	
		// Merge Cell
		exFile.MergeCell(srvDB, fmt.Sprintf("A%d",colRownum), fmt.Sprintf("B%d",colRownum))
		exFile.MergeCell(srvDB, fmt.Sprintf("C%d",colRownum), fmt.Sprintf("J%d",colRownum))

		// Style
		exFile.SetCellStyle(srvDB, fmt.Sprintf("A%d",colRownum), fmt.Sprintf("B%d",colRownum),exTitleGr)
		exFile.SetCellStyle(srvDB, fmt.Sprintf("C%d",colRownum), fmt.Sprintf("J%d",colRownum),exCont)

		exFile.SetCellValue(srvDB,fmt.Sprintf("A%d",colRownum),"Table Name")		
		exFile.SetCellValue(srvDB,fmt.Sprintf("C%d",colRownum),dvTB.table_nm)
		colRownum += 1

		// Title Style
		
		
		exFile.MergeCell(srvDB, fmt.Sprintf("A%d",colRownum), fmt.Sprintf("B%d",colRownum))
		exFile.MergeCell(srvDB, fmt.Sprintf("C%d",colRownum), fmt.Sprintf("J%d",colRownum))

		// Style
		exFile.SetCellStyle(srvDB, fmt.Sprintf("A%d",colRownum), fmt.Sprintf("B%d",colRownum),exTitleGr)
		exFile.SetCellStyle(srvDB, fmt.Sprintf("C%d",colRownum), fmt.Sprintf("J%d",colRownum),exCont)

		exFile.SetCellValue(srvDB,fmt.Sprintf("A%d",colRownum),"Description")
		if dvTB.comment.Valid {
			exFile.SetCellValue(srvDB,fmt.Sprintf("C%d",colRownum),dvTB.comment.String)
		}
		colRownum += 1
		/*
		
		--------------------- Column Sector
		
		*/
		// Column Title
		colDiv := []string{"A","B","C","D","E","F","G","H","I","J"}
		colNames := []string{"No","Column Name","Data Type","Length","Nullable","Key","Extra","Collation","Default","Comment"}

		for idx, colNm := range colNames {
			exFile.SetCellValue(srvDB, fmt.Sprintf("%s%d",colDiv[idx],colRownum),colNm)
			exFile.SetCellStyle(srvDB, fmt.Sprintf("%s%d",colDiv[idx],colRownum), fmt.Sprintf("%s%d",colDiv[idx],colRownum),exTitle)
		}
		colRownum += 1

		// Get Column
		var getColFormat string = "SELECT column_name,column_default,is_nullable,data_type,character_maximum_length AS column_length,character_set_name,collation_name,column_key,extra,column_comment FROM information_schema.COLUMNS WHERE table_name='%s' AND table_schema='%s' ORDER BY ordinal_position;"
		var getCol string = fmt.Sprintf(getColFormat,dvTB.table_nm,dvTB.schema)
		colRows, err := dbCon.Query(getCol)
		if err != nil {
			panic(err)
		}
		defer colRows.Close()		

		// Column Insert
		colNum := 0
		for colRows.Next() {
			err := colRows.Scan(
				&dvCol.col_nm,
				&dvCol.col_default,
				&dvCol.nullable,
				&dvCol.col_type,
				&dvCol.col_length,
				&dvCol.charset,
				&dvCol.collation,
				&dvCol.colkey,
				&dvCol.extra,
				&dvCol.comment)
			if err != nil {
				panic(err)
			}
			// fmt.Println("Column : ", dvCol.col_nm,dvCol.col_type,dvCol.col_length,dvCol.comment)
			colNum += 1
			exFile.SetCellValue(srvDB,fmt.Sprintf("A%d",colRownum),colNum)
			// Column Name
			exFile.SetCellValue(srvDB,fmt.Sprintf("B%d",colRownum),dvCol.col_nm)
			// Column Type
			exFile.SetCellValue(srvDB,fmt.Sprintf("C%d",colRownum),dvCol.col_type)
			// Column Length
			if dvCol.col_length.Valid {
				exFile.SetCellValue(srvDB,fmt.Sprintf("D%d",colRownum),dvCol.col_length.Int64)
			}
			// Column Nullable
			exFile.SetCellValue(srvDB,fmt.Sprintf("E%d",colRownum),dvCol.nullable)
			// Column Key
			if dvCol.colkey.Valid {
				exFile.SetCellValue(srvDB,fmt.Sprintf("F%d",colRownum),dvCol.colkey.String)
			}
			// Column Extra
			if dvCol.extra.Valid {
				exFile.SetCellValue(srvDB,fmt.Sprintf("G%d",colRownum),dvCol.extra.String)
			}
			// Column Collation
			if dvCol.collation.Valid {
				exFile.SetCellValue(srvDB,fmt.Sprintf("H%d",colRownum),dvCol.collation.String)
			}
			// Column Default
			if dvCol.col_default.Valid {
				exFile.SetCellValue(srvDB,fmt.Sprintf("I%d",colRownum),string(dvCol.col_default.String))
			}
			// Column Comment
			if dvCol.comment.Valid {
				exFile.SetCellValue(srvDB,fmt.Sprintf("J%d",colRownum),dvCol.comment.String)
			}
			// Style
			for _, colDivnm := range colDiv {
				exFile.SetCellStyle(srvDB, fmt.Sprintf("%s%d",colDivnm,colRownum), fmt.Sprintf("%s%d",colDivnm,colRownum),exCont)
			}

			colRownum += 1
		}
		/*
		
		--------------------- Constraint Sector
		
		*/
		// Constraint Title
		exFile.MergeCell(srvDB, fmt.Sprintf("A%d",colRownum), fmt.Sprintf("J%d",colRownum))
		// Style
		exFile.SetCellStyle(srvDB, fmt.Sprintf("A%d",colRownum), fmt.Sprintf("J%d",colRownum),exTitleGr)

		exFile.SetCellValue(srvDB,fmt.Sprintf("A%d",colRownum),"Foreign Key Info")
		colRownum += 1

		constDiv := []string{"A","D","E","I","J"}
		constNames := []string{"Name","Column","Referance","DELETE","UPDATE"}

		// Merge Cell
		exFile.MergeCell(srvDB, fmt.Sprintf("A%d",colRownum), fmt.Sprintf("C%d",colRownum))
		exFile.MergeCell(srvDB, fmt.Sprintf("D%d",colRownum), fmt.Sprintf("D%d",colRownum))
		exFile.MergeCell(srvDB, fmt.Sprintf("E%d",colRownum), fmt.Sprintf("H%d",colRownum))

		// Style 
		exFile.SetCellStyle(srvDB, fmt.Sprintf("A%d",colRownum), fmt.Sprintf("C%d",colRownum),exTitle)
		exFile.SetCellStyle(srvDB, fmt.Sprintf("D%d",colRownum), fmt.Sprintf("D%d",colRownum),exTitle)
		exFile.SetCellStyle(srvDB, fmt.Sprintf("E%d",colRownum), fmt.Sprintf("H%d",colRownum),exTitle)
		exFile.SetCellStyle(srvDB, fmt.Sprintf("I%d",colRownum), fmt.Sprintf("I%d",colRownum),exTitle)
		exFile.SetCellStyle(srvDB, fmt.Sprintf("J%d",colRownum), fmt.Sprintf("J%d",colRownum),exTitle)

		for idx, constNm := range constNames {
			exFile.SetCellValue(srvDB,fmt.Sprintf("%s%d",constDiv[idx],colRownum),constNm)
		}
		colRownum += 1

		// Get Constraint
		var getConstFormat string = "SELECT x.constraint_name AS constraint_key,group_concat(x.column_name) AS con_column,concat(x.referenced_table_name,'.',x.referenced_column_name) AS refer_info,y.DELETE_RULE,y.UPDATE_RULE FROM information_schema.KEY_COLUMN_USAGE x INNER JOIN information_schema.REFERENTIAL_CONSTRAINTS y ON x.constraint_name=y.constraint_name WHERE x.CONSTRAINT_SCHEMA='%s' AND x.table_name='%s' AND x.constraint_name<> 'PRIMARY' GROUP BY x.constraint_name"
		var getConst string = fmt.Sprintf(getConstFormat,dvTB.schema,dvTB.table_nm)
		constRows, err := dbCon.Query(getConst)
		if err != nil {
			panic(err)
		}
		defer constRows.Close()

		// Constraint Insert
		for constRows.Next() {
			err := constRows.Scan(&dvConst.const_key,&dvConst.const_col,&dvConst.const_refer,&dvConst.const_del,&dvConst.const_upd)
			if err != nil {
				panic(err)
			}
			if dvConst.const_refer.Valid {
				// Merge cell
				exFile.MergeCell(srvDB, fmt.Sprintf("A%d",colRownum), fmt.Sprintf("C%d",colRownum))
				exFile.MergeCell(srvDB, fmt.Sprintf("D%d",colRownum), fmt.Sprintf("D%d",colRownum))
				exFile.MergeCell(srvDB, fmt.Sprintf("E%d",colRownum), fmt.Sprintf("H%d",colRownum))

				// Style 
				exFile.SetCellStyle(srvDB, fmt.Sprintf("A%d",colRownum), fmt.Sprintf("C%d",colRownum),exCont)
				exFile.SetCellStyle(srvDB, fmt.Sprintf("D%d",colRownum), fmt.Sprintf("D%d",colRownum),exCont)
				exFile.SetCellStyle(srvDB, fmt.Sprintf("E%d",colRownum), fmt.Sprintf("H%d",colRownum),exCont)
				exFile.SetCellStyle(srvDB, fmt.Sprintf("I%d",colRownum), fmt.Sprintf("I%d",colRownum),exCont)
				exFile.SetCellStyle(srvDB, fmt.Sprintf("J%d",colRownum), fmt.Sprintf("J%d",colRownum),exCont)

				exFile.SetCellValue(srvDB,fmt.Sprintf("A%d",colRownum),dvConst.const_key)
				exFile.SetCellValue(srvDB,fmt.Sprintf("D%d",colRownum),dvConst.const_col)
				exFile.SetCellValue(srvDB,fmt.Sprintf("E%d",colRownum),dvConst.const_refer.String)
				exFile.SetCellValue(srvDB,fmt.Sprintf("I%d",colRownum),dvConst.const_del)
				exFile.SetCellValue(srvDB,fmt.Sprintf("J%d",colRownum),dvConst.const_upd)
				colRownum += 1
			}
			
			
		}
		/*
		
		--------------------- Index Info

		*/
		// Index Title
		exFile.MergeCell(srvDB, fmt.Sprintf("A%d",colRownum), fmt.Sprintf("J%d",colRownum))
		// Style
		exFile.SetCellStyle(srvDB, fmt.Sprintf("A%d",colRownum), fmt.Sprintf("J%d",colRownum),exTitleGr)

		exFile.SetCellValue(srvDB,fmt.Sprintf("A%d",colRownum),"Index Info")
		colRownum += 1

		idxDiv := []string{"A","C","G"}
		idxNames := []string{"Index Div","Index Name","Index Column"}

		// Merge Cell
		exFile.MergeCell(srvDB, fmt.Sprintf("A%d",colRownum), fmt.Sprintf("B%d",colRownum))
		exFile.MergeCell(srvDB, fmt.Sprintf("C%d",colRownum), fmt.Sprintf("F%d",colRownum))
		exFile.MergeCell(srvDB, fmt.Sprintf("G%d",colRownum), fmt.Sprintf("J%d",colRownum))

		// Style
		exFile.SetCellStyle(srvDB, fmt.Sprintf("A%d",colRownum), fmt.Sprintf("B%d",colRownum),exTitle)
		exFile.SetCellStyle(srvDB, fmt.Sprintf("C%d",colRownum), fmt.Sprintf("F%d",colRownum),exTitle)
		exFile.SetCellStyle(srvDB, fmt.Sprintf("G%d",colRownum), fmt.Sprintf("J%d",colRownum),exTitle)

		for idx, idxNm := range idxNames {
			exFile.SetCellValue(srvDB,fmt.Sprintf("%s%d",idxDiv[idx],colRownum),idxNm)
		}
		colRownum += 1

		// Get Index Info 
		var getIndexFormat string = "SELECT INDEX_NAME,NON_UNIQUE,GROUP_CONCAT(COLUMN_NAME SEPARATOR ',') AS COLUMN_NAMES FROM information_schema.STATISTICS WHERE TABLE_NAME='%s' AND TABLE_SCHEMA='%s' AND INDEX_NAME !='PRIMARY' GROUP BY TABLE_SCHEMA,TABLE_NAME,INDEX_NAME ORDER BY INDEX_NAME"
		var getIndex string = fmt.Sprintf(getIndexFormat,dvTB.table_nm,dvTB.schema)
		indexRows, err := dbCon.Query(getIndex)
		if err != nil {
			panic(err)
		}
		defer indexRows.Close()

		for indexRows.Next() {
			err := indexRows.Scan(&dvIdx.idx_nm,&dvIdx.idx_unique,&dvIdx.idx_col)
			if err != nil {
				panic(err)
			}
			// Merge Cell
			exFile.MergeCell(srvDB, fmt.Sprintf("A%d",colRownum), fmt.Sprintf("B%d",colRownum))
			exFile.MergeCell(srvDB, fmt.Sprintf("C%d",colRownum), fmt.Sprintf("F%d",colRownum))
			exFile.MergeCell(srvDB, fmt.Sprintf("G%d",colRownum), fmt.Sprintf("J%d",colRownum))

			// Style
			exFile.SetCellStyle(srvDB, fmt.Sprintf("A%d",colRownum), fmt.Sprintf("B%d",colRownum),exCont)
			exFile.SetCellStyle(srvDB, fmt.Sprintf("C%d",colRownum), fmt.Sprintf("F%d",colRownum),exCont)
			exFile.SetCellStyle(srvDB, fmt.Sprintf("G%d",colRownum), fmt.Sprintf("J%d",colRownum),exCont)
			
			if dvIdx.idx_unique == 0 {
				exFile.SetCellValue(srvDB,fmt.Sprintf("A%d",colRownum),"Unique")
			} else {
				exFile.SetCellValue(srvDB,fmt.Sprintf("A%d",colRownum),"Index")
			}			
			exFile.SetCellValue(srvDB,fmt.Sprintf("C%d",colRownum),dvIdx.idx_nm)
			exFile.SetCellValue(srvDB,fmt.Sprintf("G%d",colRownum),dvIdx.idx_col)

			colRownum += 1

		}


		/*
		
		--------------------- Table Info

		*/
		// Merge cell
		exFile.MergeCell(srvDB, fmt.Sprintf("A%d",colRownum), fmt.Sprintf("B%d",colRownum))
		exFile.MergeCell(srvDB, fmt.Sprintf("C%d",colRownum), fmt.Sprintf("D%d",colRownum))
		exFile.MergeCell(srvDB, fmt.Sprintf("E%d",colRownum), fmt.Sprintf("F%d",colRownum))
		exFile.MergeCell(srvDB, fmt.Sprintf("G%d",colRownum), fmt.Sprintf("J%d",colRownum))

		// Style
		exFile.SetCellStyle(srvDB, fmt.Sprintf("A%d",colRownum), fmt.Sprintf("B%d",colRownum),exTitle)
		exFile.SetCellStyle(srvDB, fmt.Sprintf("C%d",colRownum), fmt.Sprintf("D%d",colRownum),exCont)
		exFile.SetCellStyle(srvDB, fmt.Sprintf("E%d",colRownum), fmt.Sprintf("F%d",colRownum),exTitle)
		exFile.SetCellStyle(srvDB, fmt.Sprintf("G%d",colRownum), fmt.Sprintf("J%d",colRownum),exCont)


		exFile.SetCellValue(srvDB,fmt.Sprintf("A%d",colRownum),"Engine")
		exFile.SetCellValue(srvDB,fmt.Sprintf("C%d",colRownum),dvTB.engine)
		exFile.SetCellValue(srvDB,fmt.Sprintf("E%d",colRownum),"Row Format")
		exFile.SetCellValue(srvDB,fmt.Sprintf("G%d",colRownum),dvTB.row_format)
		colRownum += 1

		// Merge cell
		exFile.MergeCell(srvDB, fmt.Sprintf("A%d",colRownum), fmt.Sprintf("B%d",colRownum))
		exFile.MergeCell(srvDB, fmt.Sprintf("C%d",colRownum), fmt.Sprintf("D%d",colRownum))
		exFile.MergeCell(srvDB, fmt.Sprintf("E%d",colRownum), fmt.Sprintf("F%d",colRownum))
		exFile.MergeCell(srvDB, fmt.Sprintf("G%d",colRownum), fmt.Sprintf("J%d",colRownum))

		// Style
		exFile.SetCellStyle(srvDB, fmt.Sprintf("A%d",colRownum), fmt.Sprintf("B%d",colRownum),exTitle)
		exFile.SetCellStyle(srvDB, fmt.Sprintf("C%d",colRownum), fmt.Sprintf("D%d",colRownum),exCont)
		exFile.SetCellStyle(srvDB, fmt.Sprintf("E%d",colRownum), fmt.Sprintf("F%d",colRownum),exTitle)
		exFile.SetCellStyle(srvDB, fmt.Sprintf("G%d",colRownum), fmt.Sprintf("J%d",colRownum),exCont)

		exFile.SetCellValue(srvDB,fmt.Sprintf("A%d",colRownum),"Table Type")
		exFile.SetCellValue(srvDB,fmt.Sprintf("C%d",colRownum),dvTB.table_type)
		exFile.SetCellValue(srvDB,fmt.Sprintf("E%d",colRownum),"Collation")
		exFile.SetCellValue(srvDB,fmt.Sprintf("G%d",colRownum),dvTB.collation)

		colRownum += 3
	}

	
	err = exFile.SaveAs(xlsFile)
	if err != nil {
		panic(err)
	}

}