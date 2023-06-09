package cmd

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/pseudonator/yap/pkg/api"
	"github.com/pseudonator/yap/pkg/cluster"
)

var createTime = time.Unix(1500000000, 0)
var startTime = time.Unix(1600000000, 0)
var clusterType = cluster.TypeMeta()
var clusterList = &api.ClusterList{
	TypeMeta: cluster.ListTypeMeta(),
	Items: []api.Cluster{
		api.Cluster{
			TypeMeta: clusterType,
			Name:     "microk8s",
			Product:  "microk8s",
			Status: api.ClusterStatus{
				CreationTimestamp: metav1.Time{Time: createTime},
				Current:           true,
			},
		},
		api.Cluster{
			TypeMeta: clusterType,
			Name:     "kind-kind",
			Product:  "KIND",
			Status: api.ClusterStatus{
				CreationTimestamp: metav1.Time{Time: createTime},
			},
		},
	},
}

func TestDefaultPrint(t *testing.T) {
	streams, _, out, _ := genericclioptions.NewTestIOStreams()
	o := NewGetOptions()
	o.IOStreams = streams
	o.StartTime = startTime

	err := o.Print(o.transformForOutput(clusterList))
	require.NoError(t, err)
	assert.Equal(t, out.String(), `CURRENT   NAME        PRODUCT    AGE
*         microk8s    microk8s   3y
          kind-kind   KIND       3y
`)
}

func TestYAML(t *testing.T) {
	streams, _, out, _ := genericclioptions.NewTestIOStreams()
	o := NewGetOptions()
	o.IOStreams = streams
	o.StartTime = startTime

	err := o.Command().Flags().Set("output", "yaml")
	require.NoError(t, err)

	err = o.Print(o.transformForOutput(clusterList))
	require.NoError(t, err)
	assert.Equal(t, `apiVersion: yap.pseudonator.io/v1alpha1
items:
- apiVersion: yap.pseudonator.io/v1alpha1
  kind: Cluster
  name: microk8s
  product: microk8s
  status:
    creationTimestamp: "2017-07-14T02:40:00Z"
    current: true
- apiVersion: yap.pseudonator.io/v1alpha1
  kind: Cluster
  name: kind-kind
  product: KIND
  status:
    creationTimestamp: "2017-07-14T02:40:00Z"
kind: ClusterList
`, out.String())
}
