package dbexport

import (
	"database/sql"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"reflect"
	"testing"
)

func TestConverter_MarshalDbRows(t *testing.T) {
	type args struct {
		rows *sql.Rows
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "valid not null values",
			args: args{
				rows: mockRowsToSqlRows(
					sqlmock.NewRows([]string{"id", "title", "float", "bool"}).
						AddRow(
							sql.NullInt64{Int64: 1, Valid: true},
							sql.NullString{String: "title 1", Valid: true},
							sql.NullFloat64{Float64: 1.1, Valid: true},
							sql.NullBool{Bool: true, Valid: true},
						).
						AddRow(
							sql.NullInt64{Int64: 2, Valid: true},
							sql.NullString{String: "title 2", Valid: true},
							sql.NullFloat64{Float64: 2.2, Valid: true},
							sql.NullBool{Bool: false, Valid: true},
						)),
			},
			want: func() []byte {
				s := &[]struct {
					Id    int64   `json:"id"`
					Title string  `json:"title"`
					Float float64 `json:"float"`
					Bool  bool    `json:"bool"`
				}{
					{1, "title 1", 1.1, true},
					{2, "title 2", 2.2, false},
				}
				r, _ := json.Marshal(s)
				return r
			}(),
			wantErr: false,
		},
		{
			name: "valid null values",
			args: args{
				rows: mockRowsToSqlRows(
					sqlmock.NewRows([]string{"id", "title", "float", "bool"}).
						AddRow(
							sql.NullInt64{Int64: 0, Valid: false},
							sql.NullString{String: "", Valid: false},
							sql.NullFloat64{Float64: 0, Valid: false},
							sql.NullBool{Bool: false, Valid: false},
						)),
			},
			want: func() []byte {
				s := &[]struct {
					Id    *int64   `json:"id"`
					Title *string  `json:"title"`
					Float *float64 `json:"float"`
					Bool  *bool    `json:"bool"`
				}{
					{nil, nil, nil, nil},
				}
				r, _ := json.Marshal(s)
				return r
			}(),
			wantErr: false,
		},
		{
			name: "valid empty rows",
			args: args{
				rows: mockRowsToSqlRows(sqlmock.NewRows([]string{"id", "title", "float", "bool"})),
			},
			want: func() []byte {
				s := &[]struct{}{}
				r, _ := json.Marshal(s)
				return r
			}(),
			wantErr: false,
		},
		{
			name: "invalid closed rows",
			args: args{
				rows: func() *sql.Rows {
					r := mockRowsToSqlRows(sqlmock.NewRows([]string{"id", "title", "float", "bool"}))
					err := r.Close()
					if err != nil {
						return nil
					}
					return r
				}(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MarshalDbRows(tt.args.rows)
			if err != nil && tt.wantErr {
				return
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalDbRows() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			eq, err := JSONBytesEqual(got, tt.want)
			if err != nil {
				t.Errorf("MarshalDbRows() error = %v, wantErr false", err)
			}
			if !eq {
				t.Errorf("MarshalDbRows() got = %s, want %s", string(got), string(tt.want))
			}
		})
	}
}

func mockRowsToSqlRows(mockRows *sqlmock.Rows) *sql.Rows {
	db, mock, _ := sqlmock.New()
	mock.ExpectQuery("select 1").WillReturnRows(mockRows)
	rows, _ := db.Query("select 1")
	return rows
}

// JSONBytesEqual compares the JSON in two byte slices.
func JSONBytesEqual(a, b []byte) (bool, error) {
	var j, j2 interface{}
	if err := json.Unmarshal(a, &j); err != nil {
		return false, err
	}
	if err := json.Unmarshal(b, &j2); err != nil {
		return false, err
	}
	return reflect.DeepEqual(j2, j), nil
}
