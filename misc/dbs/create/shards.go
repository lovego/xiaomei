package create

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lovego/bsql"
	"github.com/lovego/config/conf"
)

func setupShard(rawDb *sql.DB, shardNo int, settings conf.ShardsSettings) error {
	if shardNo == 0 || settings.IdSeqIncrementBy == 0 {
		return nil
	}
	db := bsql.New(rawDb, 10*time.Second)

	var idSeqNames []string
	if err := db.Query(&idSeqNames, getQueryIdSeqsSql()); err != nil {
		return err
	}
	for _, seqName := range idSeqNames {
		if err := setupSequence(db, seqName, settings.IdSeqIncrementBy, shardNo); err != nil {
			return err
		}
	}
	return nil
}

func setupSequence(db *bsql.DB, seqName string, incrementBy, restartWith int) error {
	var isCalled bool
	if err := db.Query(&isCalled, `SELECT is_called from `+seqName); err != nil {
		return err
	}
	if isCalled {
		return nil
	}

	_, err := db.Exec(fmt.Sprintf(
		"ALTER SEQUENCE %s INCREMENT BY %d RESTART WITH %d", seqName, incrementBy, restartWith,
	))
	return err
}

var queryIdSeqsSql string

func getQueryIdSeqsSql() string {
	if queryIdSeqsSql == `` {
		const prefix, suffix = `nextval('`, `'::regclass)`

		queryIdSeqsSql = fmt.Sprintf(
			"SELECT substring(column_default FROM %d FOR length(column_default)-%d) AS sequence",
			len(prefix)+1, len(prefix)+len(suffix),
		) + `
    FROM information_schema.columns
    WHERE table_schema NOT IN ('pg_catalog', 'information_schema') AND
          column_name = 'id' AND column_default LIKE $$` + prefix + `%_id_seq` + suffix + `$$
    ORDER BY column_default`
	}
	return queryIdSeqsSql
}
