package integration

import (
	"github.com/knowhunger/ortoo/pkg/errors"
	"github.com/knowhunger/ortoo/pkg/model"
	"github.com/knowhunger/ortoo/pkg/ortoo"
	"github.com/stretchr/testify/require"
	"sync"
)

func (its *IntegrationTestSuite) TestProtocol() {
	its.Run("Can produce an error when key is duplicated", func() {
		key := GetFunctionName()

		config := NewTestOrtooClientConfig(its.collectionName)

		client1 := ortoo.NewClient(config, "client1")

		err := client1.Connect()
		require.NoError(its.T(), err)
		defer func() {
			_ = client1.Close()
		}()

		client2 := ortoo.NewClient(config, "client2")
		err = client2.Connect()
		require.NoError(its.T(), err)
		defer func() {
			_ = client2.Close()
		}()

		_ = client1.CreateCounter(key, ortoo.NewHandlers(
			func(dt ortoo.Datatype, old model.StateOfDatatype, new model.StateOfDatatype) {
				require.Equal(its.T(), model.StateOfDatatype_DUE_TO_CREATE, old)
				require.Equal(its.T(), model.StateOfDatatype_SUBSCRIBED, new)
			}, nil,
			func(dt ortoo.Datatype, errs ...errors.OrtooError) {
				require.NoError(its.T(), errs[0])
			}))
		require.NoError(its.T(), client1.Sync())
		wg := &sync.WaitGroup{}
		wg.Add(1)
		_ = client2.CreateCounter(key, ortoo.NewHandlers(
			nil, nil,
			func(dt ortoo.Datatype, errs ...errors.OrtooError) {
				its.ctx.L().Infof("should be duplicate error:%v", errs[0])
				require.Error(its.T(), errs[0])
				wg.Done()
			}))
		require.NoError(its.T(), client2.Sync())
		wg.Wait()
	})

	its.Run("Can produce RPC error when connect", func() {
		config := NewTestOrtooClientConfig("NOT_EXISTING")
		client1 := ortoo.NewClient(config, its.getTestName())
		err := client1.Connect()
		require.Error(its.T(), err)
		defer func() {
			_ = client1.Close()
		}()
	})
}
