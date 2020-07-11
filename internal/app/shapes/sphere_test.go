package shapes

import (
	"github.com/eriklupander/pathtracer/internal/app/geom"
	"github.com/eriklupander/pathtracer/internal/app/identity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSphereHasIdentity(t *testing.T) {
	sphere := NewSphere()
	assert.True(t, geom.Equals(sphere.Transform, identity.Matrix))
}
func TestTranslateShape(t *testing.T) {
	s := NewSphere()
	t1 := geom.Translate(2, 3, 4)
	s.SetTransform(t1)
	assert.True(t, geom.Equals(s.Transform, t1))
}
