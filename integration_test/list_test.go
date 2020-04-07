package integration

import (
	"github.com/knowhunger/ortoo/ortoo"
	"github.com/knowhunger/ortoo/ortoo/model"
	"github.com/stretchr/testify/require"
	"time"
)

func (its *OrtooIntegrationTestSuite) TestList() {
	key := GetFunctionName()

	its.Run("Can update snapshot for list", func() {
		config := NewTestOrtooClientConfig(its.collectionName)
		client1 := ortoo.NewClient(config, "listClient")

		err := client1.Connect()
		require.NoError(its.T(), err)
		defer func() {
			_ = client1.Close()
		}()

		list1 := client1.CreateList(key, ortoo.NewHandlers(
			func(dt ortoo.Datatype, old model.StateOfDatatype, new model.StateOfDatatype) {

			},
			func(dt ortoo.Datatype, opList []interface{}) {

			},
			func(dt ortoo.Datatype, errs ...error) {

			}))
		_, _ = list1.Insert(0, "a", 2, 3.141592, time.Now())
		require.NoError(its.T(), client1.Sync())
	})
}
