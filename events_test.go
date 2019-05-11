package pubsub

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"golang.org/x/xerrors"
)

func TestEventsSubscribe(t *testing.T) {
	e := NewEvents()
	_, l1 := e.Subscribe()
	_, l2 := e.Subscribe()
	defer e.Evict(l1)
	defer e.Evict(l2)
	count := e.SubscribersCount()
	require.Equal(t, 2, count, "Must be 2 subscribers")

	nowT := time.Now()
	msg := Msg{
		Scope:  "project",
		Action: "create",
		Actor: Actor{
			ID:         "project_create_event_id_1",
			Attributes: map[string]string{"id": "project_id_id"},
		},
		Time:     nowT.Unix(),
		TimeNano: nowT.UnixNano(),
	}
	e.PublishMessage(msg)
	select {
	case msgCh := <-l1:
		jmsg, ok := msgCh.(Msg)
		require.True(t, ok, "no get msg from events")
		require.Equal(t, msg.Scope, jmsg.Scope, "scope must equal")
		require.Equal(t, msg.Action, jmsg.Action, "action must equal")
	case <-time.After(1 * time.Second):
		require.NoError(t, xerrors.Errorf("Timeout waiting for broadcasted message"))
	}
	select {
	case msgCh := <-l2:
		jmsg, ok := msgCh.(Msg)
		require.True(t, ok, "no get msg from events")
		require.Equal(t, msg.Scope, jmsg.Scope, "scope must equal")
		require.Equal(t, msg.Action, jmsg.Action, "action must equal")
	case <-time.After(1 * time.Second):
		require.NoError(t, xerrors.Errorf("Timeout waiting for broadcasted message"))
	}
}

func TestEventsSubscribeTopic(t *testing.T) {
	e := NewEvents()

	noFilter := NewFilter()

	projectFilter := NewFilter()
	projectFilter.Add("scope", "project")

	projectCreateFilter := NewFilter()
	projectCreateFilter.Add("scope", "project")
	projectCreateFilter.Add("action", "create")

	_, noFilterL := e.SubscribeTopic(noFilter)
	_, projectFilterL := e.SubscribeTopic(projectFilter)
	_, projectCreateL := e.SubscribeTopic(projectCreateFilter)
	count := e.SubscribersCount()
	require.Equal(t, 3, count, "Must be 3 subscribers")

	nowT := time.Now()

	projectCreateMsg := Msg{
		Scope:  "project",
		Action: "create",
		Actor: Actor{
			ID:         "project_create_event_id_1",
			Attributes: map[string]string{"id": "project_id_id"},
		},
		Time:     nowT.Unix(),
		TimeNano: nowT.UnixNano(),
	}
	e.PublishMessage(projectCreateMsg)

	chanGetMsg := func(t *testing.T, name string, l chan interface{}, noMsg bool) {
		select {
		case msgCh := <-l:
			jmsg, ok := msgCh.(Msg)
			require.True(t, ok, "no get msg from events")
			t.Log(jmsg)
		case <-time.After(1 * time.Second):
			if noMsg {
				t.Log("name: ", name, " noMsg: ", noMsg)
				require.True(t, noMsg, "should noMsg")
			} else {
				require.NoError(t, xerrors.Errorf("Timeout waiting for broadcasted message"))
			}
		}
	}

	chanGetMsg(t, "noFilterL", noFilterL, false)
	chanGetMsg(t, "projectFilterL", projectFilterL, false)
	chanGetMsg(t, "projectCreateL", projectCreateL, false)

	projectMsg := Msg{
		Scope: "project",
		Actor: Actor{
			ID:         "project_event_id_1",
			Attributes: map[string]string{"id": "project_id_id"},
		},
		Time:     nowT.Unix(),
		TimeNano: nowT.UnixNano(),
	}
	e.PublishMessage(projectMsg)
	chanGetMsg(t, "noFilterL", noFilterL, false)
	chanGetMsg(t, "projectFilterL", projectFilterL, false)
	chanGetMsg(t, "projectCreateL", projectCreateL, true)

	companyMsg := Msg{
		Scope: "company",
		Actor: Actor{
			ID:         "project_event_id_1",
			Attributes: map[string]string{"id": "project_id_id"},
		},
		Time:     nowT.Unix(),
		TimeNano: nowT.UnixNano(),
	}
	e.PublishMessage(companyMsg)
	chanGetMsg(t, "noFilterL", noFilterL, true)
	chanGetMsg(t, "projectFilterL", projectFilterL, true)
	chanGetMsg(t, "projectCreateL", projectCreateL, true)
}
