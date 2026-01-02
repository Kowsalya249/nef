package context

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/free5gc/nef/internal/logger"
	"github.com/free5gc/openapi/models"
	"github.com/sirupsen/logrus"
)

type AfData struct {
	AfID       string
	NumSubscID uint64
	NumTransID uint64
	Subs       map[string]*AfSubscription
	PfdTrans   map[string]*AfPfdTransaction
<<<<<<< Updated upstream
=======
	QosSubs    map[string]*AfQosSubscription
>>>>>>> Stashed changes
	Mu         sync.RWMutex
	Log        *logrus.Entry
}

func (a *AfData) NewSub(numCorreID uint64, tiSub *models.NefTrafficInfluSub) *AfSubscription {
	a.NumSubscID++
	sub := AfSubscription{
		NotifCorreID: strconv.FormatUint(numCorreID, 10),
		SubID:        strconv.FormatUint(a.NumSubscID, 10),
		TiSub:        tiSub,
		Log:          a.Log.WithField(logger.FieldSubID, fmt.Sprintf("SUB:%d", a.NumSubscID)),
	}
	sub.Log.Infoln("New subscription")
	return &sub
}

func (a *AfData) NewPfdTrans() *AfPfdTransaction {
	a.NumTransID++
	pfdTr := AfPfdTransaction{
		TransID:   strconv.FormatUint(a.NumTransID, 10),
		ExtAppIDs: make(map[string]struct{}),
		Log:       a.Log.WithField(logger.FieldPfdTransID, fmt.Sprintf("PFDT:%d", a.NumTransID)),
	}
	pfdTr.Log.Infoln("New pfd transcation")
	return &pfdTr
}

func (a *AfData) IsAppIDExisted(appID string) (string, bool) {
	for _, pfdTrans := range a.PfdTrans {
		if _, ok := pfdTrans.ExtAppIDs[appID]; ok {
			return pfdTrans.TransID, true
		}
	}
	return "", false
}
<<<<<<< Updated upstream
=======

func (a *AfData) AddQosSubscription(sub *AfQosSubscription) {
	if a.QosSubs == nil {
		a.QosSubs = make(map[string]*AfQosSubscription)
	}
	a.QosSubs[sub.SubscriptionID] = sub
}

func (a *AfData) GetQosSubscription(subID string) (*AfQosSubscription, bool) {
	sub, ok := a.QosSubs[subID]
	return sub, ok
}

func (a *AfData) DeleteQosSubscription(subID string) {
	delete(a.QosSubs, subID)
}
>>>>>>> Stashed changes
