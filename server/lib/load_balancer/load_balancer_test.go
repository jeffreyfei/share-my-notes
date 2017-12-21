package load_balancer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/suite"
)

type LoadBalancerTestSuite struct {
	suite.Suite
	lb *LoadBalancer
}

func TestLoadBalancerTestSuite(t *testing.T) {
	s := new(LoadBalancerTestSuite)
	s.lb = NewLoadBalancer(1000)
	suite.Run(t, s)
}

func (s *LoadBalancerTestSuite) SetupTest() {
	s.lb.Providers = []provider{
		provider{
			"test-1",
			2,
		},
		provider{
			"test-2",
			1,
		},
		provider{
			"test-3",
			6,
		},
	}
}

func (s *LoadBalancerTestSuite) TestHasProvider() {
	assert.True(s.T(), s.lb.hasProvider("test-1"))
	assert.True(s.T(), s.lb.hasProvider("test-2"))
	assert.True(s.T(), s.lb.hasProvider("test-3"))
	assert.False(s.T(), s.lb.hasProvider("test-21"))
}

func (s *LoadBalancerTestSuite) TestComputeIdleIndexes() {
	providerJobs := make(map[string]int)
	providerJobs["test-1"] = 5
	providerJobs["test-2"] = 8
	providerJobs["test-3"] = 0
	s.lb.computeIdleIndexes(8, providerJobs)
	assert.Equal(s.T(), 3, s.lb.Providers[0].idleInd)
	assert.Equal(s.T(), 0, s.lb.Providers[1].idleInd)
	assert.Equal(s.T(), 8, s.lb.Providers[2].idleInd)
}

func (s *LoadBalancerTestSuite) TestGetNextProvider() {
	provider, err := s.lb.getNextProvider()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "test-1", provider)
	provider, err = s.lb.getNextProvider()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "test-1", provider)
	provider, err = s.lb.getNextProvider()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "test-1", provider)
	provider, err = s.lb.getNextProvider()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "test-2", provider)
	provider, err = s.lb.getNextProvider()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "test-2", provider)
	provider, err = s.lb.getNextProvider()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "test-3", provider)
}
