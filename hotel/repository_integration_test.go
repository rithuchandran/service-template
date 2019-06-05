// +build integration

package hotel

//This repository test is written using in built libraries. For more complicated use cases consider using test fixtures
//link: https://github.com/go-testfixtures/testfixtures

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type RepositoryIntegrationTestSuite struct {
	suite.Suite
	db *sql.DB
}

func (s *RepositoryIntegrationTestSuite) SetupSuite() {
	var err error
	db, err := sql.Open("postgres", fmt.Sprintf("host=localhost port=5432 user='postgres' "+
		"password='password' dbname=hotels-poc-test sslmode=disable"))
	if err != nil {
		panic(err)
	}
	s.db = db
}

func TestRepositoryIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryIntegrationTestSuite))
}

func (s *RepositoryIntegrationTestSuite) TestUpdateRegions() {
	repository := NewRepository(s.db)
	var obtainedRegion Region
	region1 := Region{Id: "1", Name: "first", Descriptor: "test region 1"}
	region2 := Region{Id: "2", Name: "second", Descriptor: "test region 2"}
	regions := Regions{"1": region1, "2": region2}
	repository.update(regions)

	var b []byte
	query := `select data from regions where name=$1`

	row := repository.db.QueryRow(query, "first")

	err := row.Scan(&b)
	fmt.Println(string(b))
	assert.Nil(s.T(), err)

	err = json.Unmarshal(b, &obtainedRegion)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), region1, obtainedRegion)

}

func (s *RepositoryIntegrationTestSuite) TestGetRegion() {
	repository := NewRepository(s.db)
	region1 := Region{Id: "2", Name: "first", Descriptor: "test region 1"}
	b, err := json.Marshal(region1)
	assert.Nil(s.T(), err)

	query := `insert into regions (id, name, data) values ($1, $2, $3)`
	_, err = repository.db.Exec(query, region1.Id, region1.Name, b)
	assert.Nil(s.T(), err)

	obtainedRegion, _ := repository.get("first")

	assert.Equal(s.T(), region1, obtainedRegion)
}

func (s *RepositoryIntegrationTestSuite) TearDownTest() {
	_, err := s.db.Exec(`delete from regions`)
	if err != nil {
		fmt.Println("tx exec error delete", err)
	}
}

func (s *RepositoryIntegrationTestSuite) TearDownSuite() {
	s.db.Close()
}
