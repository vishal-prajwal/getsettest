package consul_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/pkg/errors"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"bitbucket.org/junglee_games/getsetgo/pandora/consul"
)

var errCallbackDone = errors.New("the callback was terminated as expected")

func TestSingleStartLeaderElection(t *testing.T) {
	ctx := context.Background()
	client, err := consul.NewClient(consul.ConsulConfig{})
	require.NoError(t, err)

	err = client.StartLeaderElection(ctx, "StartLeaderElection", func(ctx context.Context) error {
		return errCallbackDone
	})

	assert.Equal(t, errCallbackDone, err)
}

func TestMultiStartLeaderElection(t *testing.T) {
	ctx := context.Background()
	client, err := consul.NewClient(consul.ConsulConfig{})
	require.NoError(t, err)

	// using wg gives `panic: sync: negative WaitGroup counter`. Even though there's wg.Add(1) and wg.Done() once?
	t.Run("instance 1 becomes a leader while instance 2 is blocked (normal instance)", func(t *testing.T) {
		var instance1 uint32
		var instance2 uint32

		go func() {
			err = client.StartLeaderElection(ctx, "MultiStartLeaderElection1", func(ctx context.Context) error {
				atomic.AddUint32(&instance1, 1)

				time.Sleep(4 * time.Second)

				return nil
			})
			require.NoError(t, err)
		}()

		time.Sleep(1 * time.Second)

		go func() {
			err = client.StartLeaderElection(ctx, "MultiStartLeaderElection1", func(ctx context.Context) error {
				atomic.AddUint32(&instance2, 1)

				return nil
			})

			require.NoError(t, err)
		}()

		assert.Eventually(t, func() bool {
			return atomic.LoadUint32(&instance1) >= uint32(1) && atomic.LoadUint32(&instance2) == uint32(0)
		}, time.Second*3, time.Millisecond*400)
	})

	t.Run("instance 1 becomes a leader then stops; while instance 2 acquire the leadership", func(t *testing.T) {
		var instance1 uint32
		var instance2 uint32

		go func() {
			err = client.StartLeaderElection(ctx, "MultiStartLeaderElection2", func(ctx context.Context) error {
				atomic.AddUint32(&instance1, 1)

				if instance1 == 1 {
					return errCallbackDone
				}
				return errCallbackDone
			})

			require.Equal(t, errCallbackDone, err)
		}()

		time.Sleep(1 * time.Second)

		assert.Eventually(t, func() bool {
			return atomic.LoadUint32(&instance1) == uint32(1) && atomic.LoadUint32(&instance2) == uint32(0)
		}, time.Second*2, time.Millisecond*400)

		go func() {
			err = client.StartLeaderElection(ctx, "MultiStartLeaderElection2", func(ctx context.Context) error {
				atomic.AddUint32(&instance2, 1)

				if instance2 == 1 {
					return errCallbackDone
				}
				return errCallbackDone
			})

			require.Equal(t, errCallbackDone, err)
		}()
		time.Sleep(3 * time.Second)

		assert.Eventually(t, func() bool {
			return atomic.LoadUint32(&instance1) == uint32(1) && atomic.LoadUint32(&instance2) == uint32(1)
		}, time.Second*5, time.Millisecond*100)
	})

	t.Run("Handling context cancelling (Process is usually shutting down)", func(t *testing.T) {
		ctx1, cancel := context.WithDeadline(ctx, time.Now().Add(500*time.Millisecond))
		defer cancel()

		var instance1 uint32
		var instance2 uint32

		go func() {
			err := client.StartLeaderElection(ctx1, "MultiStartLeaderElection4", func(ctx1 context.Context) error {
				atomic.AddUint32(&instance1, 1)
				// Simulating a long lived task, but will terminate when context Done
				timer := time.NewTimer(10 * time.Second)
				select {
				case <-timer.C:
					return errors.New("error, context timeout")
				case <-ctx1.Done():
					return errCallbackDone
				}
			})

			assert.Equal(t, errCallbackDone, err)
		}()

		time.Sleep(1 * time.Second)

		go func() {
			err := client.StartLeaderElection(ctx, "MultiStartLeaderElection4", func(ctx context.Context) error {
				atomic.AddUint32(&instance2, 1)

				if instance2 == 1 {
					return errCallbackDone
				}
				return errCallbackDone
			})

			assert.Equal(t, errCallbackDone, err)
		}()

		time.Sleep(2 * time.Second)

		assert.Eventually(t, func() bool {
			return atomic.LoadUint32(&instance1) == uint32(1) && atomic.LoadUint32(&instance2) == uint32(1)
		}, time.Second*7, time.Millisecond*500)
	})

	t.Run("instance 1 and instance 2 takes turn acquire the leadership", func(t *testing.T) {
		var instance1 uint32
		var instance2 uint32

		go func() {
			err := client.StartLeaderElection(ctx, "MultiStartLeaderElection5", func(ctx context.Context) error {
				atomic.AddUint32(&instance1, 1)
				if instance1 == 4 {
					return errCallbackDone
				}
				return nil
			})
			assert.Equal(t, errCallbackDone, err)
		}()

		go func() {
			err := client.StartLeaderElection(ctx, "MultiStartLeaderElection5", func(ctx context.Context) error {
				atomic.AddUint32(&instance2, 1)
				if instance2 == 4 {
					return errCallbackDone
				}
				return nil
			})
			assert.Equal(t, errCallbackDone, err)
		}()
		assert.Eventually(t, func() bool {
			return atomic.LoadUint32(&instance1) == 4 && atomic.LoadUint32(&instance2) == 4
		}, time.Second*10, time.Millisecond*100)
	})
}

func TestSingleStartLeaderElectionPrometheus(t *testing.T) {
	ctx := context.Background()
	client, err := consul.NewClient(consul.ConsulConfig{})
	require.NoError(t, err)

	serviceName := "StartLeaderElectionPrometheus"

	consulPrometheusInstanceCounter(t, serviceName, consul.Leader, 0)
	consulPrometheusInstanceCounter(t, serviceName, consul.NonLeader, 0)

	t.Run("prometheus test : check before being leader, check during leadership, check after loses leader role", func(t *testing.T) {
		err = client.StartLeaderElection(ctx, serviceName, func(ctx context.Context) error {
			consulPrometheusInstanceCounter(t, serviceName, consul.Leader, 1)
			consulPrometheusInstanceCounter(t, serviceName, consul.NonLeader, 0)

			return errCallbackDone
		})

		assert.Equal(t, errCallbackDone, err)

		consulPrometheusInstanceCounter(t, serviceName, consul.Leader, 0)
		consulPrometheusInstanceCounter(t, serviceName, consul.NonLeader, 0)
	})
}

// Util function to count prometheus gauge based on servicename&instancetype
func consulPrometheusInstanceCounter(t testing.TB, serviceName, instanceType string, want int) {
	t.Helper()

	res := &io_prometheus_client.Metric{}

	err := consul.ServiceInstancesGauge.WithLabelValues(serviceName, instanceType).Write(res)
	instanceOfService := res.Gauge.Value

	assert.NoError(t, err)
	assert.Equal(t, float64(want), *instanceOfService)
}
