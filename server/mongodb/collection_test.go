package mongodb_test

import (
	"context"
	"fmt"
	"github.com/knowhunger/ortoo/commons/log"
	"github.com/knowhunger/ortoo/commons/model"
	"github.com/knowhunger/ortoo/server/mongodb"
	"github.com/knowhunger/ortoo/server/mongodb/schema"
	//"log"
	"sync"
	"testing"
	"time"
)

func TestMongo(t *testing.T) {

	conf := mongodb.NewTestMongoDBConfig("ortoo_unit_test")
	mongo, err := mongodb.New(context.TODO(), conf)
	if err != nil {
		log.Logger.Fatalf("fail to create mongoDB instance:%v", err)
	}

	t.Run("Make collections simultaneously", func(t *testing.T) {
		madeCollections := make(map[uint32]*schema.CollectionDoc)
		wg := sync.WaitGroup{}
		wg.Add(10)
		for i := 0; i < 10; i++ {
			go func(idx int) {
				collection, err := mongo.InsertCollection(context.TODO(), fmt.Sprintf("hello_%d", idx))
				if err != nil {
					t.Fail()
					return
				}
				madeCollections[collection.Num] = collection
				wg.Done()
			}(i)
		}

		wg.Wait()
		if len(madeCollections) != 10 {
			t.Fail()
		}
		for _, v := range madeCollections {
			if err := mongo.DeleteCollection(context.TODO(), v.Name); err != nil {
				t.Fail()
			}
		}
	})

	t.Run("Can get clientDoc with checkpoint", func(t *testing.T) {

	})

	t.Run("Can manipulate clientDoc", func(t *testing.T) {
		c := &schema.ClientDoc{
			CUID:          "test_cuid",
			Alias:         "test_alias",
			CollectionNum: 1,
			SyncType:      "MANUAL",
			CheckPoints: map[string]*model.CheckPoint{
				"test_duid1": model.NewCheckPoint().Set(1, 2),
				"test_duid2": model.NewCheckPoint().Set(3, 4),
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := mongo.UpdateClient(context.TODO(), c); err != nil {
			t.Fatal(err)
		}
		clientWithoutCheckPoints, err := mongo.GetClientWithoutCheckPoints(context.TODO(), c.CUID)
		if err != nil {
			t.Fatal(err)
		}

		clientWithCheckPoints, err := mongo.GetClient(context.TODO(), c.CUID)
		if err != nil {
			t.Fatal(err)
		}

		if err := mongo.DeleteClient(context.TODO(), c.CUID); err != nil {
			t.Fatal(err)
		}
		if err := mongo.DeleteClient(context.TODO(), c.CUID); err == nil {
			t.Fatal(err)
		}

		log.Logger.Infof("%+v", clientWithoutCheckPoints)
		log.Logger.Infof("%+v", clientWithCheckPoints)
	})

	t.Run("Can manipulate datatypeDoc", func(t *testing.T) {
		d := &schema.DatatypeDoc{
			DUID:          "test_duid",
			Key:           "test_key",
			CollectionNum: 1,
			Type:          "test_datatype",
			Visible:       true,
			Sseq:          0,
			CreatedAt:     time.Now(),
		}
		if err := mongo.UpdateDatatype(context.TODO(), d); err != nil {
			t.Fatal(err)
		}

		datatypeDoc1, err := mongo.GetDatatype(context.TODO(), d.DUID)
		if err != nil {
			t.Fatal(err)
		}
		log.Logger.Infof("%+v", datatypeDoc1)
		datatypeDoc2, err := mongo.GetDatatype(context.TODO(), "not exist")
		if err != nil {
			t.Fatal(err)
		}
		if datatypeDoc2 != nil {
			t.FailNow()
		}
		datatypeDoc3, err := mongo.GetDatatypeByKey(context.TODO(), d.CollectionNum, d.Key)
		if err != nil {
			t.Fatal(err)
		}
		log.Logger.Infof("%+v", datatypeDoc3)
	})

	t.Run("Can manipulate operationDoc", func(t *testing.T) {
		op, err := model.NewSnapshotOperation(
			model.TypeOfDatatype_INT_COUNTER,
			model.StateOfDatatype_DUE_TO_CREATE,
			&testSnapshot{Value: 1})
		if err != nil {
			t.Fatal(err)
		}
		var operations []interface{}
		opDoc := &schema.OperationDoc{
			ID:            "test_duid:1",
			DUID:          "test_duid",
			CollectionNum: 1,
			OpType:        "Snapshot",
			Sseq:          1,
			Operation:     op,
			CreatedAt:     time.Now(),
		}
		operations = append(operations, opDoc)

		_, err = mongo.DeleteOperation(context.TODO(), opDoc.DUID, 1)
		if err != nil {
			t.Fatal(err)
		}

		if err := mongo.InsertOperations(context.TODO(), operations); err != nil {
			t.Fatal(err)
		}

		if err := mongo.GetOperations(context.TODO(), opDoc.DUID, 1); err != nil {
			t.Fatal(err)
		}

		//deletedNum,  err := mongo.DeleteOperation(context.TODO(), opDoc.DUID, 1)
		//if err != nil || deletedNum != 1{
		//	t.Fatal(err)
		//}

	})
}

type testSnapshot struct {
	Value int32 `json:"value"`
}

func (i *testSnapshot) CloneSnapshot() model.Snapshot {
	return &testSnapshot{
		Value: i.Value,
	}
}

func (i *testSnapshot) GetTypeUrl() string {
	return "github.com/knowhunger/ortoo/common/intCounterSnapshot"
}
