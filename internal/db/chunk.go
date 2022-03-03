package db

import (
	"database/sql"
	"fmt"
	"math"
	"reflect"
)

type ChunkResult struct {
	results []sql.Result
}

// LastInsertId foo
func (p ChunkResult) LastInsertId() (int64, error) {
	return p.results[len(p.results)].LastInsertId()
}

// RowsAffected foo
func (p ChunkResult) RowsAffected() (int64, error) {
	var (
		affected, i int64
		err         error
	)

	for _, v := range p.results {
		i, err = v.RowsAffected()
		affected += i
	}

	return affected, err
}

// Chunk foo
func (p *DB) Chunk(chunkSize int, query string, arg interface{}) (sql.Result, error) {
	var (
		err      error
		res      sql.Result
		max      int
		c, total int64

		cRes = ChunkResult{results: []sql.Result{}}
		v    = reflect.ValueOf(arg)
	)

	if query == "" {
		return nil, fmt.Errorf("empty query string")
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if !(v.Kind() == reflect.Array || v.Kind() == reflect.Slice) {
		return nil, fmt.Errorf("kind (%s) is not an array", v.Kind())
	}

	for i := 0; i < v.Len(); i += chunkSize {
		max = int(math.Min(float64(i+chunkSize), float64(v.Len())))

		res, err = p.NamedExec(query, v.Slice(i, max).Interface())
		cRes.results = append(cRes.results, res)
		if err != nil {
			return cRes, fmt.Errorf("%w [i:%v, max:%v]", err, i, max)
		}

		c, err = res.RowsAffected()
		if err != nil {
			return cRes, fmt.Errorf("%w [i:%v, max:%v]", err, i, max)
		}

		total += c
	}

	return cRes, nil
}
