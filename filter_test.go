package pubsub

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFilter_Len(t *testing.T) {
	filter := NewFilter()
	filter.Add("scope", "project")
	require.Equal(t, 1, filter.Len(), "only insert one  filter!")
}

func TestFilter_Include(t *testing.T) {
	filter := NewFilter()
	filter.Add("scope", "project")
	require.Equal(t, false, filter.Include(Msg{}))
	require.Equal(t, true, filter.Include(Msg{Scope: "project"}))
	require.Equal(t, false, filter.Include(Msg{Scope: "company"}))
	require.Equal(t, false, filter.Include(Msg{Scope: "company", Action: "creat"}))

	filter.Add("action", "create")
	require.Equal(t, false, filter.Include(Msg{}))
	require.Equal(t, false, filter.Include(Msg{Scope: "project"}))
	require.Equal(t, true, filter.Include(Msg{Scope: "project", Action: "create"}))
}

func TestFilter_Contains(t *testing.T) {
	f1 := Filter{filter: map[string]map[string]bool{}}
	require.Equal(t, false, f1.Contains("scope"))

	f2 := Filter{filter: map[string]map[string]bool{"scope": {"a": true}}}
	require.Equal(t, true, f2.Contains("scope"))

	f3 := Filter{filter: map[string]map[string]bool{"scope": {"a": false}}}
	require.Equal(t, true, f3.Contains("scope"))
}

func TestFilter_ExactMatch(t *testing.T) {
	f1 := Filter{filter: map[string]map[string]bool{}}
	require.Equal(t, true, f1.ExactMatch("scope", "a"))

	f2 := Filter{filter: map[string]map[string]bool{"scope": {"a": true}}}
	require.Equal(t, true, f2.ExactMatch("scope", "a"))

	f3 := Filter{filter: map[string]map[string]bool{"scope": {"a": false}}}
	require.Equal(t, false, f3.ExactMatch("scope", "a"))
}
