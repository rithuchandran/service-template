package hotel

import (
	"encoding/json"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRegion(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := NewRepository(db)

	rowString := `{"id": "1", "name": "test", "type": "", "ancestors": null, "name_full": "", 
	"descriptor": "test region 1", 
	"Descendants": null}`
	expectedRegion := Region{Id: "1", Name: "test", Descriptor: "test region 1"}

	columns := []string{"o_data"}
	mockRows := mock.NewRows(columns).AddRow(rowString)

	mock.ExpectBegin()
	mock.ExpectQuery("select data from regions where name").WithArgs("test").WillReturnRows(mockRows)
	mock.ExpectCommit()

	region, err := repo.get("test")
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err, "Expectations not met: ", err)

	assert.Equal(t, expectedRegion, region)
}

func TestGetRegionShouldReturnTxBeginError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := NewRepository(db)

	mock.ExpectBegin().WillReturnError(errors.New("tx begin error"))

	region, err := repo.get("test")
	mockErr := mock.ExpectationsWereMet()

	assert.Nil(t, mockErr, "Expectations not met: ", err)
	assert.Equal(t, Region{}, region)
	assert.Equal(t, errors.New("tx begin error"), err)
}

func TestGetRegionShouldReturnJsonScanError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := NewRepository(db)

	rowString := `{"id": "1", "name": "test", "type": "", "ancestors": null, "name_full": "", 
	"descriptor": "test region 1", 
	"Descendants": null`

	columns := []string{"o_data"}
	mockRows := mock.NewRows(columns).AddRow(rowString)

	mock.ExpectBegin()
	mock.ExpectQuery("select data from regions where name").WithArgs("test").WillReturnRows(mockRows)

	region, err := repo.get("test")

	mockErr := mock.ExpectationsWereMet()
	assert.Nil(t, mockErr, "Expectations not met: ", err)
	assert.Equal(t, Region{}, region)
	assert.EqualErrorf(t, err, "unexpected end of JSON input", "errors not e")
}

func TestGetRegionShouldReturnTxCommitError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := NewRepository(db)

	rowString := `{"id": "1", "name": "test", "type": "", "ancestors": null, "name_full": "", 
	"descriptor": "test region 1", 
	"Descendants": null}`

	columns := []string{"o_data"}
	mockRows := mock.NewRows(columns).AddRow(rowString)

	mock.ExpectBegin()
	mock.ExpectQuery("select data from regions where name").WithArgs("test").WillReturnRows(mockRows)
	mock.ExpectCommit().WillReturnError(errors.New("tx commit error"))

	region, err := repo.get("test")

	mockErr := mock.ExpectationsWereMet()
	assert.Nil(t, mockErr, "Expectations not met: ", err)
	assert.Equal(t, Region{}, region)
	assert.Equal(t, errors.New("tx commit error"), err)
}

// TODO : Add error scenario tests for update
func TestUpdateShouldInsertRegions(t *testing.T) {
	db, mock, _ := sqlmock.New()
	repo := NewRepository(db)
	region := Region{Id: "1", Name: "test"}
	regions := Regions{"1": region}
	data, _ := json.Marshal(region)

	mock.ExpectBegin()
	mock.ExpectExec("delete from regions").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec("insert into regions").WithArgs("1", "test", data).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.update(regions)
	assert.Nil(t, err)

	err = mock.ExpectationsWereMet()
	assert.Nil(t, err, "Expectations not met: ", err)
}
